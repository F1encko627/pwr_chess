package board

import (
	"errors"
	"fmt"
	"time"
	"ust_chess/internal/types"

	"github.com/rs/zerolog/log"
)

type Game struct {
	Board             [8][8]*types.Piece
	WhitePieces       map[types.Figure][]types.Piece
	BlackPieces       map[types.Figure][]types.Piece // plus empty tile
	LastMovedPiece    *types.Piece
	TurnNum           int
	IsBlackTurn       bool
	IsWhiteCantCastle bool
	IsBlackCantCastle bool
	IsKingChecked     bool
	IsCheckmate       bool
	IsPause           bool
	EnPassantPawn     *types.Piece
	LastMoveTime      time.Time
}

var classic = []types.Piece{
	types.NewPiece(types.PAWN, true, types.MustNewPos(0, 1)),
	types.NewPiece(types.PAWN, true, types.MustNewPos(1, 1)),
	types.NewPiece(types.PAWN, true, types.MustNewPos(2, 1)),
	types.NewPiece(types.PAWN, true, types.MustNewPos(3, 1)),
	types.NewPiece(types.PAWN, true, types.MustNewPos(4, 1)),
	types.NewPiece(types.PAWN, true, types.MustNewPos(5, 1)),
	types.NewPiece(types.PAWN, true, types.MustNewPos(6, 1)),
	types.NewPiece(types.PAWN, true, types.MustNewPos(7, 1)),

	types.NewPiece(types.ROOK, true, types.MustNewPos(0, 0)),
	types.NewPiece(types.KNIGHT, true, types.MustNewPos(1, 0)),
	types.NewPiece(types.BISHOP, true, types.MustNewPos(2, 0)),
	types.NewPiece(types.KING, true, types.MustNewPos(3, 0)),
	types.NewPiece(types.QUEEN, true, types.MustNewPos(4, 0)),
	types.NewPiece(types.BISHOP, true, types.MustNewPos(5, 0)),
	types.NewPiece(types.KNIGHT, true, types.MustNewPos(6, 0)),
	types.NewPiece(types.ROOK, true, types.MustNewPos(7, 0)),

	types.NewPiece(types.PAWN, false, types.MustNewPos(0, 6)),
	types.NewPiece(types.PAWN, false, types.MustNewPos(1, 6)),
	types.NewPiece(types.PAWN, false, types.MustNewPos(2, 6)),
	types.NewPiece(types.PAWN, false, types.MustNewPos(3, 6)),
	types.NewPiece(types.PAWN, false, types.MustNewPos(4, 6)),
	types.NewPiece(types.PAWN, false, types.MustNewPos(5, 6)),
	types.NewPiece(types.PAWN, false, types.MustNewPos(6, 6)),
	types.NewPiece(types.PAWN, false, types.MustNewPos(7, 6)),

	types.NewPiece(types.ROOK, false, types.MustNewPos(0, 7)),
	types.NewPiece(types.KNIGHT, false, types.MustNewPos(1, 7)),
	types.NewPiece(types.BISHOP, false, types.MustNewPos(2, 7)),
	types.NewPiece(types.KING, false, types.MustNewPos(3, 7)),
	types.NewPiece(types.QUEEN, false, types.MustNewPos(4, 7)),
	types.NewPiece(types.BISHOP, false, types.MustNewPos(5, 7)),
	types.NewPiece(types.KNIGHT, false, types.MustNewPos(6, 7)),
	types.NewPiece(types.ROOK, false, types.MustNewPos(7, 7)),
}

