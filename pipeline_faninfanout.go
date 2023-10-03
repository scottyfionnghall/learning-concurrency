package main

import (
	"fmt"
	"sync"
	"time"
)

func gen(nums ...int) <-chan int {
	out := make(chan int, len(nums))

	go func() {
		for _, n := range nums {
			out <- n
			time.Sleep(1 * time.Second)
		}
		close(out)
	}()
	return out
}

func sq(done chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()

	return out
}

func merge(done chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int, 1)

	output := func(c <-chan int) {
		defer wg.Done()

		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	done := make(chan struct{}, 2)

	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	c1 := sq(done, gen(arr...))
	c2 := sq(done, gen(arr...))

	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	for n := range merge(done, c1, c2) {
		fmt.Println(n)
	}
}
