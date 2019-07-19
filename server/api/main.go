package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/api/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `hello`)
	}))

	err := http.ListenAndServeTLS(":3001", "./localhost.pem", "./localhost-key.pem", mux)
	//err := http.ListenAndServe(":3001", mux)
	if err != nil {
		panic(err)
	}
}
