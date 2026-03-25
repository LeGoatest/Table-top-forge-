package engine

import (
	"testing"
)

func TestChess_Init(t *testing.T) {
	c := NewChess()
	if c.CurrentPlayer != White {
		t.Errorf("Expected White to start, got %v", c.CurrentPlayer)
	}
	if c.Board[0][0].Type != Rook || c.Board[0][0].Color != White {
		t.Errorf("Expected White Rook at [0][0], got %v", c.Board[0][0])
	}
	if c.Board[7][4].Type != King || c.Board[7][4].Color != BlackChess {
		t.Errorf("Expected Black King at [7][4], got %v", c.Board[7][4])
	}
}

func TestChess_SimpleMove(t *testing.T) {
	c := NewChess()
	err := c.Move(1, 0, 3, 0) // White pawn a2 to a4
	if err != nil {
		t.Fatalf("Move failed: %v", err)
	}
	if c.Board[3][0].Type != Pawn {
		t.Errorf("Expected Pawn at [3][0]")
	}
	if c.CurrentPlayer != BlackChess {
		t.Errorf("Expected Black's turn")
	}
}

func TestChess_InvalidTurn(t *testing.T) {
	c := NewChess()
	err := c.Move(6, 0, 4, 0) // Try to move Black on White's turn
	if err == nil {
		t.Error("Should not be able to move Black piece on White's turn")
	}
}
