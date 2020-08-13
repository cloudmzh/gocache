// Programmed by Mazehua (Reference groupcache)
// South China Normal University
// August 12, 2020

package singleflight

import (
	"log"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func TestSF(t *testing.T) {
	group := Group{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go group.Do("test", func() (interface{}, error) {
			defer wg.Add(-10) //once
			log.Println("sleep.....")
			time.Sleep(time.Second * 1)
			log.Println("awake.....")
			return nil, nil
		})

	}

	wg.Wait()

}
