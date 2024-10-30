package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func mutex() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("hello world")
	}()
	wg.Wait()
}

func channel() {
	ch := make(chan string)
	go func() {
		ch <- "hello world"
		close(ch)
	}()
	fmt.Println(<-ch)
}

func Atomic() {
	var done int32 = 0
	go func() {
		fmt.Println("hello world")
		atomic.StoreInt32(&done, 1)
	}()
	for atomic.LoadInt32(&done) == 0 {
		time.Sleep(3 * time.Millisecond)
	}
}

func BenchmarkMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mutex()
	}
}

func BenchmarkChannel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		channel()
	}
}

func BenchmarkAtomic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Atomic()
	}
}

func main() {
mutex()
//channel()
//Atomic()
}
