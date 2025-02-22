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

func initializeNestedPointer(ptr interface{}) {
	v := reflect.ValueOf(ptr)

	// // Use it internally and send pointers.
	// if v.Kind() != reflect.Ptr {
	// 	panic("expected a pointer")
	// }

	for v.Kind() == reflect.Ptr {
		elemType := v.Type().Elem()
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		} else if elemType.Kind() == reflect.Slice {
			v.Elem().Set(reflect.MakeSlice(elemType, 0, 0))
		}
		v = v.Elem()
	}
}

func GetPool[T any]() *Pool[T] {
	// this function should not call frequency.
	// We allow to allocated one of T to reflect type
	// and put back to pool

	newType := func() *T {
		tt := new(T)
		initializeNestedPointer(tt)
		return tt
	}

	t := newType()
	key := reflect.TypeOf(t).String()
	pool := genericTypeMap[key]

	if pool == nil {
		newPool := &Pool[T]{
			pool: &sync.Pool{
				New: func() interface{} { return newType() },
			},
		}

		genericTypeMap[key] = newPool
		pool = genericTypeMap[key]
	}

	p := pool.(*Pool[T])
	p.Put(t)

	return p
}

// Get returns an arbitrary item from the pool.
func (p *Pool[T]) Get() *T {
	itm := p.pool.Get().(*T)
	// // Try to use drainPool
	// // But we not retrieve nil
	// // So this statement we ignore
	//
	// if itm == nil {
	// 	itm = new(T)
	// }

	return itm
}

// Put places an item in the pool
func (p *Pool[T]) Put(x *T) {
	p.pool.Put(x)
}
