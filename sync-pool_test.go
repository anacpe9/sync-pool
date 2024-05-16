package syncpool

import (
	"testing"
)

func TestPool(t *testing.T) {
	type Foo struct{ x int }
	type Bar struct{ y int }

	fooPool_01 := GetPool[Foo]() // when call GetPoo[T], it's call new (T) every time
	fooPool_02 := GetPool[Foo]() // when call GetPoo[T], it's call new (T) every time
	fooPool_03 := GetPool[Foo]() // when call GetPoo[T], it's call new (T) every time

	barPool_01 := GetPool[Bar]() // when call GetPoo[T], it's call new (T) every time
	barPool_02 := GetPool[Bar]() // when call GetPoo[T], it's call new (T) every time
	barPool_03 := GetPool[Bar]() // when call GetPoo[T], it's call new (T) every time

	foo_01 := fooPool_01.Get()
	foo_02 := fooPool_02.Get()
	foo_03 := fooPool_03.Get()

	bar_01 := barPool_01.Get()
	bar_02 := barPool_02.Get()
	bar_03 := barPool_03.Get()

	// new instant all value should be 0
	if got, want := foo_01.x, 0; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
	if got, want := foo_02.x, 0; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
	if got, want := foo_03.x, 0; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	if got, want := bar_01.y, 0; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
	if got, want := bar_02.y, 0; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
	if got, want := bar_03.y, 0; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	foo_01.x = 1
	foo_02.x = 2
	foo_03.x = 3
	bar_01.y = 11
	bar_02.y = 12
	bar_03.y = 13
	val_x_01 := foo_01.x
	val_x_02 := foo_02.x
	val_x_03 := foo_03.x
	val_y_01 := bar_01.y
	val_y_02 := bar_02.y
	val_y_03 := bar_03.y

	// this prove the sync.Pool must be the same by push back
	// and pull new is should be the same one value
	fooPool_03.Put(foo_01)
	new_x_01 := fooPool_02.Get()
	if got, want := new_x_01.x, val_x_01; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	fooPool_02.Put(foo_02)
	new_x_02 := fooPool_03.Get()
	if got, want := new_x_02.x, val_x_02; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	fooPool_01.Put(foo_03)
	new_x_03 := fooPool_03.Get()
	if got, want := new_x_03.x, val_x_03; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	fooPool_04 := GetPool[Foo]() // when call GetPoo[T], it's call new (T) every time
	fooPool_04.Get()             // pop new(T)
	fooPool_02.Put(new_x_03)
	new_x_x3 := fooPool_04.Get()
	if got, want := new_x_x3.x, val_x_03; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	// ------------------------------------------------
	barPool_03.Put(bar_01)
	new_y_01 := barPool_02.Get()
	if got, want := new_y_01.y, val_y_01; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	barPool_02.Put(bar_02)
	new_y_02 := barPool_03.Get()
	if got, want := new_y_02.y, val_y_02; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	barPool_01.Put(bar_03)
	new_y_03 := barPool_03.Get()
	if got, want := new_y_03.y, val_y_03; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	barPool_04 := GetPool[Bar]() // when call GetPoo[T], it's call new (T) every time
	barPool_04.Get()             // pop new(T)
	barPool_02.Put(new_y_03)
	new_y_y3 := barPool_04.Get()
	if got, want := new_y_y3.y, val_y_03; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}

	// Should be the new
	// value should be new
	foo_x1 := fooPool_01.Get()
	if got, want := foo_x1.x, 0; got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
}
