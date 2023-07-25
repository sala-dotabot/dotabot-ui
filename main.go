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

	go func() {
		serveMux := http.NewServeMux()
		serveMux.Handle("/metrics", context.MetricsHandler)
		log.Fatal(http.ListenAndServe(":8090", serveMux))
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
