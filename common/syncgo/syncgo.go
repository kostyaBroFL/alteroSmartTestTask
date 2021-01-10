package syncgo

import (
	"sync"
)

func GoWG(wg *sync.WaitGroup, f func()) {
	wg.Add(1)
	go func() {
		f()
		wg.Done()
	}()
}
