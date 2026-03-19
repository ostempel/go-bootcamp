package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		panic("provide file as argument")
	}
	fileArg := args[1]

	log.Println("file", fileArg)

	f, err := os.Open(fileArg)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var (
		lineCount      = 0
		wordCount      = 0
		characterCount = 0
	)

	for s.Scan() {
		text := s.Text()

		lineCount++
		characterCount += len(text)

		words := strings.Fields(text)
		wordCount += len(words)
	}

	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Lines: %d\nWords: %d\nCharacters: %d\n", lineCount, wordCount, characterCount)
}
