package fanin

import (
    "sync"
)

type FanIn struct {
    inputChannels []<-chan int
    mergedChannel chan int
    wg            sync.WaitGroup
    waitOnce      sync.Once
}

func NewFanIn() *FanIn {
    return &FanIn{
        inputChannels: make([]<-chan int, 0),
        mergedChannel: make(chan int),
    }
}

func (f *FanIn) AddInputChannel(ch <-chan int) {
    f.inputChannels = append(f.inputChannels, ch)
}

func (f *FanIn) Start() {
    f.waitOnce.Do(func() {
        f.wg.Add(len(f.inputChannels))
        go func() {
            defer close(f.mergedChannel)
            for _, ch := range f.inputChannels {
                go func(input <-chan int) {
                    defer f.wg.Done()
                    for val := range input {
                        f.mergedChannel <- val
                    }
                }(ch)
            }
            f.wg.Wait()
        }()
    })
}

func (f *FanIn) MergedChannel() <-chan int {
    f.Start()
    return f.mergedChannel
}

func (f *FanIn) Close() {
    f.waitOnce = sync.Once{}
}
