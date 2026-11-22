package appname

import (
	"strings"
	"sync"
)

func FileName(path string) string {
	path = strings.TrimSpace(path)
	if i := strings.LastIndexAny(path, `\/`); i >= 0 {
		path = path[i+1:]
	}
	if len(path) > 4 && strings.EqualFold(path[len(path)-4:], ".exe") {
		path = path[:len(path)-4]
	}
	return path
}

func clean(s string) string {
	return strings.Join(strings.Fields(strings.Trim(s, "\x00")), " ")
}

func Choose(description, product, path string) string {
	if d := clean(description); d != "" {
		return d
	}
	if p := clean(product); p != "" {
		return p
	}
	return FileName(path)
}

type Cache struct {
	mu    sync.Mutex
	names map[string]string
	read  func(path string) (description, product string)
}

func NewCache(read func(path string) (string, string)) *Cache {
	return &Cache{names: map[string]string{}, read: read}
}

func (c *Cache) Name(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	key := strings.ToLower(path)
	c.mu.Lock()
	name, ok := c.names[key]
	c.mu.Unlock()
	if ok {
		return name
	}
	desc, product := "", ""
	if c.read != nil {
		desc, product = c.read(path)
	}
	name = Choose(desc, product, path)
	c.mu.Lock()
	c.names[key] = name
	c.mu.Unlock()
	return name
}
