package types

import "fmt"

// For Cartesian coordinates... skill issue
type Pos int8

func NewPos(x, y int) Pos {
	if x > -1 && y > -1 && x < 8 && y < 8 {
		return Pos(y*8 + x)
	} else {
		return Pos(-1)
	}
}

func (p Pos) String() string {
	if p.IsValid() {
		return fmt.Sprintf("[%d;%d]", p.GetX(), p.GetY())
	} else {
		return "NaN"
	}
}

func (cur *Pos) Transform(t Pos) bool {
	newPos := NewPos(cur.GetX()+t.GetX(), cur.GetY()+t.GetY())
	if newPos.IsValid() {
		*cur = newPos
		return true
	}
	return false
}

func (p Pos) IsValid() bool {
	return p > -1 && p < 64
}

func (p Pos) GetX() int {
	return int(p % 8)
}

func (p Pos) GetY() int {
	return int(p / 8)
}
