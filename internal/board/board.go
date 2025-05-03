package board

import (
	"errors"
	"fmt"
	"math"
	"time"
	"ust_chess/internal/types"

	"github.com/rs/zerolog/log"
)

type Game struct {
	Board            [8][8]*types.Piece
	WhitePieces      map[types.Type][]types.Piece
	BlackPieces      map[types.Type][]types.Piece // plus empty tile
	IsBlackTurn      bool
	IsWhiteCanCastle bool
	IsBlackCanCastle bool
	IsKingChecked    bool
	IsCheckmate      bool
	IsPause          bool
	EnPassantPawn    *types.Piece
	LastMoveTime     time.Time
	MoveFunc         *func(Game) (int, int, int, int)
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

	empty := types.GP(types.EMPTY, false, types.Pos(-1))

	g.BlackPieces[types.EMPTY] = append(g.BlackPieces[types.EMPTY], empty)

	for _, p := range start {
		if p.White {
			g.WhitePieces[p.T] = append(g.WhitePieces[p.T], p)
			g.Board[p.Pos.GetY()][p.Pos.GetX()] = &g.WhitePieces[p.T][len(g.WhitePieces[p.T])-1]
		} else {
			g.BlackPieces[p.T] = append(g.BlackPieces[p.T], p)
			g.Board[p.Pos.GetY()][p.Pos.GetX()] = &g.BlackPieces[p.T][len(g.BlackPieces[p.T])-1]
		}
	}

	for i := range 8 {
		for j := range 8 {
			if g.Board[i][j] == nil {
				g.Board[i][j] = &g.BlackPieces[types.EMPTY][0]
			} else {
				g.Board[i][j].Pos = types.NewPos(j, i)
			}
		}
	}

	return g
}

func (g Game) SetTurnSide(whiteTurn bool) Game {
	g.IsBlackTurn = !whiteTurn
	return g
}

func (g *Game) DebugRender() {
	for y := range 8 {
		fmt.Print(7-y, 8-y, " ")
		for x := range 8 {
			piece := g.Board[7-y][x]
			if piece.T == types.EMPTY {
				if (x+y)%2 != 0 {
					fmt.Print("  ")
				} else {
					fmt.Print("██")
				}
				continue
			}
			w := "B"
			if piece.White {
				w = "W"
			}
			fmt.Printf("%s%s", w, piece.T.String())
		}
		fmt.Print("\n")
	}
	fmt.Print("    ")
	for i := range 8 {
		fmt.Print(string(rune('a'+i)), " ")
	}
	fmt.Print("\n")
	fmt.Print("    ")
	for i := range 8 {
		fmt.Print(i, " ")
	}
	fmt.Print("\n White: ")
	for _, a := range g.WhitePieces {
		for _, p := range a {
			fmt.Print(p.String())
		}
	}
	fmt.Print("\n Black: ")
	for _, a := range g.BlackPieces {
		for _, p := range a {
			fmt.Print(p.String())
		}
	}
}

