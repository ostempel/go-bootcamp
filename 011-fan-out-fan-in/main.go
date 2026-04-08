package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Result struct {
	Url      string
	Status   string
	Duration time.Duration
	Err      error
}

var client = &http.Client{Timeout: time.Second * 5}

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		log.Fatalln("needs at least 2 args: worker-count and at least one url")
	}

	count, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("unable to parse worker-count: %v", err)
	}

	urls := generator(args[1:])
	results := fetch(urls, count)
	merged := merge(results...)

	for value := range merged {
		if value.Err != nil {
			fmt.Printf("%s - %s (%s)\n", value.Url, value.Err.Error(), value.Duration.String())
		} else {
			fmt.Printf("%s - %s (%s)\n", value.Url, value.Status, value.Duration.String())
		}
	}
}

func generator(input []string) <-chan string {
	urls := make(chan string)

	go func() {
		for _, in := range input {
			urls <- in
		}
		close(urls)
	}()

	return urls
}

func fetch(urls <-chan string, workerCount int) []<-chan Result {
	fetcherChannels := make([]<-chan Result, 0, workerCount)

	// fan-out: go routines read from same channel -> distribute work
	for w := 1; w <= workerCount; w++ {
		fetcherChannels = append(fetcherChannels, handleFetch(urls))
	}

	return fetcherChannels
}

func handleFetch(urls <-chan string) <-chan Result {
	results := make(chan Result)

	go func() {
		for url := range urls {
			start := time.Now()
			resp, err := client.Get(url)
			duration := time.Since(start)
			if err != nil {
				results <- Result{Url: url, Err: err, Duration: duration.Round(time.Millisecond)}
				continue
			}
			results <- Result{Url: url, Duration: duration.Round(time.Millisecond), Status: resp.Status}
			// close body here -> defer only at the end of function
			resp.Body.Close()
		}
		close(results)
	}()

	return results
}

func merge(cs ...<-chan Result) <-chan Result {
	merged := make(chan Result)
	var wg sync.WaitGroup

	output := func(results <-chan Result) {
		for r := range results {
			merged <- r
		}
		wg.Done()
	}
	wg.Add(len(cs))

	// fan-in: merge all channels into one again
	for _, c := range cs {
		go output(c)
	}

	// count active channels and wait till all are merged
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}
