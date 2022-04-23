package board

import "unicode"

type Position [2]int

type Piece struct {
	Symbol  rune
	IsWhite bool
}

func (p *Piece) String() string {
	sym := p.Symbol
	if p.IsWhite {
		sym = unicode.ToUpper(sym)
	}
	return string(sym)
}

func (p *Piece) ValidMoves(board Board) []Position {
	var moves []Position
	var pos Position
	for y, row := range board {
		for x, other := range row {
			if other == p { //Not sure if this works
				pos = Position{y, x}
				break
			}
		}
	}
	switch p.Symbol {
	case 'p':
	}
}
