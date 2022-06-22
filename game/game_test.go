package game

import (
	"testing"
)

type Position struct {
	Name  string
	Fen   string
	Depth []int
	Nodes []int
}

func TestingPositions() []Position {
	positions := make([]Position, 0, 10)
	positions = append(positions, Position{
		Name:  "Initial Position",
		Fen:   "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		Depth: []int{1, 2, 3, 5},
		Nodes: []int{20, 400, 8902, 4865609},
	})
	positions = append(positions, Position{
		Name:  "Endgame",
		Fen:   "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		Depth: []int{1, 2, 3, 5},
		Nodes: []int{14, 191, 2812, 674624},
	})
	return positions
}

// Get more positions to test here: https://gist.github.com/peterellisjones/8c46c28141c162d1d8a0f0badbc9cff9
func TestFEN(t *testing.T) {
	fenStrings := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r6r/1b2k1bq/8/8/7B/8/8/R3K2R b KQ - 3 2",
		"8/8/8/2k5/2pP4/8/B7/4K3 b - d3 0 3",
		"r1bqkbnr/pppppppp/n7/8/8/P7/1PPPPPPP/RNBQKBNR w KQkq - 2 2",
		"r3k2r/p1pp1pb1/bn2Qnp1/2qPN3/1p2P3/2N5/PPPBBPPP/R3K2R b KQkq - 3 2",
	}
	for _, fen := range fenStrings {
		g, err := FromFEN(fen)
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
		t.Log(position.Name)
		g, err := FromFEN(position.Fen)
		if err != nil {
			t.Errorf("Failed to create game with fen '%v'\n", position.Fen)
			continue
		}
		for i, depth := range position.Depth {
			if position.Nodes[i] > 10_000 {
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

func TestMovesLogged(t *testing.T) {
	// positions := TestingPositions()

	// for _, position := range positions {
	// 	t.Log(position.Name)
	// 	g, err := FromFEN(position.Fen)
	// 	if err != nil {
	// 		t.Errorf("Failed to create game with fen '%v'\n", position.Fen)
	// 		continue
	// 	}
	// 	// extra move in e2e4 and g2g4
	// 	perft := g.DividedPerft(3)
	// 	for key, val := range perft {
	// 		t.Logf("%v: %v\n", key, val)
	// 	}
	// }

	g, _ := FromFEN("rnbqkbnr/1ppppppp/8/p7/7P/8/PPPPPPP1/RNBQKBNR w KQkq a6 0 3")
	perft := g.DividedPerft(1)
	var sum int
	for key, val := range perft {
		sum += val
		t.Logf("%v: %v\n", key, val)
	}
	t.Logf("%v Nodes searched\n", sum)

}
