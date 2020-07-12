package ventiutils

import (
	"bytes"
	"encoding/binary"

	"github.com/cespedes/venti"
)

func getU8(b *[]byte) uint8 {
	r := (*b)[0]
	*b = (*b)[1:]
	return r
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

func getU48(b *[]byte) uint64 {
	c := make([]byte, 8)
	copy(c[2:8], *b)
	r := binary.BigEndian.Uint64(c)
	*b = (*b)[U48Size:]
	return r
}

func getString(b *[]byte, maxsize int) string {
	var r string
	if c := bytes.IndexByte((*b)[:maxsize], 0); c >= 0 {
		r = string((*b)[:c])
	} else {
		r = string((*b)[:maxsize])
	}
	*b = (*b)[maxsize:]
	return r
}

func getScore(b *[]byte) venti.Score {
	v := venti.Score((*b)[:VtScoreSize])
	*b = (*b)[VtScoreSize:]
	return v
}
