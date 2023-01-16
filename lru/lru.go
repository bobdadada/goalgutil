package lru

import (
	"container/list"

	cm "goalgutil/macros/cache_macro"
)

// LRU (Least recently used) 算法根据数据的历史访问记录来进行淘汰数据，
// 其核心思想是“如果数据最近被访问过，那么将来被访问的几率也更高”。
//
// 命中率：
// 当存在热点数据时，LRU的效率很好，但偶发性的、周期性的批量操作会
// 导致LRU命中率急剧下降，缓存污染情况比较严重。
//
// 复杂度：
// 实现简单
//
// 代价：
// 命中时需要遍历链表，找到命中的数据块索引，然后需要将数据移到头部。
//
// 作者：jiangmo
// 链接：https://www.jianshu.com/p/d533d8a66795
// 来源：简书
// 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
//

type LRU struct {
	MaxEntries int

	// OnEvicted optionally specifies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(k cm.Key, v cm.Value)

	ll    *list.List
	cache map[cm.Key]*list.Element
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit.
func NewLRU(maxEntries int) *LRU {
	return &LRU{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[cm.Key]*list.Element),
	}
}

// Add adds a value to the cache.
func (lru *LRU) Add(k cm.Key, v cm.Value) {
	if lru.cache == nil {
		// `make` may fail
		lru.cache = make(map[cm.Key]*list.Element)
		lru.ll = list.New()
	}

	if ee, ok := lru.cache[k]; ok {
		lru.ll.MoveToFront(ee)
		ee.Value.(*cm.Entry).V = v
		return
	}
	if (lru.MaxEntries > 0) && (lru.ll.Len() == lru.MaxEntries) {
		b := lru.ll.Back()
		k := b.Value.(*cm.Entry).K
		lru.ll.Remove(b)
		delete(lru.cache, k)
	}
	ee := lru.ll.PushFront(&cm.Entry{K: k, V: v})
	lru.cache[k] = ee
}

// Get looks up a key's value from the cache.
func (lru *LRU) Get(k cm.Key) (v any, ok bool) {
	if lru.cache == nil {
		return nil, false
	}

	if ee, hit := lru.cache[k]; hit {
		lru.ll.MoveToFront(ee)
		return ee.Value.(*cm.Entry).V, true
	}
	return nil, false
}

// Remove removes the provided key from the cache.
func (lru *LRU) Remove(k cm.Key) {
	if lru.cache == nil {
		return
	}

	if ee, hit := lru.cache[k]; hit {
		lru.ll.Remove(ee)
		delete(lru.cache, k)
	}
}

// Len returns the number of items in the cache.
func (lru *LRU) Len() int {
	if lru.cache == nil {
		return 0
	}

	return lru.ll.Len()
}

// Remove removes the provided key from the cache.
func (lru *LRU) Clear() {
	if lru.OnEvicted != nil {
		for _, e := range lru.cache {
			kv := e.Value.(*cm.Entry)
			lru.OnEvicted(kv.K, kv.V)
		}
	}
	lru.ll = nil
	lru.cache = nil
}
