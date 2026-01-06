// internal/grep/grep_test.go
package grep

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunGrep(t *testing.T) {
	// Создаём временную директорию и тестовые файлы
	tmpDir := t.TempDir()

	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")

	content1 := `line one
hello world
Line Two
HELLO there
exact match here
no match
hello again`
	if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
		t.Fatalf("failed to write test file1: %v", err)
	}

	content2 := `first line
second hello
third line
Hello World`
	if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
		t.Fatalf("failed to write test file2: %v", err)
	}

	// Таблица тестов
	tests := []struct {
		name    string
		opts    Options
		pattern string
		files   []string
		want    string
	}{
		{
			name:    "basic search",
			opts:    Options{},
			pattern: "hello",
			files:   []string{file1},
			want:    "hello world\nhello again\n",
		},
		{
			name:    "ignore case",
			opts:    Options{IgnoreCase: true},
			pattern: "hello",
			files:   []string{file1},
			want:    "hello world\nHELLO there\nhello again\n",
		},
		{
			name:    "line numbers",
			opts:    Options{LineNum: true},
			pattern: "hello",
			files:   []string{file1},
			want:    "2:hello world\n7:hello again\n",
		},
		{
			name:    "count only",
			opts:    Options{Count: true},
			pattern: "hello",
			files:   []string{file1},
			want:    "2\n",
		},
		{
			name:    "invert match",
			opts:    Options{Invert: true},
			pattern: "hello",
			files:   []string{file1},
			want:    "line one\nLine Two\nexact match here\nno match\n",
		},
		{
			name:    "context after 1",
			opts:    Options{After: 1},
			pattern: "hello",
			files:   []string{file1},
			want:    "hello world\nLine Two\n--\nhello again\n",
		},
		{
			name:    "context around 1",
			opts:    Options{Context: 1},
			pattern: "hello",
			files:   []string{file1},
			want:    "line one\nhello world\nLine Two\n--\nHELLO there\nexact match here\nno match\nhello again\n",
		},
		{
			name:    "multiple files",
			opts:    Options{},
			pattern: "hello",
			files:   []string{file1, file2},
			want: file1 + ":hello world\n" + file1 + ":hello again\n" +
				file2 + ":second hello\n" + file2 + ":Hello World\n",
		},
		{
			name:    "no matches",
			opts:    Options{},
			pattern: "xyzxyz",
			files:   []string{file1},
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			err := RunGrep(&out, tt.opts, tt.pattern, tt.files)
			if err != nil {
				t.Fatalf("RunGrep returned unexpected error: %v", err)
			}
			got := out.String()
			got = strings.ReplaceAll(got, "\r\n", "\n")
			tt.want = strings.ReplaceAll(tt.want, "\r\n", "\n")
			if got != tt.want {
				t.Errorf("%s failed:\nGOT:\n%q\nWANT:\n%q", tt.name, got, tt.want)
			}
		})
	}
}
