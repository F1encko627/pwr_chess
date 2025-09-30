package board

import (
	"fmt"
	"math"
	"ust_chess/internal/types"
)

func (g *Game) CheckValidRookMove(ix int, iy int, fx int, fy int) error {
	if ix != fx && iy != fy {
		return fmt.Errorf("rook moves horizontally or vertically [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
	}

	return g.CheckNoJumpOverPieceStraight(ix, iy, fx, fy)
}

func (g *Game) CheckValidQueenMove(ix int, iy int, fx int, fy int) error {
	if math.Abs(float64(fx-ix)) != math.Abs(float64(fy-iy)) && (ix != fx && iy != fy) {
		return fmt.Errorf("queen moves in straiht lines [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
	}
	err := g.CheckNoJumpOverPieceStraight(ix, iy, fx, fy)
	if err != nil {
		return err
	}
	err = g.CheckNoJumpOverPieceDiagonal(ix, iy, fx, fy)
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) CheckValidBishopMove(ix int, iy int, fx int, fy int) error {
	if math.Abs(float64(fx-ix)) != math.Abs(float64(fy-iy)) {
		return fmt.Errorf("bishop moves diagonally [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
	}
	return g.CheckNoJumpOverPieceDiagonal(ix, iy, fx, fy)
}

func (g *Game) CheckValidPawnMove(ix int, iy int, fx int, fy int) error {
	// TODO: EnPassant
	var delta_y int
	white := g.Board[iy][ix].White
	if white {
		delta_y = fy - iy
	} else {
		delta_y = iy - fy
	}
	if delta_y < 1 {
		return fmt.Errorf("can move only forward")
	}
	if delta_y > 2 || delta_y > 1 && (iy != 1 && iy != 6) {
		return fmt.Errorf("can't move more than one cell forward after initial move")
	}
	if ix != fx {
		if (delta_y > 1) || (math.Abs(float64(fx-ix)) > 1 && g.Board[fy][fx].T == types.EMPTY) {
			return fmt.Errorf("can move only one cell forward and to the side and only to take opponents piece")
		}
		if g.Board[fy][fx].T == types.EMPTY {
			return fmt.Errorf("can move one cell forward and to the side only to take opponents piece")
		}
	} else if white {
		iy++
		for iy <= fy {
			if g.Board[iy][fx].T != types.EMPTY {
				return fmt.Errorf("can't take piece by moving forvard")
			}
			iy++
		}
	} else {
		iy--
		for iy >= fy {
			if g.Board[iy][fx].T != types.EMPTY {
				return fmt.Errorf("can't take piece by moving forvard")
			}
			iy--
		}
	}
	return nil
}

func (g *Game) CheckValidKingMove(ix int, iy int, fx int, fy int) error {
	if math.Abs(float64(fx-ix)) > 1 || math.Abs(float64(fy-iy)) > 1 {
		return fmt.Errorf("king moves one cell per turn")
	}
	return nil
}

func (g *Game) CheckValidKnightMove(ix int, iy int, fx int, fy int) error {
	if math.Abs(float64(fx-ix))+math.Abs(float64(fy-iy)) != 3 ||
		ix == fx || iy == fy {
		return fmt.Errorf("wrong knight move [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
	}
	return nil
}

func (g *Game) CheckNoJumpOverPieceStraight(ix int, iy int, fx int, fy int) error {
	if iy == fy {
		if ix > fx {
			for i := fx + 1; i < ix; i++ {
				if g.Board[iy][i].T != types.EMPTY {
					return fmt.Errorf("can't jump over pieces")
				}
			}
		} else {
			for i := ix + 1; i < fx; i++ {
				if g.Board[iy][i].T != types.EMPTY {
					return fmt.Errorf("can't jump over pieces")
				}
			}
		}
	} else if ix == fx {
		if iy > fy {
			for i := fy + 1; i < iy; i++ {
				if g.Board[i][ix].T != types.EMPTY {
					return fmt.Errorf("can't jump over pieces")
				}
			}
		} else {
			for i := iy + 1; i < fy; i++ {
				if g.Board[i][ix].T != types.EMPTY {
					return fmt.Errorf("can't jump over piece")
				}
			}
		}
	}
	return nil
}

func (g *Game) CheckNoJumpOverPieceDiagonal(ix int, iy int, fx int, fy int) error {
	err := fmt.Errorf("can't jump over pieces diagonal")

	if ix > fx && iy > fy {
		i, j := fy+1, fx+1
		for i < iy && j < ix {
			if g.Board[i][j].T != types.EMPTY {
				return err
			}
			i++
			j++
		}
	} else if ix < fx && iy < fy {
		i, j := iy+1, ix+1
		for i < fy && j < fx {
			if g.Board[i][j].T != types.EMPTY {
				return err
			}
			i++
			j++
		}
	} else if ix < fx && iy > fy {
		i, j := iy-1, ix+1
		for i > fy && j < fx {
			if g.Board[i][j].T != types.EMPTY {
				return err
			}
			i--
			j++
		}
	} else if ix > fx && iy < fy {
		i, j := iy+1, ix-1
		for i < fy && j > fx {
			if g.Board[i][j].T != types.EMPTY {
				return err
			}
			i++
			j--
		}
	}
	return nil
}
