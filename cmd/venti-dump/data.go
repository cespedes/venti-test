package main

import (
	"github.com/cespedes/venti"
)

var vtzeroscore = venti.Score{
        0xda, 0x39, 0xa3, 0xee, 0x5e, 0x6b, 0x4b, 0x0d, 0x32, 0x55,
        0xbf, 0xef, 0x95, 0x60, 0x18, 0x90, 0xaf, 0xd8, 0x07, 0x09,
}

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
