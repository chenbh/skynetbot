package opusfile

import (
	"fmt"
	"io"
)

type OpusReader interface {
	NextPacket() ([]byte, error)
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

func (r *opusReader) NextFrame() ([]byte, error) {
	_, err := r.NextPacket()
	if err != nil {
		return nil, nil
	}

	return nil, nil
}

func (r *opusReader) NextPacket() ([]byte, error) {
	if r.currentPage == nil {
		page, err := r.ogg.NextPage()
		if err != nil {
			return nil, fmt.Errorf("reading ogg page: %v", err)
		}
		r.currentPage = &page
	}

	var buf []byte
	done := false
	for !done {
		// sanity + in case last segment was 0xff from previous page
		if r.currentPage.SegmentTable[r.currentSegment] == 0 {
			done = true
		}

		buf = append(buf, r.currentPage.Segments[r.currentSegment]...)
		r.currentSegment++

		// not a multi-segment packet -> we're done
		if r.currentPage.SegmentTable[r.currentSegment-1] != 0xff {
			done = true
		}

		// multi-page spanning packet
		if r.currentSegment >= r.currentPage.PageSegments {
			page, err := r.ogg.NextPage()
			if err != nil {
				return nil, fmt.Errorf("reading ogg page: %v", err)
			}
			r.currentPage = &page
			r.currentSegment = 0
		}
	}
	return buf, nil
}
