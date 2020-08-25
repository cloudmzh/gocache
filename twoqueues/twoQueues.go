// Programmed by Mazehua
// South China Normal University
// August 25, 2020

package twoqueues

import (
	"container/list"
)

// Cache based on 2Q algorithm
// The 2Q algorithm consists of two queues, one is LRU queue and the other is FIFO queue
type Cache struct {
	lruMaxEntries  int
	fifoMaxEntries int
	lru            *list.List
	fifo           *list.List
	cache          map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
	isLru bool
}

// New +!
func New(lMaxEntries, fMaxEntries int) *Cache {
	return &Cache{
		lruMaxEntries:  lMaxEntries,
		fifoMaxEntries: fMaxEntries,
		lru:            list.New(),
		fifo:           list.New(),
		cache:          make(map[string]*list.Element),
	}
}

// Add :If the element is not in the cache, push it into FIFO;
// if the element is in a FIFO queue, move it to the head of LRU;
// if the element is in an LRU queue, move it to the head of LRU
func (c *Cache) Add(key string, value interface{}) {
	//Initialize cache
	if c.cache == nil {
		c.lru = list.New()
		c.fifo = list.New()
		c.cache = make(map[string]*list.Element)
		c.fifoMaxEntries = 2 << 10
		c.lruMaxEntries = 2 << 10
	}

	if element, ok := c.cache[key]; ok {
		ele := element.Value.(*entry)
		ele.value = value
		if ele.isLru {
			c.lru.MoveToFront(element)
		} else {
			c.fifo.Remove(element)
			ele.isLru = true
			updateElement := c.lru.PushFront(ele)
			c.cache[key] = updateElement
		}

	} else {
		element := c.fifo.PushBack(&entry{key, value, false})
		c.cache[key] = element
	}
	//During the process of moving from a FIFO queue to an LRU queue,
	//if the element exceeds the LRU maximum, the LRU cleanup mechanism is adopted
	if c.lru.Len() > c.lruMaxEntries {
		c.RemoveLruBack()
	}
	//Check whether the number of FIFO queue elements exceeds the maximum FIFO limit
	if c.fifo.Len() > c.fifoMaxEntries {
		c.RemoveFifoFront()
	}
}

// Get ：If the element is in a FIFO queue, move it to the LRU head,
// if the element is in an LRU queue, move it to the LRU head, and finally return the element
func (c *Cache) Get(key string) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}

	if element, hit := c.cache[key]; hit {
		ele := element.Value.(*entry)
		if ele.isLru {
			c.lru.MoveToFront(element)
		} else {
			c.fifo.Remove(element)
			ele.isLru = true
			updateElement := c.lru.PushFront(ele)
			c.cache[key] = updateElement
		}
		//During the process of moving from a FIFO queue to an LRU queue,
		//if the element exceeds the LRU maximum, the LRU cleanup mechanism is adopted
		if c.lru.Len() > c.lruMaxEntries {
			c.RemoveLruBack()
		}
		return ele.value, true
	}
	return
}

// RemoveLruBack +!
func (c *Cache) RemoveLruBack() {
	if c.cache == nil {
		return
	}

	element := c.lru.Back()

	if element != nil {
		c.lru.Remove(element)
		entry := element.Value.(*entry)
		delete(c.cache, entry.key)

	}
}

// RemoveFifoFront +!
func (c *Cache) RemoveFifoFront() {
	if c.cache == nil {
		return
	}

	element := c.fifo.Front()

	if element != nil {
		c.fifo.Remove(element)
		entry := element.Value.(*entry)
		delete(c.cache, entry.key)

	}
}

// Delete +!
func (c *Cache) Delete(key string) bool {
	if c.cache == nil {
		return false
	}
	if element, hit := c.cache[key]; hit {
		ele := element.Value.(*entry)
		if ele.isLru {
			c.lru.Remove(element)
		} else {
			c.fifo.Remove(element)
		}
		delete(c.cache, key)
		return true
	}
	return false
}
