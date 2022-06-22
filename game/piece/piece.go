package piece

import "unicode"

const (
	Empty  = 0
	King   = 1
	Pawn   = 2
	Knight = 3
	Bishop = 4
	Rook   = 5
	Queen  = 6

	White = 8
	Black = 16
)

type Piece uint8

func (p Piece) String() string {
	var sym rune
	switch p % 8 {
	case Empty:
		sym = ' '
	case King:
		sym = 'k'
	case Pawn:
		sym = 'p'
	case Knight:
		sym = 'n'
	case Bishop:
		sym = 'b'
	case Rook:
		sym = 'r'
	case Queen:
		sym = 'q'
	}
	if p.IsWhite() {
		sym = unicode.ToUpper(sym)
	}
	return string(sym)
}

func (p Piece) IsWhite() bool {
	return p/White == 1
}

func FromRune(r rune) Piece {
	var p Piece
	switch unicode.ToLower(r) {
	case 'k':
		p = King
	case 'p':
		p = Pawn
	case 'n':
		p = Knight
	case 'b':
		p = Bishop
	case 'r':
		p = Rook
	case 'q':
		p = Queen
	}
	if unicode.IsUpper(r) {
		p = p | White
	} else {
		p = p | Black
	}
	return p
}
