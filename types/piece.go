package types

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
	T Type
	ID int8
	Pos Pos
	White bool
	Score uint8
}