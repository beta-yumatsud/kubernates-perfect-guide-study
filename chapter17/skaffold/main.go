package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request)  {
	_, _ = fmt.Fprintf(w, "Hello, Skaffold")
}

func main() {
	http.HandleFunc("/", handler)
	_ = http.ListenAndServe(":8080", nil)
}


