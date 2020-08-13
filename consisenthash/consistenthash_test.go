// Programmed by Mazehua (Reference groupcache)
// South China Normal University
// August 12, 2020

package consisenthash

import (
	"log"
	"testing"
)

func TestAdd(t *testing.T) {
	m := New(4, nil)
	m.Add("192.168.1.1")
	log.Println(m.hashMap)
	log.Println(m.keys)
}

func TestGet(t *testing.T) {
	m := New(50, nil)
	m.Add("192.168.1.1")
	m.Add("192.168.1.2")
	m.Add("192.168.1.3")
	log.Println(m.hashMap)
	log.Println(m.keys)
	log.Println(m.Get("mazeshua"))
}
