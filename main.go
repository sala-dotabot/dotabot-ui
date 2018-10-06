package main

import (
	"log"
	"net/http"
)

func main() {
	context, error := InitContext()
	if error != nil {
		log.Fatalf("Error while initializing context %s", error)
	}

	http.HandleFunc("/", context.Handler.Handle)
    log.Fatal(http.ListenAndServe(":8080", nil))
}