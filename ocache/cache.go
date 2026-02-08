package ocache

import (
	"container/list"
	"fmt"
	"sync"
)

type Cache struct {
	mtx sync.Mutex
	lru *lruCache
}

func (c *Cache) Get(key string) (Value, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.lru.get(key)
}

func (c *Cache) Add(key string, value Value) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.lru.add(key, value)
}

type pair struct {
	key   string
	value Value
}

type lruCache struct {
	ll       list.List                // 存*(string, Value)
	m        map[string]*list.Element // 存map[str]iterator
	capacity int                      // in bytes
	size     int
}

func (c *lruCache) get(key string) (Value, error) {
	if v, ok := c.m[key]; ok {
		p := v.Value.(*pair)
		c.ll.MoveToFront(v)
		return p.value, nil
	}
	return nil, fmt.Errorf("no such key")
}

func (c *lruCache) add(key string, value Value) error {
	if len(key)+value.Len() > c.capacity {
		return fmt.Errorf("too many bytes")
	}
	if v, ok := c.m[key]; ok {
		p := v.Value.(*pair)
		diff := value.Len() - p.value.Len()
		if c.size+diff > c.capacity {
			return fmt.Errorf("too many bytes")
		}
		p.key, p.value = key, value
		c.size += diff
		c.ll.MoveToFront(v)
		return nil
	}
	for c.size+value.Len()+len(key) > c.capacity {
		c.deleteByLru()
	}
	c.ll.PushFront(&pair{key, value})
	c.m[key] = c.ll.Front()
	c.size += value.Len() + len(key)
	return nil
}

func (c *lruCache) deleteByLru() {
	// delete tail
	if c.ll.Len() == 0 {
		return
	}
	tail := c.ll.Back()
	p := tail.Value.(*pair)
	c.size -= p.value.Len() + len(p.key)
	delete(c.m, p.key)
	c.ll.Remove(tail)
}
