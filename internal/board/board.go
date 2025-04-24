package board

import (
	"errors"
	"fmt"
	"math"
	"ust_chess/internal/types"
)

type Game struct {
	Board  [8][8]*types.Piece
	Pieces map[types.Type][]types.Piece
	State  types.State
}

var classic = []types.Piece{
	types.GP(types.PAWN, true, types.NewPos(0, 1)),
	types.GP(types.PAWN, true, types.NewPos(1, 1)),
	types.GP(types.PAWN, true, types.NewPos(2, 1)),
	types.GP(types.PAWN, true, types.NewPos(3, 1)),
	types.GP(types.PAWN, true, types.NewPos(4, 1)),
	types.GP(types.PAWN, true, types.NewPos(5, 1)),
	types.GP(types.PAWN, true, types.NewPos(6, 1)),
	types.GP(types.PAWN, true, types.NewPos(7, 1)),

	types.GP(types.ROOK, true, types.NewPos(0, 0)),
	types.GP(types.KNIGHT, true, types.NewPos(1, 0)),
	types.GP(types.BISHOP, true, types.NewPos(2, 0)),
	types.GP(types.KING, true, types.NewPos(3, 0)),
	types.GP(types.QUEEN, true, types.NewPos(4, 0)),
	types.GP(types.BISHOP, true, types.NewPos(5, 0)),
	types.GP(types.KNIGHT, true, types.NewPos(6, 0)),
	types.GP(types.ROOK, true, types.NewPos(7, 0)),

	types.GP(types.PAWN, false, types.NewPos(0, 6)),
	types.GP(types.PAWN, false, types.NewPos(1, 6)),
	types.GP(types.PAWN, false, types.NewPos(2, 6)),
	types.GP(types.PAWN, false, types.NewPos(3, 6)),
	types.GP(types.PAWN, false, types.NewPos(4, 6)),
	types.GP(types.PAWN, false, types.NewPos(5, 6)),
	types.GP(types.PAWN, false, types.NewPos(6, 6)),
	types.GP(types.PAWN, false, types.NewPos(7, 6)),

	types.GP(types.ROOK, false, types.NewPos(0, 7)),
	types.GP(types.KNIGHT, false, types.NewPos(1, 7)),
	types.GP(types.BISHOP, false, types.NewPos(2, 7)),
	types.GP(types.KING, false, types.NewPos(3, 7)),
	types.GP(types.QUEEN, false, types.NewPos(4, 7)),
	types.GP(types.BISHOP, false, types.NewPos(5, 7)),
	types.GP(types.KNIGHT, false, types.NewPos(6, 7)),
	types.GP(types.ROOK, false, types.NewPos(7, 7)),
}

