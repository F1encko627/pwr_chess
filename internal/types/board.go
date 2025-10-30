package types

import (
	"fmt"
	"slices"
)

type Cell struct {
	piece      *Piece
	targetedBy []*Piece
}

type Board struct {
	board  [8][8]*Cell
	pieces map[bool][]Piece
}

func GetBoard(initialPieces []Piece) (Board, error) {
	board := Board{}
	for _, piece := range initialPieces {
		board.pieces[piece.IsWhite()] = append(board.pieces[piece.IsWhite()], piece)
		board.GetCell(piece.position).piece =
			&board.pieces[piece.IsWhite()][len(board.pieces[piece.IsWhite()])]
	}
	return board, nil
}

func (b *Board) GetCell(pos Position) *Cell {
	return b.board[pos.GetX()][pos.GetY()]
}

func (b *Board) GetPieces(pos Position) map[bool][]Piece {
	return b.pieces
}

func (c *Cell) GetPiece() *Piece {
	return c.piece
}

func (c *Cell) RemoveFromTargetedBy(pieceId int) {
	for i, piece := range c.targetedBy {
		if piece.GetId() == pieceId {
			c.targetedBy = slices.Delete(c.targetedBy, i, i)
			return
		}
	}
}

func (b Board) DebugRender() {
	for y := range len(b.board) {
		fmt.Print(7-y, 8-y, " ")
		for x := range len(b.board[y]) {
			cell := b.GetCell(MustNewPos(x, y))
			if cell.GetPiece() == nil {
				if (x+y)%2 != 0 {
					fmt.Print("  ")
				} else {
					fmt.Print("██")
				}
				continue
			}
			color := "B"
			piece := cell.GetPiece()
			if piece.IsWhite() {
				color = "W"
			}
			fmt.Printf("%s%s", color, piece.GetType().String())
		}
		fmt.Print("\n")
	}
	fmt.Print("    ")
	for i := range 8 {
		fmt.Print(string(rune('a'+i)), " ")
	}
	fmt.Print("\n")
	fmt.Print("    ")
	for i := range 8 {
		fmt.Print(i, " ")
	}
	fmt.Print("\n")
}
