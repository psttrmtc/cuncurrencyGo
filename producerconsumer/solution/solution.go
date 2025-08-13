//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package producerconsumer

import (
	. "concurrency/producerconsumer/utils"
	"fmt"
	"sync"
	"time"
)

func producer(wg *sync.WaitGroup, stream Stream, dataCh chan<- *Tweet, done chan bool) {

}

func consumer(wg *sync.WaitGroup, dataCha <-chan *Tweet) {
	defer wg.Done()
	for t := range dataCha {
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
	doneCh := make(chan bool)
	producers := sync.WaitGroup{}
	consumers := sync.WaitGroup{}

	// Producer
	go producer(&producers, stream, dataCh, doneCh)

	// Consumer
	for i := 0; i < 10; i++ {
		consumers.Add(1)
		go consumer(&consumers, dataCh)
	}
	<-doneCh
	fmt.Printf("Process took %s\n", time.Since(start))
}
