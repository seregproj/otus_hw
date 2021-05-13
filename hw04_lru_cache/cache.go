package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.Lock()
	defer l.Unlock()

	item, ok := l.items[key]

	if ok {
		l.queue.MoveToFront(item)
		item.Value = cacheItem{value: value, key: string(key)}

		return true
	}

	if l.capacity == l.queue.Len() {
		lastItem := l.queue.Back()

		delete(l.items, Key(lastItem.Value.(cacheItem).key))
		l.queue.Remove(lastItem)
	}

	newItem := cacheItem{value: value, key: string(key)}
	addedItem := l.queue.PushFront(newItem)
	l.items[key] = addedItem

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.Lock()
	defer l.Unlock()

	item, ok := l.items[key]

	if ok {
		l.queue.MoveToFront(item)

		return item.Value.(cacheItem).value, ok
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.Lock()
	defer l.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

type cacheItem struct {
	key   string
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
