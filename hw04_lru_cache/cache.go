package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mux      sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	newValue := cacheItem{key: key, value: value}
	listItem, exists := c.items[key]
	if exists {
		listItem.Value = newValue
		c.queue.MoveToFront(listItem)
	} else {
		listItem = c.queue.PushFront(newValue)
		if c.queue.Len() > c.capacity {
			lastListItem := c.queue.Back()
			lastValue := lastListItem.Value.(cacheItem)
			delete(c.items, lastValue.key)
			c.queue.Remove(lastListItem)
		}
	}

	c.items[key] = listItem
	return exists
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	var result interface{}
	listItem, exists := c.items[key]
	if exists {
		result = listItem.Value.(cacheItem).value
		c.queue.MoveToFront(listItem)
	}

	return result, exists
}

func (c *lruCache) Clear() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
