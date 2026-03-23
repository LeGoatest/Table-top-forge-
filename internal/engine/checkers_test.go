package engine

import (
	"context"
	"testing"
)

func TestCheckers_Reset(t *testing.T) {
	c := NewCheckers()
	if c.CurrentPlayer != Red {
		t.Errorf("Expected current player to be Red, got %v", c.CurrentPlayer)
	}

	// Verify pieces on board
	if c.Board[0][1] != Red {
		t.Errorf("Expected Red piece at [0][1], got %v", c.Board[0][1])
	}
	if c.Board[7][0] != Black {
		t.Errorf("Expected Black piece at [7][0], got %v", c.Board[7][0])
	}
	if c.Board[3][3] != Empty {
		t.Errorf("Expected Empty at [3][3], got %v", c.Board[3][3])
	}
}

func TestCheckers_Move(t *testing.T) {
	c := NewCheckers()
	params := map[string][]string{
		"fromRow": {"2"},
		"fromCol": {"1"},
		"toRow":   {"3"},
		"toCol":   {"0"},
	}

	err := c.HandleMove(context.Background(), params)
	if err != nil {
		t.Fatalf("HandleMove failed: %v", err)
	}

	if c.Board[2][1] != Empty {
		t.Errorf("Source cell should be empty, got %v", c.Board[2][1])
	}
	if c.Board[3][0] != Red {
		t.Errorf("Target cell should have Red piece, got %v", c.Board[3][0])
	}
	if c.CurrentPlayer != Black {
		t.Errorf("Expected current player to be Black, got %v", c.CurrentPlayer)
	}
}

func TestCheckers_InvalidMove(t *testing.T) {
	c := NewCheckers()
	// Trying to move an empty cell
	params := map[string][]string{
		"fromRow": {"3"},
		"fromCol": {"3"},
		"toRow":   {"4"},
		"toCol":   {"4"},
	}

	err := c.HandleMove(context.Background(), params)
	if err == nil {
		t.Fatal("Expected error when moving empty cell, got nil")
	}
}
