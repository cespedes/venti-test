package ventiutils

import (
	"fmt"
	"github.com/cespedes/venti"
)

const (
	VtRootSize       = 300
	VtRootVersion    = 2
	vtRootVersionBig = 1 << 15
)

type VtRoot struct {
	Name      string
	Type      string
	Score     venti.Score
	BlockSize uint16
	PrevScore venti.Score
}

func VtRootUnpack(b []byte) (*VtRoot, error) {
	vers := getU16(&b)
	if vers != VtRootVersion {
		return nil, fmt.Errorf("unknown root version")
	}
	root := new(VtRoot)

	root.Name = getString(&b, 128)
	root.Type = getString(&b, 128)
	root.Score = getScore(&b)
	root.BlockSize = getU16(&b)
	root.PrevScore = getScore(&b)
	return root, nil
}
