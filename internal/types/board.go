package types

type Board struct {
	board  [8][8]*Piece
	pieces map[bool][]Piece
}

func GetBoard() (Board, error) {
	return Board{}, nil
}
