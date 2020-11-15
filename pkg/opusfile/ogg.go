package opusfile

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type OggHeader struct {
	Version     uint8
	IsContinued bool
	IsFirstPage bool
	IsLastPage  bool

	GranulePosition uint64
	BitstreamSerial uint32
	PageSequence    uint32
	CrcChecksum     uint32

	PageSegments uint8
	SegmentTable []uint8
}

type OggPage struct {
	OggHeader

	Segments [][]byte

	// Size of all segments in bytes
	SegmentTotal int
}

func (p OggPage) Bytes(includeChecksum bool) []byte {
	totalSize := 27 + int(p.PageSegments) + p.SegmentTotal
	buf := bytes.NewBuffer(make([]byte, 0, totalSize))
	var b []byte

	buf.WriteString("OggS")
	buf.WriteByte(byte(p.Version))

	headerType := uint8(0x0)
	if p.IsContinued {
		headerType = headerType | 0x1
	}
	if p.IsFirstPage {
		headerType = headerType | 0x2
	}
	if p.IsLastPage {
		headerType = headerType | 0x4
	}
	buf.WriteByte(byte(headerType))

	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, p.GranulePosition)
	buf.Write(b)

	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, p.BitstreamSerial)
	buf.Write(b)

	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, p.PageSequence)
	buf.Write(b)

	b = make([]byte, 4)
	if includeChecksum {
		binary.LittleEndian.PutUint32(b, p.CrcChecksum)
	}
	buf.Write(b)

	buf.WriteByte(byte(p.PageSegments))
	for _, s := range p.Segments {
		buf.Write(s)
	}

	return buf.Bytes()
}

type OggReader interface {
	NextPage() (OggPage, error)
}

type oggReader struct {
	r io.Reader
}

func NewOggReader(in io.Reader) OggReader {
	// TODO: verify that it at least has the OggS header
	return &oggReader{
		r: in,
	}
}

func (o *oggReader) NextPage() (OggPage, error) {
	header, err := o.parseHeader()
	if err != nil {
		return OggPage{}, err
	}

	var totalBytes int
	for _, s := range header.SegmentTable {
		totalBytes += int(s)
	}

	buf := make([]byte, totalBytes)
	n, err := o.r.Read(buf)
	if err != nil {
		return OggPage{}, err
	}
	if n != int(totalBytes) {
		return OggPage{}, errors.New("invalid file")
	}

	segments := make([][]byte, header.PageSegments)
	var idx int
	for i, s := range header.SegmentTable {
		segments[i] = buf[idx : idx+int(s)]
		idx += int(s)
	}

	// TODO: verify crc

	return OggPage{
		OggHeader:    header,
		Segments:     segments,
		SegmentTotal: totalBytes,
	}, nil
}

// https://tools.ietf.org/html/rfc3533#section-6
func (o *oggReader) parseHeader() (OggHeader, error) {
	header := make([]byte, 27)
	n, err := o.r.Read(header)
	if err != nil {
		return OggHeader{}, err
	}
	if n != 27 {
		return OggHeader{}, errors.New("invalid file")
	}

	magicNumber := header[0:4]
	version := uint8(header[4])
	headerType := header[5]
	if string(magicNumber) != "OggS" {
		return OggHeader{}, errors.New("invalid header")
	}

	granulePosition := binary.LittleEndian.Uint64(header[6:14])
	bitstreamSerial := binary.LittleEndian.Uint32(header[14:18])
	pageSequence := binary.LittleEndian.Uint32(header[18:22])
	checksum := binary.LittleEndian.Uint32(header[22:26])

	pageSegments := uint8(header[26])
	segmentTable := make([]uint8, pageSegments)
	n, err = o.r.Read(segmentTable)
	if err != nil {
		return OggHeader{}, err
	}
	if n != int(pageSegments) {
		return OggHeader{}, errors.New("invalid file")
	}

	return OggHeader{
		Version:     version,
		IsContinued: headerType&0x1 == 1,
		IsFirstPage: headerType&0x2 == 1,
		IsLastPage:  headerType&0x4 == 1,

		GranulePosition: granulePosition,
		BitstreamSerial: bitstreamSerial,
		PageSequence:    pageSequence,
		CrcChecksum:     checksum,

		PageSegments: pageSegments,
		SegmentTable: segmentTable,
	}, nil
}
