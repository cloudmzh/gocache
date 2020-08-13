// Programmed by Mazehua (Reference groupcache)
// South China Normal University
// August 12, 2020

package lru

import (
	"log"
	"testing"
)

func TestAddAndGet(t *testing.T) {
	c := New(0)
	c.Add("mzh", 24)
	c.Add("lyj", 23)
	c.Add("mzh", 22)
	log.Println(c.Get("mzh"))
	log.Println(c.Get("lihan"))
}

func TestLen(t *testing.T) {
	c := New(0)
	c.Add("mzh", 24)
	c.Add("lyj", 23)
	c.Add("mzh", 22)
	log.Println(c.Len())
}
