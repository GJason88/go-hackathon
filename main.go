package main

import (
	"fmt"
	"log"
	"net/http"
)

func baseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "The time is %s", time)
}

func main() {
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/", baseHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// This is a comment.
