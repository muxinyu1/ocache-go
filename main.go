package main

import (
	"fmt"
	"hash/crc32"
	"math/rand/v2"
	"net/http"
	"ocache/impl"
	"sync"
)

const PEER_NUM = 2
const REPLICAS = 4
const PORT_BASE = 1024

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func stringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func makeDb(groupNames []string, keys []string) *map[string]*[][]string {
	m := make(map[string]*[][]string)
	for _, groupName := range groupNames {
		kv := make([][]string, 0)
		for _, key := range keys {
			kv = append(kv, []string{key, stringWithCharset(32)})
		}
		m[groupName] = &kv
	}
	return &m
}

func main() {

	groupNames := []string{"Scores", "Tsinghua", "ByteDance", "学生"}
	keys := []string{"mxy", "mqs", "sheeep", "random", "主席"}
	l := len(groupNames)
	groupBytes := make([]int, 0, l)
	for range l {
		groupBytes = append(groupBytes, 256) // 先都设置成256字节
	}

	// make db
	db := makeDb(groupNames, keys)
	for k, v := range *db {
		fmt.Printf("Group: %s\n", k)
		for _, kv := range *v {
			fmt.Printf("\t%s: %s\n", kv[0], kv[1])
		}
	}
	getter := impl.NewDbGetter(db)

	var wg sync.WaitGroup
	for i := range PEER_NUM {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			p := impl.NewPeer(
				groupNames,
				groupBytes,
				idx,
				PEER_NUM,
				PORT_BASE,
				REPLICAS,
				crc32.ChecksumIEEE,
				getter,
			)
			err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", PORT_BASE+idx), p)
			if err != nil {
				panic(err)
			}
		}(i)
	}
	wg.Wait()
}
