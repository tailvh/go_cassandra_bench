package component

import "sync"

// Slice type that can be safely shared between goroutines
type ConcurrentSlice struct {
  sync.RWMutex
  items []interface{}
}

// Concurrent slice item
type ConcurrentSliceItem struct {
  Index int
  Value interface{}
}

// Appends an item to the concurrent slice
func (cs *ConcurrentSlice) Append(item interface{}) {
  cs.Lock()
  defer cs.Unlock()

  cs.items = append(cs.items, item)
}

func (cs *ConcurrentSlice) Len() int {
  cs.Lock()
  defer cs.Unlock()

  value := len(cs.items)
  return value
}

// Iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builin range keyword
func (cs *ConcurrentSlice) Iter() <-chan ConcurrentSliceItem {
  c := make(chan ConcurrentSliceItem)

  f := func() {
    cs.Lock()
    defer cs.Lock()
    for index, value := range cs.items {
      c <- ConcurrentSliceItem{index, value}
    }
    close(c)
  }
  go f()

  return c
}

type ConcurrentMap struct {
  sync.RWMutex
  items map[string]interface{}
}

// Concurrent map item
type ConcurrentMapItem struct {
  Key   string
  Value interface{}
}

// Sets a key in a concurrent map
func (cm *ConcurrentMap) Set(key string, value interface{}) {
  cm.Lock()
  defer cm.Unlock()

  cm.items[key] = value
}

// Gets a key from a concurrent map
func (cm *ConcurrentMap) Get(key string) (interface{}, bool) {
  cm.Lock()
  defer cm.Unlock()

  value, ok := cm.items[key]

  return value, ok
}

// Iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (cm *ConcurrentMap) Iter() <-chan ConcurrentMapItem {
  c := make(chan ConcurrentMapItem)

  f := func() {
    cm.Lock()
    defer cm.Unlock()

    for k, v := range cm.items {
      c <- ConcurrentMapItem{k, v}
    }
    close(c)
  }
  go f()

  return c
}
