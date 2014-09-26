package main

import (
	"fmt"
	"net/http"
	"regexp"
)

type route struct {
	pattern *regexp.Regexp
	verb    string
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) Handler(pattern *regexp.Regexp, verb string, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, verb, handler})
}

func (h *RegexpHandler) HandleFunc(r string, v string, handler func(http.ResponseWriter, *http.Request)) {
	re := regexp.MustCompile(r)
	h.routes = append(h.routes, &route{re, v, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) && route.verb == r.Method {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	server := &Server{}

	reHandler := new(RegexpHandler)

	reHandler.HandleFunc("/todos/[0-9]+$", "GET", server.showTodo)
	reHandler.HandleFunc("/todos$", "GET", server.todoIndex)

	reHandler.HandleFunc("/", "GET", server.homepage)

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", reHandler)
}

type Server struct{}

func (s *Server) homepage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "<h1>HomePage</h1>")
}

func (s *Server) todoIndex(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "<h2>Array of Todos JSON</h2>")
}

func (s *Server) showTodo(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "<h2>Todo JSON</h2>")
}
