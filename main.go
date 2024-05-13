package main

import (
	"fmt"
	future "patterns/libs/future"
	pool "patterns/libs/pool"
	"time"
)

type Task struct {
	ID int
}


func (w *Task) execute() {
	w.ID += 100
	fmt.Println("Executing task 1...")
	delay := 5 * time.Second
	fmt.Printf("Waiting for %s before executing tasks...\n", delay)
	time.Sleep(delay)
}


type Future struct {
	ID int
}

func (w *Future) execute() int {
	w.ID += 100
	fmt.Println("Executing task 1...")
	delay := 5 * time.Second
	fmt.Printf("Waiting for %s before executing tasks...\n", delay)
	time.Sleep(delay)
	return 1
}



func main() {

	//////////////////////////////////////////////////////////////////////////////////

	// Worker-pool
	var taskItems []Task

	for i := 0; i < 10; i++ {
		taskItems = append(taskItems, Task{ID: i})
	}

	task1 := pool.NewTask(taskItems[0].execute)
	task2 := pool.NewTask(taskItems[1].execute)

	task3 := pool.NewTask(func(){
		fmt.Println("Executing task 3...")
	})

	pool.ProcessTasks(10, []pool.Task{task1, task2, task3})


	//////////////////////////////////////////////////////////////////////////////////

	// Futures
	var futures []Future
	for i := 0; i < 10; i++ {
		futures = append(futures, Future{ID: i})
	}
	future1 := future.NewFuture(futures[0].execute)
	future1.Submit()

	future2 := future.NewFuture(futures[1].execute)
	future2.Submit()

	future2.Submit()
	result := future1.Get()

	fmt.Println(result)

	//////////////////////////////////////////////////////////////////////////////////

	// Futures
}
