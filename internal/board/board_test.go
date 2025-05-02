package board_test

import (
	"errors"
	"fmt"
	"testing"
	"ust_chess/internal/board"
	"ust_chess/internal/types"
)

type Test struct {
	Title        string
	InitialState board.Game
	Moves        []TestMove
}
type TestMove struct {
	From      types.Pos
	To        types.Pos
	Validator GameValidator
}

type GameValidator func(*board.Game, types.Pos, types.Pos) error

var tests = []Test{
	{
		"move order",
		board.NewGame([]types.Piece{}),
		[]TestMove{
			{
				types.NewPos(1, 1),
				types.NewPos(1, 2),
				ValidateOnError(false, "white move denied", false),
			},
			{
				types.NewPos(2, 1),
				types.NewPos(2, 2),
				ValidateOnError(true, "white moved twice", false),
			},
			{
				types.NewPos(1, 6),
				types.NewPos(1, 5),
				ValidateOnError(false, "black move denied", false),
			},
			{
				types.NewPos(2, 6),
				types.NewPos(2, 5),
				ValidateOnError(true, "black moved twice", false),
			},
		},
	},
	{
		"can't move empty cell",
		board.NewGame([]types.Piece{}),
		[]TestMove{
			{
				types.NewPos(3, 5),
				types.NewPos(3, 4),
				ValidateOnError(true, "empty cell moved as piece", false),
			},
		},
	},
	{
		"can't take own piece",
		board.NewGame([]types.Piece{}),
		[]TestMove{
			{
				types.NewPos(0, 0),
				types.NewPos(1, 0),
				ValidateOnError(true, "white rook took white piece", false),
			},
			{
				types.NewPos(3, 1),
				types.NewPos(3, 2),
				ValidateOnError(false, "white move denied", false),
			},
			{
				types.NewPos(0, 7),
				types.NewPos(1, 7),
				ValidateOnError(true, "black rook took white piece", false),
			},
		},
	},
	{
		"bishop сan't jump over pieces clear",
		board.NewGame([]types.Piece{
			types.GP(types.PAWN, false, types.NewPos(1, 2)),
			types.GP(types.PAWN, false, types.NewPos(5, 2)),
			types.GP(types.PAWN, false, types.NewPos(1, 6)),
			types.GP(types.PAWN, false, types.NewPos(5, 6)),

			types.GP(types.BISHOP, true, types.NewPos(3, 4)),
		}),
		[]TestMove{
			{
				types.NewPos(3, 4),
				types.NewPos(0, 1),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(6, 1),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(0, 7),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(6, 7),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
		},
	},
	{
		"bishop сan't jump over pieces obstructed",
		board.NewGame([]types.Piece{
			types.GP(types.PAWN, false, types.NewPos(5, 2)),
			types.GP(types.PAWN, false, types.NewPos(1, 6)),
			types.GP(types.PAWN, false, types.NewPos(1, 2)),
			types.GP(types.PAWN, false, types.NewPos(5, 6)),

			types.GP(types.PAWN, false, types.NewPos(3, 5)),
			types.GP(types.PAWN, false, types.NewPos(3, 3)),
			types.GP(types.PAWN, false, types.NewPos(4, 4)),
			types.GP(types.PAWN, false, types.NewPos(2, 4)),

			types.GP(types.BISHOP, true, types.NewPos(3, 4)),
		}),
		[]TestMove{
			{
				types.NewPos(3, 4),
				types.NewPos(0, 1),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(6, 1),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(0, 7),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(6, 7),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
		},
	},
	{
		"bishop moves straight clear",
		board.NewGame([]types.Piece{
			types.GP(types.PAWN, false, types.NewPos(1, 2)),
			types.GP(types.PAWN, false, types.NewPos(5, 2)),
			types.GP(types.PAWN, false, types.NewPos(1, 6)),
			types.GP(types.PAWN, false, types.NewPos(5, 6)),

			types.GP(types.BISHOP, true, types.NewPos(3, 4)),
		}),
		[]TestMove{
			{
				types.NewPos(3, 4),
				types.NewPos(1, 2),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(5, 2),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(1, 6),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(5, 6),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(2, 3),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(4, 3),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(2, 5),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(4, 5),
				ValidateOnError(false, "bishop not moved", true),
			},
		},
	},
	{
		"bishop moves straight obstructed",
		board.NewGame([]types.Piece{
			types.GP(types.PAWN, false, types.NewPos(5, 2)),
			types.GP(types.PAWN, false, types.NewPos(1, 2)),
			types.GP(types.PAWN, false, types.NewPos(1, 6)),
			types.GP(types.PAWN, false, types.NewPos(5, 6)),

			types.GP(types.PAWN, false, types.NewPos(3, 5)),
			types.GP(types.PAWN, false, types.NewPos(3, 3)),
			types.GP(types.PAWN, false, types.NewPos(4, 4)),
			types.GP(types.PAWN, false, types.NewPos(2, 4)),

			types.GP(types.BISHOP, true, types.NewPos(3, 4)),
		}),
		[]TestMove{
			{
				types.NewPos(3, 4),
				types.NewPos(1, 2),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(5, 2),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(1, 6),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(5, 6),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(2, 3),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(4, 3),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(2, 5),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.NewPos(3, 4),
				types.NewPos(4, 5),
				ValidateOnError(false, "bishop not moved", true),
			},
		},
	},
	{
		"black king gets checked by everyone",
		board.NewGame([]types.Piece{
			types.GP(types.KING, false, types.NewPos(0, 0)),

			types.GP(types.QUEEN, true, types.NewPos(2, 6)),
			types.GP(types.ROOK, true, types.NewPos(3, 6)),
			types.GP(types.BISHOP, true, types.NewPos(4, 6)),
			types.GP(types.PAWN, true, types.NewPos(1, 2)),
			types.GP(types.KNIGHT, true, types.NewPos(5, 4)),
		}),
		[]TestMove{
			{
				types.NewPos(2, 6),
				types.NewPos(2, 0),
				ValidateKingChecked(false, "no queen horizontal check", true),
			},
			{
				types.NewPos(2, 6),
				types.NewPos(4, 4),
				ValidateKingChecked(false, "no queen diagonal check", true),
			},
			{
				types.NewPos(3, 6),
				types.NewPos(3, 0),
				ValidateKingChecked(false, "no rook check", true),
			},
			{
				types.NewPos(4, 6),
				types.NewPos(5, 5),
				ValidateKingChecked(false, "no bishop check", true),
			},
			{
				types.NewPos(1, 2),
				types.NewPos(1, 1),
				ValidateKingChecked(false, "no pawn check", true),
			},
			{
				types.NewPos(5, 4),
				types.NewPos(2, 1),
				ValidateKingChecked(false, "no knight check", true),
			},
		},
	},
}

func TestEverything(t *testing.T) {
	for _, test := range tests {
		t.Log(test.Title)
		for _, move := range test.Moves {
			if err := move.Validator(&test.InitialState, move.From, move.To); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func ValidateOnError(expectError bool, errorText string, resetAfter bool) GameValidator {
	return func(Game *board.Game, From types.Pos, To types.Pos) error {
		init := *Game
		err := Game.MovePiece(From.GetX(), From.GetY(), To.GetX(), To.GetY())
		if (err != nil) != expectError {
			Game.DebugRender()
			if err != nil {
				return fmt.Errorf("%s: %s", errorText, err)
			} else {
				return errors.New(errorText)
			}
		}
		if resetAfter {
			*Game = init
		}
		return nil
	}
}

func ValidateKingChecked(expectCheck bool, errorText string, resetAfter bool) GameValidator {
	return func(Game *board.Game, From types.Pos, To types.Pos) error {
		init := *Game
		err := Game.MovePiece(From.GetX(), From.GetY(), To.GetX(), To.GetY())
		if err != nil {
			Game.DebugRender()
			return fmt.Errorf("move error: %s", err)
		}
		if Game.KingChecked != expectCheck {
			Game.DebugRender()
			return errors.New(errorText)
		}
		if resetAfter {
			*Game = init
		}
		return nil
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
		game.DebugRender()
		t.Fatal("white took it's own piece")
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
	game.BlackTurn = false
	game.DebugRender()
	err := game.MovePiece(3, 4, 0, 1)
	if err == nil {
		game.DebugRender()
		t.Fatal("queen jumped over piece")
	}
	game.BlackTurn = false
	err = game.MovePiece(3, 4, 3, 0)
	if err == nil {
		t.Fatal("queen jumped over piece")
	}
	game.BlackTurn = false
	err = game.MovePiece(3, 4, 7, 1)
	if err == nil {
		t.Fatal("queen jumped over piece")
	}
	game.BlackTurn = false
	err = game.MovePiece(3, 4, 7, 4)
	if err == nil {
		t.Fatal("queen jumped over piece")
	}
	game.BlackTurn = false
	err = game.MovePiece(3, 4, 5, 6)
	if err == nil {
		t.Fatal("queen jumped over piece")
	}
	game.BlackTurn = false
	err = game.MovePiece(3, 4, 3, 6)
	if err == nil {
		t.Fatal("queen jumped over piece")
	}
	game.BlackTurn = false
	err = game.MovePiece(3, 4, 1, 6)
	if err != nil {
		t.Fatal("queen not moved")
	}
	game.BlackTurn = false
	err = game.MovePiece(1, 6, 1, 3)
	if err != nil {
		t.Fatal("queen not moved")
	}
	game.BlackTurn = false
	err = game.MovePiece(1, 3, 2, 4)
	if err != nil {
		t.Fatal("queen not moved")
	}
	game.BlackTurn = false
	err = game.MovePiece(2, 4, 4, 4)
	if err != nil {
		t.Fatal("queen not moved")
	}
	game.BlackTurn = false
	err = game.MovePiece(4, 4, 4, 5)
	if err != nil {
		t.Fatal("queen not moved")
	}
	game.BlackTurn = false
	err = game.MovePiece(4, 5, 1, 2)
	if err != nil {
		t.Fatal("queen not moved")
	}
	game.BlackTurn = false
	err = game.MovePiece(1, 2, 5, 2)
	if err != nil {
		t.Fatal("queen not moved")
	}
	game.BlackTurn = false
	err = game.MovePiece(5, 2, 5, 4)
	if err != nil {
		t.Fatal("queen not moved")
	}
}
