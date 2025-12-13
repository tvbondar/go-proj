package sorter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

const maxChunkBytes = 100 * 1024 * 1024

type sorterContext struct {
	cfg *Config
}

var ctx sorterContext

func Sort(r io.Reader, w io.Writer, cfg *Config) error {
	ctx.cfg = cfg
	files, err := createSortedChunks(r)
	if err != nil {
		return err
	}
	defer cleanup(files)
	return mergeWithHeap(files, w)
}

func cleanup(files []*os.File) {
	for _, f := range files {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close temp file: %v\n", err)
		}
		if err := os.Remove(f.Name()); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to remove temp file %s: %v\n", f.Name(), err)
		}
	}
}

func createSortedChunks(r io.Reader) ([]*os.File, error) {
	sc := bufio.NewScanner(r)
	var chunk []string
	var files []*os.File
	bytesInChunk := 0

	for {
		chunk = chunk[:0]
		bytesInChunk = 0

		for bytesInChunk < maxChunkBytes && sc.Scan() {
			line := sc.Text()
			chunk = append(chunk, line)
			bytesInChunk += len(line) + 1
		}
		if len(chunk) == 0 {
			break
		}

		sortChunk(chunk)
		f, err := os.CreateTemp(ctx.cfg.TempDir, "mysort-*.tmp")
		if err != nil {
			return nil, err
		}

		for _, l := range chunk {
			if _, err := fmt.Fprintln(f, l); err != nil {
				_ = f.Close()
				_ = os.Remove(f.Name())
				return nil, err
			}
		}
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			_ = f.Close()
			_ = os.Remove(f.Name())
			return nil, err
		}
		files = append(files, f)
	}
	return files, sc.Err()
}

func sortChunk(lines []string) {
	sort.SliceStable(lines, func(i, j int) bool {
		cmp := CompareKeys(ExtractKey(lines[i], ctx.cfg), ExtractKey(lines[j], ctx.cfg), ctx.cfg)
		if cmp != 0 {
			return cmp < 0
		}
		if ctx.cfg.Reverse {
			return lines[i] > lines[j]
		}
		return lines[i] < lines[j]
	})
}
