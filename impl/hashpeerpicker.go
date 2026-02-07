package impl

import (
	"ocache/ocache"
	"slices"
	"sort"
	"strconv"
	"sync"
)

type Hasher func([]byte) uint32

type HashRingPeerPicker struct {
	mtx    sync.RWMutex
	nodes  []uint32
	m      map[uint32]ocache.Client
	hasher Hasher
}

// peers和clients长度应当相等并且一一对应
//
// 将本身对应的client设置为nil
func NewHashPeerPicker(replicas int, peers []string, clients []ocache.Client, hasher Hasher) *HashRingPeerPicker {

	nodes := make([]uint32, 0)
	m := make(map[uint32]ocache.Client)
	for i, peer := range peers {
		for j := range replicas {
			hashed := hasher([]byte(peer + strconv.Itoa(j)))
			nodes = append(nodes, hashed)
			m[hashed] = clients[i]
		}
	}
	slices.Sort(nodes)
	return &HashRingPeerPicker{
		nodes:  nodes,
		m:      m,
		hasher: hasher,
	}
}

func (p *HashRingPeerPicker) PickPeer(key string) ocache.Client {
	hashed := p.hasher([]byte(key))
	l := len(p.nodes)
	idx := sort.Search(l, func(i int) bool { return p.nodes[i] > hashed })
	if idx == l {
		idx = 0
	}
	return p.m[p.nodes[idx]]
}
