package ocache

import "container/list"

type Group struct {
	groupName  string
	cache      Cache
	peerPicker PeerPicker // TODO thread safe
	getter     Getter     // TODO thread safe
	table      Table
}

func NewGroup(groupName string, maxBytes int, peerPicker PeerPicker, getter Getter) *Group {
	return &Group{
		groupName: groupName,
		cache: Cache{
			lru: lruCache{
				m:        make(map[string]*list.Element),
				capacity: maxBytes,
			},
		},
		peerPicker: peerPicker,
		getter:     getter,
	}
}

func (g *Group) getLocally(string, key string) (Value, error) {
	v, err := g.getter.Get(g.groupName, key)
	if err != nil {
		return nil, err
	}
	g.cache.Add(key, v)
	return v, nil
}

func (g *Group) Get(key string) (Value, error) {
	// cache
	if v, err := g.cache.Get(key); err == nil {
		return v, nil
	}
	// peer
	if g.peerPicker != nil {
		cli := g.peerPicker.PickPeer(key)
		if cli != nil {
			return g.table.Do(g.groupName, key, cli.Get)
		}
	}
	// db
	return g.table.Do(g.groupName, key, g.getLocally)
}
