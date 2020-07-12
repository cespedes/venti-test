// +build !windows

package main

import (
	"os"
	"syscall"
)

func SysFileToEntry(e *Entry, f os.FileInfo) {
	stat, ok := f.Sys().(*syscall.Stat_t)
	if ok {
		e.UID = stat.Uid
		e.GID = stat.Gid
		e.Major = uint(stat.Rdev/256)
		e.Minor = uint(stat.Rdev%256)
	}
}
