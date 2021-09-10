package patterns

import (
	"crypto/sha1"
	"sync"
)

type ShardedMap struct {
	shards []*sync.Map
}

func NewShardedMap(n int) *ShardedMap {
	m := &ShardedMap{shards: make([]*sync.Map, 0, n)}
	for i := 0; i < n; i++ {
		m.shards = append(m.shards, new(sync.Map))
	}
	return m
}

func (m *ShardedMap) Set(key string, value interface{}) {
	m.shards[m.destinationShard(key)].Store(key, value)
}

func (m *ShardedMap) Get(key string) (value interface{}, ok bool) {
	value, ok = m.shards[m.destinationShard(key)].Load(key)
	return
}

func (m *ShardedMap) Delete(key string) {
	m.shards[m.destinationShard(key)].Delete(key)
}

func (m *ShardedMap) destinationShard(key string) int {
	sum := sha1.Sum([]byte(key))
	return int(sum[sha1.Size/2-1]) % len(m.shards)
}
