package audio

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
	"github.com/chenbh/skynetbot/pkg/opusfile"
)

func (s *state) clip(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	for k, c := range clips {
		f, err := os.Create(fmt.Sprintf("%v.opus", k))
		if err != nil {
			return err
		}

		out, err := opusfile.NewOpusWriter(f)
		if err != nil {
			return err
		}
		defer out.Close()

		fmt.Printf("writting %v packets\n", len(c.packets))
		for _, p := range c.packets {
			out.WritePacket([][]byte{p.Opus}, uint64(p.Timestamp))
		}
		fmt.Println("done")
	}
	return nil
}

func record(vc *discordgo.VoiceConnection, done <-chan struct{}) {
	for {
		select {
		case p := <-vc.OpusRecv:
			ssrc := int(p.SSRC)

			if c, found := clips[ssrc]; found {
				c.record(p)
			}
		case <-done:
			return
		}
	}
}

func voiceHandler(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
	fmt.Println(vs)
	if vs.Speaking == true {
		fmt.Printf("generated clip for %v\n", vs.UserID)
		clips[vs.SSRC] = &clip{
			user:    vs.UserID,
			packets: make([]*discordgo.Packet, 0),
		}
	}
}
