package sorter

import (
	"bufio"
	"os"
)

type Config struct {
	KeyField             int
	Numeric              bool
	Reverse              bool
	Unique               bool
	Month                bool
	IgnoreTrailingBlanks bool
	CheckIfSorted        bool
	HumanNumeric         bool
	MaxMemoryMB          int64
	TempDir              string
}

type Key struct {
	Value   any // string | float64 | int
	Raw     string
	Trimmed string
}

type mergeItem struct {
	line  string
	sc    *bufio.Scanner
	file  *os.File
	index int
}

type mergeHeap []*mergeItem
