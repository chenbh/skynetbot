package audio

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
	"github.com/chenbh/skynetbot/pkg/opusfile"
)

var (
	filler         = []byte{0xfc, 0xff, 0xfe} // 20ms of Opus encoded silence
	recordDuration = 60
)

type clip struct {
	mu        *sync.Mutex
	startTime uint64 // timestamp (unix) for start of this ssrc

	packets []*discordgo.Packet
}

func (c *clip) record(p *discordgo.Packet) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.packets) == 0 {
		c.startTime = uint64(time.Now().Unix())
	}
	c.packets = append(c.packets, p)
}

func (s *state) clip(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	s.trimOldAudio()

	files := make([]*os.File, 0)
	for k, c := range s.recordings {
		if len(c.packets) == 0 { // don't clip people who haven't talked
			continue
		}

		user, err := b.Session().User(k)
		if err != nil {
			return err
		}
		f, err := os.Create(path.Join(os.TempDir(), user.Username+".opus"))
		if err != nil {
			return err
		}
		out, err := opusfile.NewOpusWriter(f)
		if err != nil {
			return err
		}

		files = append(files, f)

		fmt.Printf("writting %v packets\n", len(c.packets))
		writePackets(out, c)
		out.Finish()
		fmt.Println("done")
	}

	attachments := make([]*discordgo.File, len(files))
	for i, f := range files {
		f.Seek(0, 0)
		attachments[i] = &discordgo.File{
			Name:        path.Base(f.Name()),
			ContentType: "audio/opus",
			Reader:      f,
		}
	}

	_, err := b.Session().ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("last %v seconds of audio:", recordDuration),
		Files:   attachments,
	})
	fmt.Println(err)

	for _, f := range files {
		f.Close()
	}

	return err
}

func writePackets(out opusfile.OpusWriter, c *clip) {
	lastTs := uint64(c.packets[0].Timestamp) // use the rtp timestamp of the first packet as base
	for i, p := range c.packets {
		if i > 0 {
			delta := uint64(p.Timestamp) - lastTs - 960 // 48 khz * 20ms = 960 samples per packet

			// backfill missing packets with silence
			for j := uint64(960); j <= delta; j += 960 {
				out.WritePacket([][]byte{filler}, uint64(lastTs+j))
			}
		}

		lastTs = uint64(p.Timestamp)
		out.WritePacket([][]byte{p.Opus}, uint64(p.Timestamp))
	}
}

func (s *state) receiveAudio() {
	for {
		if s.vc == nil {
			return
		}

		select {
		case p := <-s.vc.OpusRecv:
			user, foundUser := s.ssrc[int(p.SSRC)]
			if foundUser {
				if c, foundClip := s.recordings[user]; foundClip {
					c.record(p)
				}
			}
		case _, closed := <-s.doneRecording:
			if closed {
				return
			}
		}
	}
}

func (s *state) voiceHandler() discordgo.VoiceSpeakingUpdateHandler {
	return func(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
		// let's not listen to ourselves
		if vs.UserID == s.vc.UserID {
			return
		}

		if vs.Speaking == true {
			now := uint64(time.Now().Unix())
			pkts := make([]*discordgo.Packet, 0)

			// storing packets per user instead of ssrc makes it nicer for us later on
			// at the expense of having to map ssrc -> user everytime a packet comes in
			s.ssrc[vs.SSRC] = vs.UserID
			if c, found := s.recordings[vs.UserID]; !found {
				s.recordings[vs.UserID] = &clip{
					mu:        &sync.Mutex{},
					startTime: now,
					packets:   pkts,
				}
			} else {
				c.startTime = now
				c.packets = pkts
			}
		} else {
			// TODO: does this capture disconnect events?
			fmt.Printf("%v stopped speaking\n", vs.UserID)
		}
	}
}

func (s *state) trimOldAudio() {
	fmt.Println("starting gc")

	now := time.Now()
	for _, c := range s.recordings {
		c.mu.Lock()
		defer c.mu.Unlock()

		// TODO timestamp is generated randomly per ssrc and increases at rate of 960
		// per packet (20ms @48khz)
		if expired(c.startTime, now) {
			if len(c.packets) == 0 {
				continue
			}

			idx := 0
			base := c.packets[0].Timestamp
			var duration uint64

			fmt.Println(c.startTime, base)

			// figure out index of first non-expired packet
			for i, p := range c.packets {
				sampleDelta := p.Timestamp - base        // how many samples has elapsed since the start
				timeDelta := uint64(sampleDelta / 48000) // 48khz = 48000 samples / second

				if expired(c.startTime+timeDelta, now) {
					idx = i
					duration = timeDelta
				}
			}

			fmt.Printf("trimmed %v packets\n", idx)
			c.packets = c.packets[idx+1:]
			if len(c.packets) != 0 {
				c.startTime = c.startTime + duration
			}

			fmt.Printf("%v packets left\n", len(c.packets))
			fmt.Println(c.startTime)
		}
	}
	fmt.Println("done gc")
}

func (s *state) gcAudio() {
	timer := time.NewTicker(30 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			s.trimOldAudio()
		case _, closed := <-s.doneRecording:
			if closed {
				return
			}
		}
	}
}

func expired(ts uint64, now time.Time) bool {
	return ts < uint64(now.Add(-60*time.Second).Unix())
}
