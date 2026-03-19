package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var client = &http.Client{Timeout: 5 * time.Second}

func main() {
	log.Println("Starting healthcheck!")
	urls := os.Args[1:]

	// create waitgroup of specific length
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go checkHealth(url, &wg)
	}

	// wait for every func finished
	wg.Wait()

	log.Println("Healthcheck done!")
}

func checkHealth(url string, wg *sync.WaitGroup) {
	// always defer at beginning -> panic in func can cause errors
	defer wg.Done()

	start := time.Now()
	res, err := client.Get(url)

	duration := time.Since(start)

	if err != nil {
		fmt.Printf("%s - FAIL (%s) - %s\n", url, err.Error(), duration)
	} else {
		// close resp body --> leaking body (http-connections stays open)
		defer res.Body.Close()
		fmt.Printf("%s - (%s) - %s\n", url, res.Status, duration)
	}
}
