package board

import (
	"errors"
	"fmt"
	"math"
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
				White: i < 3,
				Score: uint8(0),
			}
			g.Pieces[p.T] = p
			g.Board[i][k] = &p
		}
	}
	return g
}
/*
	TODO:
	- Математическая нотация
	- Проверить не быстрее ли использовать одномерный массив вместо двумерного

	- Проверка на валидность хода для каждой фигуры
*/
func (g *Game) MovePiece(ix int, iy int, fx int, fy int) error {
	// Game over
	if g.State == types.BLACK_CHECKMATE ||
	g.State == types.WHITE_CHECKMATE || 
	g.State == types.STALEMATE {
		return errors.New("game over: can't make moves")
	}
	// Outside of the board
	if ix < 0 || ix > 8 || iy < 0 || iy > 8 {
		return fmt.Errorf("out of bounds position [%d;%d]", ix, iy)
	} else if  fx < 0 || fx > 8 || fy < 0 || fy > 8 {
		return fmt.Errorf("out of bounds move [%d;%d]", fx, fy)
	}
	// Move to same cell
	if ix == fx && iy == fy {
		return fmt.Errorf("same cell move (no move) [%d;%d]", ix, iy)
	}
	// Moving nothing
	if g.Board[ix][iy].T == types.EMPTY {
		return fmt.Errorf("no piece to move [%d;%d]", ix, iy)
	}
	// Wrong color move
	if (g.State == types.BLACK_TURN ||
		g.State == types.BLACK_CHECK) == g.Board[ix][iy].White ||
	(g.State == types.WHITE_TURN ||
		g.State == types.WHITE_CHECK) != g.Board[ix][iy].White {
		return fmt.Errorf("wrong color move [%d;%d]", ix, iy)
	}
	// Move to own piece
	if g.Board[fx][fy].T != ' ' &&
	(g.Board[fx][fy].White == g.Board[ix][iy].White) {
		return fmt.Errorf("can't beat same color [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
	}
	// Piece move validation
	switch g.Board[ix][iy].T {
	case types.KING:
	case types.QUEEN:
		if math.Abs(float64(fx-ix)) != math.Abs(float64(fy-iy)) || (ix != fx && iy != fy) {
			return fmt.Errorf("queen moves in straiht lines [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
		}
		// TODO: проверять что не прыгаем через фигуры
	case types.BISHOP:
		if math.Abs(float64(fx-ix)) != math.Abs(float64(fy-iy)) {
			return fmt.Errorf("bishop moves diagonally [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
		}
		// TODO: проверять что не прыгаем через фигуры
	case types.ROOK:
		if ix != fx && iy != fy {
			return fmt.Errorf("rook moves horizontally or vertically [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
		}
		// TODO: проверять что не прыгаем через фигуры
	case types.KNIGHT:
		if math.Abs(float64(fx-ix)) + math.Abs(float64(fy-iy)) != 3 ||
		ix == fx || iy == fy {
			return fmt.Errorf("wrong knight move [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
		}
	case types.PAWN:
	}
	
	g.Board[fx][fy] = g.Board[ix][iy]
	g.Board[ix][iy] = &types.Piece{T: types.EMPTY}
	if g.State == types.WHITE_TURN {
		g.State = types.BLACK_TURN
	} else {
		g.State = types.WHITE_TURN
	}
	return nil
}