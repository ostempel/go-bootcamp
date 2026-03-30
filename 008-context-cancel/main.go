package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Result struct {
	Url      string
	Status   string
	Duration time.Duration
	Err      error
}

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		log.Fatal("commands needs at least 2 args: [delay, urls...]")
	}

	delay, err := time.ParseDuration(args[0])
	if err != nil {
		log.Fatalf("delay %q is not a valid duration: %s", args[0], err)
	}
	urls := args[1:]

	// create context with timeout and propagate to all requests
	ctx, cancel := context.WithTimeout(context.Background(), delay)
	defer cancel()

	// channel for receiving results
	messages := make(chan Result, len(urls))
	fmt.Printf("initialized channel - size %d\n", len(urls))

	// go routine spawning for concurrent reqs
	for _, url := range urls {
		go fetchAPI(ctx, messages, url)
	}

	succeededCount := 0
	totalCount := 0

	// loop over length of channel an receive messages
	for range len(urls) {
		msg := <-messages
		totalCount++
		if msg.Err != nil {
			fmt.Printf("%s - %s (%s)\n", msg.Url, msg.Err, msg.Duration)
		} else {
			succeededCount++
			fmt.Printf("%s - %s (%s)\n", msg.Url, msg.Status, msg.Duration)
		}
	}
	fmt.Printf("Total: %d/%d successful\n", succeededCount, totalCount)
}

func fetchAPI(ctx context.Context, messages chan Result, url string) {
	// create request with timeout
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		// don't panic -> all go routines will crash
		messages <- Result{Url: url, Err: err}
	}

	start := time.Now()
	// send timeout
	res, err := http.DefaultClient.Do(req)
	duration := time.Since(start)

	result := Result{
		Url:      url,
		Duration: duration.Round(time.Millisecond),
	}

	if err != nil {
		result.Err = err
	} else {
		defer res.Body.Close()
		result.Status = res.Status
	}

	// send result over channel
	messages <- result
}
