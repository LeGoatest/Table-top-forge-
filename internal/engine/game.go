package engine

import (
	"context"
	"io"
)

// Game defines the interface for all board games in the engine.
type Game interface {
	// Name returns the name of the game.
	Name() string

	// HandleMove processes a move from the player.
	// It returns an error if the move is invalid.
	HandleMove(ctx context.Context, params map[string][]string) error

	// RenderHTML writes the game's current state as HTML to the provided writer.
	RenderHTML(ctx context.Context, w io.Writer) error

	// GetState returns a serializable state for persistence.
	GetState() (any, error)

	// LoadState restores the game from a serializable state.
	LoadState(state any) error

	// Reset resets the game to its initial state.
	Reset() error
}
