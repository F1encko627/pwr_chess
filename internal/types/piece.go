package types

import "fmt"

type Type rune

const (
	KING Type = '♚' + iota
	QUEEN
	ROOK
	BISHOP
	KNIGHT
	PAWN
	EMPTY = ' '
)

type Piece struct {
	T     Type
	ID    int8
	Pos   Pos
	White bool
	Score uint8
}

var seq = int8(-1)

func (t *Type) String() string {
	return string(*t)
}

func (p *Piece) String() string {
	return fmt.Sprint(p.T.String(), p.White, p.Pos.String(), p.Score, p.ID)
}

// "NewPiece" or "GetPiece" piece with correct ID and other fields
func GP(T Type, White bool, pos Pos) Piece {
	if T == EMPTY {
		return Piece{
			T:     EMPTY,
			ID:    -1,
			Pos:   Pos(-1),
			White: false,
			Score: 0,
		}
	}
	seq++
	return Piece{
		T:     T,
		ID:    seq,
		Pos:   pos,
		White: White,
		Score: 0,
	}
}
