package types

type State uint8

const (
	WHITE_TURN State = iota
	WHITE_CHECK
	BLACK_TURN
	BLACK_CHECK
	WHITE_CHECKMATE
	BLACK_CHECKMATE
	STALEMATE
	// PAUSE // Reconsider for blitz
)