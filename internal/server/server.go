package server

import "net/http"

type Server struct {
	Port int
	Mux http.ServeMux
	Handlers map[string][]http.HandlerFunc
}

func Serve() {
	
}

func GetRoot(w http.ResponseWriter, r *http.Request) {
	
}

func MovePiece(w http.ResponseWriter, r *http.Request) {
	
}

func GetBoard(w http.ResponseWriter, r *http.Request) {
	
}