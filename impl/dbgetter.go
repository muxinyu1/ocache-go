package impl

import (
	"fmt"
	"ocache/ocache"
)

type DbGetter struct {
	db *map[string]*[][]string
}

func (g *DbGetter) Get(groupName string, key string) (ocache.Value, error) {
	v, ok := (*(g.db))[groupName]
	if !ok {
		// db 没有 group
		return nil, fmt.Errorf("No such group: \"%s\"", groupName)
	}
	for _, pair := range *v {
		if key == pair[0] {
			return FromString(pair[1]), nil
		}
	}
	return nil, fmt.Errorf("No such key: \"%s\"", key)
}

func NewDbGetter(db *map[string]*[][]string) *DbGetter {
	return &DbGetter{
		db: db,
	}
}
