// Programmed by Mazehua (Reference groupcache)
// South China Normal University
// August 13, 2020

package gocache

import (
	"fmt"
	"gocache/consisenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

//
// DefaultBasePath     HTTP default path
// DefaultReplicas     The default number of virtual nodes for a hash ring
const (
	DefaultBasePath = "/_gocache/"
	DefaultReplicas = 50
)

// HTTPPool :HTTP contains all Peer <Peer units>
// PeerPicker interface is implemented
type HTTPPool struct {
	// The basic URL of Peer
	self     string
	basePath string
	mu       sync.Mutex
	// Place the Peer node in a hash ring
	peers *consisenthash.Map
	// Map the remote node to the corresponding httpGetter
	httpGetters map[string]*httpGetter
}

// NewHTTPPool +
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: DefaultBasePath,
	}
}

// Set :Register Peers into cluster (Hash ring)
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//Initializes the node hash ring
	p.peers = consisenthash.New(DefaultReplicas, nil)
	p.peers.Add(peers...)
	// Initializes the list of remote mappings
	p.httpGetters = make(map[string]*httpGetter)
	// eg.    192.168.0.1 -> 192.168.0.1/_gocache/
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{
			baseURL: peer + p.basePath,
		}
	}
}

// PickPeer :Select the Peer node based on the key value
func (p *HTTPPool) PickPeer(key string) (peer PeerGetter, ok bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//consisenthash.Get(...)
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

// Log !+
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// The cache is listening as a server
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.RequestURI() == "/favicon.ico" {
		return
	}

	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	p.Log("%s %s", r.Method, r.URL.Path)

	parts := strings.Split(r.URL.Path[len(p.basePath):], "/")
	if len(parts) < 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	//Parts [0] is a group name
	groupName := parts[0]
	group := GetGroup(groupName)

	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	// .../_gocache/[groupsName]/_del/[keys]   DELETE
	//.../_gocache/[groupsName]/[keys]   GET
	if parts[1] != "_del" {
		view, err := group.Get(parts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteSlice())

	} else if parts[1] == "_del" {
		group.Delete(parts[2])
	}

}

// ----------------------------------------
// httpGetter :Implement client function
// PeerGetter interface is implemented
type httpGetter struct {
	baseURL string
}

// Get :Gets the cache remotely as a client
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	// .../_gocache/[groupsName]/[keys]   GET
	u := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))

	res, err := http.Get(u)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	//The response code is not 200
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Server return code: %d", res.StatusCode)
	}

	//Read all bytes
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Reading response body: %v", err)
	}
	return bytes, nil
}

// Delete :Deletes the cache remotely as a client
func (h *httpGetter) Delete(group string, key string) error {
	// .../_gocache/[groupsName]/_del/[keys]   DELETE
	u := fmt.Sprintf("%v%v/_del/%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))

	res, err := http.Get(u)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	//The response code is not 200
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Reading response body: %v", err)
	}

	return nil
}
