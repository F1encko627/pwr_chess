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
	KING Figure = '♚' + iota
	QUEEN
	ROOK
	BISHOP
	KNIGHT
	PAWN
	CONST_FIGURE_LIST_LENGTH
	EMPTY Figure = ' '
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
	ErrFigureNotSupported = errors.New("figure not in provided constants")   // Rare ocasion during board build or piece transformation
	ErrSameColorPiece     = errors.New("can't take own piece")               // Maybe should wrap piece position
	ErrWrongMovePattern   = errors.New("wrong move pattern")                 // "piece can't move like that"
	ErrMoveNotPossibleNow = errors.New("move not possible now")              // Wrong prerequisites for a potentially valid move
	ErrCantJumpOverPieces = errors.New("can't jump over pieces")             // Wrapper for piece position
	ErrEnPassantMove      = errors.New("pawn avaliable for en passant")      // Signal for Game Service
	ErrEnPassantTake      = errors.New("check for possible en passant take") // Signal for Game Service
	ErrCastleMove         = errors.New("check for possible castle")          // Signal for Game Service
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

func (p *Piece) Take() {
	p.isTaken = true
}

// ID generation
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

func (p *Piece) MakeMove(move Move, board *Board) error {
	switch p.figure {
	case KING:
		return checkMoveKing(move, board)
	case QUEEN:
		return checkMoveQueen(move, board)
	case ROOK:
		return checkMoveRook(move, board)
	case BISHOP:
		return checkMoveBishop(move, board)
	case KNIGHT:
		return checkMoveKnight(move, board)
	case PAWN:
		return checkMovePawn(move, board)
	default:
		return errors.Join(ErrFigureNotSupported, fmt.Errorf("%s", p.figure))
	}
}

func checkMoveKnight(move Move, board *Board) error {
	if !checkMovePatternKnight(move) {
		return errors.Join(ErrWrongMovePattern, fmt.Errorf("%s", move))
	}
	if !checkSameColorTake(move, board) {
		return ErrSameColorPiece
	}
	return nil
}

func checkMoveRook(move Move, board *Board) error {
	if !checkMovePatternStraight(move) {
		return errors.Join(ErrWrongMovePattern, fmt.Errorf("%s", move))
	}
	if obsticle := checkJumpOverPieceStraightOrDiagonal(move, board); obsticle != nil {
		return errors.Join(ErrCantJumpOverPieces, obsticle)
	}
	if !checkSameColorTake(move, board) {
		return ErrSameColorPiece
	}
	return nil
}

func checkMoveBishop(move Move, board *Board) error {
	if !checkMovePatternDiagonal(move) {
		return errors.Join(ErrWrongMovePattern, fmt.Errorf("%s", move))
	}
	if obsticle := checkJumpOverPieceStraightOrDiagonal(move, board); obsticle != nil {
		return errors.Join(ErrCantJumpOverPieces, obsticle)
	}
	if !checkSameColorTake(move, board) {
		return ErrSameColorPiece
	}

	return nil
}

func checkMoveQueen(move Move, board *Board) error {
	if !checkMovePatternStraight(move) && !checkMovePatternDiagonal(move) {
		return errors.Join(ErrWrongMovePattern, fmt.Errorf("%s", move))
	}
	if obsticle := checkJumpOverPieceStraightOrDiagonal(move, board); obsticle != nil {
		return errors.Join(ErrCantJumpOverPieces, obsticle)
	}
	if !checkSameColorTake(move, board) {
		return ErrSameColorPiece
	}
	return nil
}

func checkMovePawn(move Move, board *Board) error {
	piece := board.GetCell(move.GetInitial()).GetPiece()
	ix := move.GetInitial().GetX()
	iy := move.GetInitial().GetY()
	fy := move.GetFinal().GetY()
	fx := move.GetFinal().GetX()
	dx := iAbs(ix - fx)
	dy := iAbs(iy - fy)

	if dy > 2 || dx > 1 || dx == 1 && dy != 1 {
		return ErrWrongMovePattern
	}
	if dy == 2 {
		if piece.IsWhite() && iy != 1 || !piece.IsWhite() && iy != 6 {
			return ErrMoveNotPossibleNow
		}
		if move.GetDirection() == UP {
			pos := MustNewPos(ix, iy+1)
			if board.GetCell(pos).GetPiece() != nil {
				return errors.Join(ErrCantJumpOverPieces, fmt.Errorf("%s", pos))
			}
		} else {
			pos := MustNewPos(ix, iy-1)
			if board.GetCell(pos).GetPiece() != nil {
				errors.Join(ErrCantJumpOverPieces, fmt.Errorf("%s", pos))
			}
		}
		left, err := NewPos(ix-1, fy)
		if err == nil {
			if l_piece := board.GetCell(left).GetPiece(); l_piece != nil {
				return ErrEnPassantMove
			}
		}
		right, err := NewPos(ix-1, fy)
		if err == nil {
			if r_piece := board.GetCell(right).GetPiece(); r_piece != nil {
				return ErrEnPassantMove
			}
		}
	} else if dx == 1 && board.GetCell(move.GetFinal()).GetPiece() == nil && board.GetCell(MustNewPos(iy, iy)).GetPiece() != nil && (piece.IsWhite() && board.GetCell(MustNewPos(ix, iy-1)).GetPiece() != nil || board.GetCell(MustNewPos(ix, iy+1)).GetPiece() != nil) {
		return ErrEnPassantTake
	} else {
		front := board.GetCell(move.GetFinal()).GetPiece()
		if front != nil {
			return ErrMoveNotPossibleNow
		}
	}
	return nil
}

