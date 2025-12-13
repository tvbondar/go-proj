package sorter

import (
	"math"
	"strconv"
	"strings"
	"unicode"
)

// map с месяцами для сортировки по месяцам
var MonthOrder = map[string]int{
	"JANUARY": 1, "FEBRUARY": 2, "MARCH": 3, "APRIL": 4, "MAY": 5, "JUNE": 6,
	"JULY": 7, "AUGUST": 8, "SEPTEMBER": 9, "OCTOBER": 10, "NOVEMBER": 11, "DECEMBER": 12,
}

// проверяет наличие флагов в команде
func extractField(line string, fieldNum int) string {
	if fieldNum <= 0 {
		return line
	}
	fields := strings.Split(line, "\t")
	index := fieldNum - 1
	if index >= len(fields) {
		return ""
	}
	return fields[index]
}

// если есть хвостовые пробелы, игнорируем их
func TrimTrailingBlanks(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// сортировка по числовому значению с учетом суффиксов
func ParseHumanNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	mult := 1.0
	last := s[len(s)-1]

	switch last {
	case 'K', 'k':
		mult = 1e3
		s = s[:len(s)-1]
	case 'M', 'm':
		mult = 1e6
		s = s[:len(s)-1]
	case 'G', 'g':
		mult = 1e9
		s = s[:len(s)-1]
	case 'T', 't':
		mult = 1e12
		s = s[:len(s)-1]
	}

	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, false
	}
	return f * mult, true

}

// извлечение ключа из команды и применение нужной функции
func ExtractKey(line string, cfg *Config) Key {
	field := extractField(line, cfg.KeyField)
	trimmed := field
	if cfg.IgnoreTrailingBlanks {
		trimmed = TrimTrailingBlanks(field)
	}

	if cfg.Month {
		s := strings.TrimSpace(trimmed)
		if len(s) < 3 {
		} else {
			upper := strings.ToUpper(s[:3])
			if m, ok := MonthOrder[upper]; ok {
				return Key{Value: m, Raw: field, Trimmed: trimmed}
			}
		}
	}

	if cfg.Numeric || cfg.HumanNumeric {
		var num float64
		var ok bool
		if cfg.HumanNumeric {
			num, ok = ParseHumanNumber(field)
		} else {
			num, _ = strconv.ParseFloat(strings.TrimSpace(trimmed), 64)
			ok = true
		}
		if ok {
			return Key{Value: num, Raw: field, Trimmed: trimmed}
		}
	}

	return Key{Value: trimmed, Raw: field, Trimmed: trimmed}
}

// Сравнение ключей
func CompareKeys(a, b Key, cfg *Config) int {
	if a.Value == nil {
		if b.Value == nil {
			return 0
		}
		return -1
	}
	if b.Value == nil {
		return 1
	}

	var cmp int
	switch av := a.Value.(type) {
	case string:
		cmp = strings.Compare(av, b.Value.(string))
	case float64:
		bv := b.Value.(float64)
		if math.IsNaN(av) {
			cmp = -1
		} else if math.IsNaN(bv) {
			cmp = 1
		} else if av < bv {
			cmp = -1
		} else if av > bv {
			cmp = 1
		}
	case int:
		cmp = av - b.Value.(int)
	}

	if cfg.Reverse {
		cmp = -cmp
	}
	return cmp
}
