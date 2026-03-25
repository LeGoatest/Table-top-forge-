package engine

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"time"
)

type BackgammonPlayer int

const (
	Player1 BackgammonPlayer = 1
	Player2 BackgammonPlayer = 2
)

type Backgammon struct {
	Points       [24]int // Positive for Player1, negative for Player2
	Bar          [2]int  // [0] for Player1, [1] for Player2
	Off          [2]int  // [0] for Player1, [1] for Player2
	CurrentPlayer BackgammonPlayer
	Dice         []int
}

func NewBackgammon() *Backgammon {
	b := &Backgammon{}
	b.Reset()
	return b
}

func (b *Backgammon) Name() string {
	return "Backgammon"
}

func (b *Backgammon) Reset() error {
	b.Points = [24]int{}
	// Initial setup
	b.Points[0] = 2   // Player 1
	b.Points[5] = -5  // Player 2
	b.Points[7] = -3  // Player 2
	b.Points[11] = 5  // Player 1
	b.Points[12] = -5 // Player 2
	b.Points[16] = 3  // Player 1
	b.Points[18] = 5  // Player 1
	b.Points[23] = -2 // Player 2

	b.Bar = [2]int{0, 0}
	b.Off = [2]int{0, 0}
	b.CurrentPlayer = Player1
	b.RollDice()
	return nil
}

func (b *Backgammon) RollDice() {
	rand.Seed(time.Now().UnixNano())
	d1 := rand.Intn(6) + 1
	d2 := rand.Intn(6) + 1
	if d1 == d2 {
		b.Dice = []int{d1, d1, d1, d1}
	} else {
		b.Dice = []int{d1, d2}
	}
}

func (b *Backgammon) Move(from, dieIdx int) error {
	if dieIdx < 0 || dieIdx >= len(b.Dice) {
		return errors.New("invalid die index")
	}

	// Movement logic to be implemented via TDD
	return nil
}

func (b *Backgammon) HandleMove(ctx context.Context, params map[string][]string) error {
	return nil
}

func (b *Backgammon) RenderHTML(ctx context.Context, w io.Writer) error { return nil }
func (b *Backgammon) GetState() (any, error) { return b, nil }
func (b *Backgammon) LoadState(state any) error { return nil }
func (b *Backgammon) ResetInternal() error { return b.Reset() }
