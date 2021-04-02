package safego

import (
	"context"
	"fmt"
	"sync"
)

func f1() {
	fmt.Println("f1")
	panic("f1")
}

func f1WithParams(a int, b string) func() {
	return func() {
		fmt.Println(a, b)
		fmt.Println("f1WithParams")
		panic("f1WithParams")
	}
}

func Demo(ctx context.Context) {
	fmt.Println(ctx)
	g := NewPanicGroup()
	g.Go(f1)
	g.Go(f1WithParams(123, "abc"))
	err := g.Wait(ctx)
	//// 也可以链式调用
	//err := NewPanicGroup().Go(f1).Go(f2).Wait(ctx)
	if err != nil {
		fmt.Println(err)
	}
}

func f2(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	panic("f2")
	fmt.Println("f2")
}

func Demo2() {
	var wg sync.WaitGroup
	wg.Add(1)
	go f2(&wg)
	wg.Wait()
}
