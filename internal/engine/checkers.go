package engine

import (
	"context"
	"errors"
	"io"
	"strconv"
)

// PieceType represents the type of piece.
type PieceType int

const (
	Empty PieceType = iota
	Red
	Black
	RedKing
	BlackKing
)

// Position represents a cell on the board.
type Position struct {
	Row int
	Col int
}

// Checkers represents the logic for a checkers game.
type Checkers struct {
	Board         [8][8]PieceType
	CurrentPlayer PieceType // Red or Black
	SelectedPiece *Position // For handling multi-jumps
	JumpActive    bool      // True if the current player is in the middle of a multi-jump
}

// NewCheckers creates a new instance of a checkers game.
func NewCheckers() *Checkers {
	c := &Checkers{}
	c.Reset()
	return c
}

// Name returns the name of the game.
func (c *Checkers) Name() string {
	return "Checkers"
}

// Reset initializes the board for a new game.
func (c *Checkers) Reset() error {
	c.Board = [8][8]PieceType{}
	c.CurrentPlayer = Red
	c.SelectedPiece = nil
	c.JumpActive = false

	// Initialize Red pieces (rows 0, 1, 2)
	for row := 0; row < 3; row++ {
		for col := 0; col < 8; col++ {
			if (row+col)%2 != 0 {
				c.Board[row][col] = Red
			}
		}
	}

	// Initialize Black pieces (rows 5, 6, 7)
	for row := 5; row < 8; row++ {
		for col := 0; col < 8; col++ {
			if (row+col)%2 != 0 {
				c.Board[row][col] = Black
			}
		}
	}

	return nil
}

// HandleMove processes a move from the player.
func (c *Checkers) HandleMove(ctx context.Context, params map[string][]string) error {
	fromRowStr := getFirstParam(params, "fromRow")
	fromColStr := getFirstParam(params, "fromCol")
	toRowStr := getFirstParam(params, "toRow")
	toColStr := getFirstParam(params, "toCol")

	if fromRowStr == "" || fromColStr == "" || toRowStr == "" || toColStr == "" {
		return errors.New("missing move parameters")
	}

	fromRow, _ := strconv.Atoi(fromRowStr)
	fromCol, _ := strconv.Atoi(fromColStr)
	toRow, _ := strconv.Atoi(toRowStr)
	toCol, _ := strconv.Atoi(toColStr)

	return c.Move(fromRow, fromCol, toRow, toCol)
}

// Move executes a move from (fromRow, fromCol) to (toRow, toCol).
func (c *Checkers) Move(fromRow, fromCol, toRow, toCol int) error {
	if !c.isValidCoordinate(fromRow, fromCol) || !c.isValidCoordinate(toRow, toCol) {
		return errors.New("coordinates out of bounds")
	}

	piece := c.Board[fromRow][fromCol]
	if piece == Empty {
		return errors.New("no piece at source location")
	}

	if !c.isCurrentPlayerPiece(piece) {
		return errors.New("not your piece")
	}

	if c.Board[toRow][toCol] != Empty {
		return errors.New("target location is not empty")
	}

	if c.JumpActive && (c.SelectedPiece.Row != fromRow || c.SelectedPiece.Col != fromCol) {
		return errors.New("must continue jumping with the same piece")
	}

	rowDiff := toRow - fromRow
	colDiff := toCol - fromCol
	absRowDiff := abs(rowDiff)
	absColDiff := abs(colDiff)

	if absRowDiff == 1 && absColDiff == 1 && !c.JumpActive {
		if !c.canMoveInDirection(piece, rowDiff) {
			return errors.New("invalid move direction")
		}
		c.executeMove(fromRow, fromCol, toRow, toCol)
		c.switchPlayer()
		return nil
	}

	if absRowDiff == 2 && absColDiff == 2 {
		midRow := (fromRow + toRow) / 2
		midCol := (fromCol + toCol) / 2
		capturedPiece := c.Board[midRow][midCol]

		if capturedPiece == Empty || c.isCurrentPlayerPiece(capturedPiece) {
			return errors.New("no opponent piece to jump over")
		}

		if !c.canMoveInDirection(piece, rowDiff) {
			return errors.New("invalid move direction")
		}

		c.executeJump(fromRow, fromCol, toRow, toCol, midRow, midCol)

		if c.canJumpAgain(toRow, toCol) {
			c.JumpActive = true
			c.SelectedPiece = &Position{Row: toRow, Col: toCol}
		} else {
			c.switchPlayer()
		}
		return nil
	}

	return errors.New("invalid move distance")
}

func (c *Checkers) executeMove(fromRow, fromCol, toRow, toCol int) {
	c.Board[toRow][toCol] = c.Board[fromRow][fromCol]
	c.Board[fromRow][fromCol] = Empty
	c.promoteIfNecessary(toRow, toCol)
}

func (c *Checkers) executeJump(fromRow, fromCol, toRow, toCol, midRow, midCol int) {
	c.Board[toRow][toCol] = c.Board[fromRow][fromCol]
	c.Board[fromRow][fromCol] = Empty
	c.Board[midRow][midCol] = Empty
	c.promoteIfNecessary(toRow, toCol)
}

func (c *Checkers) canJumpAgain(row, col int) bool {
	piece := c.Board[row][col]
	directions := []int{-2, 2}

	for _, dr := range directions {
		for _, dc := range directions {
			tr, tc := row+dr, col+dc
			mr, mc := row+dr/2, col+dc/2

			if c.isValidCoordinate(tr, tc) && c.Board[tr][tc] == Empty {
				midPiece := c.Board[mr][mc]
				if midPiece != Empty && !c.isCurrentPlayerPiece(midPiece) {
					if c.canMoveInDirection(piece, dr) {
						return true
					}
				}
			}
		}
	}
	return false
}

func (c *Checkers) isValidCoordinate(row, col int) bool {
	return row >= 0 && row < 8 && col >= 0 && col < 8
}

func (c *Checkers) isCurrentPlayerPiece(piece PieceType) bool {
	if c.CurrentPlayer == Red {
		return piece == Red || piece == RedKing
	}
	return piece == Black || piece == BlackKing
}

func (c *Checkers) canMoveInDirection(piece PieceType, rowDiff int) bool {
	if piece == Red {
		return rowDiff > 0
	}
	if piece == Black {
		return rowDiff < 0
	}
	return true
}

func (c *Checkers) promoteIfNecessary(row, col int) {
	piece := c.Board[row][col]
	if piece == Red && row == 7 {
		c.Board[row][col] = RedKing
	} else if piece == Black && row == 0 {
		c.Board[row][col] = BlackKing
	}
}

func (c *Checkers) switchPlayer() {
	if c.CurrentPlayer == Red {
		c.CurrentPlayer = Black
	} else {
		c.CurrentPlayer = Red
	}
	c.JumpActive = false
	c.SelectedPiece = nil
}

func (c *Checkers) RenderHTML(ctx context.Context, w io.Writer) error { return nil }

func (c *Checkers) GetState() (any, error) { return c, nil }

func (c *Checkers) LoadState(state any) error { return nil }

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func getFirstParam(params map[string][]string, key string) string {
	values, ok := params[key]
	if !ok || len(values) == 0 {
		return ""
	}
	return values[0]
}
