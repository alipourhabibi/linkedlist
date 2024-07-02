package main

import (
	v1 "linkedlist/api/v1"
	v2 "linkedlist/api/v2"
	"net/http"
)

func main() {
	v1 := v1.V1()
	v2, err := v2.V2()
	if err != nil {
		// TODO
	}

	mux := http.NewServeMux()
	mux.Handle("/v1/", http.StripPrefix("/v1", v1))
	mux.Handle("/v2/", http.StripPrefix("/v2", v2))

	http.ListenAndServe(":8080", mux)
}
