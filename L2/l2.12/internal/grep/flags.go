// flags.go
package grep

import (
	"flag"
	"os"
)

type Options struct {
	After      int  // -A N
	Before     int  // -B N
	Context    int  //-C N
	Count      bool // -c
	IgnoreCase bool // -i
	Invert     bool // -v
	Fixed      bool // -F
	LineNum    bool // -n
}

func ParseFlags() (Options, []string) {
	var opts Options

	flag.IntVar(&opts.After, "A", 0, "print N lines after match")
	flag.IntVar(&opts.Before, "B", 0, "print N lines before match")
	flag.IntVar(&opts.Context, "C", 0, "print N lines of context around match")
	flag.BoolVar(&opts.Count, "c", false, "print only count of matching lines")
	flag.BoolVar(&opts.IgnoreCase, "i", false, "ignore case")
	flag.BoolVar(&opts.Invert, "v", false, "invert match (select non-matching lines)")
	flag.BoolVar(&opts.Fixed, "F", false, "pattern is a fixed string, not regex")
	flag.BoolVar(&opts.LineNum, "n", false, "print line numbers")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// Если указан -C, он переопределяет -A и -B
	if opts.Context > 0 {
		opts.After = opts.Context
		opts.Before = opts.Context
	}

	return opts, args
}
