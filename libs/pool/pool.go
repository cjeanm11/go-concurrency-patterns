package workerpool

import (
	"sync"
)

type Item interface{ }

type TaskCommand interface {
    invoke()
}

type Task struct {
	execute func()
}

func NewTask(execute func()) Task {
    return Task{execute: execute}
}


func (t Task) invoke() {
	t.execute()
}


func ProcessTasks[T TaskCommand](numWorkers int, items []T) {
	var wg sync.WaitGroup
	taskChannel := make(chan T, numWorkers)
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go workerTask(&wg, taskChannel)
	}

	for _, item := range items {
		taskChannel <- item
	}
	close(taskChannel)
	wg.Wait()
}

func ProcessItems[T Item](numWorkers int, items []T, task func(T)) {
	var wg sync.WaitGroup
	taskChannel := make(chan T, numWorkers)
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go workerItem(&wg, taskChannel, task)
	}

	for _, item := range items {
		taskChannel <- item
	}
	close(taskChannel)
	wg.Wait()
}

func workerItem[T Item](wg *sync.WaitGroup, items <-chan T, task func(T)) {
	defer wg.Done()
	for item := range items {
		task(item)
	}
}

func workerTask[T TaskCommand](wg *sync.WaitGroup, items <-chan T) {
	defer wg.Done()
	for item := range items {
		item.invoke()
	}
}
