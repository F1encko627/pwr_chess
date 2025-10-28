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
	From      types.Position
	To        types.Position
	Validator GameValidator
}

type GameValidator func(*board.Game, types.Position, types.Position) error

var tests = []Test{
	{
		"move order",
		board.NewGame([]types.Piece{}),
		[]TestMove{
			{
				types.MustNewPos(1, 1),
				types.MustNewPos(1, 2),
				ValidateOnError(false, "white move denied", false),
			},
			{
				types.MustNewPos(2, 1),
				types.MustNewPos(2, 2),
				ValidateOnError(true, "white moved twice", false),
			},
			{
				types.MustNewPos(1, 6),
				types.MustNewPos(1, 5),
				ValidateOnError(false, "black move denied", false),
			},
			{
				types.MustNewPos(2, 6),
				types.MustNewPos(2, 5),
				ValidateOnError(true, "black moved twice", false),
			},
		},
	},
	{
		"can't move empty cell",
		board.NewGame([]types.Piece{}),
		[]TestMove{
			{
				types.MustNewPos(3, 5),
				types.MustNewPos(3, 4),
				ValidateOnError(true, "empty cell moved as piece", false),
			},
		},
	},
	{
		"can't take own piece",
		board.NewGame([]types.Piece{}),
		[]TestMove{
			{
				types.MustNewPos(0, 0),
				types.MustNewPos(1, 0),
				ValidateOnError(true, "white rook took white piece", false),
			},
			{
				types.MustNewPos(3, 1),
				types.MustNewPos(3, 2),
				ValidateOnError(false, "white move denied", false),
			},
			{
				types.MustNewPos(0, 7),
				types.MustNewPos(1, 7),
				ValidateOnError(true, "black rook took white piece", false),
			},
		},
	},
	{
		"bishop сan't jump over pieces clear",
		board.NewGame([]types.Piece{
			types.NewPiece(types.PAWN, false, types.MustNewPos(1, 2)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(5, 2)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(1, 6)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(5, 6)),

			types.NewPiece(types.BISHOP, true, types.MustNewPos(3, 4)),
		}),
		[]TestMove{
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(0, 1),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(6, 1),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(0, 7),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(6, 7),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
		},
	},
	{
		"bishop сan't jump over pieces obstructed",
		board.NewGame([]types.Piece{
			types.NewPiece(types.PAWN, false, types.MustNewPos(5, 2)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(1, 6)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(1, 2)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(5, 6)),

			types.NewPiece(types.PAWN, false, types.MustNewPos(3, 5)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(3, 3)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(4, 4)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(2, 4)),

			types.NewPiece(types.BISHOP, true, types.MustNewPos(3, 4)),
		}),
		[]TestMove{
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(0, 1),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(6, 1),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(0, 7),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(6, 7),
				ValidateOnError(true, "bishop jumped over piece", false),
			},
		},
	},
	{
		"bishop moves straight clear",
		board.NewGame([]types.Piece{
			types.NewPiece(types.PAWN, false, types.MustNewPos(1, 2)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(5, 2)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(1, 6)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(5, 6)),

			types.NewPiece(types.BISHOP, true, types.MustNewPos(3, 4)),
		}),
		[]TestMove{
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(1, 2),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(5, 2),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(1, 6),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(5, 6),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(2, 3),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(4, 3),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(2, 5),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(4, 5),
				ValidateOnError(false, "bishop not moved", true),
			},
		},
	},
	{
		"bishop moves straight obstructed",
		board.NewGame([]types.Piece{
			types.NewPiece(types.PAWN, false, types.MustNewPos(5, 2)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(1, 2)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(1, 6)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(5, 6)),

			types.NewPiece(types.PAWN, false, types.MustNewPos(3, 5)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(3, 3)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(4, 4)),
			types.NewPiece(types.PAWN, false, types.MustNewPos(2, 4)),

			types.NewPiece(types.BISHOP, true, types.MustNewPos(3, 4)),
		}),
		[]TestMove{
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(1, 2),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(5, 2),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(1, 6),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(5, 6),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(2, 3),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(4, 3),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(2, 5),
				ValidateOnError(false, "bishop not moved", true),
			},
			{
				types.MustNewPos(3, 4),
				types.MustNewPos(4, 5),
				ValidateOnError(false, "bishop not moved", true),
			},
		},
	},
	{
		"black king gets checked by everyone",
		board.NewGame([]types.Piece{
			types.NewPiece(types.KING, false, types.MustNewPos(0, 7)),

			types.NewPiece(types.QUEEN, true, types.MustNewPos(2, 1)),
			types.NewPiece(types.ROOK, true, types.MustNewPos(3, 1)),
			types.NewPiece(types.BISHOP, true, types.MustNewPos(4, 1)),
			types.NewPiece(types.PAWN, true, types.MustNewPos(1, 5)),
			types.NewPiece(types.KNIGHT, true, types.MustNewPos(4, 5)),
		}),
		[]TestMove{
			{
				types.MustNewPos(2, 1),
				types.MustNewPos(2, 7),
				ValidateKingChecked(true, "no queen horizontal check", true),
			},
			{
				types.MustNewPos(2, 1),
				types.MustNewPos(4, 3),
				ValidateKingChecked(true, "no queen diagonal check", true),
			},
			{
				types.MustNewPos(3, 1),
				types.MustNewPos(3, 7),
				ValidateKingChecked(true, "no rook check", true),
			},
			{
				types.MustNewPos(4, 1),
				types.MustNewPos(5, 2),
				ValidateKingChecked(true, "no bishop check", true),
			},
			{
				types.MustNewPos(1, 5),
				types.MustNewPos(1, 6),
				ValidateKingChecked(true, "no pawn check", true),
			},
			{
				types.MustNewPos(4, 5),
				types.MustNewPos(2, 6),
				ValidateKingChecked(true, "no knight check", true),
			},
		},
	},
	{
		"queen can't jump over pieces",
		board.NewGame([]types.Piece{
			types.NewPiece(types.PAWN, true, types.MustNewPos(3, 1)),
			types.NewPiece(types.PAWN, true, types.MustNewPos(1, 2)),
			types.NewPiece(types.PAWN, true, types.MustNewPos(5, 2)),
			types.NewPiece(types.PAWN, true, types.MustNewPos(5, 4)),
			types.NewPiece(types.PAWN, true, types.MustNewPos(4, 5)),
			types.NewPiece(types.PAWN, true, types.MustNewPos(3, 5)),
			types.NewPiece(types.PAWN, true, types.MustNewPos(0, 4)),
			types.NewPiece(types.PAWN, true, types.MustNewPos(0, 7)),

			types.NewPiece(types.QUEEN, false, types.MustNewPos(3, 4)),
		}),
		[]TestMove{},
	},
}

func TestEverything(t *testing.T) {
	for _, test := range tests {
		t.Log(test.Title)
		for _, move := range test.Moves {
			if err := move.Validator(&test.InitialState, move.From, move.To); err != nil {
				test.InitialState.DebugRender()
				t.Fatal(err)
			}
		}
	}
}

func ValidateOnError(expectError bool, errorText string, resetAfter bool) GameValidator {
	return func(Game *board.Game, From types.Position, To types.Position) error {
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
	return func(Game *board.Game, From types.Position, To types.Position) error {
		init := *Game
		err := Game.MovePiece(From.GetX(), From.GetY(), To.GetX(), To.GetY())
		if err != nil {
			Game.DebugRender()
			return fmt.Errorf("move error: %s", err)
		}
		if Game.IsKingChecked != expectCheck {
			Game.DebugRender()
			return errors.New(errorText)
		}
		if resetAfter {
			*Game = init
		}
		return nil
	}
}
