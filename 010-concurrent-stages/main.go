package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	numbers := generator(args)
	squared := square(numbers)
	filtered := filter(squared, 100)

	for n := range filtered {
		fmt.Println(n)
	}
}

func generator(input []string) <-chan int {
	numbers := make(chan int)

	go func() {
		for _, value := range input {
			num, err := strconv.Atoi(value)
			if err != nil {
				log.Printf("unable to parse input %q: %v\n", value, err)
				continue
			}

			numbers <- num
		}
		close(numbers)
	}()

	return numbers
}

func square(numbers <-chan int) <-chan int {
	squared := make(chan int)

	go func() {
		for value := range numbers {
			squared <- value * value
		}
		close(squared)
	}()

	return squared
}

func filter(squared <-chan int, threshold int) <-chan int {
	filtered := make(chan int)

	go func() {
		for value := range squared {
			if value >= threshold {
				filtered <- value
			}
		}
		close(filtered)
	}()

	return filtered
}
