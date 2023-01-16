package lru

import (
	"container/list"

	cm "goalgutil/macros/cache_macro"
)

// Two queues（以下使用2Q代替）算法类似于LRU-2，不同点在于2Q将LRU-2算法中的访问历史队列（注意这不是缓存数据的）改为一个FIFO缓存队列，
// 即：2Q算法有两个缓存队列，一个是FIFO队列，一个是LRU队列。 当数据第一次访问时，2Q算法将数据缓存在FIFO队列里面，当数据第二次被访问时，
// 则将数据从FIFO队列移到LRU队列里面，两个队列各自按照自己的方法淘汰数据。
//
// 命中率：
// 2Q算法的命中率要高于LRU。
//
// 复杂度：
// 需要两个队列，但两个队列本身都比较简单。
//
// 代价：
// FIFO和LRU的代价之和。2Q算法和LRU-2算法命中率类似，内存消耗也比较接近，但对于最后缓存的数据来说，
// 2Q会减少一次从原始存储读取数据或者计算数据的操作。
//

type LRU2Q struct {
	MaxEntries int

	// OnEvicted optionally specifies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(k cm.Key, v cm.Value)

	ll     *list.List
	fifo   *list.List
	cache  map[cm.Key]*list.Element
	qcount map[cm.Key]*list.Element
}

// New creates a new Cache. maxEntries must be larger than zero.
func NewLRU2Q(maxEntries int) *LRU2Q {
	if maxEntries <= 0 {
		panic("maxEntries must be larger than 0!")
	}

	return &LRU2Q{
		MaxEntries: maxEntries,
		ll:         list.New(),
		fifo:       list.New(),
		cache:      make(map[cm.Key]*list.Element),
		qcount:     make(map[cm.Key]*list.Element),
	}
}

// Add adds a value to the cache.
func (lru2q *LRU2Q) Add(k cm.Key, v cm.Value) {
	if lru2q.cache == nil {
		// `make` may fail
		lru2q.cache = make(map[cm.Key]*list.Element)
		lru2q.ll = list.New()
	}

	if lru2q.qcount == nil {
		lru2q.qcount = make(map[cm.Key]*list.Element)
		lru2q.fifo = list.New()
	}

	// key exists in LRU cache
	if ee, ok := lru2q.cache[k]; ok {
		lru2q.ll.MoveToFront(ee)
		ee.Value.(*cm.Entry).V = v
		return
	}

	// key exists in FIFO
	if ee, ok := lru2q.qcount[k]; ok {

		kv := ee.Value.(*cm.Entry)

		// delete the element in FIFO
		lru2q.fifo.Remove(ee)
		delete(lru2q.qcount, k)

		// add the element into LRU
		if lru2q.ll.Len() == lru2q.MaxEntries {
			b := lru2q.ll.Back()
			k := b.Value.(*cm.Entry).K
			lru2q.ll.Remove(b)
			delete(lru2q.cache, k)
		}
		lru2q.cache[k] = lru2q.ll.PushFront(kv)

		return
	}

	// add key into FIFO
	if lru2q.fifo.Len() == lru2q.MaxEntries {
		b := lru2q.fifo.Back()
		k := b.Value.(*cm.Entry).K
		lru2q.fifo.Remove(b)
		delete(lru2q.qcount, k)
	}
	lru2q.qcount[k] = lru2q.fifo.PushFront(&cm.Entry{K: k, V: v})
}

// Get looks up a key's value from the cache.
func (lru2q *LRU2Q) Get(k cm.Key) (v cm.Value, ok bool) {

	if lru2q.cache != nil {
		if ee, hit := lru2q.cache[k]; hit {
			lru2q.ll.MoveToFront(ee)
			return ee.Value.(*cm.Entry).V, true
		}
	}

	if lru2q.qcount != nil {
		if ee, hit := lru2q.qcount[k]; hit {
			// delete the element in FIFO
			lru2q.fifo.Remove(ee)
			delete(lru2q.qcount, k)

			// make LRU
			if lru2q.cache == nil {
				lru2q.cache = make(map[cm.Key]*list.Element)
				lru2q.ll = list.New()
			}

			// add the element into LRU
			if lru2q.ll.Len() == lru2q.MaxEntries {
				b := lru2q.ll.Back()
				k := b.Value.(*cm.Entry).K
				lru2q.ll.Remove(b)
				delete(lru2q.cache, k)
			}
			kv := ee.Value.(*cm.Entry)
			lru2q.cache[k] = lru2q.ll.PushFront(kv)

			return kv.V, true
		}
	}

	return nil, false
}

// Remove removes the provided key from the cache.
func (lru2q *LRU2Q) Remove(k cm.Key) {
	if lru2q.cache != nil {
		if ee, hit := lru2q.cache[k]; hit {
			lru2q.ll.Remove(ee)
			delete(lru2q.cache, k)
		}
	}

	if lru2q.qcount != nil {
		if ee, hit := lru2q.qcount[k]; hit {
			lru2q.fifo.Remove(ee)
			delete(lru2q.qcount, k)
		}
	}
}

// Len returns the number of items in the cache.
func (lru2q *LRU2Q) Len() int {
	var n int = 0

	if lru2q.cache != nil {
		n += lru2q.ll.Len()
	}

	if lru2q.qcount != nil {
		n += lru2q.fifo.Len()
	}

	return n
}

// Remove removes the provided key from the cache.
func (lru2q *LRU2Q) Clear() {
	if lru2q.OnEvicted != nil {
		for _, e := range lru2q.cache {
			kv := e.Value.(*cm.Entry)
			lru2q.OnEvicted(kv.K, kv.V)
		}

		for _, e := range lru2q.qcount {
			kv := e.Value.(*cm.Entry)
			lru2q.OnEvicted(kv.K, kv.V)
		}
	}

	lru2q.ll = nil
	lru2q.qcount = nil
	lru2q.fifo = nil
	lru2q.cache = nil
}
