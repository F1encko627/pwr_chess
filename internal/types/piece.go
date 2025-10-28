package types

import "fmt"

type Figure rune

const (
	EMPTY Figure = ' '
	KING         = '♚' + iota
	QUEEN
	ROOK
	BISHOP
	KNIGHT
	PAWN
)

type Piece struct {
	T        Figure
	ID       int
	Pos      Position
	IsInGame bool
	White    bool
	Score    uint8
}

var sequence = int(-1)

func (t *Figure) String() string {
	return string(*t)
}

func (p *Piece) String() string {
	return fmt.Sprint(p.T.String(), p.White, p.Pos.String(), p.Score, p.ID)
}

func NewPiece(T Figure, White bool, pos Position) Piece {
	sequence++
	return Piece{
		T:        T,
		ID:       sequence,
		Pos:      pos,
		IsInGame: true,
		White:    White,
		Score:    0,
	}
}
