package types

import (
	"errors"
	"fmt"
)

type Direction int

const (
	SAME_SQUARE = iota
	UP          // ix == fx && iy < fy
	LEFT_UP     // ix > fx && iy < fy
	LEFT        // iy == fy && ix > fx
	LEFT_DOWN   // ix > fx && iy > fy
	DOWN        // ix == fx && iy > fy
	RIGHT_DOWN  // ix < fx && iy > fy
	RIGHT       // iy == fy && ix < fx
	RIGHT_UP    // ix < fx && iy < fy
)

type Move struct {
	posInit   Position
	posFinal  Position
	Direction Direction
}

var (
	ErrSameSquaremove = errors.New("same quare move")
	ErrInitialPos     = errors.New("invalid initial position")
	ErrFinalPos       = errors.New("invalid final position")
)

func GetMove(ix, iy, fx, fy int) (Move, error) {
	posInit, err := NewPos(ix, iy)
	if err != nil {
		return Move{}, errors.Join(ErrInitialPos, err)
	}

	posFinal, err := NewPos(fx, fy)
	if err != nil {
		return Move{}, errors.Join(ErrFinalPos, err)
	}

	var move = Move{posInit, posFinal, getMoveDirection(ix, iy, fx, fy)}

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

func getMoveDirection(ix, iy, fx, fy int) Direction {
	/*
		+-----+-----+
		|x0;y1|x1;y1|
		+-----+-----+
		|x0;y0|x1;y0|
		+-----+-----+
	*/
	if ix == fx && iy < fy {
		return UP
	}
	if ix == fx && iy > fy {
		return DOWN
	}
	if iy == fy && ix < fx {
		return RIGHT
	}
	if iy == fy && ix > fx {
		return LEFT
	}
	if ix > fx && iy < fy {
		return LEFT_UP
	}
	if ix > fx && iy > fy {
		return LEFT_DOWN
	}
	if ix < fx && iy < fy {
		return RIGHT_UP
	}
	if ix < fx && iy > fy {
		return RIGHT_DOWN
	}
	return SAME_SQUARE
}
