package board

import (
	"errors"
	"fmt"
	"math"
	"time"
	"ust_chess/internal/types"
)

type Game struct {
	Board         [8][8]*types.Piece
	WhitePieces   map[types.Type][]types.Piece
	BlackPieces   map[types.Type][]types.Piece // plus empty tile
	BlackTurn     bool
	KingChecked   bool
	Checkmate     bool
	Pause         bool
	EnPassantPawn *types.Piece
	LastMoveTime  time.Time
	MoveFunc      *func(Game) (int, int, int, int)
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

	g.WhitePieces = make(map[types.Type][]types.Piece)
	g.BlackPieces = make(map[types.Type][]types.Piece)

	empty := types.GP(types.EMPTY, false, types.Pos(0))

	g.BlackPieces[types.EMPTY] = append(g.BlackPieces[types.EMPTY], empty)

	for _, p := range start {
		if p.White {
			g.WhitePieces[p.T] = append(g.WhitePieces[p.T], p)
		} else {
			g.BlackPieces[p.T] = append(g.BlackPieces[p.T], p)
		}
		g.Board[p.Pos.GetY()][p.Pos.GetX()] = &p
	}

	for i := range 8 {
		for j := range 8 {
			if g.Board[i][j] == nil {
				g.Board[i][j] = &g.BlackPieces[types.EMPTY][0]
			}
		}
	}

	return g
}

func (g Game) SetTurnSide(whiteTurn bool) Game {
	g.BlackTurn = !whiteTurn
	return g
}

func (g *Game) DebugRender() {
	for y := range 8 {
		for x := range 8 {
			piece := g.Board[y][x]
			w := "B"
			if piece.White {
				w = "W"
			}
			fmt.Printf("| y:%d x:%d %s%s ", y, x, w, string(piece.T))
		}
		fmt.Print("|\n")
	}
}

/*
TODO:
- Математическая нотация
- Проверить не быстрее ли использовать одномерный массив вместо двумерного

- Проверка на валидность хода для каждой фигуры
*/
func (g *Game) MovePiece(ix int, iy int, fx int, fy int) error {
	// Pause
	if g.Pause {
		return errors.New("game paused: can't make moves")
	}
	// Game over
	if g.Checkmate {
		return errors.New("game over: can't make moves")
	}
	// Outside of the board
	if ix < 0 || ix > 8 || iy < 0 || iy > 8 {
		return fmt.Errorf("out of bounds position [%d;%d]", ix, iy)
	} else if fx < 0 || fx > 8 || fy < 0 || fy > 8 {
		return fmt.Errorf("out of bounds move [%d;%d]", fx, fy)
	}
	// Same cell move
	if ix == fx && iy == fy {
		return fmt.Errorf("same cell move (no move) [%d;%d]", ix, iy)
	}
	// Moving nothing
	if g.Board[iy][ix].T == types.EMPTY {
		return fmt.Errorf("no piece to move [%d;%d]", ix, iy)
	}
	// Wrong color move
	if g.BlackTurn == g.Board[iy][ix].White {
		return fmt.Errorf("wrong color move [%d;%d]", ix, iy)
	}
	// Trying take same color piece
	if g.Board[fy][fx].T != ' ' &&
		(g.Board[fy][fx].White == g.Board[iy][ix].White) {
		return fmt.Errorf("can't beat same color [%d;%d] -> [%d;%d]", ix, iy, fx, fy)
	}
	// Per piece move validation
	switch g.Board[iy][ix].T {
	case types.KING:
		if err := g.CheckValidKingMove(ix, iy, fx, fy); err != nil {
			return err
		}
	case types.QUEEN:
		if err := g.CheckValidQueenMove(ix, iy, fx, fy); err != nil {
			return err
		}
	case types.BISHOP:
		if err := g.CheckValidBishopMove(ix, iy, fx, fy); err != nil {
			return err
		}
	case types.ROOK:
		if err := g.CheckValidRookMove(ix, iy, fx, fy); err != nil {
			return err
		}
	case types.KNIGHT:
		if err := g.CheckValidKnightMove(ix, iy, fx, fy); err != nil {
			return err
		}
	case types.PAWN:
		if err := g.CheckValidPawnMove(ix, iy, fx, fy); err != nil {
			return err
		}
	}

	var ownKingPos, enemyKingPos = types.Pos(-1), types.Pos(-1)
	if _, ok := g.WhitePieces[types.KING]; ok {
		ownKingPos = g.WhitePieces[types.KING][0].Pos
	}
	if _, ok := g.BlackPieces[types.KING]; ok {
		enemyKingPos = g.BlackPieces[types.KING][0].Pos
	}
	if g.BlackTurn {
		ownKingPos, enemyKingPos = enemyKingPos, ownKingPos
	}

	game_save := *g
	g.KingChecked = false

	// General move execution
	g.Board[fy][fx] = g.Board[iy][ix]
	g.Board[iy][ix] = &g.BlackPieces[types.EMPTY][0]

	if ownKingPos.IsValid() && g.IsKingChecked(ownKingPos.GetX(), ownKingPos.GetY(), !g.BlackTurn) {
		// Revert illigal move
		*g = game_save
		return fmt.Errorf("invalid turn: own king still checked")
	}

	// enemy king status
	if enemyKingPos.IsValid() {
		g.KingChecked = g.IsKingChecked(enemyKingPos.GetX(), enemyKingPos.GetY(), g.BlackTurn)
		if g.KingChecked {
			// check for checkmate
		}
	}

	// Other color turn
	g.BlackTurn = !g.BlackTurn

	return nil
}

