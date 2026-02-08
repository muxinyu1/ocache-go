package bloom

import (
	"iter"
	"math"
)

// 对于固定value 和 key，hash结果固定
type MapBytes2Uint64 func(value []byte, seed uint64) uint64

// READONLY, thus thread safe
type BloomFilter struct {
	bits   []byte
	m      uint64 // 多少位
	k      uint64 // 多少个hash函数
	hasher MapBytes2Uint64
}

type Generator interface {
	All() iter.Seq[[]byte]
}

func NewBloomFilter(n uint64, p float64, gen Generator, hasher MapBytes2Uint64) *BloomFilter {
	m := uint64(math.Ceil(-(float64(n) * math.Log(p)) / (math.Ln2 * math.Ln2)))
	k := uint64(math.Ceil(math.Ln2 * float64(m) / float64(n)))
	byteNum := (m + 7) / 8

	bits := make([]byte, byteNum)

	// 在构造后只读，在这里填充bits
	for value := range gen.All() {
		for i := range k {
			idx := hasher(value, i) % m
			bits[idx/8] |= (1 << (idx % 8))
		}
	}

	return &BloomFilter{
		bits:   bits,
		m:      m,
		k:      k,
		hasher: hasher,
	}
}

func (b *BloomFilter) Test(key []byte) bool {
	for i := range b.k {
		hashed := b.hasher(key, i)
		idx := hashed % b.m
		if b.bits[idx/8]&(1<<(idx%8)) == 0 {
			return false
		}
	}
	return true
}
