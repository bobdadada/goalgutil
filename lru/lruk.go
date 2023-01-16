package lru

import (
	"container/list"

	cm "goalgutil/macros/cache_macro"
)

// LRU-K中的K代表最近使用的次数，因此LRU可以认为是LRU-1。
// LRU-K的主要目的是为了解决LRU算法“缓存污染”的问题，
// 其核心思想是将“最近使用过1次”的判断标准扩展为“最近使用过K次”。
//
// LRU-K具有LRU的优点，同时能够避免LRU的缺点，实际应用中LRU-2是综合各种因素后最优的选择，
// LRU-3或者更大的K值命中率会高，但适应性差，需要大量的数据访问才能将历史访问记录清除掉。
// LRU-K具有LRU的优点，同时能够避免LRU的缺点，实际应用中LRU-2是综合各种因素后最优的选择，
// LRU-3或者更大的K值命中率会高，但适应性差，需要大量的数据访问才能将历史访问记录清除掉。
//
// 命中率：
// LRU-K降低了“缓存污染”带来的问题，命中率比LRU要高。
//
// 复杂度：
// LRU-K降低了“缓存污染”带来的问题，命中率比LRU要高。
//
// 代价：
// 由于LRU-K还需要记录那些被访问过、但还没有放入缓存的对象，
// 因此内存消耗会比LRU要多；当数据量很大的时候，内存消耗会比较可观。
//

type LRUK struct {
	MaxEntries int
	MaxHitting int

	// OnEvicted optionally specifies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(key cm.Key, value cm.Value)

	ll    *list.List
	count map[cm.Key]int
	cache map[cm.Key]*list.Element
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit.
func NewLRUK(maxEntries, maxHitting int) *LRUK {
	if maxHitting <= 0 {
		panic("MaxHitting must be larger than 0!")
	}
	return &LRUK{
		MaxEntries: maxEntries,
		MaxHitting: maxHitting,
		ll:         list.New(),
		count:      make(map[cm.Key]int),
		cache:      make(map[cm.Key]*list.Element),
	}
}

// Add adds a value to the cache.
func (lruk *LRUK) Add(k cm.Key, v cm.Value) {
	if lruk.cache == nil {
		lruk.cache = make(map[cm.Key]*list.Element)
		lruk.ll = list.New()
		lruk.count = make(map[cm.Key]int)
	}

	if ee, ok := lruk.cache[k]; ok {
		lruk.ll.MoveToFront(ee)
		ee.Value.(*cm.Entry).V = v
		return
	}

	if _, ok := lruk.count[k]; !ok {
		lruk.count[k] = 0
	}
	lruk.count[k] += 1
	if lruk.count[k] < lruk.MaxHitting {
		return
	}

	delete(lruk.count, k)

	if (lruk.MaxEntries > 0) && (lruk.ll.Len() == lruk.MaxEntries) {
		b := lruk.ll.Back()
		k := b.Value.(*cm.Entry).K
		lruk.ll.Remove(b)
		delete(lruk.cache, k)
	}
	ee := lruk.ll.PushFront(&cm.Entry{K: k, V: v})
	lruk.cache[k] = ee
}

// Get looks up a key's value from the cache.
func (lruk *LRUK) Get(k cm.Key) (v cm.Value, ok bool) {
	if lruk.cache == nil {
		return nil, false
	}

	if ee, hit := lruk.cache[k]; hit {
		lruk.ll.MoveToFront(ee)
		return ee.Value.(*cm.Entry).V, true
	}

	if _, ok := lruk.count[k]; !ok {
		lruk.count[k] = 0
	}
	lruk.count[k] += 1

	return nil, false
}

// Remove removes the provided key from the cache.
func (lruk *LRUK) Remove(k cm.Key) {
	if lruk.cache == nil {
		return
	}

	if ee, hit := lruk.cache[k]; hit {
		lruk.ll.Remove(ee)
		delete(lruk.cache, k)
	}
}

// Len returns the number of items in the cache.
func (lruk *LRUK) Len() int {
	if lruk.cache == nil {
		return 0
	}

	return lruk.ll.Len()
}

// Remove removes the provided key from the cache.
func (lruk *LRUK) Clear() {
	if lruk.OnEvicted != nil {
		for _, e := range lruk.cache {
			kv := e.Value.(*cm.Entry)
			lruk.OnEvicted(kv.K, kv.V)
		}
	}

	lruk.ll = nil
	lruk.count = nil

	lruk.cache = nil
}
