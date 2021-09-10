package patterns

import (
	"testing"
)

func TestShardedMap(t *testing.T) {
	const key = "key"
	const value = "value"
	m := NewShardedMap(4)

	m.Set(key, value)
	if v, ok := m.shards[m.destinationShard(key)].Load(key); !ok || v.(string) != value {
		t.Errorf("v, ok == %v, %v, want %v, %v", v, ok, key, value)
	}
	if v, ok := m.Get(key); !ok || v.(string) != value {
		t.Errorf("v, ok == %v, %v, want %v, %v", v, ok, key, value)
	}

	m.Delete(key)
	if _, ok := m.Get(key); ok {
		t.Errorf("ok == %v, want %v", ok, false)
	}
}
