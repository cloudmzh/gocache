// Programmed by Mazehua (Reference groupcache)
// South China Normal University
// August 12, 2020

package lru

import "container/list"

// Cache :LRU Cache
type Cache struct {
	// MaxEntries :the maximum number of cached entries, and 0 represents no limit
	MaxEntries int
	ll         *list.List
	cache      map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
}

// New +!
func New(maxEntries int) *Cache {
	return &Cache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

// Add :add key and value to the cache
func (c *Cache) Add(key string, value interface{}) {
	if c.cache == nil {
		c.ll = list.New()
		c.cache = make(map[string]*list.Element)
	}

	// If the element exists, put the element in the header of the linked list,
	// point to the *entry type, and update the value
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		element.Value.(*entry).value = value
		return
	}

	//Create a new * Entry node element and place it in the list header,map the key to the list nodes
	element := c.ll.PushFront(&entry{key, value})
	c.cache[key] = element
	//Exceeds the limit, clears the last element
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.RemoveOldest()
	}
}

// Get +!
func (c *Cache) Get(key string) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	//If the cache value exists, the linked list node is placed in the list header and the linked list element is returned
	if element, hit := c.cache[key]; hit {
		c.ll.MoveToFront(element)
		return element.Value.(*entry).value, true
	}
	return
}

// RemoveOldest Clear the last element
func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}

	element := c.ll.Back()

	// Deletes the last element of the list and deletes the cache map
	if element != nil {
		c.ll.Remove(element)
		entry := element.Value.(*entry)
		delete(c.cache, entry.key)

	}
}

// Delete +!
func (c *Cache) Delete(key string) bool {
	if c.cache == nil {
		return false
	}
	// Deletes the element of the list and deletes the cache map
	if element, hit := c.cache[key]; hit {
		c.ll.Remove(element)
		delete(c.cache, key)
		return true
	}
	return false
}

// Len LRU Cache length
func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}
