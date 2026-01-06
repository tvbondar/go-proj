// internal/grep/grep.go
package grep

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

func RunGrep(w io.Writer, opts Options, pattern string, filenames []string) error {
	// читаем из stdin, если не указаны файлы
	if len(filenames) == 0 {
		filenames = []string{"-"}
	}

	multipleFiles := len(filenames) > 1
	// Компилируем шаблон поиска
	var matcher func(string) bool
	if opts.Fixed {
		//-F: фиксированная подстрока
		target := pattern
		if opts.IgnoreCase {
			target = strings.ToLower(pattern)
		}
		matcher = func(line string) bool {
			s := line
			if opts.IgnoreCase {
				s = strings.ToLower(s)
			}
			return strings.Contains(s, target)
		}
	} else {
		//Обычный режим с регулярными выражениями
		flags := ""
		if opts.IgnoreCase {
			flags = "(?i)" //игнорируем регистр
		}
		re, err := regexp.Compile(flags + pattern)
		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}
		matcher = re.MatchString
	}

	// Обрабатываем каждый файл
	for _, filename := range filenames {
		lines, err := readLines(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", filename, err)
			continue
		}
		// Находим индексы строк, которые соответствуют шаблону
		matches := []int{}
		for i, line := range lines {
			match := matcher(line)
			if opts.Invert {
				match = !match // -v: инвертируем совпадение
			}
			if match {
				matches = append(matches, i)
			}
		}

		// C: подсчет количества совпадений
		if opts.Count {
			prefix := ""
			if multipleFiles {
				prefix = filename + ":"
			}
			if _, err := fmt.Fprintln(w, prefix+fmt.Sprint(len(matches))); err != nil {
				return err
			}
			continue
		}
		if len(matches) == 0 {
			continue
		}
		// Формируем диапазоны строк для вывода с учетом контекста
		var merged [][]int

		if opts.Before == 0 && opts.After == 0 {
			// Без контекста — выводим только совпадающие строки, каждая отдельно
			for _, m := range matches {
				merged = append(merged, []int{m, m})
			}
		} else {
			// С контекстом — строим диапазоны и объединяем пересекающиеся
			ranges := buildRanges(matches, len(lines), opts.Before, opts.After)
			merged = mergeRanges(ranges)
		}
		// Префикс с именем файла, если их несколько
		filePrefix := ""
		if multipleFiles {
			filePrefix = filename + ":"
		}

		//Выводим все группы совпадений с разделителями
		for gi, rg := range merged {
			// Разделитель между группами
			if gi > 0 && (opts.Before > 0 || opts.After > 0) {
				if _, err := fmt.Fprintln(w, "--"); err != nil {
					return err
				}
			}
			// Проходим по строкам в диапазоне
			for i := rg[0]; i <= rg[1]; i++ {
				// При нулевом контексте пропускаем несовпадающие строки
				if opts.Before == 0 && opts.After == 0 && !isMatch(i, matches) {
					continue
				}
				// Формируем префикс строки
				prefix := ""
				if opts.LineNum {
					prefix = fmt.Sprintf("%d:", i+1)
				}
				isMatching := isMatch(i, matches)

				//Логика префикса для нескольких файлов
				if multipleFiles {
					//Все строки получают имя файла при нескольких файлах
					if isMatching {
						prefix = filePrefix + prefix
					} else if !isMatching {
						//Только несовпадающие строки получают имя файла если файл один
						prefix = "-" + prefix
					}
				}
				// Выводим строку с префиксом
				if _, err := fmt.Fprintf(w, "%s%s\n", prefix, lines[i]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// читает все строки из файла или stdin
func readLines(filename string) ([]string, error) {
	var scanner *bufio.Scanner

	if filename == "-" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer func() {
			if cerr := f.Close(); cerr != nil {
				fmt.Fprintf(os.Stderr, "warning: error closing file %s: %v\n", filename, cerr)
			}
		}()
		scanner = bufio.NewScanner(f)
	}

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// Строит диапазоны строк для вывода с учетом контекста
func buildRanges(matches []int, totalLines, before, after int) [][]int {
	ranges := make([][]int, len(matches))
	for i, m := range matches {
		start := max(0, m-before)
		end := min(totalLines-1, m+after)
		ranges[i] = []int{start, end}
	}
	return ranges
}

// Объединяет пересекающиеся диапазоны строк
func mergeRanges(ranges [][]int) [][]int {
	if len(ranges) == 0 {
		return nil
	}
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i][0] < ranges[j][0]
	})

	merged := [][]int{ranges[0]}
	for _, r := range ranges[1:] {
		last := merged[len(merged)-1]
		if r[0] <= last[1]+1 {
			if r[1] > last[1] {
				last[1] = r[1]
			}
		} else {
			merged = append(merged, r)
		}
	}
	return merged
}

// Проверяет, есть ли строка в списке совпадений
func isMatch(line int, matches []int) bool {
	i := sort.SearchInts(matches, line)
	return i < len(matches) && matches[i] == line
}
