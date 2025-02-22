package syncpool

import (
	"testing"
)

const logTpl = "got: %d, want: %d"

func TestPool(t *testing.T) {
	type Foo struct{ x int }
	type Bar struct{ y int }

	fooPool01 := GetPool[Foo]() // when call GetPoo[T], it's call new (T) every time
	fooPool02 := GetPool[Foo]() // when call GetPoo[T], it's call new (T) every time
	fooPool03 := GetPool[Foo]() // when call GetPoo[T], it's call new (T) every time

	barPool01 := GetPool[Bar]() // when call GetPoo[T], it's call new (T) every time
	barPool02 := GetPool[Bar]() // when call GetPoo[T], it's call new (T) every time
	barPool03 := GetPool[Bar]() // when call GetPoo[T], it's call new (T) every time

	foo01 := fooPool01.Get()
	foo02 := fooPool02.Get()
	foo03 := fooPool03.Get()

	bar01 := barPool01.Get()
	bar02 := barPool02.Get()
	bar03 := barPool03.Get()

	// new instant all value should be 0
	if got, want := foo01.x, 0; got != want {
		t.Errorf(logTpl, got, want)
	}
	if got, want := foo02.x, 0; got != want {
		t.Errorf(logTpl, got, want)
	}
	if got, want := foo03.x, 0; got != want {
		t.Errorf(logTpl, got, want)
	}

	if got, want := bar01.y, 0; got != want {
		t.Errorf(logTpl, got, want)
	}
	if got, want := bar02.y, 0; got != want {
		t.Errorf(logTpl, got, want)
	}
	if got, want := bar03.y, 0; got != want {
		t.Errorf(logTpl, got, want)
	}

	foo01.x = 1
	foo02.x = 2
	foo03.x = 3
	bar01.y = 11
	bar02.y = 12
	bar03.y = 13
	valX01 := foo01.x
	valX02 := foo02.x
	valX03 := foo03.x
	valY01 := bar01.y
	valY02 := bar02.y
	valY03 := bar03.y

	// this prove the sync.Pool must be the same by push back
	// and pull new is should be the same one value
	fooPool03.Put(foo01)
	newX01 := fooPool02.Get()
	if got, want := newX01.x, valX01; got != want {
		t.Errorf(logTpl, got, want)
	}

	fooPool02.Put(foo02)
	newX02 := fooPool03.Get()
	if got, want := newX02.x, valX02; got != want {
		t.Errorf(logTpl, got, want)
	}

	fooPool01.Put(foo03)
	newX03 := fooPool03.Get()
	if got, want := newX03.x, valX03; got != want {
		t.Errorf(logTpl, got, want)
	}

	fooPool04 := GetPool[Foo]() // when call GetPoo[T], it's call new (T) every time
	fooPool04.Get()             // pop new(T)
	fooPool02.Put(newX03)
	newXx3 := fooPool04.Get()
	if got, want := newXx3.x, valX03; got != want {
		t.Errorf(logTpl, got, want)
	}

	// ------------------------------------------------
	barPool03.Put(bar01)
	newY01 := barPool02.Get()
	if got, want := newY01.y, valY01; got != want {
		t.Errorf(logTpl, got, want)
	}

	barPool02.Put(bar02)
	newY02 := barPool03.Get()
	if got, want := newY02.y, valY02; got != want {
		t.Errorf(logTpl, got, want)
	}

	barPool01.Put(bar03)
	newY03 := barPool03.Get()
	if got, want := newY03.y, valY03; got != want {
		t.Errorf(logTpl, got, want)
	}

	barPool04 := GetPool[Bar]() // when call GetPoo[T], it's call new (T) every time
	barPool04.Get()             // pop new(T)
	barPool02.Put(newY03)
	newYy3 := barPool04.Get()
	if got, want := newYy3.y, valY03; got != want {
		t.Errorf(logTpl, got, want)
	}

	// Should be the new
	// value should be new
	fooX1 := fooPool01.Get()
	if got, want := fooX1.x, 0; got != want {
		t.Errorf(logTpl, got, want)
	}
}

// // Function to drain the pool
// func drainPool(pool *sync.Pool) {
// 	for {
// 		if pool.Get() == nil {
// 			break
// 		}
// 	}
// }

// func TestSyncPoolReturnsNil(t *testing.T) {
// 	type Dude struct{ x int }

// 	dudePool01 := GetPool[Dude]()
// 	pool := dudePool01.pool

// 	// Ensure the pool is drained
// 	drainPool(pool)

// 	// Internal pull should be empty and return nil
// 	// Then we must create new instead
// 	item := dudePool01.Get()
// 	if item == nil {
// 		t.Fatalf("Expected nil, but got %v", item)
// 	}
// }
