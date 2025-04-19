package types

type Pos int8

func NewPos(x, y int) Pos {
	return Pos(y * 8 + x)
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