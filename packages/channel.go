import (
	"fmt"
)

func Channel() {
	ch := make(chan string)

	go func() {
		ch <- "hello world"
		close(ch)
	}()

	fmt.Println(<-ch)
}