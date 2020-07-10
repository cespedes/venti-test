package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cespedes/venti"
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

func dump(c *venti.Client, indent int, score venti.Score, typ venti.Type) error {
	r, err := c.Read(typ, score, 4096)
	_ = r
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	r.Close()
	for i := 0; i < indent; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("%s ", p(score))
	switch typ {
	case venti.VtRoot:
		fmt.Println("root")
	case venti.VtData:
		fmt.Printf("data n=%d\n", len(b))
	case venti.VtDir:
		if len(b)%40 != 0 {
			return fmt.Errorf("wrong size for directory: %d", len(b))
		}
		fmt.Printf("dir n=%d\n", len(b)/40)
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
