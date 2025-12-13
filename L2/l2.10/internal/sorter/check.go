package sorter

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func CheckSorted(r io.Reader, cfg *Config) error {
	sc := bufio.NewScanner(r)
	var prev string
	var prevKey *Key
	lineNum := 0

	for sc.Scan() {
		lineNum++
		line := sc.Text()
		key := ExtractKey(line, cfg)

		if prevKey != nil {
			if CompareKeys(key, *prevKey, cfg) < 0 ||
				(CompareKeys(key, *prevKey, cfg) == 0 && ((cfg.Reverse && line < prev) || (!cfg.Reverse && line > prev))) {
				fmt.Fprintf(os.Stderr, "mysort: disorder on line %d\n", lineNum)
				return fmt.Errorf("not sorted")
			}
		}
		prev, prevKey = line, &key
	}
	return nil
}
