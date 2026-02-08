package impl

import (
	"fmt"
	"net/http"
	"ocache/bloom"
	"ocache/ocache"
	"strings"
)

type Peer struct {
	groups map[string]*ocache.Group
}

func NewPeer(
	groupNames []string,
	groupBytes []int,
	idx int,
	total int,
	portBase int,
	replicas int,
	hasher Hasher,
	getter ocache.Getter,
	bf *bloom.BloomFilter) *Peer {
	m := make(map[string]*ocache.Group)
	clients := make([]ocache.Client, 0, total)
	peers := make([]string, 0, total)
	for i := range total {
		peer := fmt.Sprintf("http://127.0.0.1:%d", portBase+i)
		peers = append(peers, peer)
		if i != idx {
			clients = append(clients, &HttpClient{
				baseUrl: peer,
			})
			continue
		}
		clients = append(clients, nil)
	}
	peerPicker := NewHashPeerPicker(
		replicas,
		peers,
		clients,
		hasher,
	)
	for i, groupName := range groupNames {
		m[groupName] = ocache.NewGroup(
			groupName,
			groupBytes[i],
			peerPicker,
			getter,
			bf,
		)
	}
	return &Peer{
		groups: m,
	}
}

func (p *Peer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 3 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Bad request"))
		return
	}
	groupName, key := parts[1], parts[2]
	group, ok := p.groups[groupName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Group \"%s\" not found", groupName)
		return
	}
	value, err := group.Get(key)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			w.WriteHeader(http.StatusNotFound)
			w.Write(value.AsBytes())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(value.AsBytes())
}
