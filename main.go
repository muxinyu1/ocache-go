package main

import (
	"fmt"
	"hash/crc32"
	"math/rand/v2"
	"net/http"
	"ocache/bloom"
	"ocache/docker"
	"ocache/impl"
	"ocache/ocache"
	"sync"
)

type stringClient string

func (s stringClient) Get(group string, key string) (ocache.Value, error) {
	return nil, nil
}

func BloomHash(value []byte, seed uint64) uint64 {
	// 简单的多重哈希：将 seed 混入数据中
	h := crc32.NewIEEE()
	fmt.Fprintf(h, "%d", seed)
	h.Write(value)
	return uint64(h.Sum32())
}

const PEER_NUM = 2
const REPLICAS = 4
const PORT_BASE = 1024
const DATA_NUM = 10000
const P = 0.001

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
	// db := makeDb(groupNames, keys)
	// for k, v := range *db {
	// 	fmt.Printf("Group: %s\n", k)
	// 	for _, kv := range *v {
	// 		fmt.Printf("\t%s: %s\n", kv[0], kv[1])
	// 	}
	// }
	getter, err := docker.NewMysqlGetter("root:root@tcp(127.0.0.1:3306)/ocache")
	if err != nil {
		panic(err)
	}

	// 验证 Scores group key 分布
	fmt.Println("Key distribution for group 'Scores':")
	pickerPeers := make([]string, 0, PEER_NUM)
	pickerClients := make([]ocache.Client, 0, PEER_NUM)
	for i := range PEER_NUM {
		addr := fmt.Sprintf("http://127.0.0.1:%d", PORT_BASE+i)
		pickerPeers = append(pickerPeers, addr)
		pickerClients = append(pickerClients, stringClient(addr))
	}
	picker := impl.NewHashPeerPicker(REPLICAS, pickerPeers, pickerClients, crc32.ChecksumIEEE)
	for _, key := range keys {
		cli := picker.PickPeer(key)
		fmt.Printf("Key: %s -> Peer: %s\n", key, cli)
	}

	// bloom filter
	_ = bloom.NewBloomFilter(DATA_NUM, P, getter, BloomHash)

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
				nil,
			)
			err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", PORT_BASE+idx), p)
			if err != nil {
				panic(err)
			}
		}(i)
	}
	wg.Wait()
}
