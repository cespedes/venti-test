package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, `Usage: venti-usage <host:port>`)
		os.Exit(1)
	}
	resp, err := http.Get(os.Args[1] + "/storage")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	if scanner.Scan() {
		fmt.Printf("index: %s\n", scanner.Text())
	}
	if scanner.Scan() {
		fmt.Printf("totalArenas,activeArenas: %s\n", scanner.Text())
	}
	if scanner.Scan() {
		fmt.Printf("totalSpace,usedSpace: %s\n", scanner.Text())
	}
	for scanner.Scan() {
		fmt.Printf("= %s =\n", scanner.Text())
	}
}
