package main

import (
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	if len(port) == 0 {
		port = "80"
	}
	http.ListenAndServe(":"+port, http.DefaultServeMux)
}
