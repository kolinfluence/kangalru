package lruxbytes

import (
	"testing"

	cx "github.com/cloudxaas/gocx"
	"github.com/phuslu/lru"
	"github.com/maypok86/otter"
)

func BenchmarkOtterSet(b *testing.B) {
	cache, err := otter.MustBuilder[string, string](1024 * 100).
		CollectStats().
		Cost(func(key string, value string) uint32 {
			return 1
		}).
		Build()
	if err != nil {
		b.Fatal(err)
	}

	keys := make([]string, 100000)
	values := make([]string, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
		values[i] = "value"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(keys[i%100000], values[i%100000])
	}
}

func BenchmarkOtterGet(b *testing.B) {
	cache, err := otter.MustBuilder[string, string](1024 * 100).
		CollectStats().
		Cost(func(key string, value string) uint32 {
			return 1
		}).
		Build()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 100000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), "value")
	}
	keys := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = fmt.Sprintf("key%d", i%100000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(keys[i])
	}
}

func BenchmarkOtterDelete(b *testing.B) {
	cache, err := otter.MustBuilder[string, string](1024 * 100).
		CollectStats().
		Cost(func(key string, value string) uint32 {
			return 1
		}).
		Build()
	if err != nil {
		b.Fatal(err)
	}

	keys := make([]string, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
		cache.Set(keys[i], "value")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Delete(keys[i%100000])
	}
}

// FNV-1a hash function for byte slices
func FNV1aHash(key []byte) uint32 {
	const (
		offset32 uint32 = 2166136261
		prime32  uint32 = 16777619
	)
	var hash uint32 = offset32
	for _, c := range key {
		hash ^= uint32(c)
		hash *= prime32
	}
	return hash
}

func BenchmarkPhusluLRUSet(b *testing.B) {
	cache := lru.NewLRUCache[string, []byte](1024 * 100)
	keys := make([][]byte, 100000)
	values := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		values[i] = make([]byte, 1024) // 1 KB values
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(cx.B2s(keys[i%1000]), values[i%1000])
	}
}

func BenchmarkPhusluLRUGet(b *testing.B) {
	cache := lru.NewLRUCache[string, []byte](1024 * 100)
	for i := 0; i < 100000; i++ {
		cache.Set(cx.B2s([]byte{byte(i)}), make([]byte, 1024)) // 1 KB values
	}
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = []byte{byte(i % 100000)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(cx.B2s(keys[i]))
	}
}

func BenchmarkPhusluLRUDelete(b *testing.B) {
	cache := lru.NewLRUCache[string, []byte](1024 * 100)
	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		cache.Set(cx.B2s(keys[i]), make([]byte, 1024)) // 1 KB values
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Delete(cx.B2s(keys[i%100000]))
	}
}

func BenchmarkCXLRUBytesSet(b *testing.B) {
	cache := NewLRUCache(1024*100, 1, FNV1aHash)
	keys := make([][]byte, 100000)
	values := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		values[i] = make([]byte, 1024) // 1 KB values
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(keys[i%100000], values[i%100000])
	}
}

func BenchmarkCXLRUBytesGet(b *testing.B) {
	cache := NewLRUCache(1024*100, 1, FNV1aHash)
	for i := 0; i < 100000; i++ {
		cache.Set([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
	}
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = []byte{byte(i % 100000)}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(keys[i])
	}
}

func BenchmarkCXLRUBytesDel(b *testing.B) {
	cache := NewLRUCache(1024*100, 1, FNV1aHash)
	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		cache.Set(keys[i], make([]byte, 1024)) // 1 KB values
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Del(keys[i%100000])
	}
}

func BenchmarkCXLRUBytesSetParallel(b *testing.B) {
	cache := NewLRUCache(1024*100, 1, FNV1aHash)
	keys := make([][]byte, 100000)
	values := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		values[i] = make([]byte, 1024) // 1 KB values
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			cache.Set(keys[i%100000], values[i%100000])
		}
	})
}

func BenchmarkCXLRUBytesGetParallel(b *testing.B) {
	cache := NewLRUCache(1024*100, 1, FNV1aHash)
	for i := 0; i < 100000; i++ {
		cache.Set([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
	}

	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			_, _ = cache.Get(keys[i%100000])
		}
	})
}

func BenchmarkCXLRUBytesDelParallel(b *testing.B) {
	cache := NewLRUCache(1024*100, 1, FNV1aHash)
	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		cache.Set(keys[i], make([]byte, 1024)) // 1 KB values
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			cache.Del(keys[i%100000])
		}
	})
}

func BenchmarkCXLRUBytesShardedSet(b *testing.B) {
	cache := NewShardedCache(16, 1024*100, 1, FNV1aHash)
	keys := make([][]byte, 100000)
	values := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		values[i] = make([]byte, 1024) // 1 KB values
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(keys[i%100000], values[i%100000])
	}
}

func BenchmarkCXLRUBytesShardedGet(b *testing.B) {
	cache := NewShardedCache(16, 1024*100, 1, FNV1aHash)
	for i := 0; i < 100000; i++ {
		cache.Set([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
	}
	keys := make([][]byte, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(keys[i])
	}
}

func BenchmarkCXLRUBytesShardedDel(b *testing.B) {
	cache := NewShardedCache(16, 1024*100, 1, FNV1aHash)
	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		cache.Set(keys[i], make([]byte, 1024)) // 1 KB values
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Del(keys[i%100000])
	}
}

func BenchmarkCXLRUBytesShardedSetParallel(b *testing.B) {
	cache := NewShardedCache(16, 1024*100, 1, FNV1aHash)
	keys := make([][]byte, 100000)
	values := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		values[i] = make([]byte, 1024) // 1 KB values
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			cache.Set(keys[i%100000], values[i%100000])
		}
	})
}

func BenchmarkCXLRUBytesShardedGetParallel(b *testing.B) {
	cache := NewShardedCache(16, 1024*100, 1, FNV1aHash)
	for i := 0; i < 100000; i++ {
		cache.Set([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
	}

	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			_, _ = cache.Get(keys[i%100000])
		}
	})
}

func BenchmarkCXLRUBytesShardedDelParallel(b *testing.B) {
	cache := NewShardedCache(16, 1024*100, 1, FNV1aHash)
	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		cache.Set(keys[i], make([]byte, 1024)) // 1 KB values
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			cache.Del(keys[i%100000])
		}
	})
}
