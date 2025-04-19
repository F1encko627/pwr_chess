package board

import (
	"ust_chess/types"
)

type Game struct {
	Board [8][8]*types.Piece
	Pieces map[types.Type]types.Piece
	State types.State
}

func NewGame() Game {
	g := Game{}
	start := [][]types.Type{
		{'♜', '♞', '♝', '♛', '♚', '♝', '♞', '♜'},
		{'♟', '♟', '♟', '♟', '♟', '♟', '♟', '♟'},
		{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		{'♟', '♟', '♟', '♟', '♟', '♟', '♟', '♟'},
		{'♜', '♞', '♝', '♚', '♛', '♝', '♞', '♜'},
	}
	g.Pieces = make(map[types.Type]types.Piece, 16)
	g.State = types.WHITE_TURN
	for i := range 8 {
		for k := range 8 {
			p := types.Piece{
				T: start[i][k],
				ID: int8(i+k),
				Pos: types.NewPos(k, i),
				White: i > 3,
				Score: uint8(0),
			}
			g.Pieces[p.T] = p
			g.Board[i][k] = &p
		}
	}
	return g
}
/*
	TODO
	- Математическая нотация
	- Проверить не быстрее ли использовать одномерный массив вместо двумерного

	- Проверка на валидность хода для каждой фигуры
*/
func (g *Game) MovePiece(ix int, iy int, fx int, fy int) error {
	// Move outside of the board
	if ix < 0 || ix > 8 || iy < 0 || iy > 8 || fx < 0 || fx > 8 || fy < 0 || fy > 8 {
		return nil
	}
	// To same cell
	if ix == fx && iy == ix {
		return nil
	}
	// Move nothing
	if g.Board[ix][iy].T == ' ' {
		return nil
	}
	// To own piece
	if g.Board[fx][fy].T != ' ' && (g.Board[fx][fy].White == g.Board[ix][iy].White) {
		return nil;
	}

	switch g.Board[fx][fy].T {
	case types.KING:
	case types.QUEEN:
	case types.BISHOP:
	case types.ROOK:
	case types.KNIGHT:
	case types.PAWN:
	}
	
	g.Board[fx][fy] = g.Board[ix][iy]
	g.Board[ix][iy] = &types.Piece{T: ' '}
	return nil
}