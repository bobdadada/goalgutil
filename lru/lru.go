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

package lru

import (
	"container/list"
)

type Key any

type entry struct {
	key   Key
	value any
}

type LRU struct {
	MaxEntries int
	ll         *list.List
	cache      map[Key]*list.Element
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit.
func NewLRU(maxEntries int) *LRU {
	return &LRU{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[Key]*list.Element),
	}
}

// Add adds a value to the cache.
func (lru *LRU) Add(k Key, v any) {
	if ee, ok := lru.cache[k]; ok {
		lru.ll.MoveToFront(ee)
		ee.Value.(*entry).value = v
		return
	}
	if (lru.MaxEntries > 0) && (lru.ll.Len() == lru.MaxEntries) {
		b := lru.ll.Back()
		k := b.Value.(*entry).key
		lru.ll.Remove(b)
		delete(lru.cache, k)
	}
	ee := lru.ll.PushFront(&entry{k, v})
	lru.cache[k] = ee
}

// Get looks up a key's value from the cache.
func (lru *LRU) Get(k Key) (v any, ok bool) {
	if ee, hit := lru.cache[k]; hit {
		lru.ll.MoveToFront(ee)
		return ee.Value.(*entry).value, true
	}
	return nil, false
}

// Remove removes the provided key from the cache.
func (lru *LRU) Remove(k Key) {
	if ee, hit := lru.cache[k]; hit {
		lru.ll.Remove(ee)
		delete(lru.cache, k)
	}
}

// Len returns the number of items in the cache.
func (lru *LRU) Len() int {
	return lru.ll.Len()
}

// Remove removes the provided key from the cache.
func (lru *LRU) Clear() {
	lru.ll = list.New()
	lru.cache = make(map[Key]*list.Element)
}

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
	K          int
	ll         *list.List
	cache      map[Key]*list.Element
	count      map[Key]int
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit.
func NewLRUK(maxEntries, k int) *LRUK {
	if k <= 0 {
		panic("k must be larger than 0!")
	}
	return &LRUK{
		MaxEntries: maxEntries,
		K:          k,
		ll:         list.New(),
		cache:      make(map[Key]*list.Element),
		count:      make(map[Key]int),
	}
}

// Add adds a value to the cache.
func (lruk *LRUK) Add(k Key, v any) {
	if ee, ok := lruk.cache[k]; ok {
		lruk.ll.MoveToFront(ee)
		ee.Value.(*entry).value = v
		return
	}

	if _, ok := lruk.count[k]; !ok {
		lruk.count[k] = 0
	}
	lruk.count[k] += 1
	if lruk.count[k] < lruk.K {
		return
	}

	delete(lruk.count, k)

	if (lruk.MaxEntries > 0) && (lruk.ll.Len() == lruk.MaxEntries) {
		b := lruk.ll.Back()
		k := b.Value.(*entry).key
		lruk.ll.Remove(b)
		delete(lruk.cache, k)
	}
	ee := lruk.ll.PushFront(&entry{k, v})
	lruk.cache[k] = ee
}

// Get looks up a key's value from the cache.
func (lruk *LRUK) Get(k Key) (v any, ok bool) {
	if ee, hit := lruk.cache[k]; hit {
		lruk.ll.MoveToFront(ee)
		return ee.Value.(*entry).value, true
	}

	if _, ok := lruk.count[k]; !ok {
		lruk.count[k] = 0
	}
	lruk.count[k] += 1

	return nil, false
}

// Remove removes the provided key from the cache.
func (lruk *LRUK) Remove(k Key) {
	if ee, hit := lruk.cache[k]; hit {
		lruk.ll.Remove(ee)
		delete(lruk.cache, k)
	}
}

// Len returns the number of items in the cache.
func (lruk *LRUK) Len() int {
	return lruk.ll.Len()
}

// Remove removes the provided key from the cache.
func (lruk *LRUK) Clear() {
	lruk.ll = list.New()
	lruk.cache = make(map[Key]*list.Element)
	lruk.count = make(map[Key]int)
}
