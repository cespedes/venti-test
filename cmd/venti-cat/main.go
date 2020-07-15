package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/cespedes/venti"
	"github.com/cespedes/venti-utils"
)

func main() {
	var score venti.Score
	var t, size, dsize, psize uint64
	var err error

	if len(os.Args) != 6 {
		fmt.Fprintf(os.Stderr, "Usage: %s <score> <type> <size> <dsize> <psize>\n", os.Args[0])
		os.Exit(1)
	}
	addr := os.Getenv("venti")
	score, err = venti.ParseScore(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	t, err = strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	size, err = strconv.ParseUint(os.Args[3], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	dsize, err = strconv.ParseUint(os.Args[4], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	psize, err = strconv.ParseUint(os.Args[5], 10, 64)
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

	d, err := ventiutils.OpenData(client, score, venti.Type(t), size, uint(dsize), uint(psize))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, err = io.Copy(os.Stdout, d)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// fmt.Printf("d: %v\n", d)
}
