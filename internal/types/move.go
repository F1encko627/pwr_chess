package types

import (
	"errors"
	"fmt"
)

type Direction int

const (
	SAME_SQUARE = iota
	UP          // same x, asc y
	LEFT_UP     // desc x, asc y
	LEFT        // desc x, same y
	LEFT_DOWN   // desc x, desc y
	DOWN        // same x, desc y
	RIGHT_DOWN  // asc x, desc y
	RIGHT       // asc x, same y
	RIGHT_UP    // asc x, asc y
)

func getmoveDirection(ix, iy, fx, fy int) Direction {
	if ix == fx && iy < fy {
		return UP
	}
	if ix == fx {
		return DOWN
	}
	if iy == fy && ix < fx {
		return RIGHT
	}
	if iy == fy {
		return LEFT
	}
	if ix < fx && iy < fy {
		return LEFT_UP
	}
	if ix < fx {
		return LEFT_DOWN
	}
	if ix > fx && iy < fy {
		return RIGHT_UP
	}
	if ix > fx {
		return RIGHT_DOWN
	}
	return SAME_SQUARE
}

type Move struct {
	posInit   Position
	posFinal  Position
	Direction Direction
}

var (
	ErrSameSquaremove = errors.New("same quare Move")
	ErrInitialPos     = errors.New("invalid initial position")
	ErrFinalPos       = errors.New("invalid final position")
)

func Getmove(ix, iy, fx, fy int) (Move, error) {
	posInit, err := NewPos(ix, iy)
	if err != nil {
		return Move{}, errors.Join(ErrInitialPos, err)
	}

	posFinal, err := NewPos(fx, fy)
	if err != nil {
		return Move{}, errors.Join(ErrFinalPos, err)
	}

	var move = Move{posInit, posFinal, getmoveDirection(ix, iy, fx, fy)}

	if move.GetDirection() == SAME_SQUARE {
		return Move{},
			errors.Join(ErrSameSquaremove, errors.New(move.String()))
	}
	return move, nil
}

func (m Move) String() string {
	return fmt.Sprintf("%s->%s", m.GetInitial(), m.GetFinal())
}

func (m Move) GetInitial() Position {
	return m.posInit
}

func (m Move) GetFinal() Position {
	return m.posFinal
}

func (m Move) GetDirection() Direction {
	return m.Direction
}
