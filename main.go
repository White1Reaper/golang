package main

import (
	m "packages/mutex"
	c "packages/channel"
	a "packages/atomic"
)

func main() {
	m.Mutex()
	c.Channel()
	a.Atomic()
}
