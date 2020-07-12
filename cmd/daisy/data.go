package main

import (
	"time"

	"github.com/cespedes/venti"
)

const (
	TypeRegular = iota
	TypeDir
	TypeSymlink
	TypeBlock
	TypeChar
	TypeNamedPipe
	TypeSocket
	TypeUnknown
)

type Entry struct {
	Name string
	Size uint64
	Type int
	Mode     uint32 // chmod-compatible
	ModTime  time.Time
	Score    venti.Score

	UID      uint32 // only meaningful in UNIX-like systems
	GID      uint32
	DestLink string
	Major, Minor uint
}
