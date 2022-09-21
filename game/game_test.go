package game_test

import (
	"fmt"
	"testing"

	"bareman.net/chess-engine/game"
)

type Position struct {
	Name  string
	Fen   string
	Depth []int
	Nodes []int
}

// Testing positions from the Chess Programming Wiki: https://www.chessprogramming.org/Perft_Results
func TestingPositions() []Position {
	return []Position{
		{
			Name:  "Initial Position",
			Fen:   "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			Depth: []int{1, 2, 3, 5},
			Nodes: []int{20, 400, 8902, 4_865_609},
		},
		{
			Name:  "Kiwipete",
			Fen:   "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			Depth: []int{1, 2, 3, 5},
			Nodes: []int{48, 2039, 97_862, 1_936_900_690},
		},
		{
			Name:  "Endgame",
			Fen:   "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			Depth: []int{1, 2, 3, 5},
			Nodes: []int{14, 191, 2812, 67_4624},
		},
		{
			Name:  "Middle Game",
			Fen:   "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			Depth: []int{1, 2, 3, 5},
			Nodes: []int{6, 264, 9467, 15_833_292},
		},
		{
			Name:  "Talkchess",
			Fen:   "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
			Depth: []int{1, 2, 3, 5},
			Nodes: []int{44, 1486, 62_379, 89_941_194},
		},
		{
			Name:  "Edwards 2",
			Fen:   "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
			Depth: []int{1, 2, 3, 5},
			Nodes: []int{46, 2_079, 89_890, 164_075_551},
		},
	}
}

func TestFEN(t *testing.T) {
	fenStrings := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r6r/1b2k1bq/8/8/7B/8/8/R3K2R b KQ - 3 2",
		"8/8/8/2k5/2pP4/8/B7/4K3 b - d3 0 3",
		"r1bqkbnr/pppppppp/n7/8/8/P7/1PPPPPPP/RNBQKBNR w KQkq - 2 2",
		"r3k2r/p1pp1pb1/bn2Qnp1/2qPN3/1p2P3/2N5/PPPBBPPP/R3K2R b KQkq - 3 2",
	}
	for _, fen := range fenStrings {
		g, err := game.FromFEN(fen)
		if err != nil {
			t.Errorf("Failed to create game from FEN string: %s\n", err)
			continue
		}
		newFen := g.ToFEN()
		if newFen != fen {
			t.Errorf("Failed to match output fen string to input.\n Input: %v\n Output: %v\n", fen, newFen)
		}
	}
}

func TestMoves(t *testing.T) {
	positions := TestingPositions()

	for _, position := range positions {
		t.Logf("Testing %v\n", position.Name)
		g, err := game.FromFEN(position.Fen)
		if err != nil {
			t.Errorf("Failed to create game with fen '%v'\n", position.Fen)
			continue
		}
		for i, depth := range position.Depth {
			if position.Nodes[i] > 90_000 {
				t.Logf("Skipping depth %v. Too slow\n", depth)
				break
			}
			calculatedNodes := g.Perft(depth)
			expectedNodes := position.Nodes[i]
			t.Logf("Depth %v: Expected %v, Got %v\n", depth, expectedNodes, calculatedNodes)
			if calculatedNodes != expectedNodes {
				t.Errorf("%v nodes off\n", calculatedNodes-expectedNodes)
			}
		}
	}
}

func TestIncrementalHash(t *testing.T) {
	positions := TestingPositions()

	for _, position := range positions {
		g, err := game.FromFEN(position.Fen)
		if err != nil {
			t.Errorf("Failed to create game with fen '%v'\n", position.Fen)
			continue
		}
		if g.Hash != game.Hash(g) {
			t.Errorf("Initial Hashes do not match\n")
		}
		mvs := g.AllLegalMoves()
		for _, m := range mvs {
			g.Make(m)
			if g.Hash != game.Hash(g) {
				t.Errorf("Hashes do not match after making move\n")
			}
			g.Unmake()
			if g.Hash != game.Hash(g) {
				t.Errorf("Hashes do not match after Unmaking move\n")
			}
		}

	}
}

func BenchmarkMoveGeneration(b *testing.B) {
	b.StopTimer()
	positions := TestingPositions()

	for _, position := range positions {
		g, err := game.FromFEN(position.Fen)
		if err != nil {
			b.Errorf("Failed to create game with fen '%v'\n", position.Fen)
			continue
		}

		b.StartTimer()
		b.Run(position.Name, func(b *testing.B) { g.Perft(3) })
		b.StopTimer()
	}
}

func BenchmarkInitialPosition(b *testing.B) {
	b.StopTimer()
	b.ResetTimer()
	g := game.Default()
	b.StartTimer()
	for i := 1; i < 6; i++ {
		b.Run(fmt.Sprintf("Depth_%v", i), func(b *testing.B) { g.Perft(i) })
	}
}

func BenchmarkLegalMoves(b *testing.B) {
	b.StopTimer()
	b.ResetTimer()
	positions := TestingPositions()
	for _, p := range positions {
		g, err := game.FromFEN(p.Fen)
		if err != nil {
			b.Errorf("Failed to build game from Fen %v\n", p.Fen)
			continue
		}
		b.StartTimer()
		b.Run(p.Name, func(b *testing.B) { g.AllLegalMoves() })
		b.StopTimer()
	}
}

func BenchmarkPseudolegalMoves(b *testing.B) {
	b.StopTimer()
	b.ResetTimer()
	positions := TestingPositions()
	for _, p := range positions {
		g, err := game.FromFEN(p.Fen)
		if err != nil {
			b.Errorf("Failed to build game from Fen %v\n", p.Fen)
			continue
		}
		b.StartTimer()
		b.Run(p.Name, func(b *testing.B) { g.PseudoLegalMoves("") })
		b.StopTimer()
	}
}

func BenchmarkMake(b *testing.B) {
	b.StopTimer()
	b.ResetTimer()
	positions := TestingPositions()
	for _, p := range positions {
		g, err := game.FromFEN(p.Fen)
		if err != nil {
			b.Errorf("Failed to build game from Fen %v\n", p.Fen)
			continue
		}
		moves := g.AllLegalMoves()
		b.StartTimer()
		b.Run(p.Name, func(b *testing.B) {
			for _, mv := range moves {
				g.Make(mv)
			}
		})
		b.StopTimer()
	}
}
