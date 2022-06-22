package move

import "testing"

func TestMove(t *testing.T) {
	move, err := EmptyMove("f4e3")
	if err != nil {
		t.Logf("Failed to create move: %v", err)
	}
	t.Log(move)
	t.Logf("Moving from index %v to %v\n", move.OriginIndex(), move.DestIndex())
}
