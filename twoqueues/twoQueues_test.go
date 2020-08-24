package twoqueues

import ( //明天添加
	"log"
	"testing"
)

func TestAddAndGet(t *testing.T) {
	c := New(10, 10)
	log.Println("start")
	c.Add("mzh", "16")
	log.Println(c.lru.Len(), c.fifo.Len())

	c.Add("mzh", "17")
	log.Println(c.lru.Len(), c.fifo.Len())
	log.Println(c.Get("mzh"))

	c.Add("xqh", "55")
	log.Println(c.lru.Len(), c.fifo.Len())
	log.Println(c.Get("xqh"))
	log.Println(c.lru.Len(), c.fifo.Len())
}

func TestRemoveLruAndFifo(t *testing.T) {
	c := New(3, 3)
	c.Add("mzh", "17")
	c.Add("xqh", "17")
	c.Add("cx", "17")
	log.Println(c.lru.Len(), c.fifo.Len()) // 0  3
	c.Add("lyj", "17")
	log.Println(c.lru.Len(), c.fifo.Len()) // 0  3
	c.Get("xqh")
	log.Println(c.lru.Len(), c.fifo.Len()) // 1  2
	c.Add("sam", "22")
	log.Println(c.lru.Len(), c.fifo.Len()) // 1  3

	c.Get("mzh")                           //invaluable
	log.Println(c.lru.Len(), c.fifo.Len()) // 1  3
	c.Get("lyj")
	c.Add("nwj", "22")
	log.Println(c.lru.Len(), c.fifo.Len()) // 2  3
	c.Get("sam")
	log.Println(c.lru.Len(), c.fifo.Len()) // 3  2
	c.Get("nwj")
	log.Println(c.lru.Back().Value)
	c.Get("cx")
	log.Println(c.lru.Len(), c.fifo.Len()) // 3  1
}

func TestDelete(t *testing.T) {
	c := New(3, 3)
	c.Add("mzh", "17")
	c.Add("xqh", "17")
	c.Add("cx", "17")
	c.Get("cx")
	log.Println(c.lru.Len(), c.fifo.Len())
	c.Delete("cx")
	c.Delete("xqh")
	log.Println(c.lru.Len(), c.fifo.Len())
}
