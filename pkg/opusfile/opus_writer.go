package opusfile

import (
	"encoding/binary"
	"io"
)

type OpusWriter interface {
	WritePacket(packet [][]byte, timestamp uint64) error
	Close() error
}

type opusWriter struct {
	ogg OggWriter

	pageIndex        uint32
	prevPagePosition uint64
	prevTimestamp    uint64
}

func NewOpusWriter(out io.WriteCloser) (OpusWriter, error) {
	oggWriter := NewOggWriter(out)

	writer := &opusWriter{
		ogg: oggWriter,

		prevPagePosition: 1,
	}

	err := writer.writeHeaders()
	if err != nil {
		return nil, err
	}

	return writer, nil
}

func (w *opusWriter) writeHeaders() error {
	idHeader := make([]byte, 19)
	copy(idHeader[0:], opusIdSig)
	idHeader[8] = 1
	idHeader[9] = 2

	binary.LittleEndian.PutUint16(idHeader[10:], 3840)  // pre-skip
	binary.LittleEndian.PutUint32(idHeader[12:], 48000) // sample rate
	binary.LittleEndian.PutUint16(idHeader[16:], 0)     // output gain
	idHeader[18] = 0                                    // mono or stereo

	idPage := w.ogg.NewPage([][]byte{idHeader}, 0, w.pageIndex)
	idPage.IsFirstPage = true
	err := w.ogg.WritePage(idPage)
	if err != nil {
		return err
	}
	w.pageIndex++

	commentHeader := make([]byte, 25)
	copy(commentHeader[0:], opusCommentSig)
	binary.LittleEndian.PutUint32(commentHeader[8:], 9)  // vendor name length
	copy(commentHeader[12:], "skynetbot")                // vendor name
	binary.LittleEndian.PutUint32(commentHeader[21:], 0) // comment list Length

	commentPage := w.ogg.NewPage([][]byte{commentHeader}, 0, w.pageIndex)
	err = w.ogg.WritePage(commentPage)
	if err == nil {
		w.pageIndex++
	}
	return err
}

func (w *opusWriter) WritePacket(p [][]byte, timestamp uint64) error {
	if w.prevTimestamp != 0 {
		w.prevPagePosition += timestamp - w.prevTimestamp
	}
	w.prevTimestamp = timestamp
	page := w.ogg.NewPage(p, w.prevPagePosition, w.pageIndex)
	w.pageIndex++

	return w.ogg.WritePage(page)
}

func (w *opusWriter) Close() error {
	return w.ogg.Close(w.prevPagePosition, w.pageIndex)
}
