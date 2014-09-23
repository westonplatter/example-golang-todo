package main

import (
	"fmt"
	"net/http"
	"regexp"
)

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *RegexpHandler) HandleFunc(s string, handler func(http.ResponseWriter, *http.Request)) {
	rex := regexp.MustCompile(s)
	h.routes = append(h.routes, &route{rex, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	server := &Server{}

	reHandler := new(RegexpHandler)
	reHandler.HandleFunc("/todos/[0-9]+$", server.showTodo)
	reHandler.HandleFunc("/todos$", server.todoIndex)

	reHandler.HandleFunc("/", server.homepage)

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