/*
TODO:
- Конвертация из математической нотации
- Проверить быстрее ли / проще ли проверки в одномерном массиве. Перейти полностью на types.Pos вместо ix iy fx fy
- Взятие на проходе
- Рокировка (просто двигать короля к финальному положению по клетке за раз, и проверять не находится ли он под шахом)
- Проверка на мат
 1. Король не может уйти с битого поля (+ взятие ближайшей фигуры соперника + проверка на шах после взятия)
 2. Нападающая фигура не может быть взята любым ходом (взять готовую логику валидности хода. ВЗЯТИЕ ФИГУРЫ НЕ ГАРАНТИРУЕТ ОТСУТСТВИЕ ШАХА НА СЛЕДУЮЩЕМ ХОДЕ)
 3. ЛИНЯЯ ШАХА НЕ МОЖЕТ БЫТЬ ПЕРЕКРЫТА ХОДОМ ЛЮБОЙ ФИГУРЫ (НЕ РАБОТАЕТ ДЛЯ КОНЯ И ЕСЛИ ФИГУРА ВПРИТЫК. ПЕРЕКРЫТИЕ ЛИНИИ ШАХА НЕ ГАРАНТИРУЕТ ОТСУТСТВИЕ ШАХА НА СЛЕДУЮЩЕМ ХОДЕ)
*/
func (g *Game) MovePiece(ix int, iy int, fx int, fy int) error {
	// Pause
	if g.IsPause {
		return errors.New("game paused: can't make moves")
	}
	// Game over
	if g.IsCheckmate {
		return errors.New("game over: can't make moves")
	}
	// Same cell move
	if ix == fx && iy == fy {
		return fmt.Errorf("same cell move (no move) [%d;%d]", ix, iy)
	}
	// Outside of the board
	if !types.NewPos(ix, iy).IsValid() {
		return fmt.Errorf("out of bounds position [%d;%d]", ix, iy)
	} else if !types.NewPos(fx, fy).IsValid() {
		return fmt.Errorf("out of bounds move [%d;%d]", fx, fy)
	}
	// Wrong color move
	if g.IsBlackTurn == g.Board[iy][ix].White {
		return fmt.Errorf("wrong color move [%d;%d]", ix, iy)
	}
	// Trying to take same color piece
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
	default:
		return fmt.Errorf("no piece to move [%d;%d]", ix, iy)
	}

	game_save := *g

	// General move execution
	g.Board[fy][fx] = g.Board[iy][ix]
	g.Board[iy][ix] = &g.BlackPieces[types.EMPTY][0]
	g.Board[fy][fx].Pos = types.NewPos(fx, fy)

	var ownKingPos, enemyKingPos = types.Pos(-1), types.Pos(-1)
	if _, ok := g.WhitePieces[types.KING]; ok {
		ownKingPos = g.WhitePieces[types.KING][0].Pos
	}
	if _, ok := g.BlackPieces[types.KING]; ok {
		enemyKingPos = g.BlackPieces[types.KING][0].Pos
	}
	if g.IsBlackTurn {
		ownKingPos, enemyKingPos = enemyKingPos, ownKingPos
	}

	if ownKingPos.IsValid() && g.CheckKingChecked(ownKingPos.GetX(), ownKingPos.GetY(), !g.IsBlackTurn) {
		// Revert illigal move
		*g = game_save
		return fmt.Errorf("invalid turn: own king still checked")
	}

	if enemyKingPos.IsValid() {
		g.IsKingChecked = g.CheckKingChecked(enemyKingPos.GetX(), enemyKingPos.GetY(), g.IsBlackTurn)
		if g.IsKingChecked {
			// TODO: check for checkmate
		}
	} else {
		g.IsKingChecked = false
	}

	// Other color turn
	g.IsBlackTurn = !g.IsBlackTurn

	return nil
}

