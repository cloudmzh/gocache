// Programmed by Mazehua (Reference groupcache)
// South China Normal University
// August 12, 2020

package singleflight

import "sync"

// call :an ongoing or completed request.
// Use the Sync.waitGroup lock to avoid reentrancy.
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group :main data structure, managing requests for different keys
type Group struct {
	mu sync.Mutex // protects m
	m  map[string]*call
}

// Do :For the specified parameter (key),
// no matter how many gocoroutines and no matter how many calls,
// the function (fn) is run only once before returning the result
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	//Found that the function is running, blocking the gocoroutine
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	g.m[key] = c //add key to m
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key) //delete key from m
	g.mu.Unlock()

	return c.val, c.err

}
