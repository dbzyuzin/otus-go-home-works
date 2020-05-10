package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	results := make(chan error, 1)
	tasksChan := make(chan Task, N)
	quit := make(chan struct{})

	var tnum int
	for tnum = 0; tnum < N && tnum < len(tasks); tnum++ {
		tasksChan <- tasks[tnum]
	}
	//go pushTasks(tasks, tasksChan, quit)

	for i := 0; i < N; i++ {
		go worker(tasksChan, results, quit)
	}

	errCounter := 0
	for {
		if <-results != nil {
			errCounter++
		}
		if errCounter >= M {
			close(tasksChan)
			return ErrErrorsLimitExceeded
		} else {
			if tnum >= len(tasks) {
				break
			}
			tasksChan <- tasks[tnum]
			tnum++
		}
	}

	return nil
}
func worker(in <-chan Task, out chan<- error, quit chan struct{}) {
	for {
		var task Task
		select {
		case <-quit:
			return
		case task = <-in:
		}

		if task == nil {
			return
		}
		err := task()

		select {
		case <-quit:
			return
		case out <- err:
		}
	}
}

func pushTasks(tasks []Task, tasksChan chan<- Task, quit chan struct{}) {
	defer close(tasksChan)
	for _, task := range tasks {
		select {
		case <-quit:
			return
		case tasksChan <- task:
		}
	}
}
