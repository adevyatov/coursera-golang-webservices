package main

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// ExecutePipeline сюда писать код
func ExecutePipeline(jobs ...job) {

	var in chan interface{}
	var channels []chan interface{}

	for i := 0; i < len(jobs); i++ {
		channels = append(channels, make(chan interface{}))

		if i == 0 {
			in = make(chan interface{}, 5)
		} else {
			in = channels[i-1]
		}

		out := channels[i]
		go jobs[i](in, out)

		runtime.Gosched()
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go closeChannels(wg, 340, channels[len(channels)-1:])
	wg.Wait()
}

func closeChannels(wg *sync.WaitGroup, timeout int, channels []chan interface{}) {
	time.Sleep(time.Millisecond * time.Duration(timeout))
	for _, channel := range channels {
		close(channel)
	}
	wg.Done()
}

func SingleHash(in, out chan interface{}) {
	value := fmt.Sprintf("%v", <-in)

	out <- fmt.Sprintf("%v~%v", DataSignerCrc32(value), DataSignerCrc32(DataSignerMd5(value)))
}

func MultiHash(in chan interface{}, out chan interface{}) {
	var result string
	value := fmt.Sprintf("%v", <-in)
	fmt.Printf("Got %v\n", value)

	for i := 0; i < 6; i++ {
		result += fmt.Sprintf("%d%v", i, DataSignerCrc32(value))
	}

	out <- result
}

func CombineResults(in chan interface{}, out chan interface{}) {
	var values = make([]string, 10)

	for value := range in {
		values = append(values, fmt.Sprintf("%v", value))
	}

	sort.Strings(values)

	out <- strings.Join(values[:], "_")
}
