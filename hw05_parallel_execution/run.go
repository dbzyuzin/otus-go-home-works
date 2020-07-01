package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	if len(tasks) < n {
		n = len(tasks)
	}
	tasksChan := make(chan Task, 1)
	results := make(chan error)
	done := make(chan struct{})

	wg := StartWorkers(n, tasksChan, results, done)
	defer wg.Wait()

	errCount := 0
	checkRes := func(res error) error {
		if res != nil {
			errCount++
		}
		if errCount >= m {
			close(done)
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

	close(tasksChan)
	for res := range results {
		res = checkRes(res)
		if res != nil {
			return res
		}
	}

	close(done)
	return nil
}

func StartWorkers(n int, tasksChan chan Task, results chan error, done chan struct{}) *sync.WaitGroup {
	workersWg := &sync.WaitGroup{}

	workersWg.Add(n)
	for i := 0; i < n; i++ {
		go worker(tasksChan, results, done, workersWg)
	}

	go func() {
		workersWg.Wait()
		close(results)
	}()

	return workersWg
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
