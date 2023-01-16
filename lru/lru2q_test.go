package lru_test

import (
	"testing"

	"goalgutil/lru"
)

func TestLRU2QGet(t *testing.T) {
	type simpleStruct struct {
		int
		string
	}

	type complexStruct struct {
		int
		simpleStruct
	}

	getTests := []struct {
		name       string
		maxEntries int
		count      int
		keyToAdd   interface{}
		keyToGet   interface{}
		expectedOk bool
	}{
		{"string_hit", 1, 1, "myKey", "myKey", true},
		{"string_miss", 1, 1, "myKey", "nonsense", false},
		{"simple_struct_hit", 2, 3, simpleStruct{1, "two"}, simpleStruct{1, "two"}, true},
		{"simple_struct_miss", 2, 2, simpleStruct{1, "two"}, simpleStruct{0, "noway"}, false},
		{"complex_struct_hit", 2, 2, complexStruct{1, simpleStruct{2, "three"}},
			complexStruct{1, simpleStruct{2, "three"}}, true},
	}
	for _, tt := range getTests {
		lru2q := lru.NewLRU2Q(tt.maxEntries)
		for i := 0; i < tt.count; i++ {
			lru2q.Add(tt.keyToAdd, 1234)
		}
		val, ok := lru2q.Get(tt.keyToGet)
		if ok != tt.expectedOk {
			t.Fatalf("%s: K = %v; count = %v; cache hit = %v; want %v", tt.name, tt.maxEntries, tt.count, ok, !ok)
		} else if ok && val != 1234 {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}
	}
}

func TestLRU2QRemove(t *testing.T) {
	lru2q := lru.NewLRU2Q(4)
	lru2q.Add("myKey", 1234)
	if val, ok := lru2q.Get("myKey"); !ok {
		t.Fatal("TestLRUKRemove returned no match")
	} else if val != 1234 {
		t.Fatalf("TestLRUKRemove failed.  Expected %d, got %v", 1234, val)
	}

	lru2q.Remove("myKey")
	if _, ok := lru2q.Get("myKey"); ok {
		t.Fatal("TestLRUKRemove returned a removed entry")
	}
}
