// Programmed by Mazehua (Reference groupcache)
// South China Normal University
// August 13, 2020

package gocache

// PeerPicker +!
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter +!
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
	Delete(group string, key string) error
}
