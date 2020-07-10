package venti

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// Useful constants:
const (
	VtScoreSize = 20

	U8Size  = 1
	U16Size = 2
	U32Size = 4
	U64Size = 8

	IBucketSize = U32Size + U16Size
	IEntrySize  = U64Size + U32Size + 2*U16Size + 2*U8Size + VtScoreSize

	PartBlank  = 256 * 1024 // untouched section at beginning of partition
	ANameSize  = 64         // maximum length for names
	HeadSize   = 512        // size of a header after PartBlank
	IsectMagic = 0xd15c5ec7 // magic number in isects, at the beginning of header
)

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

func getU16(b *[]byte) uint16 {
	r := binary.BigEndian.Uint16((*b)[:U16Size])
	*b = (*b)[U16Size:]
	return r
}

func getU32(b *[]byte) uint32 {
	r := binary.BigEndian.Uint32((*b)[:U32Size])
	*b = (*b)[U32Size:]
	return r
}

func getString(b *[]byte) string {
	var r string
	if c := bytes.IndexByte((*b)[:ANameSize], 0); c >= 0 {
		r = string((*b)[:c])
	} else {
		r = string((*b)[:ANameSize])
	}
	*b = (*b)[ANameSize:]
	return r
}

// u64log2 returns floor(log2(v))
func u64log2(v uint64) int {
	var i int
	for i = 0; i < 64; i++ {
		if (v >> i) <= 1 {
			break
		}
	}
	return i
}

func ParseIsect(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	b := make([]byte, HeadSize)
	n, err := file.ReadAt(b, PartBlank)
	if err != nil || n != HeadSize {
		panic(err)
	}
	if m := getU32(&b); m != IsectMagic {
		panic(fmt.Sprintf("index section has wrong magic number: 0x%08x (expected 0x%08x)", m, IsectMagic))
	}
	var isect Isect
	isect.version = getU32(&b)
	isect.name = getString(&b)
	isect.index = getString(&b)
	isect.blocksize = getU32(&b)
	isect.blockbase = getU32(&b)
	isect.blocks = getU32(&b)
	isect.start = getU32(&b)
	isect.stop = getU32(&b)
	if isect.version == 2 {
		isect.bucketmagic = getU32(&b)
	}

	fmt.Printf("version = %d\n", isect.version)
	fmt.Printf("name = %q\n", isect.name)
	fmt.Printf("index = %q\n", isect.index)
	fmt.Printf("blocksize = %d\n", isect.blocksize)
	fmt.Printf("blockbase = %d\n", isect.blockbase)
	fmt.Printf("blocks = %d\n", isect.blocks)
	fmt.Printf("start = %d\n", isect.start)
	fmt.Printf("stop = %d\n", isect.stop)
	fmt.Printf("bucketmagic = 0x%08x\n", isect.bucketmagic)
	var i uint32
	for i = 0; i < isect.blocks; i++ {
		// IBucket:
		b1 := make([]byte, IBucketSize)
		n, err := file.ReadAt(b1, int64(isect.blockbase+i*isect.blocksize))
		if err != nil || n != IBucketSize {
			panic(err)
		}
		len := getU16(&b1)
		magic := getU32(&b1)
		if isect.version == 2 && magic != isect.bucketmagic {
			len = 0
		}
		fmt.Printf("bucket %d: %d entries\n", i, len)
	}
}
