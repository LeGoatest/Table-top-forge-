package engine

import (
	"context"
	"io"
	"math/rand"
	"time"
)

// SlotMachine represents a simple slot machine game.
type SlotMachine struct {
	Reels   []string
	Credits int
	Symbols []string
	LastWin int
	IsSpinning bool
}

// NewSlotMachine creates a new slot machine instance.
func NewSlotMachine() *SlotMachine {
	s := &SlotMachine{
		Credits: 100,
		Symbols: []string{"cherry", "lemon", "grape", "bell", "diamond", "seven"},
	}
	s.Reset()
	return s
}

func (s *SlotMachine) Name() string {
	return "Slot Machine"
}

func (s *SlotMachine) Reset() error {
	s.Reels = []string{"seven", "seven", "seven"}
	s.LastWin = 0
	s.IsSpinning = false
	return nil
}

func (s *SlotMachine) Spin() {
	rand.Seed(time.Now().UnixNano())
	if s.Credits < 10 {
		return
	}
	s.Credits -= 10
	s.IsSpinning = true

	for i := 0; i < 3; i++ {
		s.Reels[i] = s.Symbols[rand.Intn(len(s.Symbols))]
	}

	s.calculateWin()
	s.IsSpinning = false
}

func (s *SlotMachine) calculateWin() {
	if s.Reels[0] == s.Reels[1] && s.Reels[1] == s.Reels[2] {
		// Jackpot!
		win := 0
		switch s.Reels[0] {
		case "cherry": win = 50
		case "lemon": win = 100
		case "grape": win = 150
		case "bell": win = 250
		case "diamond": win = 500
		case "seven": win = 1000
		}
		s.LastWin = win
		s.Credits += win
	} else if s.Reels[0] == s.Reels[1] || s.Reels[1] == s.Reels[2] || s.Reels[0] == s.Reels[2] {
		s.LastWin = 20
		s.Credits += 20
	} else {
		s.LastWin = 0
	}
}

func (s *SlotMachine) HandleMove(ctx context.Context, params map[string][]string) error {
	action := getFirstParam(params, "action")
	if action == "spin" {
		s.Spin()
	}
	return nil
}

func (s *SlotMachine) RenderHTML(ctx context.Context, w io.Writer) error { return nil }
func (s *SlotMachine) GetState() (any, error) { return s, nil }
func (s *SlotMachine) LoadState(state any) error { return nil }
func (s *SlotMachine) ResetInternal() error { return s.Reset() }
