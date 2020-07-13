package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cespedes/venti"
	utils "github.com/cespedes/venti-utils"
)

func p(score venti.Score) string {
	switch len(score) {
	case 0:
		return "nil"
	case 20:
		s := ""
		for i := 0; i < 20; i++ {
			s += fmt.Sprintf("%02x", score[i])
		}
		return s
	default:
		panic(fmt.Sprintf("unknown score: %v", score))
	}
}

func printIndent(indent int) {
	for i := 0; i < 4*indent; i++ {
		fmt.Print(" ")
	}
}

func dump(c *venti.Client, indent int, score venti.Score, typ venti.Type) error {
	var b []byte
	if !bytes.Equal([]byte(score), []byte(utils.VtZeroScore)) {
		r, err := c.Read(typ, score, utils.VtMaxLumpSize)
		if err != nil {
			return err
		}
		b, err = ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		r.Close()
	}
	printIndent(indent)
	fmt.Printf("%s ", p(score))
	switch {
	case typ == venti.VtRoot:
		root, err := utils.VtRootUnpack(b)
		if err != nil {
			return err
		}
		fmt.Printf("root name=%q type=%q prev=%s bsize=%d\n", root.Name, root.Type, p(root.PrevScore), root.BlockSize)
		err = dump(c, indent+1, root.Score, venti.VtDir)
		if err != nil {
			return err
		}
	case typ == venti.VtDir:
		size := len(b)
		if size%utils.VtEntrySize != 0 {
			return fmt.Errorf("wrong size for directory: %d", size)
		}
		size /= utils.VtEntrySize
		fmt.Printf("dir n=%d\n", size)
		for i := 0; i < size; i++ {
			entry, err := utils.VtEntryUnpack(b[i*utils.VtEntrySize:])
			if err != nil {
				panic(err)
			}
			printIndent(indent + 1)
			fmt.Printf("%d: gen=%d psize=%d dsize=%d type=%d flags=0x%02x size=%d score=%s\n", i,
				entry.Gen, entry.PSize, entry.DSize, entry.Type, entry.Flags, entry.Size, p(entry.Score))
			err = dump(c, indent+2, entry.Score, entry.Type)
			if err != nil {
				return err
			}
		}
	case typ == venti.VtData:
		size := len(b)
		fmt.Printf("data n=%d\n", size)
		i := 0
		for i < size {
			printIndent(indent + 1)
			s := ""
			for j := 0; j < 16; j++ {
				if i < size {
					fmt.Printf(" %02x", b[i])
					if b[i] >= 32 && b[i] <= 126 {
						s += string(b[i])
					} else {
						// s += "�"
						// s += "_"
						s += "…"
					}
				} else {
					fmt.Print("   ")
				}
				i++
			}
			fmt.Printf("  %s\n", s)
		}
	case typ > venti.VtData && typ < venti.VtDir:
		size := len(b)
		if size%20 != 0 {
			return fmt.Errorf("wrong size for pointer to data: %d", size)
		}
		size /= 20
		fmt.Printf("data+%d n=%d\n", typ-venti.VtData, size)
		for i := 0; i < size; i++ {
			err := dump(c, indent+1, venti.Score(b[20*i:20*(i+1)]), typ-1)
			if err != nil {
				return err
			}
		}
	case typ > venti.VtDir && typ < venti.VtRoot:
		size := len(b)
		if size%20 != 0 {
			return fmt.Errorf("wrong size for pointer to dir: %d", size)
		}
		size /= 20
		fmt.Printf("dir+%d n=%d\n", typ-venti.VtDir, size)
		for i := 0; i < size; i++ {
			err := dump(c, indent+1, venti.Score(b[20*i:20*(i+1)]), typ-1)
			if err != nil {
				return err
			}
		}
	default:
		fmt.Printf("To-Do (type=%d) :-)\n", typ)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <score>\n", os.Args[0])
		os.Exit(1)
	}
	addr := os.Getenv("venti")
	score, err := venti.ParseScore(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	client, err := venti.Dial(addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer client.Close()

	// We don't know which type is the block, so we will try all of them:
	for i := venti.VtData; i <= venti.VtRoot; i++ {
		if r, err := client.Read(i, score, 32768); err == nil {
			r.Close()
			err = dump(client, 0, score, i)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
	}
	fmt.Fprintf(os.Stderr, "cannot find block %s\n", score)
	os.Exit(1)
}
