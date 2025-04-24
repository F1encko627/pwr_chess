package board_test

import (
	"testing"
	"ust_chess/internal/board"
	"ust_chess/internal/types"
)

func TestWalkOrder(t *testing.T) {
	game := board.NewGame([]types.Piece{})
	err := game.MovePiece(1, 1, 1, 2)
	if err != nil {
		t.Error("white move denied")
	}

	err = game.MovePiece(2, 1, 2, 2)
	if err == nil {
		t.Error("white moved twice")
	}

	err = game.MovePiece(6, 1, 5, 1)
	if err != nil {
		t.Error("black move denied")
	}

	err = game.MovePiece(2, 6, 2, 5)
	if err == nil {
		t.Error("black moved twice")
	}
}

func TestCantMoveEmpty(t *testing.T) {
	game := board.NewGame([]types.Piece{})
	err := game.MovePiece(2, 2, 3, 2)
	if err == nil {
		t.Error("empty space moved as piece")
	}
}

func TestCantTakeOwnPiece(t *testing.T) {
	game := board.NewGame([]types.Piece{})
	err := game.MovePiece(4, 0, 4, 1)
	if err == nil {
		t.Error("white took it's own piece")
	}

	_ = game.MovePiece(1, 1, 2, 1)
	err = game.MovePiece(4, 7, 4, 6)
	if err == nil {
		t.Error("black took it's own piece")
	}
}

func TestRookWalksStraigt(t *testing.T) {
	game := board.NewGame([]types.Piece{
		types.GP(types.PAWN, true, types.NewPos(3, 1)),
		types.GP(types.PAWN, true, types.NewPos(1, 2)),
		types.GP(types.PAWN, true, types.NewPos(5, 2)),
		types.GP(types.PAWN, true, types.NewPos(5, 4)),
		types.GP(types.PAWN, true, types.NewPos(4, 5)),
		types.GP(types.PAWN, true, types.NewPos(3, 5)),
		types.GP(types.PAWN, true, types.NewPos(0, 4)),
		types.GP(types.PAWN, true, types.NewPos(0, 7)),

		types.GP(types.ROOK, false, types.NewPos(3, 4)),
	})
	err := game.MovePiece(4, 0, 4, 1)
	if err == nil {
		t.Error("white took it's own piece")
	}
}

func TestQueenCantJumpOverPiece(t *testing.T) {
	game := board.NewGame([]types.Piece{
		types.GP(types.PAWN, false, types.NewPos(3, 1)),
		types.GP(types.PAWN, false, types.NewPos(1, 2)),
		types.GP(types.PAWN, false, types.NewPos(5, 2)),
		types.GP(types.PAWN, false, types.NewPos(5, 4)),
		types.GP(types.PAWN, false, types.NewPos(4, 5)),
		types.GP(types.PAWN, false, types.NewPos(3, 5)),
		types.GP(types.PAWN, false, types.NewPos(0, 4)),
		types.GP(types.PAWN, false, types.NewPos(0, 7)),

		types.GP(types.QUEEN, true, types.NewPos(3, 4)),
	})
	game.State = types.WHITE_TURN
	err := game.MovePiece(3, 4, 0, 1)
	if err != nil {
		t.Error("queen jumped over piece")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(3, 4, 3, 0)
	if err != nil {
		t.Error("queen jumped over piece")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(3, 4, 7, 1)
	if err != nil {
		t.Error("queen jumped over piece")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(3, 4, 7, 4)
	if err != nil {
		t.Error("queen jumped over piece")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(3, 4, 5, 6)
	if err != nil {
		t.Error("queen jumped over piece")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(3, 4, 3, 6)
	if err != nil {
		t.Error("queen jumped over piece")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(3, 4, 1, 6)
	if err != nil {
		t.Error("queen not moved")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(1, 6, 1, 3)
	if err != nil {
		t.Error("queen not moved")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(1, 3, 2, 4)
	if err != nil {
		t.Error("queen not moved")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(2, 4, 4, 4)
	if err != nil {
		t.Error("queen not moved")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(4, 4, 4, 5)
	if err != nil {
		t.Error("queen not moved")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(4, 5, 1, 2)
	if err != nil {
		t.Error("queen not moved")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(1, 2, 5, 2)
	if err != nil {
		t.Error("queen not moved")
	}
	game.State = types.WHITE_TURN
	err = game.MovePiece(5, 2, 5, 4)
	if err != nil {
		t.Error("queen not moved")
	}
}

func TestBishopWalksStraigt(t *testing.T) {
	t.Error("not implemented")
}

func TestQueenWalksStraigt(t *testing.T) {
	t.Error("not implemented")
}