func checkMoveKing(move Move, board *Board) error {
	return nil
}

func checkMovePatternStraight(move Move) bool {
	dir := move.GetDirection()
	return (dir == UP || dir == DOWN || dir == LEFT || dir == RIGHT)
}

func checkMovePatternDiagonal(move Move) bool {
	return (iAbs(move.GetFinal().GetX()-move.GetInitial().GetX()) ==
		iAbs(move.GetFinal().GetY()-move.GetInitial().GetY()))
}

func checkMovePatternKnight(move Move) bool {
	dx := iAbs(move.GetFinal().GetX() - move.GetInitial().GetX())
	dy := iAbs(move.GetFinal().GetY() - move.GetInitial().GetY())
	return (dx == 2 && dy == 1) || (dx == 1 && dy == 2)
}

func checkJumpOverPieceStraightOrDiagonal(move Move, board *Board) error {
	ix := move.GetInitial().GetX()
	fx := move.GetFinal().GetX()
	iy := move.GetInitial().GetY()
	fy := move.GetFinal().GetY()
	switch move.GetDirection() {
	case UP:
		iy++
		for iy < fy {
			pos := MustNewPos(ix, iy)
			if piece := board.GetCell(pos).GetPiece(); piece != nil {
				return fmt.Errorf("%s", pos)
			}
			iy++
		}
	case DOWN:
		iy--
		for iy > fy {
			pos := MustNewPos(ix, iy)
			if piece := board.GetCell(pos).GetPiece(); piece != nil {
				return fmt.Errorf("%s", pos)
			}
			iy--
		}
	case LEFT:
		ix--
		for ix > fx {
			pos := MustNewPos(ix, iy)
			if piece := board.GetCell(pos).GetPiece(); piece != nil {
				return fmt.Errorf("%s", pos)
			}
			ix--
		}
	case RIGHT:
		ix++
		for ix < fx {
			pos := MustNewPos(ix, iy)
			if piece := board.GetCell(pos).GetPiece(); piece != nil {
				return fmt.Errorf("%s", pos)
			}
			ix++
		}
	case LEFT_UP:
		ix--
		iy++
		for ix > fx && iy < fy {
			pos := MustNewPos(ix, iy)
			if piece := board.GetCell(pos).GetPiece(); piece != nil {
				return fmt.Errorf("%s", pos)
			}
			ix--
			iy++
		}
	case LEFT_DOWN:
		ix--
		iy--
		for ix > fx && iy > fy {
			pos := MustNewPos(ix, iy)
			if piece := board.GetCell(pos).GetPiece(); piece != nil {
				return fmt.Errorf("%s", pos)
			}
			ix--
			iy--
		}
	case RIGHT_UP:
		ix++
		iy++
		for ix < fx && iy < fy {
			pos := MustNewPos(ix, iy)
			if piece := board.GetCell(pos).GetPiece(); piece != nil {
				return fmt.Errorf("%s", pos)
			}
			ix++
			iy++
		}
	case RIGHT_DOWN:
		ix++
		iy--
		for ix < fx && iy > fy {
			pos := MustNewPos(ix, iy)
			if piece := board.GetCell(pos).GetPiece(); piece != nil {
				return fmt.Errorf("%s", pos)
			}
			ix++
			iy--
		}
	}

	return nil
}

func checkSameColorTake(move Move, board *Board) bool {
	targetedPiece := board.GetCell(move.GetFinal()).GetPiece()
	if targetedPiece == nil {
		return true
	}
	return targetedPiece.IsWhite() != board.GetCell(move.GetInitial()).GetPiece().IsWhite()
}

func iAbs(n int) int {
	return n * ((2*n + 1) % 2)
}
