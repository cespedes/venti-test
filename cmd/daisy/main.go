package main

import (
	"fmt"
	"os"
)

func FileToEntry(filename string) (Entry, error) {
	f, err := os.Lstat(filename)
	if err != nil {
		return Entry{}, err
	}
	var e Entry
	e.Name = f.Name()
	e.Size = uint64(f.Size())
	e.ModTime = f.ModTime()
	e.DestLink, _ = os.Readlink(filename) // ignore errors

	SysFileToEntry(&e, f)

	mode := f.Mode()

	switch {
	case mode&os.ModeType == 0:
		e.Type = TypeRegular
	case mode&os.ModeDir != 0:
		e.Type = TypeDir
	case mode&os.ModeSymlink != 0:
		e.Type = TypeSymlink
	case mode&(os.ModeDevice|os.ModeCharDevice) == os.ModeDevice:
		e.Type = TypeBlock
	case mode&(os.ModeDevice|os.ModeCharDevice) == os.ModeDevice|os.ModeCharDevice:
		e.Type = TypeChar
	case mode&os.ModeNamedPipe == os.ModeNamedPipe:
		e.Type = TypeNamedPipe
	case mode&os.ModeSocket == os.ModeSocket:
		e.Type = TypeSocket
	default:
		e.Type = TypeUnknown
	}

	e.Mode = (uint32)(mode & os.ModePerm)
	if mode&os.ModeSetuid != 0 {
		e.Mode |= 0o4000
	}
	if mode&os.ModeSetgid != 0 {
		e.Mode |= 0o2000
	}
	if mode&os.ModeSticky != 0 {
		e.Mode |= 0o1000
	}

	fmt.Printf("Name: %s\n", e.Name)
	fmt.Printf("Size: %d\n", e.Size)
	fmt.Printf("Type: %d\n", e.Type)
	fmt.Printf("Mode: %03o\n", e.Mode)
	fmt.Printf("ModTime: %s\n", e.ModTime)
	fmt.Printf("UID: %d\n", e.UID)
	fmt.Printf("GID: %d\n", e.GID)
	fmt.Printf("DestLink: %s\n", e.DestLink)

	return e, nil
}

func PrintEntry(e Entry) {
}

func main() {
	for i := 1; i < len(os.Args); i++ {
		entry, err := FileToEntry(os.Args[i])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		PrintEntry(entry)
	}
}
