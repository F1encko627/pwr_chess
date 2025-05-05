package server

import "net/http"

type Server struct {
	Port      int
	Mux       http.ServeMux
	Templates map[string][]http.HandlerFunc
}

func Serve() {

}

func Index(w http.ResponseWriter, r *http.Request) {

}

func Play(w http.ResponseWriter, r *http.Request) {

}

func EnterRoom(w http.ResponseWriter, r *http.Request) {

}

func Create(w http.ResponseWriter, r *http.Request) {

}
