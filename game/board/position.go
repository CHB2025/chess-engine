package board

import (
	"fmt"
	"strconv"
)

type Position struct {
	row int
	col rune
}

func (this Position) String() string {
	return string(this.col) + fmt.Sprint(this.row)
}

func PositionFromString(str string) (Position, error) {
	if len(str) != 2 {
		return Position{}, fmt.Errorf("Length of string provided must be 2. Received %v", len(str))
	}
	row, err := strconv.Atoi(string(str[1]))
	if err != nil {
		return Position{}, err
	}
	return Position{row, rune(str[0])}, nil
}
