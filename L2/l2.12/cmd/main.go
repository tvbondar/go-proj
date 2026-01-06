package main

import (
	"fmt"
	"os"

	"l2.12/internal/grep"
)

func main() {
	opts, args := grep.ParseFlags()

	pattern := args[0]
	filenames := args[1:]

	if err := grep.RunGrep(os.Stdout, opts, pattern, filenames); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
