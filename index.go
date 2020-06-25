package venti

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	// PartBlank is the untouched section at beginning of partition
	PartBlank = 256 * 1024

	// HeadSize is the size of a header after PartBlank
	HeadSize   = 512
	IsectMagic = 0xd15c5ec7
)

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
	if m := binary.BigEndian.Uint32(b[0:4]); m != IsectMagic {
		panic(fmt.Sprintf("index section has wrong magic number: 0x%08x (expected 0x%08x)", m, IsectMagic))
	}
}
