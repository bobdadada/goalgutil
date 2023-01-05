package lru_test

import (
	"testing"

	"goalgutil/lru"
)

type simpleStruct struct {
	int
	string
}

type complexStruct struct {
	int
	simpleStruct
}

func TestLRUGet(t *testing.T) {
	getTests := []struct {
		name       string
		keyToAdd   interface{}
		keyToGet   interface{}
		expectedOk bool
	}{
		{"string_hit", "myKey", "myKey", true},
		{"string_miss", "myKey", "nonsense", false},
		{"simple_struct_hit", simpleStruct{1, "two"}, simpleStruct{1, "two"}, true},
		{"simple_struct_miss", simpleStruct{1, "two"}, simpleStruct{0, "noway"}, false},
		{"complex_struct_hit", complexStruct{1, simpleStruct{2, "three"}},
			complexStruct{1, simpleStruct{2, "three"}}, true},
	}
	for _, tt := range getTests {
		lru := lru.NewLRU(0)
		lru.Add(tt.keyToAdd, 1234)
		val, ok := lru.Get(tt.keyToGet)
		if ok != tt.expectedOk {
			t.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
		} else if ok && val != 1234 {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}
	}
}

func TestLRURemove(t *testing.T) {
	lru := lru.NewLRU(0)
	lru.Add("myKey", 1234)
	if val, ok := lru.Get("myKey"); !ok {
		t.Fatal("TestLRURemove returned no match")
	} else if val != 1234 {
		t.Fatalf("TestLRURemove failed.  Expected %d, got %v", 1234, val)
	}

	lru.Remove("myKey")
	if _, ok := lru.Get("myKey"); ok {
		t.Fatal("TestLRURemove returned a removed entry")
	}
}

func TestLRUKGet(t *testing.T) {
	getTests := []struct {
		name       string
		k          int
		count      int
		keyToAdd   interface{}
		keyToGet   interface{}
		expectedOk bool
	}{
		{"string_hit", 1, 1, "myKey", "myKey", true},
		{"string_hit", 2, 1, "myKey", "myKey", false},
		{"string_miss", 1, 1, "myKey", "nonsense", false},
		{"simple_struct_hit", 2, 3, simpleStruct{1, "two"}, simpleStruct{1, "two"}, true},
		{"simple_struct_miss", 2, 2, simpleStruct{1, "two"}, simpleStruct{0, "noway"}, false},
		{"complex_struct_hit", 2, 2, complexStruct{1, simpleStruct{2, "three"}},
			complexStruct{1, simpleStruct{2, "three"}}, true},
	}
	for _, tt := range getTests {
		lruk := lru.NewLRUK(0, tt.k)
		for i := 0; i < tt.count; i++ {
			lruk.Add(tt.keyToAdd, 1234)
		}
		val, ok := lruk.Get(tt.keyToGet)
		if ok != tt.expectedOk {
			t.Fatalf("%s: K = %v; count = %v; cache hit = %v; want %v", tt.name, tt.k, tt.count, ok, !ok)
		} else if ok && val != 1234 {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}
	}
}

func TestLRUKRemove(t *testing.T) {
	lruk := lru.NewLRUK(0, 1)
	lruk.Add("myKey", 1234)
	if val, ok := lruk.Get("myKey"); !ok {
		t.Fatal("TestLRUKRemove returned no match")
	} else if val != 1234 {
		t.Fatalf("TestLRUKRemove failed.  Expected %d, got %v", 1234, val)
	}

	lruk.Remove("myKey")
	if _, ok := lruk.Get("myKey"); ok {
		t.Fatal("TestLRUKRemove returned a removed entry")
	}
}
