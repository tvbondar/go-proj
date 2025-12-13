package sorter

import (
	"flag"
	"os"
)

func ParseFlags() *Config {
	cfg := &Config{
		MaxMemoryMB: 256,
		TempDir:     os.TempDir(),
	}

	flag.IntVar(&cfg.KeyField, "k", 0, "sort by field number (1-based)")
	flag.BoolVar(&cfg.Numeric, "n", false, "numeric sort")
	flag.BoolVar(&cfg.Reverse, "r", false, "reverse order")
	flag.BoolVar(&cfg.Unique, "u", false, "unique lines")
	flag.BoolVar(&cfg.Month, "M", false, "month sort (Jan, Feb, ...)")
	flag.BoolVar(&cfg.IgnoreTrailingBlanks, "b", false, "ignore trailing blanks")
	flag.BoolVar(&cfg.CheckIfSorted, "c", false, "check if already sorted")
	flag.BoolVar(&cfg.HumanNumeric, "h", false, "human numeric sort (10K, 5M...)")
	flag.Int64Var(&cfg.MaxMemoryMB, "S", 256, "max memory in MB")
	flag.StringVar(&cfg.TempDir, "T", os.TempDir(), "temp directory")

	flag.Parse()
	return cfg
}
