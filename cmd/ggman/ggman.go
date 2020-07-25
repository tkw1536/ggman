package main

import (
	"os"

	"github.com/tkw1536/ggman"
)

func main() {
	code, err := ggman.Main(os.Args[1:])
	if code != 0 && err != "" {
		os.Stderr.WriteString(err + "\n")
	}
	os.Exit(code)
}
