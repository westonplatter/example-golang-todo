package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
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

// todo "Object"
type Todo struct {
	Id         int    `json:"Id"`
	Title      string `json:"Title"`
	Category   string `json:"Category"`
	Dt_created int64  `json:"int64"`
	Dt_updated int64  `json:"int64"`
	State      string `json:"State"`
}

// store "context" values and connections in the server struct
type Server struct {
	db *sql.DB
}

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang_todo_dev")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(100)
	defer db.Close()

	server := &Server{db: db}

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

// simple HTML/JS/CSS pages

func (s *Server) homepage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./index.html")
}

func (s *Server) assets(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, req.URL.Path[1:])
}

// Todo CRUD

func (s *Server) todoIndex(res http.ResponseWriter, req *http.Request) {
	var todos []*Todo

	rows, err := s.db.Query("SELECT Id, Title, Category, Dt_created, Dt_updated, State FROM Todo")
	error_check(res, err)
	for rows.Next() {
		todo := &Todo{}
		rows.Scan(&todo.Id, &todo.Title, &todo.Category, &todo.Dt_created, &todo.Dt_updated, &todo.State)
		todos = append(todos, todo)
	}
	rows.Close()

	jsonResponse(res, todos)
}

func (s *Server) todoCreate(res http.ResponseWriter, req *http.Request) {
	todo := &Todo{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&todo)
	if err != nil {
		fmt.Println("ERROR decoding JSON - ", err)
		return
	}

	result, err := s.db.Exec("INSERT INTO Todo(Title, Category, State, Dt_created) VALUES(?, ?, ?, ?)", todo.Title, todo.Category, todo.State, todo.Dt_created)

	if err != nil {
		fmt.Println("ERROR saving to db - ", err)
	}

	newId64, err := result.LastInsertId()
	newId := int(newId64)
	todo = &Todo{Id: newId}

	s.db.
		QueryRow("SELECT State, Title, Category, Dt_created, Dt_updated FROM Todo WHERE Id=?", todo.Id).
		Scan(&todo.State, &todo.Title, &todo.Category, &todo.Dt_created, &todo.Dt_updated)

	jsonResponse(res, todo)
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

func jsonResponse(res http.ResponseWriter, data interface{}) {

	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	payload, err := json.Marshal(data)
	if error_check(res, err) {
		return
	}

	fmt.Fprintf(res, string(payload))
}

func error_check(res http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
