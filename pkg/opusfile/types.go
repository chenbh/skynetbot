package opusfile

const (
	oggSig         = "OggS"
	opusIdSig      = "OpusHead"
	opusCommentSig = "OpusTags"
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

type OpusPacket []byte