func NewGame(start []types.Piece) Game {
	if len(start) == 0 {
		start = classic
	}

	g := Game{}

	g.WhitePieces = make(map[types.Figure][]types.Piece)
	g.BlackPieces = make(map[types.Figure][]types.Piece)

	empty := types.NewPiece(types.EMPTY, false, types.Pos(-1))

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
				g.Board[i][j].Pos = types.MustNewPos(j, i)
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

func (g *Game) BasicMoveChecks(ix, iy, fx, fy int) error {
	// Same cell move
	if ix == fx && iy == fy {
		return fmt.Errorf("same cell move (no move) [%d;%d]", ix, iy)
	}
	// Outside of the board
	if !types.MustNewPos(ix, iy).IsValid() {
		return fmt.Errorf("out of bounds position [%d;%d]", ix, iy)
	} else if !types.MustNewPos(fx, fy).IsValid() {
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
	return nil
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
 3. ЛИНЯЯ ШАХА НЕ МОЖЕТ БЫТЬ ПЕРЕКРЫТА ХОДОМ ЛЮБОЙ ФИГУРЫ (НЕ РАБОТАЕТ ДЛЯ КОНЯ И ЕСЛИ ФИГУРА ВПРИТЫК. ПЕРЕКРЫТИЕ ЛИНИИ ШАХА НЕ ГАРАНТИРУЕТ ОТСУТСТВИЕ ШАХА НА СЛЕДУЮЩЕМ ХОДЕ (СКРЫТЫЙ ШАХ, БЛОКИРУЕМЫЙ ФИГУРОЙ, КОТОРАЯ ПЫТАЕТСЯ ЗАБЛОКИРОВАТЬ ШАХ))
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
	if err := g.BasicMoveChecks(ix, iy, fx, fy); err != nil {
		return err
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
	g.Board[fy][fx].Pos = types.Pos(-1)
	g.Board[iy][ix].Pos = types.MustNewPos(fx, fy)
	g.Board[fy][fx] = g.Board[iy][ix]
	g.Board[iy][ix] = &g.BlackPieces[types.EMPTY][0]

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

	if ownKingPos.IsValid() && g.CheckKingChecked(ownKingPos.GetX(), ownKingPos.GetY(), !g.IsBlackTurn).IsValid() {
		// Revert illigal move
		*g = game_save
		return fmt.Errorf("invalid turn: own king still checked")
	}

	if enemyKingPos.IsValid() {
		CheckPos := g.CheckKingChecked(enemyKingPos.GetX(), enemyKingPos.GetY(), g.IsBlackTurn)
		g.IsKingChecked = CheckPos.IsValid()
		if g.IsKingChecked {
			g.IsCheckmate = g.CheckForCheckMate(g.Board[CheckPos.GetY()][CheckPos.GetX()])
			// TODO: check for checkmate
		}
	} else {
		g.IsKingChecked = false
	}

	// Other color turn
	g.IsBlackTurn = !g.IsBlackTurn
	g.TurnNum++

	return nil
}

// 1. Король не может отойти в сторону.
// 2. Король не может взять атакующую фигуру.
// 3.а. Ни одна другая фигура не может взять атакующую (шах от пешки или коня).
// 3.б. Ни одна другая фигура не может своим ходом перекрыть линию шаха.
// Проверить на возможность отойти, если нельзя - дальше
// Собрать список атакующих фигур
// Если фигур больше двух - мат.
// Если фигуру нельзя забрать - дальше.
// Если фигура конь - мат.
// Проверить что линию атаки нельзя перекрыть
func (g *Game) CheckForCheckMate(lastCheckedPiece *types.Piece) bool {

	switch lastCheckedPiece.T {
	case types.PAWN:
		fallthrough
	case types.KNIGHT:
		// TODO: просто проверить что нельзя взять фигуру - это единственный способ.
	default:
		// TODO: проверить что нельзя взять фигуру, проверить что нельзя закрыться.
	}
	return false
}

func (g *Game) CheckKingChecked(x, y int, kingIsWhite bool) types.Pos {
	log.Trace().Str("pos", types.MustNewPos(x, y).String()).Bool("kingIsWhite", kingIsWhite).Msg("CheckKingChecked")
	// up
	for i := y + 1; i < 8; i++ {
		if g.Board[i][x].White != kingIsWhite &&
			(g.Board[i][x].T == types.ROOK ||
				g.Board[i][x].T == types.QUEEN) {
			log.Trace().Bool("isKingWhite", kingIsWhite).Str("CheckedBy", g.Board[i][x].String()).Msg("king attacked")
			return types.MustNewPos(x, i)
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
			return types.MustNewPos(x, i)
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
			return types.MustNewPos(i, y)
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
			return types.MustNewPos(i, y)
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
			return types.MustNewPos(j, i)
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
			return types.MustNewPos(j, i)
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
			return types.MustNewPos(j, i)
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
			return types.MustNewPos(j, i)
		} else if g.Board[i][j].T != types.EMPTY {
			break
		}
		i++
		j--
	}

	// KNIGHT CHECK
	possibleKnight := []types.Pos{
		types.MustNewPos(x+2, y+1),
		types.MustNewPos(x+1, y+2),
		types.MustNewPos(x-2, y-1),
		types.MustNewPos(x-1, y-2),
		types.MustNewPos(x-2, y+1),
		types.MustNewPos(x-1, y+2),
		types.MustNewPos(x+2, y-1),
		types.MustNewPos(x+1, y-2),
	}
	for _, position := range possibleKnight {
		if !position.IsValid() {
			continue
		}
		if kingIsWhite {
			for _, knight := range g.BlackPieces[types.KNIGHT] {
				if knight.Pos == position {
					log.Trace().Bool("isKingWhite", kingIsWhite).Str("CheckedBy", knight.String()).Msg("king attacked")
					return knight.Pos
				}
			}
		} else {
			for _, knight := range g.WhitePieces[types.KNIGHT] {
				if knight.Pos == position {
					log.Trace().Bool("isKingWhite", kingIsWhite).Str("CheckedBy", knight.String()).Msg("king attacked")
					return knight.Pos
				}
			}
		}
	}

	return types.Pos(-1)
}
