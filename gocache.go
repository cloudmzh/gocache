// Programmed by Mazehua (Reference groupcache)
// South China Normal University
// August 13, 2020

package gocache

import (
	"fmt"
	"gocache/singleflight"
	"gocache/twoqueues"
	"log"
	"sync"
)

//---------------------cache------------------
// The Cache is a structure that encapsulates
//the LRU cache and supports concurrent access

// cache :The LRU is packaged
type cache struct {
	mu             sync.Mutex
	tq             *twoqueues.Cache
	maxLruEntries  int
	maxFifoEntries int
}

// k-v: string-ByteView

// add :Concurrent secure LRU add
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.tq == nil {
		c.tq = twoqueues.New(c.maxLruEntries, c.maxFifoEntries) //Delayed initialization
	}
	c.tq.Add(key, value)
}

// get :Concurrent secure LRU get
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.tq == nil {
		return
	}
	if v, ok := c.tq.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

// delete :Concurrent secure LRU delete
func (c *cache) delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.tq == nil {
		return false
	}
	return c.tq.Delete(key)
}

//-----------------Getter---------------------
//A Getter is a callback function interface
//for accessing data sources outside the cache

// Getter :Callback function interface
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunction type
type GetterFunction func(key string) ([]byte, error)

// Get +!
func (f GetterFunction) Get(key string) ([]byte, error) {
	return f(key)
}

//-----------------Group-------------------
//The Group is the cache namespace

// Group +!
type Group struct {
	name   string
	getter Getter //Getters are used to obtain the local data source
	data   cache
	peers  PeerPicker //The cluster to which a group belongs
	loader *singleflight.Group
}

var (
	mu     sync.RWMutex //protects groups
	groups = make(map[string]*Group)
)

// NewGroup +!
func NewGroup(name string, lMaxNum int, fMaxNum int, getter Getter) *Group {
	if getter == nil {
		panic("the Group getter is nil .")
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:   name,
		getter: getter,
		data:   cache{maxLruEntries: lMaxNum, maxFifoEntries: fMaxNum},
		loader: &singleflight.Group{},
	}

	groups[name] = g
	return g
}

// GetGroup +!
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

// Get :The cache is first fetched from the peer
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("Group err: key is \"\"")
	}
	if value, ok := g.data.get(key); ok {
		log.Printf("[Group %s] Hit cache,key is %s\n", g.name, key)
		return value, nil
	}
	return g.load(key)
}

// load :Look for cache lookups in other peers
func (g *Group) load(key string) (value ByteView, err error) {

	// Run the function only once before the result is returned
	bv, err := g.loader.Do(key, func() (interface{}, error) {

		if g.peers != nil {
			// Based on the key value, the target host is found on the hash ring
			if peer, ok := g.peers.PickPeer(key); ok {
				// Search for cache on remote target host
				if bytes, err := peer.(*httpGetter).Get(g.name, key); err == nil {
					return ByteView{b: bytes}, nil
				}
				log.Println("[Group] Failed to get from peer", g.name)
			}
		}
		return g.getLocally(key)
	})

	if err == nil {
		return bv.(ByteView), nil
	}
	return
}

// getLocally :Look for data from a local data source
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	c := make([]byte, len(bytes))
	copy(c, bytes)
	value := ByteView{b: c}
	g.Set(key, value)
	return value, nil
}

// Set !+
func (g *Group) Set(key string, value ByteView) {
	g.data.add(key, value)
	log.Printf("[Group %s] Add cache,key is %s\n", g.name, key)
}

// Delete +!
func (g *Group) Delete(key string) {
	//First remove the local cache value
	g.data.delete(key)

	if g.peers != nil {
		// Based on the key value, the target host is found on the hash ring
		if peer, ok := g.peers.PickPeer(key); ok {
			// Delete the value of the host on the hash ring
			err := peer.Delete(g.name, key)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// RegisterPeers :Register Peer in HTTPPool
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}
