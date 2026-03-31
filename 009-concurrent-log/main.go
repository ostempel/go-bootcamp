package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	LogLevel  string
	Message   string
	Error     error
}

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		log.Fatal("command needs 2 args: <input.log> <count-workers>")
	}

	filepath := args[0]
	workerCounter, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("unable to parse worker count", "worker-count", args[1])
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("unable to read log file", "filepath", filepath)
	}

	s := bufio.NewScanner(file)

	jobs := make(chan string)
	results := make(chan LogEntry)

	// waitgroup to wait for all done workers before
	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(workerCounter)

	for w := 1; w <= workerCounter; w++ {
		// spawn workers handling jobs
		go handleLogEntry(w, &wg, jobs, results)
	}

	// send jobs in separate goroutine to avoid deadlock with results collection
	go func() {
		for s.Scan() {
			text := s.Text()
			jobs <- text
		}
		// close in goroutine -> outside it would immediately close
		close(jobs)
		err = s.Err()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// wait till all workers finished -> close results channel -> read entries of channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// if results not closed -> deadlock (waits for more messages)
	for entry := range results {
		log.Println(entry.Message)
	}
}

func handleLogEntry(id int, wg *sync.WaitGroup, jobs <-chan string, results chan<- LogEntry) {
	defer wg.Done()
	log.Printf("spawned worker %d\n", id)

	for j := range jobs {
		log.Printf("worker %d starting with new job", id)
		fields := strings.Fields(j)

		t := strings.Join(fields[0:2], " ")
		timestamp, err := time.Parse(time.DateTime, t)
		if err != nil {
			log.Printf("worker %d errored with error: %v", id, err)
			results <- LogEntry{Error: err, Message: j}
			continue
		}

		logLevel := fields[2]
		message := strings.Join(fields[3:], " ")
		log.Printf("worker %d finished with job", id)
		results <- LogEntry{Timestamp: timestamp, LogLevel: logLevel, Message: message}
	}

	log.Printf("destroyed worker %d\n", id)
}
