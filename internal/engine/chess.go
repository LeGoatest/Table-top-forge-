package engine

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type ChessPieceType string

const (
	Pawn   ChessPieceType = "p"
	Rook   ChessPieceType = "r"
	Knight ChessPieceType = "n"
	Bishop ChessPieceType = "b"
	Queen  ChessPieceType = "q"
	King   ChessPieceType = "k"
	None   ChessPieceType = ""
)

type ChessColor string

const (
	White ChessColor = "w"
	BlackChess ChessColor = "b"
)

type ChessPiece struct {
	Type  ChessPieceType
	Color ChessColor
}

type Chess struct {
	Board         [8][8]ChessPiece
	CurrentPlayer ChessColor
}

func NewChess() *Chess {
	c := &Chess{}
	c.Reset()
	return c
}

func (c *Chess) Name() string {
	return "Chess"
}

func (c *Chess) Reset() error {
	c.Board = [8][8]ChessPiece{}
	c.CurrentPlayer = White

	// Setup Pawns
	for i := 0; i < 8; i++ {
		c.Board[1][i] = ChessPiece{Type: Pawn, Color: White}
		c.Board[6][i] = ChessPiece{Type: Pawn, Color: BlackChess}
	}

	// Setup pieces
	backRank := []ChessPieceType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for i, t := range backRank {
		c.Board[0][i] = ChessPiece{Type: t, Color: White}
		c.Board[7][i] = ChessPiece{Type: t, Color: BlackChess}
	}

	return nil
}

func (c *Chess) HandleMove(ctx context.Context, params map[string][]string) error {
	fromRow, _ := strconv.Atoi(getFirstParam(params, "fromRow"))
	fromCol, _ := strconv.Atoi(getFirstParam(params, "fromCol"))
	toRow, _ := strconv.Atoi(getFirstParam(params, "toRow"))
	toCol, _ := strconv.Atoi(getFirstParam(params, "toCol"))

	return c.Move(fromRow, fromCol, toRow, toCol)
}

func (c *Chess) Move(fr, fc, tr, tc int) error {
	if !c.isValidCoordinate(fr, fc) || !c.isValidCoordinate(tr, tc) {
		return errors.New("out of bounds")
	}

	piece := c.Board[fr][fc]
	if piece.Type == None || piece.Color != c.CurrentPlayer {
		return errors.New("invalid selection")
	}

	target := c.Board[tr][tc]
	if target.Type != None && target.Color == c.CurrentPlayer {
		return errors.New("cannot capture your own piece")
	}

	// Simple move logic for demonstration (TDD will refine this)
	c.Board[tr][tc] = piece
	c.Board[fr][fc] = ChessPiece{Type: None}

	if c.CurrentPlayer == White {
		c.CurrentPlayer = BlackChess
	} else {
		c.CurrentPlayer = White
	}

	return nil
}

func (c *Chess) isValidCoordinate(r, col int) bool {
	return r >= 0 && r < 8 && col >= 0 && col < 8
}

func (c *Chess) RenderHTML(ctx context.Context, w io.Writer) error { return nil }
func (c *Chess) GetState() (any, error) { return c, nil }
func (c *Chess) LoadState(state any) error { return nil }
func (c *Chess) ResetInternal() error { return c.Reset() }