func NewGame(start []types.Piece) Game {
	if len(start) == 0 {
		start = classic
	}
	g := Game{}
	g.Pieces = make(map[types.Type][]types.Piece)
	g.State = types.WHITE_TURN
	empty := types.GP(types.EMPTY, false, types.Pos(0))
	g.Pieces[types.EMPTY] = append(g.Pieces[types.EMPTY], empty)
	for _, p := range start {
		g.Pieces[p.T] = append(g.Pieces[p.T], p)
		g.Board[p.Pos.GetY()][p.Pos.GetX()] = &p
	}
	for i := range 8 {
		for j := range 8 {
			if g.Board[i][j] == nil {
				g.Board[i][j] = &g.Pieces[types.EMPTY][0]
			}
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
	} else if fx < 0 || fx > 8 || fy < 0 || fy > 8 {
		return fmt.Errorf("out of bounds move [%d;%d]", fx, fy)
	}
	// Move to same cell
	if ix == fx && iy == fy {
		return fmt.Errorf("same cell move (no move) [%d;%d]", ix, iy)
	}
	// Moving nothing
	if g.Board[iy][ix].T == types.EMPTY {
		return fmt.Errorf("no piece to move [%d;%d]", ix, iy)
	}
	// Wrong color move
	if (g.State == types.BLACK_TURN ||
		g.State == types.BLACK_CHECK) == g.Board[iy][ix].White ||
		(g.State == types.WHITE_TURN ||
			g.State == types.WHITE_CHECK) != g.Board[iy][ix].White {
		return fmt.Errorf("wrong color move [%d;%d]", ix, iy)
	}
	// Move to own piece
	if g.Board[fy][fx].T != ' ' &&
		(g.Board[fy][fx].White == g.Board[iy][ix].White) {
		return fmt.Errorf("can't beat same color [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
	}
	// Piece move validation
	switch g.Board[iy][ix].T {
	case types.KING:
		if math.Abs(float64(fx-ix)) > 1 || math.Abs(float64(fy-iy)) > 1 {
			return fmt.Errorf("king moves one cell per turn")
		}
	case types.QUEEN:
		if math.Abs(float64(fx-ix)) != math.Abs(float64(fy-iy)) && (ix != fx && iy != fy) {
			return fmt.Errorf("queen moves in straiht lines [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
		}
		if iy == fy {
			if ix > fx {
				for i := fx + 1; i < ix; i++ {
					if g.Board[iy][i].T != types.EMPTY {
						return fmt.Errorf("queen can't jump over pieces")
					}
				}
			} else {
				for i := ix + 1; i < fx; i++ {
					if g.Board[iy][i].T != types.EMPTY {
						return fmt.Errorf("queen can't jump over pieces")
					}
				}
			}
		} else if ix == fx {
			if iy > fy {
				for i := fy + 1; i < iy; i++ {
					if g.Board[i][ix].T != types.EMPTY {
						return fmt.Errorf("queen can't jump over pieces")
					}
				}
			} else {
				for i := iy + 1; i < fy; i++ {
					if g.Board[i][ix].T != types.EMPTY {
						return fmt.Errorf("queen can't jump over pieces")
					}
				}
			}
		} else if ix > fx && iy > fy {
			i, j := fy+1, fx+1
			for i < iy && j < ix {
				if g.Board[i][j].T != types.EMPTY {
					return fmt.Errorf("queen can't jump over pieces")
				}
				i++
				j++
			}
		} else if ix < fx && iy < fy {
			i, j := iy+1, ix+1
			for i < fy && j < fx {
				if g.Board[i][j].T != types.EMPTY {
					return fmt.Errorf("queen can't jump over pieces")
				}
				i++
				j++
			}
		} else if ix < fx && iy > fy {
			i, j := iy-1, ix+1
			for i > fy && j < fx {
				if g.Board[i][j].T != types.EMPTY {
					return fmt.Errorf("queen can't jump over pieces")
				}
				i--
				j++
			}
		} else if ix > fx && iy < fy {
			i, j := iy+1, ix-1
			for i < fy && j > fx {
				if g.Board[i][j].T != types.EMPTY {
					return fmt.Errorf("queen can't jump over pieces")
				}
				i++
				j--
			}
		}
	case types.BISHOP:
		if math.Abs(float64(fx-ix)) != math.Abs(float64(fy-iy)) {
			return fmt.Errorf("bishop moves diagonally [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
		}
		if ix > fx && iy > fy {
			i, j := fy+1, fx+1
			for i < iy && j < ix {
				if g.Board[i][j].T != types.EMPTY {
					return fmt.Errorf("bishop can't jump over pieces")
				}
				i++
				j++
			}
		} else if ix < fx && iy < fy {
			i, j := iy+1, ix+1
			for i < fy && j < fx {
				if g.Board[i][j].T != types.EMPTY {
					return fmt.Errorf("bishop can't jump over pieces")
				}
				i++
				j++
			}
		} else if ix < fx && iy > fy {
			i, j := iy-1, ix+1
			for i > fy && j < fx {
				if g.Board[i][j].T != types.EMPTY {
					return fmt.Errorf("bishop can't jump over pieces")
				}
				i--
				j++
			}
		} else if ix > fx && iy < fy {
			i, j := iy+1, ix-1
			for i < fy && j > fx {
				if g.Board[i][j].T != types.EMPTY {
					return fmt.Errorf("bishop can't jump over pieces")
				}
				i++
				j--
			}
		}
	case types.ROOK:
		if ix != fx && iy != fy {
			return fmt.Errorf("rook moves horizontally or vertically [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
		}
		if iy == fy {
			if ix > fx {
				for i := fx + 1; i < ix; i++ {
					if g.Board[iy][i].T != types.EMPTY {
						return fmt.Errorf("rook can't jump over pieces")
					}
				}
			} else {
				for i := ix + 1; i < fx; i++ {
					if g.Board[iy][i].T != types.EMPTY {
						return fmt.Errorf("rook can't jump over pieces")
					}
				}
			}
		} else if ix == fx {
			if iy > fy {
				for i := fy + 1; i < iy; i++ {
					if g.Board[i][ix].T != types.EMPTY {
						return fmt.Errorf("rook can't jump over pieces")
					}
				}
			} else {
				for i := iy + 1; i < fy; i++ {
					if g.Board[i][ix].T != types.EMPTY {
						return fmt.Errorf("rook can't jump over pieces")
					}
				}
			}
		}
	case types.KNIGHT:
		if math.Abs(float64(fx-ix))+math.Abs(float64(fy-iy)) != 3 ||
			ix == fx || iy == fy {
			return fmt.Errorf("wrong knight move [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
		}
	case types.PAWN:
	}

	g.Board[fy][fx] = g.Board[iy][ix]
	g.Board[iy][ix] = &types.Piece{T: types.EMPTY}
	if g.State == types.WHITE_TURN {
		g.State = types.BLACK_TURN
	} else {
		g.State = types.WHITE_TURN
	}
	return nil
}
