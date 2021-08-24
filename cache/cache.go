package cache


import 	"sync"


// A thread-safe map
type Cache struct {
	dict map[int]bool
	mux  sync.Mutex
}

// Constructor
func NewCache() Cache {
	return Cache{dict: make(map[int]bool)}
}


// check if id is in map
func (c *Cache) CheckVisited(id int) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.dict[id]
	if ok == false {
		c.dict[id] = true
		return false
	}
	return true
}