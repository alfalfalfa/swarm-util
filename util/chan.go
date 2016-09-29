package util

import "sync"

func MargeChan(list ...<-chan []byte) <-chan []byte {
	out := make(chan []byte, 0)
	go func() {
		defer func() {
			close(out)
		}()
		wg := sync.WaitGroup{}
		f := func(c <-chan []byte) {
			defer wg.Done()
			for {
				b := <-c
				if b == nil {
					return
				}
				out <- b
			}
		}
		wg.Add(len(list))
		for _, c := range list {
			go f(c)
		}
		wg.Wait()
	}()
	return out
}
