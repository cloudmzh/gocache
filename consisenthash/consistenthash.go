// Programmed by Mazehua (Reference groupcache)
// South China Normal University
// August 12, 2020

package consisenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash :Converts the byte string to uint32
type Hash func(data []byte) uint32

// Map :A hash table for storing virtual nodes
type Map struct {
	hash     Hash           //custom hash function,Default:CRC32
	replicas int            //multiples of virtual nodes
	keys     []int          //hash ring
	hashMap  map[int]string //virtual nodes map to real nodes
}

// New +!
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE //By default, CRC32 is used to calculate the Hash value
	}
	return m
}

// Add :add nodes to the hash ring
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// virtual nodes
			//1192.168.0.1、2192.168.0.1、3192.168.0.1
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			//add nodes to the hash ring
			m.keys = append(m.keys, hash)
			//mapping to real nodes
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get :get the nearest node in the hash ring (clockwise)
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	length := len(m.keys)
	return m.hashMap[m.keys[idx%length]]
}
