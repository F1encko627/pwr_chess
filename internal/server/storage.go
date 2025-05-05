package server

import (
	"ust_chess/internal/board"
)

type Room struct {
	ID         int
	Game       board.Game
	White      User
	Black      User
	Spectators []User
}

type User struct {
	ID     int
	Name   string
	Tocken string
}
