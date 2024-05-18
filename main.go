package main

import (
	"fmt"
	fanin "patterns/libs/fanin"
	future "patterns/libs/future"
	p "patterns/libs/pipeline"
	pool "patterns/libs/workerpool"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////

// Worker Pool

type Task struct {
	ID int
}

func (w *Task) execute() any {
	w.ID += 100
	fmt.Println("Executing task 1...")
	delay := 5 * time.Second
	fmt.Printf("Waiting for %s before executing tasks...\n", delay)
	time.Sleep(delay)
	return 1
}

//////////////////////////////////////////////////////////////////////////////////

// Future

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

//////////////////////////////////////////////////////////////////////////////////

func execute1(data any) any {
	return data.(int) * 2
}

func execute2(data any) any {
	return data.(int) + 1
}

func execute3(data any) any {
	return data.(int) * data.(int)
}

//////////////////////////////////////////////////////////////////////////////////func sendData(ch chan<- int, start, end int) {

func sendData(ch chan<- int, start, end int) {
	defer close(ch)
	for i := start; i < end; i++ {
		ch <- i
		time.Sleep(time.Millisecond * 500)
	}
}

func main() {

	//////////////////////////////////////////////////////////////////////////////////

	// Worker Pool

	var taskItems []Task
	for i := 0; i < 10; i++ {
		taskItems = append(taskItems, Task{ID: i})
	}

	var itemFuncs []func() any
	for _, item := range taskItems {
		itemFuncs = append(itemFuncs, item.execute)
	}

	result := pool.ProcessTasks(10, itemFuncs)

	fmt.Print("Worker pool :")
	fmt.Println(result...)
	//////////////////////////////////////////////////////////////////////////////////

	// Completable Futures

	var futures []Future
	for i := 0; i < 10; i++ {
		futures = append(futures, Future{ID: i})
	}
	future1 := future.NewFuture(futures[0].execute)
	future1.Submit()

	future2 := future.NewFuture(futures[1].execute)
	future2.Submit()

	future3 := future.NewFuture(futures[2].execute)
	future3.Submit()

	result2 := future2.Get()
	result3 := future3.Get()
	result1 := future1.Get()

	fmt.Println("Future result: ", result1+result2+result3)
	//////////////////////////////////////////////////////////////////////////////////

	// Pipeline

	builder := p.NewPipelineBuilder()
	builder.AddStage(execute1)
	builder.AddStage(execute2)
	builder.AddStage(execute3)
	pipeline := builder.Build()

	in := make(chan interface{})
	go func() {
		defer close(in)
		for _, n := range []int{1, 4} {
			in <- n
		}
	}()

	for result := range p.ExecutePipeline(in, pipeline) {
		fmt.Println("Pipeline result: ", result)
	}
	//////////////////////////////////////////////////////////////////////////////////

	// Fan-in

	ch1 := make(chan int)
	ch2 := make(chan int)
	
	f := fanin.NewFanIn()
	
	go sendData(ch1, 0, 5)
	go sendData(ch2, 5, 10)
	
	f.AddInputChannel(ch1)
	f.AddInputChannel(ch2)
	
	f.Start()
	
	merged := f.MergedChannel()
	
	for val := range merged {
		fmt.Println("Merged:", val)
	}
	
}
