//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	. "concurrency/producerconsumer/utils"
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, dataCh chan<- *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(dataCh)
			return
		}
		dataCh <- tweet
	}
}

func consumer(dataCh <-chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range dataCh {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	dataCh := make(chan *Tweet)
	var wg sync.WaitGroup

	// Producer
	wg.Add(1)
	go producer(stream, dataCh, &wg)

	// Consumer
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go consumer(dataCh, &wg)
	}
	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))
}