func (g *Game) CheckKingChecked(x, y int, kingIsWhite bool) bool {
	log.Trace().Str("pos", types.NewPos(x, y).String()).Bool("kingIsWhite", kingIsWhite).Msg("CheckKingChecked")
	// up
	for i := y + 1; i < 8; i++ {
		if g.Board[i][x].White != kingIsWhite &&
			(g.Board[i][x].T == types.ROOK ||
				g.Board[i][x].T == types.QUEEN) {
			log.Trace().Bool("isKingWhite", kingIsWhite).Str("CheckedBy", g.Board[i][x].String()).Msg("king attacked")
			return true
		} else if g.Board[i][x].T != types.EMPTY {
			break
		}
	}
	// down
	for i := y - 1; i > -1; i-- {
		if g.Board[i][x].White != kingIsWhite &&
			(g.Board[i][x].T == types.ROOK ||
				g.Board[i][x].T == types.QUEEN) {
			log.Trace().Bool("isKingWhite", kingIsWhite).Str("CheckedBy", g.Board[i][x].String()).Msg("king attacked")
			return true
		} else if g.Board[i][x].T != types.EMPTY {
			break
		}
	}
	// left
	for i := x + 1; i < 8; i++ {
		if g.Board[y][i].White != kingIsWhite &&
			(g.Board[y][i].T == types.ROOK ||
				g.Board[y][i].T == types.QUEEN) {
			log.Trace().Bool("isKingWhite", kingIsWhite).Str("CheckedBy", g.Board[y][i].String()).Msg("king attacked")
			return true
		} else if g.Board[y][i].T != types.EMPTY {
			break
		}
	}
	// right
	for i := x - 1; i > -1; i-- {
		if g.Board[y][i].White != kingIsWhite &&
			(g.Board[y][i].T == types.ROOK ||
				g.Board[y][i].T == types.QUEEN) {
			log.Trace().Bool("isKingWhite", kingIsWhite).Str("CheckedBy", g.Board[y][i].String()).Msg("king attacked")
			return true
		} else if g.Board[y][i].T != types.EMPTY {
			break
		}
	}
	// up right
	for i, j := y-1, x+1; i > -1 && j < 8; {
		if g.Board[i][j].White != kingIsWhite &&
			(g.Board[i][j].T == types.BISHOP ||
				g.Board[i][j].T == types.QUEEN ||
				(g.Board[i][j].T == types.PAWN && g.Board[i][j].White)) {
			log.Trace().Bool("isKingWhite", kingIsWhite).Str("CheckedBy", g.Board[i][j].String()).Msg("king attacked")
			return true
		} else if g.Board[i][j].T != types.EMPTY {
			break
		}
		i--
		j++
	}

	// up left
	for i, j := y-1, x-1; i > -1 && j > -1; {
		if g.Board[i][j].White != kingIsWhite &&
			(g.Board[i][j].T == types.BISHOP ||
				g.Board[i][j].T == types.QUEEN ||
				g.Board[i][j].T == types.PAWN && g.Board[i][j].White) {
			log.Trace().Bool("isKingWhite", kingIsWhite).Str("CheckedBy", g.Board[i][j].String()).Msg("king attacked")
			return true
		} else if g.Board[i][j].T != types.EMPTY {
			break
		}
		i--
		j--
	}

	// down right
	for i, j := y+1, x+1; i < 8 && j < 8; {
		if g.Board[i][j].White != kingIsWhite &&
			(g.Board[i][j].T == types.BISHOP ||
				g.Board[i][j].T == types.QUEEN ||
				g.Board[i][j].T == types.PAWN && !g.Board[i][j].White) {
			log.Trace().Bool("isKingWhite", kingIsWhite).Str("CheckedBy", g.Board[i][j].String()).Msg("king attacked")
			return true
		} else if g.Board[i][j].T != types.EMPTY {
			break
		}
		i++
		j++
	}

	// down left
	for i, j := y+1, x-1; i < 8 && j > -1; {
		if g.Board[i][j].White != kingIsWhite &&
			(g.Board[i][j].T == types.BISHOP ||
				g.Board[i][j].T == types.QUEEN ||
				g.Board[i][j].T == types.PAWN && !g.Board[i][j].White) {
			log.Trace().Bool("isKingWhite", kingIsWhite).Str("CheckedBy", g.Board[i][j].String()).Msg("king attacked")
			return true
		} else if g.Board[i][j].T != types.EMPTY {
			break
		}
		i++
		j--
	}

	// KNIGHT CHECK
	possibleKnight := []types.Pos{
		types.NewPos(x+2, y+1),
		types.NewPos(x+1, y+2),
		types.NewPos(x-2, y-1),
		types.NewPos(x-1, y-2),
		types.NewPos(x-2, y+1),
		types.NewPos(x-1, y+2),
		types.NewPos(x+2, y-1),
		types.NewPos(x+1, y-2),
	}
	for _, position := range possibleKnight {
		if !position.IsValid() {
			continue
		}
		if kingIsWhite {
			for _, knight := range g.BlackPieces[types.KNIGHT] {
				if knight.Pos == position {
					fmt.Println("isWhiteKing", kingIsWhite, knight.String())
					return true
				}
			}
		} else {
			for _, knight := range g.WhitePieces[types.KNIGHT] {
				if knight.Pos == position {
					fmt.Println("isWhiteKing", kingIsWhite, knight.String())
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
