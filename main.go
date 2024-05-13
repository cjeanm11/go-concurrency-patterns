package main

import (
	"fmt"
	wp "patterns/libs/worker-pool"
	"time"
)

type Worker struct {
	ID   int
	Desc string
	Sum  *int
}

func taskForItem(task Worker) {
	fmt.Printf("Processing task %d: %s\n", task.ID, task.Desc)
}

func task1Execution(item *Worker) {
	item.ID += 100
	fmt.Println("Executing task 1...")
	delay := 5 * time.Second
	fmt.Printf("Waiting for %s before executing tasks...\n", delay)
	time.Sleep(delay)
}

func task2Execution(item *Worker) {
	item.ID += 100
	fmt.Println("Executing task 2...")
	delay := 2 * time.Second
	fmt.Printf("Waiting for %s before executing tasks...\n", delay)
	time.Sleep(delay)
}

func main() {
	// Worker-pool
	var taskItems []Worker

	for i := 0; i < 10; i++ {
		taskItems = append(taskItems, Worker{ID: i, Desc: fmt.Sprintf("Task %d", i)})
	}

	task1 := wp.NewTask(func() {
		task1Execution(&taskItems[0])
	})

	task2 := wp.NewTask(func() {
		task2Execution(&taskItems[1])
	})

	wp.ProcessTasks(10, []wp.Task{task1, task2})
	wp.ProcessItems(2, taskItems, taskForItem)
}
