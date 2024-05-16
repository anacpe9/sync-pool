package syncpool

// inspired by:
// https://github.com/noxiouz/golang-generics-util/blob/main/xsync/pool.go
// https://github.com/mkmik/syncpool/blob/main/pool.go
import (
	"reflect"
	"sync"
)

// global sync.Pool map
var genericTypeMap = make(map[string]interface{})

// Pool is a strongly typed version of sync.Pool from go standard library
type Pool[T any] struct {
	pool *sync.Pool
}

func GetPool[T any]() *Pool[T] {
	// this function should not call frequency.
	// We allow to allocated one of T to reflect type
	// and put back to pool
	t := new(T)
	key := reflect.TypeOf(t).String()
	pool := genericTypeMap[key]

	if pool != nil {
		p := pool.(*Pool[T])
		p.Put(t)

		return p
	}

	newPool := &Pool[T]{
		pool: &sync.Pool{
			New: func() interface{} { return new(T) },
		},
	}

	newPool.Put(t)
	genericTypeMap[key] = newPool
	return newPool
}

// Get returns an arbitrary item from the pool.
func (p *Pool[T]) Get() *T {
	return p.pool.Get().(*T)
}

// Put places an item in the pool
func (p *Pool[T]) Put(x *T) {
	p.pool.Put(x)
}
