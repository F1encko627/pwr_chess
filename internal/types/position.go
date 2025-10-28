package types

import (
	"errors"
	"fmt"
)

// For Cartesian coordinates... skill issue
type Position struct {
	x, y int
}

var ErrOutOfBounds = errors.New("out of bounds position")

func NewPos(x, y int) (Position, error) {
	if x < 0 || x > 8 {
		return Position{},
			errors.Join(ErrOutOfBounds, errors.New(Position{x, y}.String()))
	}
	return Position{x, y}, nil
}

func MustNewPos(x, y int) Position {
	var position, err = NewPos(x, y)
	if err != nil {
		panic(err)
	}
	return position
}

func (p Position) String() string {
	return fmt.Sprintf("(%d; %d)", p.GetX(), p.GetY())
}

func (p Position) GetX() int {
	return p.x
}

func (p Position) GetY() int {
	return p.y
}
