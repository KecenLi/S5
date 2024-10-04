package main

import (
	"fmt"
	"github.com/ChrisGora/semaphore"
	"math/rand"
	"sync"
	"time"
)

type buffer struct {
	b                 []int
	size, read, write int
}

func newBuffer(size int) buffer {
	return buffer{
		b:     make([]int, size),
		size:  size,
		read:  0,
		write: 0,
	}
}

func (buffer *buffer) get() int {
	x := buffer.b[buffer.read]
	fmt.Println("Get\t", x, "\t", buffer)
	buffer.read = (buffer.read + 1) % len(buffer.b)
	return x
}

func (buffer *buffer) put(x int) {
	buffer.b[buffer.write] = x
	fmt.Println("Put\t", x, "\t", buffer)
	buffer.write = (buffer.write + 1) % len(buffer.b)
}

func producer(buffer *buffer, spaceAvailable, workAvailable semaphore.Semaphore, mutex *sync.Mutex, start, delta int) {
	x := start
	for {
		spaceAvailable.Wait() // 等待有空位喵~
		mutex.Lock()          // 加锁保护共享资源喵~
		buffer.put(x)         // 生产数据喵~
		mutex.Unlock()        // 解锁喵~
		workAvailable.Post()  // 通知消费者有新数据了喵~
		x = x + delta
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}

func consumer(buffer *buffer, spaceAvailable, workAvailable semaphore.Semaphore, mutex *sync.Mutex) {
	for {
		workAvailable.Wait()  // 等待有数据可用喵~
		mutex.Lock()          // 加锁保护共享资源喵~
		_ = buffer.get()      // 消费数据喵~
		mutex.Unlock()        // 解锁喵~
		spaceAvailable.Post() // 通知生产者有空位了喵~
		time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
	}
}

func main() {
	buffer := newBuffer(5) // 保持不变，返回的是值类型喵~
	mutex := &sync.Mutex{} // 初始化互斥锁喵~

	spaceAvailable := semaphore.Init(5, 5) // 初始化空闲空间信号量，开始时有 5 个空位喵~
	workAvailable := semaphore.Init(5, 0)  // 初始化工作信号量，开始时没有可用数据喵~

	// 传递缓冲区的地址给 producer 和 consumer 喵~
	go producer(&buffer, spaceAvailable, workAvailable, mutex, 1, 1)
	go producer(&buffer, spaceAvailable, workAvailable, mutex, 1000, -1)

	consumer(&buffer, spaceAvailable, workAvailable, mutex)
}
