package main

import (
	"fmt"
	"io"
	"os"

	"github.com/cespedes/venti"
)

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
	//	fmt.Printf("* venti server is at %q\n", addr)
	//	fmt.Println("* Dialing...")
	client, err := venti.Dial(addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	//	err = client.Ping()
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "Ping() error: %v\n", err)
	//	}
	//	err = client.Ping()
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "Ping() error: %v\n", err)
	//	}
	//	err = client.Close()
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "Close() error: %v\n", err)
	//		os.Exit(1)
	//	}
	//	fmt.Printf("venti.Dial() returned %v\n", client)
	//	fmt.Printf("* reading venti packet with score %s...\n", score)
	r, err := client.Read(venti.VtData, score, 4096)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err := io.Copy(os.Stdout, r); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
