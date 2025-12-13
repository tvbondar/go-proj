package anagram_test

import (
	"reflect"
	"testing"

	"l2.11/internal/anagram"
)

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  map[string][]string
	}{
		{
			name:  "basic example",
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"},
			want: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			name:  "with duplicates and mixed case",
			input: []string{"Пятак", "пятка", "Тяпка", "пятак", "листок", "Листок"},
			want: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок"},
			},
		},
		{
			name:  "no anagrams",
			input: []string{"кот", "дом", "мир"},
			want:  map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := anagram.FindAnagrams(tt.input)

			if len(got) != len(tt.want) {
				t.Errorf("wrong number of groups: got %d, want %d", len(got), len(tt.want))
				return
			}

			for key, wantSlice := range tt.want {
				gotSlice, exists := got[key]
				if !exists {
					t.Errorf("missing key: %s", key)
					continue
				}

				// Проверяем, что слайсы равны (сортировка гарантирована)
				if !reflect.DeepEqual(gotSlice, wantSlice) {
					t.Errorf("for key %s: got %v, want %v", key, gotSlice, wantSlice)
				}
			}
		})
	}
}
