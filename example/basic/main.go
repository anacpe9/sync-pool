package main

import (
	syncpool "github.com/anacpe9/sync-pool"
)

func main() {
	type Foo struct{ x int }
	type Bar struct{ y int }

	fooPool_01 := syncpool.GetPool[Foo]()
	fooPool_02 := syncpool.GetPool[Foo]()
	barPool_01 := syncpool.GetPool[Bar]()

	foo_01 := fooPool_01.Get()
	foo_02 := fooPool_02.Get()
	bar_01 := barPool_01.Get()

	fooPool_01.Put(foo_01)
	fooPool_01.Put(foo_02)

	barPool_01.Put(bar_01)
}
