package main

import (
	"os"

	"l2.10/internal/sorter"
)

func main() {
	cfg := sorter.ParseFlags()

	if cfg.CheckIfSorted {
		if err := sorter.CheckSorted(os.Stdin, cfg); err != nil {
			os.Exit(1)
		}
		return
	}

	if err := sorter.Sort(os.Stdin, os.Stdout, cfg); err != nil {
		os.Exit(1)
	}
}
