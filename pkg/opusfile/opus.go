package opusfile

import (
	"fmt"
	"io"
)

type OpusPacket struct {
}

type OpusReader interface {
	NextPacket() (OpusPacket, error)
	NextPacketRaw() ([]byte, error)
	NextFrame() ([]byte, error)
}

type opusReader struct {
	ogg OggReader

	currentPage    *OggPage
	currentSegment uint8
}

func NewOpusReader(in io.Reader) (OpusReader, error) {
	oggReader := NewOggReader(in)

	_, err := oggReader.NextPage()
	if err != nil {
		return nil, fmt.Errorf("invalid id header: %v", err)
	}
	_, err = oggReader.NextPage()
	if err != nil {
		return nil, fmt.Errorf("invalid comment header: %v", err)
	}

	// TODO: verify ID header and comment header

	return &opusReader{
		ogg: oggReader,
	}, nil
}

func (o *opusReader) NextFrame() ([]byte, error) {
	_, err := o.NextPacketRaw()
	if err != nil {
		return nil, nil
	}

	return nil, nil
}

func (o *opusReader) NextPacket() (OpusPacket, error) {
	b, err := o.NextPacketRaw()
	if err != nil {
		return OpusPacket{}, nil
	}

	return parsePacket(b)
}

func (o *opusReader) NextPacketRaw() ([]byte, error) {
	if o.currentPage == nil {
		page, err := o.ogg.NextPage()
		if err != nil {
			return nil, fmt.Errorf("reading ogg page: %v", err)
		}
		o.currentPage = &page
	}

	var buf []byte
	done := false
	for !done {
		// sanity + in case last segment was 0xff from previous page
		if o.currentPage.SegmentTable[o.currentSegment] == 0 {
			done = true
		}

		buf = append(buf, o.currentPage.Segments[o.currentSegment]...)
		o.currentSegment++

		// not a multi-segment packet -> we're done
		if o.currentPage.SegmentTable[o.currentSegment-1] != 0xff {
			done = true
		}

		// multi-page spanning packet
		if o.currentSegment >= o.currentPage.PageSegments {
			page, err := o.ogg.NextPage()
			if err != nil {
				return nil, fmt.Errorf("reading ogg page: %v", err)
			}
			o.currentPage = &page
			o.currentSegment = 0
		}
	}
	return buf, nil
}

// TODO: do we really care about parsing the packet?
func parsePacket([]byte) (OpusPacket, error) {
	return OpusPacket{}, nil
}
