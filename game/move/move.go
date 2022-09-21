package move

import (
	"fmt"
	"regexp"
	"strconv"

	"bareman.net/chess-engine/game/piece"
)

const (
	PositionRegex = `^[a-hA-H][0-8]$`
	MoveRegex     = `^([a-hA-H][0-8]){2}[qrbnQRBN]?$`
)

type Move struct {
	Origin     string
	Dest       string
	Capture    piece.Piece
	Promotion  piece.Piece
	EnPassant  bool
	Castle     bool
	BoardState struct {
		WQCastle bool
		WKCastle bool
		BQCastle bool
		BKCastle bool
		EPTarget int
	}
}

func (m *Move) String() string {
	return m.Origin + m.Dest
}

func (m *Move) OriginIndex() int {
	row, _ := strconv.Atoi(string(m.Origin[1]))
	col := int(m.Origin[0] - 'a')
	return (row-1)<<3 + col
}

func (m *Move) DestIndex() int {
	row, _ := strconv.Atoi(string(m.Dest[1]))
	col := int(m.Dest[0] - 'a')
	return (row-1)<<3 + col
}

func FullMove(mv string, Capture piece.Piece, EnPassant, Castle bool) (*Move, error) {
	m, err := EmptyMove(mv)
	if err != nil {
		return nil, err
	}

	m.Capture, m.EnPassant, m.Castle = Capture, EnPassant, Castle

	return m, nil

}

func EmptyMove(mv string) (*Move, error) {
	reg := regexp.MustCompile(MoveRegex)
	if !reg.MatchString(mv) {
		return nil, fmt.Errorf("Invalid Move given. Received %v\n", mv)
	}

	var promote piece.Piece = piece.Empty
	if len(mv) == 5 {
		promote = piece.FromRune(rune(mv[4]))
	}

	return &Move{
		Origin:    mv[:2],
		Dest:      mv[2:4],
		Promotion: promote,
	}, nil
}
