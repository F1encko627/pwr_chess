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
	T     Type
	ID    int8
	Pos   Pos
	White bool
	Score uint8
}

var seq = int8(-1)

// "NewPiece" or "GetPiece" piece with correct ID and other fields
func GP(T Type, White bool, pos Pos) Piece {
	if T == EMPTY {
		return Piece{
			T:     EMPTY,
			ID:    -1,
			Pos:   NewPos(-1, -1),
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
