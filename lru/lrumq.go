package lru

import (
	"container/list"

	cm "goalgutil/macros/cache_macro"
)

// MQ算法根据访问频率将数据划分为多个队列，不同的队列具有不同的访问优先级，其核心思想是：优先缓存访问次数多的数据。
// MQ算法将缓存划分为多个LRU队列，每个队列对应不同的访问优先级，访问优先级是根据访问次数计算出来的。
// 例如：
// Q0，Q1....Qk代表不同的优先级队列。Q-history代表从缓存中淘汰数据，但记录了数据的索引和引用次数的队列。
//    1. 新插入的数据放入Q0；
//    2. 每个队列按照LRU管理数据；
//    3. 当数据的访问次数达到一定次数，需要提升优先级时，将数据从当前队列删除，加入到高一级队列的头部；
//    4. 为了防止高优先级数据永远不被淘汰，当数据在指定的时间里访问没有被访问时，需要降低优先级，将数据从当前队列删除，加入到低一级的队列头部；
//    5. 需要淘汰数据时，从最低一级队列开始按照LRU淘汰；每个队列淘汰数据时，将数据从缓存中删除，将数据索引加入Q-history头部；
//    6. 如果数据在Q-history中被重新访问，则重新计算其优先级，移到目标队列的头部；
//    7. Q-history按照LRU淘汰数据的索引。
//
// 命中率:
// MQ降低了“缓存污染”带来的问题，命中率比LRU要高。
//
// 复杂度:
// MQ需要维护多个队列，且需要维护每个数据的访问时间，复杂度比LRU高。
//
// 代价:
// MQ需要记录每个数据的访问时间，需要定时扫描所有队列，代价比LRU要高。
//
// 注：虽然MQ的队列看起来数量比较多，但由于所有队列之和受限于缓存容量的大小，因此这里多个队列长度之和和一个LRU队列是一样的，因此队列扫描性能也相近。
//
// 作者：jiangmo
// 链接：https://www.jianshu.com/p/d533d8a66795
// 来源：简书
// 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

type LRUMQ struct {
	MaxEntries int
	NumQueues  int

	// OnEvicted optionally specifies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(k cm.Key, v cm.Value)

	ll     *list.List
	fifo   *list.List
	cache  map[cm.Key]*list.Element
	qcount map[cm.Key]*list.Element

	qhistory map[cm.Key]*list.Element
}
