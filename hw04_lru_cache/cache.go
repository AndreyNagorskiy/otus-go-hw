package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) findItem(key Key) (*ListItem, bool) {
	item, ok := l.items[key]
	return item, ok
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if l.capacity <= 0 {
		return false
	}

	cache := cacheItem{key, value}

	item, exists := l.findItem(key)

	if exists {
		item.Value = cache
		l.queue.MoveToFront(item)

		return true
	}

	newListItem := l.queue.PushFront(cache)

	if l.queue.Len() > l.capacity {
		lastQueueItem := l.queue.Back()
		delete(l.items, lastQueueItem.Value.(cacheItem).key)
		l.queue.Remove(lastQueueItem)
	}

	l.items[key] = newListItem

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	item, exists := l.findItem(key)

	if !exists {
		return nil, false
	}

	l.queue.MoveToFront(item)

	return item.Value.(cacheItem).value, true
}

func (l *lruCache) Clear() {
	l.items = make(map[Key]*ListItem, l.capacity)
	l.queue = NewList()
}
