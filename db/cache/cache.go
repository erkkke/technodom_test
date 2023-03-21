package cache

import (
	"sync"
)

type Cache interface {
	// Add Метод для добавление новой записи в кэш или обновления уже существующей
	Add(key, value string)
	// Get Метод для получения значения из кэша по ключу. Если значения нет, вернуть пустую строку и false
	Get(key string) (value string, ok bool)
	// Len Метод для получения количества ключей в кэше
	Len() int
}

type InMemoryCache struct {
	capacity int
	items    map[string]string
	keys     []string
	mu       sync.Mutex
}

func NewInMemoryCache(capacity int) Cache {
	return &InMemoryCache{
		capacity: capacity,
		items:    make(map[string]string),
		keys:     make([]string, 0, capacity),
	}
}

func (c *InMemoryCache) Add(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[key]; !ok {
		if len(c.items) >= c.capacity {
			delete(c.items, c.keys[0])
			c.keys = c.keys[1:]
		}
		c.keys = append(c.keys, key)
	}

	c.items[key] = value
}

func (c *InMemoryCache) Get(key string) (value string, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok = c.items[key]
	return
}

func (c *InMemoryCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.items)
}
