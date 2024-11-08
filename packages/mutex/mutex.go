package mutex

import (
	"fmt"
	"sync"
)

func Mutex() {
	var m sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.Lock()
		fmt.Println("hello world")
		m.Unlock()
	}()

	wg.Wait()
}
