package game

import (
	"fmt"
	"math/rand"

	"bareman.net/chess-engine/game/move"
	"bareman.net/chess-engine/game/piece"
)

const (
	BTMHashIndex           = 768
	WKCastleHashIndex      = 769
	WQCastleHashIndex      = 770
	BKCastleHashIndex      = 771
	BQCastleHashIndex      = 772
	EPTargetHashIndexStart = 773
)

func (g *Game) InitializeHash() {
	var seed int64 = rand.Int63()
	source := rand.NewSource(seed)
	r := rand.New(source)
	fmt.Printf("Hashing using seed %v\n", seed)
	for i := 0; i < 781; i++ {
		g.hashKeys[i] = r.Uint64()
	}
}

// Must be done after making/before unmaking to work properly
func (g *Game) incrementHash(m *move.Move, p piece.Piece) {
	g.Hash ^= g.hashKeys[hashIndex(p, m.OriginIndex())]
	if m.Promotion == piece.Empty {
		g.Hash ^= g.hashKeys[hashIndex(p, m.DestIndex())]
	} else {
		g.Hash ^= g.hashKeys[hashIndex(m.Promotion, m.DestIndex())]
	}

	if m.Capture != piece.Empty && !m.EnPassant {
		g.Hash ^= g.hashKeys[hashIndex(m.Capture, m.DestIndex())]
	}
	if m.EnPassant {
		file := m.Dest[0]
		rank := m.Origin[1]
		index := indexFromPosition(string(file) + string(rank))
		g.Hash ^= g.hashKeys[hashIndex(m.Capture, index)]
	}

	if m.Castle {
		rank := m.Origin[1]
		file := byte(int(m.Origin[0]) + (int(m.Dest[0])-int(m.Origin[0]))/2)
		rePos := string(file) + string(rank)
		rsPos := "a" + string(rank)
		if file > m.Origin[0] {
			rsPos = "h" + string(rank)
		}
		reIndex := indexFromPosition(rePos)
		rsIndex := indexFromPosition(rsPos)

		g.Hash ^= g.hashKeys[hashIndex(piece.Rook|p.Color(), reIndex)]
		g.Hash ^= g.hashKeys[hashIndex(piece.Rook|p.Color(), rsIndex)]
	}

	g.Hash ^= g.hashKeys[BTMHashIndex]
	if g.WKCastle != m.BoardState.WKCastle {
		g.Hash ^= g.hashKeys[WKCastleHashIndex]
	}
	if g.WQCastle != m.BoardState.WQCastle {
		g.Hash ^= g.hashKeys[WQCastleHashIndex]
	}
	if g.BKCastle != m.BoardState.BKCastle {
		g.Hash ^= g.hashKeys[BKCastleHashIndex]
	}
	if g.BQCastle != m.BoardState.BQCastle {
		g.Hash ^= g.hashKeys[BQCastleHashIndex]
	}
	if g.EPTarget != m.BoardState.EPTarget {
		if g.EPTarget != -1 {
			_, col := coordinates(g.EPTarget)
			g.Hash ^= g.hashKeys[BQCastleHashIndex+col]
		}
		if m.BoardState.EPTarget != -1 {
			_, col := coordinates(m.BoardState.EPTarget)
			g.Hash ^= g.hashKeys[BQCastleHashIndex+col]
		}
	}
}

func hashIndex(p piece.Piece, index int) int {
	return (int(p.Type()-1)<<1+int(p)>>4)<<6 + index
}

func Hash(g *Game) uint64 {
	// hash key index is 2*PieceType*64 + 64 if black + index
	// 2*p.Type()*64 + 64*p >> 4 + i
	var hash uint64
	for i, p := range g.Board {
		if p != piece.Empty {
			hashKeyIndex := hashIndex(p, i)
			hash ^= g.hashKeys[hashKeyIndex]
		}
	}
	// q on square 63 would be index 767
	if !g.WhiteToMove {
		hash ^= g.hashKeys[BTMHashIndex]
	}
	if g.WKCastle {
		hash ^= g.hashKeys[WKCastleHashIndex]
	}
	if g.WQCastle {
		hash ^= g.hashKeys[WQCastleHashIndex]
	}
	if g.BKCastle {
		hash ^= g.hashKeys[BKCastleHashIndex]
	}
	if g.BQCastle {
		hash ^= g.hashKeys[BQCastleHashIndex]
	}
	if g.EPTarget != -1 {
		_, col := coordinates(g.EPTarget)
		hash ^= g.hashKeys[BQCastleHashIndex+col]
	}

	return hash
}
