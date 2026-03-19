package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func handleHello(writer http.ResponseWriter, req *http.Request) {
	io.WriteString(writer, "Hello World!\n")
}

func handleTime(writer http.ResponseWriter, req *http.Request) {
	io.WriteString(writer, time.Now().String()+"\n")
}

// interceptor pattern -> able to chain interceptors
func interceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("method", req.Method, "path", req.URL.Path)
		next.ServeHTTP(w, req)
	})
}

func main() {
	log.Println("Starting server...")

	helloHandler := interceptor(http.HandlerFunc(handleHello))
	timeHandler := interceptor(http.HandlerFunc(handleTime))

	http.Handle("/hello", helloHandler)
	http.Handle("/time", timeHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
