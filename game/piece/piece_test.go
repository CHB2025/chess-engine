package piece_test

import (
	"testing"

	"bareman.net/chess-engine/game/piece"
)

func TestCreation(t *testing.T) {
	p := piece.FromRune('R')
	if p != 5|8 {
		t.Errorf("Expected R, got %v\n", p)
	}
}

func TestIsWhite(t *testing.T) {
	var p piece.Piece = 5 | 8
	if !p.IsWhite() {
		t.Errorf("Expected IsWhite to be true")
	}
}

func TestColor(t *testing.T) {
	var p piece.Piece = 5 | 8
	if p.Color() != piece.White {
		t.Errorf("Expected Piece to be white")
	}
}

func TestType(t *testing.T) {
	var p piece.Piece = 5 | 8
	if p.Type() != piece.Rook {
		t.Errorf("Expected Piece to be Rook")
	}
}
