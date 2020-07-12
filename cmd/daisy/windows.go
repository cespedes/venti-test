// +build windows

package main

import (
	"os"
)

func SysFileToEntry(e *Entry, f os.FileInfo) {
	// nothing to do in Windows (no UID, no GID, no file devices)
}
