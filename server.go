package main

import (
	"fmt"
	"net/http"
	"regexp"
)

// net/http based router

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

// store "context" values and connections in the server struct

type Server struct{}

func main() {
	server := &Server{}

	reHandler := new(RegexpHandler)

	reHandler.HandleFunc("/todos/$", "GET", server.todoIndex)
	reHandler.HandleFunc("/todos/$", "POST", server.todoCreate)
	reHandler.HandleFunc("/todos/[0-9]+$", "GET", server.todoShow)
	reHandler.HandleFunc("/todos/[0-9]+$", "PUT", server.todoUpdate)
	reHandler.HandleFunc("/todos/[0-9]+$", "DELETE", server.todoDelete)

	reHandler.HandleFunc(".*.[js|css|png|eof|svg|ttf|woff]", "GET", server.assets)
	reHandler.HandleFunc("/", "GET", server.homepage)

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", reHandler)
}

func (s *Server) homepage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./index.html")
}

func (s *Server) assets(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, req.URL.Path[1:])
}

func (s *Server) todoIndex(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Array of todo json")
}

func (s *Server) todoCreate(res http.ResponseWriter, req *http.Request) {
	fmt.Println(res, "Created Todo. Send back todo json")
	fmt.Println("Created Todo. Send back todo json")
}

func (s *Server) todoShow(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Render todo json")
}

func (s *Server) todoUpdate(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Updated todo. Render todo json")
}

func (s *Server) todoDelete(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Deleted todo. Render todo json")
}
