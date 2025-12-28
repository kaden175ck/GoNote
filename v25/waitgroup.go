package v25

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	sum int64
)

const P = 10

func Sum() {
	sum = 0
	var wg sync.WaitGroup
	for i := 0; i < P; i++ {
		wg.Go(func() {
			atomic.AddInt64(&sum, 1)
		})
	}
	wg.Wait()
	fmt.Println(sum)
}
