package main

import (
	"github.com/cespedes/venti"
)

const (
	vtEntryActive = 1<<0		// entry is in use
	_vtEntryDir = 1<<1		// a directory
	_vtEntryDepthShift = 2		// shift for pointer depth
	_vtEntryDepthMask = 7<<2	// mask for pointer depth
	vtEntryLocal = 1<<5		// for local storage only
	_vtEntryBig = 1<<6		// dsize and psize are encoded differently
	vtEntryNoArchive = 1<<7		// for local storage only
)

func vtEntryUnpack(b []byte) (*VtEntry, error) {
	entry := new(VtEntry)

	entry.Gen = getU32(&b)
	entry.PSize = getU16(&b)
	entry.DSize = getU16(&b)
	entry.Flags = getU8(&b)

	if (entry.Flags & _vtEntryBig) != 0 {
		entry.PSize = (entry.PSize >> 5) << (entry.PSize & 31)
		entry.DSize = (entry.DSize >> 5) << (entry.DSize & 31)
	}
	if (entry.Flags & _vtEntryDir) != 0 {
		entry.Type = venti.VtDir
	} else {
		entry.Type = venti.VtData
	}
	entry.Type += venti.Type((entry.Flags & _vtEntryDepthMask) >> _vtEntryDepthShift)
	entry.Flags &= ^byte(_vtEntryDir | _vtEntryDepthMask | _vtEntryBig)

	b = b[5:]
	entry.Size = getU48(&b)
	entry.Score = getScore(&b)

	return entry, nil
}

