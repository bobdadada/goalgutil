package cache_macro

type Key any
type Value any

type Entry struct {
	K Key
	V Value
}

type Cache interface {
	Add(k Key, v Value)
	Get(k Key) (v Value, ok bool)
	Remove(k Key)
	Len() int
	Clear()
}
