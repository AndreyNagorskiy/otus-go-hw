package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded          = errors.New("errors limit exceeded")
	ErrWorkersMustBeGreaterThanZero = errors.New("workers must be greater than zero")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	switch {
	case n <= 0:
		return ErrWorkersMustBeGreaterThanZero
	case m <= 0:
		return ErrErrorsLimitExceeded
	case len(tasks) == 0:
		return nil
	}

	doneCh := make(chan struct{})
	tCh := make(chan Task, n)
	errCh := make(chan struct{}, m)

	var wg sync.WaitGroup
	var closeOnce sync.Once

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(doneCh, tCh, errCh, &wg, m, &closeOnce)
	}

	go func() {
		defer close(tCh)
		for _, task := range tasks {
			select {
			case <-doneCh:
				return
			case tCh <- task:
			}
		}
	}()

	wg.Wait()

	if len(errCh) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(
	doneCh chan struct{},
	tCh <-chan Task,
	errCh chan<- struct{},
	wg *sync.WaitGroup,
	m int,
	closeOnce *sync.Once,
) {
	defer wg.Done()

	for {
		select {
		case <-doneCh:
			return
		case task, ok := <-tCh:
			if !ok {
				return
			}

			if err := task(); err != nil {
				select {
				case errCh <- struct{}{}:
					if len(errCh) >= m {
						closeOnce.Do(func() { close(doneCh) })
					}
				case <-doneCh:
					return
				}
			}
		}
	}
}
