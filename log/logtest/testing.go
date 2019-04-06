package logtest

import (
	"fmt"
	"sync"
	"testing"

	"github.com/fhofherr/golf/log"
)

// GenerateKEYVALs generates a fixed amount of key value pairs.
func GenerateKEYVALs(tb testing.TB, n int) []interface{} {
	if n < 1 {
		tb.Fatalf("expected n >= 1; got %d", n)
	}
	var kvs []interface{}
	for i := 0; i < n; i++ {
		kvs = append(kvs, fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	return kvs
}

// StressTestLogger exercises the logger returned by factory using multiple Go routines.
//
// Each Go routine writes nMSGs before it terminates.
func StressTestLogger(tb testing.TB, factory func() log.Logger, nGoRoutines, nMSGs int) {
	if nGoRoutines < 1 {
		tb.Fatalf("expected nGoRoutines >= 1; got %d", nGoRoutines)
	}
	if nMSGs < 1 {
		tb.Fatalf("expected nMSGs >= 1; got %d", nMSGs)
	}
	start := make(chan struct{})
	errc := make(chan error, nGoRoutines)
	wg := &sync.WaitGroup{}
	wg.Add(nGoRoutines)
	for i := 0; i < nGoRoutines; i++ {
		go func(i int) {
			defer wg.Done()
			// Block until the start channel is closed, then write the messages
			<-start
			logger := factory()
			for j := 0; j < nMSGs; j++ {
				err := logger.Log("go routine", i+1, "msg no", j+1)
				if err != nil {
					errc <- err
					return
				}
			}

		}(i)
	}
	close(start) // Trigger the loggers
	wg.Wait()
	for {
		select {
		case err := <-errc:
			tb.Error(err)
		default:
			return
		}
	}
}
