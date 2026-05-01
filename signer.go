package main

import (
	"fmt"
	"sync"
)

func ExecutePipeline(jbs ...job) {
	var wg sync.WaitGroup
	wg.Add(len(jbs))

	firstch := make(chan interface{})
	in := firstch

	for _, jb := range jbs {
		nextChan := make(chan interface{})
		go func(j job, in, out chan interface{}) {
			defer wg.Done()
			j(in, out)
			close(out)
		}(jb, in, nextChan)
		in = nextChan
	}

	// go func() {
	// 	for i := 1; i <= 5; i++ {
	// 		fmt.Printf("Отправляем: %d\n", i)
	// 		firstch <- i
	// 	}
	// 	close(firstch)
	// }()

	wg.Wait()
}

func main() {
	double := func(in, out chan interface{}) {
		for val := range in {
			result := val.(int) * 2
			fmt.Printf("Удваиваем: %d -> %d\n", val, result)
			out <- result
		}
	}

	print := func(in, out chan interface{}) {
		for val := range in {
			fmt.Println("Результат:", val)
		}
	}

	ExecutePipeline(double, print)
}
