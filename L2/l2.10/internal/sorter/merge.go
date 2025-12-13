package sorter

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"os"
)

func (h mergeHeap) Len() int { return len(h) }
func (h mergeHeap) Less(i, j int) bool {
	a, b := ExtractKey(h[i].line, ctx.cfg), ExtractKey(h[j].line, ctx.cfg)
	if cmp := CompareKeys(a, b, ctx.cfg); cmp != 0 {
		return cmp < 0
	}
	return h[i].index < h[j].index
}
func (h mergeHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i]; h[i].index, h[j].index = j, i }
func (h *mergeHeap) Push(x any)   { *h = append(*h, x.(*mergeItem)) }
func (h *mergeHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

func mergeWithHeap(files []*os.File, w io.Writer) error {
	if len(files) == 0 {
		return nil
	}
	if len(files) == 1 {
		_, err := io.Copy(w, files[0])
		return err
	}

	h := make(mergeHeap, 0, len(files))
	for i, f := range files {
		sc := bufio.NewScanner(f)
		if sc.Scan() {
			heap.Push(&h, &mergeItem{line: sc.Text(), sc: sc, file: f, index: i})
		}
	}

	var prevKey *Key
	for h.Len() > 0 {
		item := heap.Pop(&h).(*mergeItem)
		key := ExtractKey(item.line, ctx.cfg)

		if !ctx.cfg.Unique || prevKey == nil || CompareKeys(*prevKey, key, ctx.cfg) != 0 {
			if _, err := fmt.Fprintln(w, item.line); err != nil {
				return err // пишем в stdout — ошибка критична
			}
			prevKey = &key
		}

		if item.sc.Scan() {
			item.line = item.sc.Text()
			heap.Push(&h, item)
		}
	}
	return nil
}
