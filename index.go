package ventiutils

import (
	"fmt"
	"os"
)

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
	isect.name = getString(&b, ANameSize)
	isect.index = getString(&b, ANameSize)
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
