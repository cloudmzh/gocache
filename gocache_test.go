// Programmed by Mazehua (Reference groupcache)
// South China Normal University
// August 13, 2020
package gocache

import (
	"fmt"
	"log"
	"testing"
)

//
var db = map[string]string{
	"mzh": "scnu",
	"cx":  "scnu",
} //db source

func getFromDB(key string) ([]byte, error) {
	if v, ok := db[key]; ok {
		return []byte(v), nil
	}
	return nil, fmt.Errorf("DB ERR: key is null")
}

func TestGet(t *testing.T) {
	getter := GetterFunction(getFromDB)
	g := NewGroup("new namespace", 2<<8, 2<<8, getter)
	g = GetGroup("new namespace")

	bv, err := g.Get("mzh")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(bv.b))
	bv, err = g.Get("mzh")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(bv.b))
	bv = ByteView{b: []byte("cccccccccccccc")}
	g.Set("cx", bv)
	bv, err = g.Get("cx")
	if err != nil {
		log.Println(err)
	}
}
