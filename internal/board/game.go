package board

import (
	"errors"
	"fmt"
	"time"
	"ust_chess/internal/types"
)

var (
	ErrDiscoveredCheck  = errors.New("discovered check")
	ErrGamePaused       = errors.New("game paused")
	ErrGameEnded        = errors.New("game ended")
	ErrKingCheckedStill = errors.New("own king checked still")
	ErrNoPieceToMove    = errors.New("no piece to move")
	ErrIlligalMove      = errors.New("illigal move")
	ErrOpponentsTurn    = errors.New("opponent's turn")
)

type Game struct {
	Board             types.Board
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
	History           error
	Error             string
}

var classic = []types.Piece{
	types.MustNewPiece(types.PAWN, true, types.MustNewPos(0, 1)),
	types.MustNewPiece(types.PAWN, true, types.MustNewPos(1, 1)),
	types.MustNewPiece(types.PAWN, true, types.MustNewPos(2, 1)),
	types.MustNewPiece(types.PAWN, true, types.MustNewPos(3, 1)),
	types.MustNewPiece(types.PAWN, true, types.MustNewPos(4, 1)),
	types.MustNewPiece(types.PAWN, true, types.MustNewPos(5, 1)),
	types.MustNewPiece(types.PAWN, true, types.MustNewPos(6, 1)),
	types.MustNewPiece(types.PAWN, true, types.MustNewPos(7, 1)),

	types.MustNewPiece(types.ROOK, true, types.MustNewPos(0, 0)),
	types.MustNewPiece(types.KNIGHT, true, types.MustNewPos(1, 0)),
	types.MustNewPiece(types.BISHOP, true, types.MustNewPos(2, 0)),
	types.MustNewPiece(types.KING, true, types.MustNewPos(3, 0)),
	types.MustNewPiece(types.QUEEN, true, types.MustNewPos(4, 0)),
	types.MustNewPiece(types.BISHOP, true, types.MustNewPos(5, 0)),
	types.MustNewPiece(types.KNIGHT, true, types.MustNewPos(6, 0)),
	types.MustNewPiece(types.ROOK, true, types.MustNewPos(7, 0)),

	types.MustNewPiece(types.PAWN, false, types.MustNewPos(0, 6)),
	types.MustNewPiece(types.PAWN, false, types.MustNewPos(1, 6)),
	types.MustNewPiece(types.PAWN, false, types.MustNewPos(2, 6)),
	types.MustNewPiece(types.PAWN, false, types.MustNewPos(3, 6)),
	types.MustNewPiece(types.PAWN, false, types.MustNewPos(4, 6)),
	types.MustNewPiece(types.PAWN, false, types.MustNewPos(5, 6)),
	types.MustNewPiece(types.PAWN, false, types.MustNewPos(6, 6)),
	types.MustNewPiece(types.PAWN, false, types.MustNewPos(7, 6)),

	types.MustNewPiece(types.ROOK, false, types.MustNewPos(0, 7)),
	types.MustNewPiece(types.KNIGHT, false, types.MustNewPos(1, 7)),
	types.MustNewPiece(types.BISHOP, false, types.MustNewPos(2, 7)),
	types.MustNewPiece(types.KING, false, types.MustNewPos(3, 7)),
	types.MustNewPiece(types.QUEEN, false, types.MustNewPos(4, 7)),
	types.MustNewPiece(types.BISHOP, false, types.MustNewPos(5, 7)),
	types.MustNewPiece(types.KNIGHT, false, types.MustNewPos(6, 7)),
	types.MustNewPiece(types.ROOK, false, types.MustNewPos(7, 7)),
}

/*
 1. Проверить что ход возможен (конец игры, пауза, ...)
    1.1 Вернуть ошибку.
 2. Проверить что ходит фигура нужного цвета.
    2.1 Вернуть ошибку.
 3. Проверить что фигура так ходит.
    3.1 Вернуть ошибку.
 4. Сохранить состояние игры.
 5. Проверить на шах своему королю.
    5.1 Откатить состояние игры, вернуть ошибку.
 6. Засчитать очки.
 7. Проверить на шах другому королю.
 8. Сменить ходящую сторону.
    8.1. Шаха нет - return.
 9. Проверить на мат:
    а) Убрать короля с доски и dummy фигурой проверить позиции вокруг на шах (возможность отойти или взять фигуру, если она рядом).
    б) вернуть короля на доску.
    в) Проверить что атакующую фигуру может взять фигура кроме короля (ОПАСНОСТЬ СКРЫТОГО ШАХА!!!)
    г) Проверить что линию атаки фигуры можно перекрыть фигурой кроме короля (если это не конь)
*/
func (g *Game) MakeMove(move types.Move) error {
	if g.IsPause {
		return ErrGamePaused
	}
	if g.IsCheckmate {
		return ErrGameEnded
	}
	piece := g.Board.GetCell(move.GetInitial()).GetPiece()
	if piece == nil {
		return errors.Join(ErrNoPieceToMove, fmt.Errorf("%s", move))
	}
	if piece.IsWhite() == g.IsBlackTurn {
		return ErrOpponentsTurn
	}
	if err := piece.MakeMove(move, &g.Board); err != nil {
		switch {
		case errors.Is(err, types.ErrEnPassantMove):
			g.EnPassantPawn = g.Board.GetCell(move.GetInitial()).GetPiece()
		case errors.Is(err, types.ErrEnPassantTake):
			pos := types.MustNewPos(0, 0)
			piece := g.Board.GetCell(pos).GetPiece()
			piece.Take()
			piece = nil
		case errors.Is(err, types.ErrCastleMove):
			if err := checkForValidCastle(move, &g.Board); err != nil {
				return errors.Join(ErrIlligalMove, err)
			}
		default:
			return errors.Join(ErrIlligalMove, err)
		}
	}

	g.Board.MakeMove(move)

	g.EnPassantPawn = nil
	g.IsBlackTurn = !g.IsBlackTurn

	return nil
}

func checkForValidCastle(move types.Move, board *types.Board) error {
	return nil
}

func checkForValidEnpassantTake(move types.Move, board *types.Board) error {
	return nil
}

type GameOutDto struct {
	IsBlackTurn   bool
	IsKingChecked bool
	IsCheckmate   bool
	Board         [][]PieceOutDto
	Error         string
}

type PieceOutDto struct {
	T     string
	White bool
	Pos   string
}

func (g *Game) GetForRender() GameOutDto {
	pieces := make([][]PieceOutDto, 8)
	for y := range 8 {
		pieces[y] = make([]PieceOutDto, 8)
		for x := range 8 {
			raw_piece := g.Board.GetCell(types.MustNewPos(x, y)).GetPiece()
			if raw_piece != nil {
				pieces[y][x] = PieceOutDto{
					T:     raw_piece.GetType().String(),
					White: raw_piece.IsWhite(),
					Pos:   raw_piece.GetPosition().String(),
				}
				continue
			}
			pieces[y][x] = PieceOutDto{
				T:     " ",
				White: false,
				Pos:   types.MustNewPos(x, y).String(),
			}
		}
	}
	return GameOutDto{
		IsBlackTurn:   g.IsBlackTurn,
		IsKingChecked: g.IsKingChecked,
		IsCheckmate:   g.IsCheckmate,
		Board:         pieces,
		Error:         g.Error,
	}
}

/*
Проверять состояние игры на шах и мат сразу при создании.
*/
func NewGame(pieces []types.Piece) Game {
	board, err := types.GetBoard(classic)
	if err != nil {
		panic(err)
	}
	return Game{Board: board}
}
