package ventiutils

import (
	"bytes"
	"fmt"
	"io"

	"github.com/cespedes/venti"
)

var VtZeroScore = venti.Score{
	0xda, 0x39, 0xa3, 0xee, 0x5e, 0x6b, 0x4b, 0x0d, 0x32, 0x55,
	0xbf, 0xef, 0x95, 0x60, 0x18, 0x90, 0xaf, 0xd8, 0x07, 0x09,
}

var vtZeroScore2 = venti.Score(make([]byte, 20))

// Useful constants:
const (
	U8Size  = 1
	U16Size = 2
	U32Size = 4
	U48Size = 6
	U64Size = 8

	VtMaxLumpSize = 65535

	VtScoreSize = 20

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

type Data struct {
	client *venti.Client
	score  venti.Score
	t      venti.Type
	size   uint64    // total size of data
	offset uint64    // next offset to read from
	dsize  uint      // size of VtData blocks
	psize  uint      // size of VtData+N blocks
	blocks [8][]byte // cache of already downloaded blocks
	pos    [8]uint64 // number of block downloaded in "blocks"
}

func readBlock(b []byte, c *venti.Client, score venti.Score, t venti.Type) error {
	if bytes.Equal([]byte(score), []byte(VtZeroScore)) || bytes.Equal([]byte(score), []byte(vtZeroScore2)) {
		for i := range b {
			b[i] = 0
		}
		return nil
	}
	r, err := c.Read(t, score, VtMaxLumpSize)
	if err != nil {
		return err
	}
	defer r.Close()

	n := 0
	for n < len(b) && err == nil {
		var nn int
		nn, err = r.Read(b[n:])
		n += nn
	}
	if n < len(b) && err != io.EOF {
		return err
	}
	for i := n; i < len(b); i++ {
		b[i] = 0
	}
	return nil

}

func OpenData(c *venti.Client, s venti.Score, t venti.Type, size uint64, dsize uint, psize uint) (*Data, error) {
	if c == nil {
		return nil, fmt.Errorf("OpenData: nil client")
	}
	if t < venti.VtData || t > venti.VtData+7 {
		return nil, fmt.Errorf("OpenData: invalid type %d", t)
	}
	if dsize == 0 {
		return nil, fmt.Errorf("OpenData: dsize cannot be zero")
	}
	if psize == 0 {
		return nil, fmt.Errorf("OpenData: psize cannot be zero")
	}
	var d Data
	d.client = c
	d.score = s
	d.t = t
	d.size = size
	d.offset = 0
	d.dsize = dsize
	d.psize = psize

	for i := 1; i <= int(t-venti.VtData); i++ {
		d.blocks[i] = make([]byte, psize)
	}
	d.blocks[0] = make([]byte, dsize)

	if err := readBlock(d.blocks[int(t-venti.VtData)], c, s, t); err != nil {
		return nil, err
	}
	for i := 0; i < 8; i++ {
		d.pos[i] = 1<<64 - 1
	}
	d.getBlocks()

	return &d, nil
}

func (d *Data) getBlocks() error {
	var offset = d.offset
	var newpos [8]uint64
	// newpos[0] = offset % uint64(d.dsize)
	offset /= uint64(d.dsize)
	for i := 0; i < int(d.t-venti.VtData)-1; i++ {
		newpos[i] = offset % uint64(d.psize)
		offset /= uint64(d.psize)
	}
	if d.t > venti.VtData {
		newpos[int(d.t-venti.VtData)] = offset
	}
	// newpos contains the blocks I should have.  We have to download them if they are different from d.pos
	for i := int(d.t-venti.VtData) - 1; i >= 0; i-- {
		if d.pos[i] != newpos[i] {
			b := d.blocks[i+1][newpos[i]:]
			score := getScore(&b)

			if err := readBlock(d.blocks[i], d.client, score, venti.Type(d.t)+venti.VtData); err != nil {
				return err
			}
			d.pos[i] = newpos[i]
		}
	}

	return nil
}

func (d *Data) Read(p []byte) (n int, err error) {
	d.getBlocks()
	if d.offset >= d.size {
		return 0, io.EOF
	}
	off := d.offset % uint64(d.dsize)
	end := uint64(d.dsize)
	if d.size-d.offset < end-off {
		end = off + d.size - d.offset
	}
	n = copy(p, d.blocks[0][off:end])
	d.offset += uint64(n)
	return n, nil
}

func (d *Data) Close() error {
	return nil
}
