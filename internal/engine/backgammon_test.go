package engine

import (
	"testing"
)

func TestBackgammon_Init(t *testing.T) {
	b := NewBackgammon()
	if b.Points[0] != 2 {
		t.Errorf("Expected 2 checkers at point 0, got %d", b.Points[0])
	}
	if len(b.Dice) < 2 {
		t.Errorf("Expected at least 2 dice, got %d", len(b.Dice))
	}
}

func TestBackgammon_Dice(t *testing.T) {
	b := NewBackgammon()
	// Test double roll
	// We can't easily test random, but we can test the logic if we mock or just run many times
	foundDouble := false
	for i := 0; i < 100; i++ {
		b.RollDice()
		if len(b.Dice) == 4 {
			foundDouble = true
			if b.Dice[0] != b.Dice[1] || b.Dice[1] != b.Dice[2] || b.Dice[2] != b.Dice[3] {
				t.Errorf("Invalid double dice: %v", b.Dice)
			}
		}
	}
	if !foundDouble {
		t.Log("No double found in 100 rolls, statistically possible but unlikely")
	}
}
