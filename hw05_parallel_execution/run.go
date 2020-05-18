package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	if len(tasks) < N {
		N = len(tasks)
	}
	tasksChan := make(chan Task, 1)
	results := make(chan error)
	done := make(chan struct{})

	workersWg := &sync.WaitGroup{}

	workersWg.Add(N)
	for i := 0; i < N; i++ {
		go worker(tasksChan, results, done, workersWg)
	}

	errCount := 0
	checkRes := func(res error) error {
		if res != nil {
			errCount++
		}
		if errCount >= M {
			close(done)
			workersWg.Wait()
			return ErrErrorsLimitExceeded
		}
		return nil
	}

	for i := 0; i < len(tasks); {
		select {
		case tasksChan <- tasks[i]:
			i++
		case res := <-results:
			res = checkRes(res)
			if res != nil {
				return res
			}
		}
	}
	res := checkRes(<-results)
	if res != nil {
		return res
	}

	close(tasksChan)
	close(done)
	workersWg.Wait()
	return nil
}

func worker(in <-chan Task, out chan<- error, done chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range in {
		err := task()

		select {
		case <-done:
			return
		case out <- err:
		}
	}
}
