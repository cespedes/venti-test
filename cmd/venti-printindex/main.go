package main

import (
	"os"

	venti "github.com/cespedes/venti-utils"
)

func main() {
	venti.ParseIsect(os.Args[1])
}
