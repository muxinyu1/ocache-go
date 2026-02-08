package ocache

import (
	"container/list"
	"fmt"
	"math/rand/v2"
	"ocache/bloom"
)

type Group struct {
	groupName  string
	mainCache  *Cache
	hotCache   *Cache
	peerPicker PeerPicker // TODO thread safe
	getter     Getter     // TODO thread safe
	table      Table
	bf         *bloom.BloomFilter
}

func NewGroup(groupName string, maxBytes int, peerPicker PeerPicker, getter Getter, bf *bloom.BloomFilter) *Group {
	return &Group{
		groupName: groupName,
		mainCache: &Cache{
			lru: &lruCache{
				m:        make(map[string]*list.Element),
				capacity: maxBytes,
			},
		},
		hotCache: &Cache{
			lru: &lruCache{
				m:        make(map[string]*list.Element),
				capacity: maxBytes * 2,
			},
		},
		peerPicker: peerPicker,
		getter:     getter,
		bf:         bf,
	}
}

func (g *Group) getLocally(groupName string, key string) (Value, error) {
	v, err := g.getter.Get(g.groupName, key)
	if err != nil {
		return nil, err
	}
	g.mainCache.Add(key, v)
	return v, nil
}

func (g *Group) Get(key string) (Value, error) {
	// main cache
	if v, err := g.mainCache.Get(key); err == nil {
		return v, nil
	}
	// hot cache
	if v, err := g.hotCache.Get(key); err == nil {
		return v, nil
	}

	// bloom
	if g.bf != nil && !g.bf.Test([]byte(key)) {
		// fmt.Println("no such data found by bloom")
		return nil, fmt.Errorf("key not found by bloom filter")
	}

	// peer
	if g.peerPicker != nil {
		cli := g.peerPicker.PickPeer(key)
		if cli != nil {
			v, err := g.table.Do(g.groupName, key, cli.Get)
			if err == nil && rand.IntN(10) == 0 {
				// 1 / 10概率进二级缓存
				// 热点key 二级缓存
				g.hotCache.Add(key, v)
			}
			return v, err
		}
	}
	// db
	return g.table.Do(g.groupName, key, g.getLocally)
}
