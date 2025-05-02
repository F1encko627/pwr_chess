package types

type State uint8

const (
	NORMAL State = iota
	PAUSE
	WHITE_CHECKMATE
	BLACK_CHECKMATE
	STALEMATE
)
