package asm

// Filling methods for instructions
const (
	fillNone = iota
	fillLink // for jumps
	fillLow  // for immediate instructions
	fillHigh // for addui
	fillLabel
)
