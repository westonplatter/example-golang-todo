package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/todos", TodoHandler)
	r.HandleFunc("/", HomepageHandler)

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}

func HomepageHandler(req http.ResponseWriter, res *http.Request) {
	fmt.Fprintf(req, "<h1>HomePage</h1>")
}

func TodoHandler(req http.ResponseWriter, res *http.Request) {
	fmt.Fprintf(req, "<h2>Todos</h2>")
}
