package main

import (
	"fmt"
	"sync"
)

func main() {
	sum := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			mu.Lock() // 加锁喵
			sum = sum + 1
			mu.Unlock() // 解锁喵
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(sum)
}
