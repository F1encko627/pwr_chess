package types

import (
	"errors"
	"fmt"
)

type IPiece interface {
	CanTransform() []Figure
	ClearTargetedCells(board *Board)
	GetId() int
	GetPosition() Position
	GetScore() int
	GetTargetedCells() []Position
	GetType() Figure
	IsTaken() bool
	IsWhite() bool
	MarkTargetedCells(board *Board)
	String() string
}

type Figure rune

const (
	EMPTY Figure = ' '
	KING         = '♚' + iota
	QUEEN
	ROOK
	BISHOP
	KNIGHT
	PAWN
	CONST_FIGURE_LIST_LENGTH
)

func (f Figure) Name() string {
	switch f {
	case KING:
		return "king"
	case QUEEN:
		return "queen"
	case KNIGHT:
		return "knight"
	case ROOK:
		return "rook"
	case BISHOP:
		return "bishop"
	case PAWN:
		return "pawn"
	default:
		return "???"
	}
}

func (t Figure) String() string {
	return string(t)
}

var (
	ErrFigureNotSupported = errors.New("figure not in provided constants")
)

type Piece struct {
	id            int
	figure        Figure
	position      Position
	score         int
	isTaken       bool
	isWhite       bool
	targetedCells []Position
}

var sequence = 0

func NewPiece(pieceType Figure, isWhite bool, position Position) (Piece, error) {
	if pieceType < KING || pieceType >= CONST_FIGURE_LIST_LENGTH {
		return Piece{}, errors.Join(ErrFigureNotSupported, fmt.Errorf("%s", pieceType))
	}
	sequence++
	return Piece{
		figure:        pieceType,
		id:            sequence,
		position:      position,
		isTaken:       false,
		isWhite:       isWhite,
		score:         0,
		targetedCells: make([]Position, 2),
	}, nil
}

func MustNewPiece(pieceType Figure, isWhite bool, position Position) Piece {
	piece, err := NewPiece(pieceType, isWhite, position)
	if err != nil {
		panic(err)
	}
	return piece
}

func (p Piece) GetId() int {
	return p.id
}

func (p Piece) GetScore() int {
	return p.score
}

func (p Piece) IsWhite() bool {
	return p.isWhite
}

func (p Piece) IsTaken() bool {
	return p.isTaken
}

func (p Piece) GetPosition() Position {
	return p.position
}

func (p Piece) GetType() Figure {
	return p.figure
}

func (p Piece) String() string {
	return fmt.Sprint(p.GetType(), p.isWhite, p.position.String(), p.score, p.id)
}

func (p Piece) CanTransform() []Figure {
	if p.GetType() != PAWN {
		return []Figure{}
	}

	if p.IsWhite() && p.GetPosition().GetY() == 7 ||
		!p.IsWhite() && p.GetPosition().GetY() == 0 {
		return []Figure{PAWN, QUEEN, BISHOP, ROOK, KNIGHT}
	}
	return []Figure{}
}

func (p Piece) GetTargetedCells() []Position {
	return p.targetedCells
}

func (p *Piece) ClearTargetedCells(board *Board) {
	for _, position := range p.targetedCells {
		board.GetCell(position).RemoveFromTargetedBy(p.GetId())
	}
}

func (p *Piece) MarkTargetedCells(board *Board) {
	switch p.figure {
	case KING:

	case QUEEN:

	case ROOK:

	case BISHOP:

	case KNIGHT:

	case PAWN:

	default:
		panic(fmt.Sprintf("unknown piece type: %s", p))
	}
}
