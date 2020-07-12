package ventiutils

import (
	"github.com/cespedes/venti"
)

var VtZeroScore = venti.Score{
        0xda, 0x39, 0xa3, 0xee, 0x5e, 0x6b, 0x4b, 0x0d, 0x32, 0x55,
        0xbf, 0xef, 0x95, 0x60, 0x18, 0x90, 0xaf, 0xd8, 0x07, 0x09,
}

// Useful constants:
const (
        U8Size  = 1
        U16Size = 2
        U32Size = 4
        U48Size = 6
        U64Size = 8

	VtMaxLumpSize = 65535

        VtScoreSize = 20
        VtEntrySize = 40

        VtRootSize = 300
        VtRootVersion = 2
        vtRootVersionBig = 1<<15

	IBucketSize = U32Size + U16Size
	IEntrySize  = U64Size + U32Size + 2*U16Size + 2*U8Size + VtScoreSize

	PartBlank  = 256 * 1024 // untouched section at beginning of partition
	ANameSize  = 64         // maximum length for names
	HeadSize   = 512        // size of a header after PartBlank
	IsectMagic = 0xd15c5ec7 // magic number in isects, at the beginning of header
)

type VtRoot struct {
	Name      string
	Type      string
	Score     venti.Score
	BlockSize uint16
	PrevScore venti.Score
}

type VtEntry struct {
	Gen   uint32
	PSize uint16
	DSize uint16
	Type  venti.Type
	Flags byte
	Size  uint64
	Score venti.Score
}

type Isect struct {
	// Fields stored on disk:
	version     uint32
	name        string // text label
	index       string // index owning the section
	blocksize   uint32 // size of hash buckets in index
	blockbase   uint32 // address of start of on disk index table
	blocks      uint32 // total blocks on disk; some may be unused
	start       uint32 // first bucket in this section
	stop        uint32 // limit of buckets in this section
	bucketmagic uint32

	// Computed values:
	blocklog int    // log2(blocksize)
	buckmax  int    // max. entries in a index bucket
	tabbase  uint32 // base address of index config table on disk
	tabsize  uint32 // max. bytes in index config
}
