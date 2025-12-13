package anagram

import (
	"sort"
	"strings"
)

// функция, находящая все множества анаграмм в словаре
func FindAnagrams(words []string) map[string][]string {
	lower := make([]string, len(words))
	for i, word := range words {
		lower[i] = strings.ToLower(word)
	}

	groups := make(map[string][]int)
	for i, word := range lower {
		//создаем сигнатуру - отсортированные map
		r := []rune(word)
		sort.Slice(r, func(i, j int) bool {
			return r[i] < r[j]
		})
		sig := string(r)
		groups[sig] = append(groups[sig], i)
	}
	res := make(map[string][]string)
	for _, indices := range groups {
		if len(indices) < 2 {
			continue //берем только группы из 2 и более слов
		}
		// ищем слово с наименьшим индексом в группе
		minIndex := indices[0]
		for _, index := range indices {
			if index < minIndex {
				minIndex = index
			}
		}
		key := lower[minIndex]
		//собираем уникальные слова из группы
		anagrams := make([]string, 0, len(indices))
		seen := make(map[string]bool)
		for _, index := range indices {
			word := lower[index]
			if !seen[word] {
				seen[word] = true
				anagrams = append(anagrams, word)
			}
		}
		//сортируем
		sort.Strings(anagrams)
		res[key] = anagrams
	}
	return res
}
