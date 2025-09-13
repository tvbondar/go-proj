package main

import (
	"fmt"
	"sync"
)

// структура map
type Map struct {
	mu   sync.Mutex
	data map[string]int
}

// конструктор
func NewMap() *Map {
	return &Map{
		data: make(map[string]int),
	}
}

// set - безопасная запись
func (m *Map) Set(key string, value int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// get - безопасное чтение
func (m *Map) Get(key string) (int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, exists := m.data[key]
	return value, exists
}

func main() {
	m := NewMap()
	var wgSet sync.WaitGroup
	var wgGet sync.WaitGroup
	//параллельная запись
	for i := 0; i < 10; i++ {
		wgSet.Add(1)
		go func(i int) {
			defer wgSet.Done()
			key := fmt.Sprintf("key%d", i)
			m.Set(key, i)
		}(i)
	}

	wgSet.Wait()

	//параллельное чтение
	for i := 0; i < 10; i++ {
		wgGet.Add(1)
		go func(i int) {
			defer wgGet.Done()
			key := fmt.Sprintf("key%d", i)
			if value, exists := m.Get(key); exists {
				fmt.Printf("Найдено: %s -> %d\n", key, value)
			} else {
				fmt.Printf("Не найдено: %s\n", key)
			}
		}(i)
	}

	wgGet.Wait()
	fmt.Println("Все операции завершены")
}
