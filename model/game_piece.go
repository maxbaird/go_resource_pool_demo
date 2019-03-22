package model

import "fmt"

type (
	// GamePiece is abstraction for a cube/word/piece
	GamePiece struct {
		Id int
		Field2 string
		Field3 [40*1024*1024]int
	}
)

func (gp *GamePiece) String() string {
	if gp == nil {
		return "GamePiece<nil>"
	}
	return fmt.Sprintf("GamePiece-id: %v", gp.Id)
}