func (g *Game) IsKingChecked(x, y int, kingIsWhite bool) (is_check bool) {
	// up
	for i := y + 1; i < 8; i++ {
		if g.Board[i][x].White != kingIsWhite &&
			(g.Board[i][x].T == types.ROOK ||
				g.Board[i][x].T == types.QUEEN) {
			return true
		} else if g.Board[y][i].T != types.EMPTY {
			break
		}
	}
	// down
	for i := y - 1; i > 8; i-- {
		if g.Board[i][x].White != kingIsWhite &&
			(g.Board[i][x].T == types.ROOK ||
				g.Board[i][x].T == types.QUEEN) {
			return true
		} else if g.Board[y][i].T != types.EMPTY {
			break
		}
	}
	// left
	for i := x + 1; i < 8; i++ {
		if g.Board[y][i].White != kingIsWhite &&
			(g.Board[y][i].T == types.ROOK ||
				g.Board[y][i].T == types.QUEEN) {
			return true
		} else if g.Board[y][i].T != types.EMPTY {
			break
		}
	}
	// right
	for i := x - 1; i > 8; i-- {
		if g.Board[y][i].White != kingIsWhite &&
			(g.Board[y][i].T == types.ROOK ||
				g.Board[y][i].T == types.QUEEN) {
			return true
		} else if g.Board[y][i].T != types.EMPTY {
			break
		}
	}
	// up right
	for i, j := y-1, x+1; i > 0 && j < 8; {
		if g.Board[y][i].White != kingIsWhite &&
			(g.Board[y][i].T == types.BISHOP ||
				g.Board[y][i].T == types.QUEEN) {
			return true
		} else if g.Board[y][i].T != types.EMPTY {
			break
		}
		i--
		j++
	}

	// up left
	for i, j := y-1, x-1; i > 0 && j > 0; {
		if g.Board[y][i].White != kingIsWhite &&
			(g.Board[y][i].T == types.BISHOP ||
				g.Board[y][i].T == types.QUEEN) {
			return true
		} else if g.Board[y][i].T != types.EMPTY {
			break
		}
		i--
		j--
	}

	// down right
	for i, j := y+1, x+1; i < 8 && j < 8; {
		if g.Board[y][i].White != kingIsWhite &&
			(g.Board[y][i].T == types.BISHOP ||
				g.Board[y][i].T == types.QUEEN) {
			return true
		} else if g.Board[y][i].T != types.EMPTY {
			break
		}
		i++
		j++
	}

	// down left
	for i, j := y+1, x-1; i < 8 && j > 0; {
		if g.Board[y][i].White != kingIsWhite &&
			(g.Board[y][i].T == types.BISHOP ||
				g.Board[y][i].T == types.QUEEN) {
			return true
		} else if g.Board[y][i].T != types.EMPTY {
			break
		}
		i++
		j--
	}

	// KNIGHT CHECK
	k := []types.Pos{
		types.NewPos(x+3, y+2),
		types.NewPos(x+2, y+3),
		types.NewPos(x-3, y-2),
		types.NewPos(x-2, y-3),
		types.NewPos(x-3, y+2),
		types.NewPos(x-2, y+3),
		types.NewPos(x+3, y-2),
		types.NewPos(x+2, y-3),
	}
	for _, p := range k {
		if kingIsWhite {
			for _, k := range g.BlackPieces[types.KNIGHT] {
				if k.Pos == p {
					return true
				}
			}
		} else {
			for _, k := range g.WhitePieces[types.KNIGHT] {
				if k.Pos == p {
					return true
				}
			}
		}
	}

	return false
}

func (g *Game) CheckValidRookMove(ix int, iy int, fx int, fy int) error {
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
	return nil
}

func (g *Game) CheckValidQueenMove(ix int, iy int, fx int, fy int) error {
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
	return nil
}

func (g *Game) CheckValidBishopMove(ix int, iy int, fx int, fy int) error {
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
	return nil
}

func (g *Game) CheckValidPawnMove(ix int, iy int, fx int, fy int) error {
	// TODO: EnPassant
	var diff int
	white := g.Board[iy][ix].White
	if white {
		diff = fy - iy
	} else {
		diff = iy - fy
	}
	if diff < 1 {
		return fmt.Errorf("can move only forward")
	}
	if diff > 2 || diff > 1 && (iy != 1 && iy != 6) {
		return fmt.Errorf("can't move more than one cell forward after initial move")
	}
	if ix != fx {
		if (diff > 1) || (math.Abs(float64(fx-ix)) > 1 && g.Board[fy][fx].T == types.EMPTY) {
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
