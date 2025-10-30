package board

import (
	"time"
	"ust_chess/internal/types"
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
