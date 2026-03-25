package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"syscall/js"

	"github.com/a-h/templ"
	"github.com/user/boardgame-engine/internal/db"
	"github.com/user/boardgame-engine/internal/engine"
	"github.com/user/boardgame-engine/templates"
)

var (
	game       *engine.Checkers
	slots      *engine.SlotMachine
	blackjack  *engine.Blackjack
	dominoes   *engine.Dominoes
	chess      *engine.Chess
	database   *db.Database
	activeGame string

	chessSelected *engine.Position
)

func main() {
	fmt.Println("WASM Go Initialized")

	var err error
	database, err = db.NewDatabase("boardgames")
	if err != nil {
		fmt.Printf("Failed to init DB: %v\n", err)
	}

	game = engine.NewCheckers()
	slots = engine.NewSlotMachine()
	blackjack = engine.NewBlackjack()
	dominoes = engine.NewDominoes()
	chess = engine.NewChess()
	activeGame = "checkers"

	loadState()

	js.Global().Set("handleRequest", js.FuncOf(handleRequest))

	select {}
}

func loadState() {
	if database == nil {
		return
	}

	val, err := database.Load("activeGame")
	if err == nil && !val.IsNull() && !val.IsUndefined() {
		activeGame = val.String()
	}

	val, err = database.Load("checkers")
	if err == nil && !val.IsNull() && !val.IsUndefined() {
		json.Unmarshal([]byte(val.String()), &game)
	}

	val, err = database.Load("slots")
	if err == nil && !val.IsNull() && !val.IsUndefined() {
		json.Unmarshal([]byte(val.String()), &slots)
	}

	val, err = database.Load("blackjack")
	if err == nil && !val.IsNull() && !val.IsUndefined() {
		json.Unmarshal([]byte(val.String()), &blackjack)
	}

	val, err = database.Load("dominoes")
	if err == nil && !val.IsNull() && !val.IsUndefined() {
		json.Unmarshal([]byte(val.String()), &dominoes)
	}

	val, err = database.Load("chess")
	if err == nil && !val.IsNull() && !val.IsUndefined() {
		json.Unmarshal([]byte(val.String()), &chess)
	}
}

func saveState() {
	if database == nil {
		return
	}

	database.Save("activeGame", activeGame)

	checkersJSON, _ := json.Marshal(game)
	database.Save("checkers", string(checkersJSON))

	slotsJSON, _ := json.Marshal(slots)
	database.Save("slots", string(slotsJSON))

	blackjackJSON, _ := json.Marshal(blackjack)
	database.Save("blackjack", string(blackjackJSON))

	dominoesJSON, _ := json.Marshal(dominoes)
	database.Save("dominoes", string(dominoesJSON))

	chessJSON, _ := json.Marshal(chess)
	database.Save("chess", string(chessJSON))
}

func handleRequest(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		return "Invalid arguments"
	}

	method := args[0].String()
	pathURL := args[1].String()
	isHX := args[2].Bool()

	u, err := url.Parse(pathURL)
	if err != nil {
		return fmt.Sprintf("Error parsing URL: %v", err)
	}

	ctx := context.Background()
	var buf bytes.Buffer

	render := func(title string, component templ.Component) error {
		if !isHX {
			return templates.Layout(title, component).Render(ctx, &buf)
		}
		return component.Render(ctx, &buf)
	}

	switch u.Path {
	case "/":
		activeGame = "checkers"
		err = render("Checkers", templates.CheckersBoard(game))
	case "/slots":
		activeGame = "slots"
		params := u.Query()
		if params.Get("action") == "spin" {
			slots.Spin()
			saveState()
		}
		err = render("Slots", templates.SlotMachineGame(slots))
	case "/blackjack":
		activeGame = "blackjack"
		params := u.Query()
		action := params.Get("action")
		if action != "" {
			blackjack.HandleMove(ctx, params)
			saveState()
		}
		err = render("Blackjack", templates.BlackjackGame(blackjack))
	case "/dominoes":
		activeGame = "dominoes"
		params := u.Query()
		action := params.Get("action")
		if action == "play" {
			idx, _ := strconv.Atoi(params.Get("idx"))
			side := params.Get("side")
			dominoes.Play(0, idx, side)
			saveState()
		} else if len(dominoes.Hands[0]) == 0 && len(dominoes.Board) == 0 {
			dominoes.Deal(2)
		}
		err = render("Dominoes", templates.DominoesGame(dominoes))
	case "/chess":
		activeGame = "chess"
		err = render("Chess", templates.ChessBoard(chess))
	case "/chess/select":
		activeGame = "chess"
		params := u.Query()
		row, _ := strconv.Atoi(params.Get("row"))
		col, _ := strconv.Atoi(params.Get("col"))

		if chessSelected == nil {
			p := chess.Board[row][col]
			if p.Type != engine.None && p.Color == chess.CurrentPlayer {
				chessSelected = &engine.Position{Row: row, Col: col}
			}
		} else {
			if chessSelected.Row == row && chessSelected.Col == col {
				chessSelected = nil
			} else {
				moveErr := chess.Move(chessSelected.Row, chessSelected.Col, row, col)
				if moveErr == nil {
					chessSelected = nil
					saveState()
				} else {
					p := chess.Board[row][col]
					if p.Type != engine.None && p.Color == chess.CurrentPlayer {
						chessSelected = &engine.Position{Row: row, Col: col}
					}
				}
			}
		}
		err = render("Chess", templates.ChessBoard(chess))
	case "/select":
		activeGame = "checkers"
		params := u.Query()
		row, _ := strconv.Atoi(params.Get("row"))
		col, _ := strconv.Atoi(params.Get("col"))

		if game.SelectedPiece == nil {
			p := game.Board[row][col]
			if p != engine.Empty && ((game.CurrentPlayer == engine.Red && (p == engine.Red || p == engine.RedKing)) ||
				(game.CurrentPlayer == engine.Black && (p == engine.Black || p == engine.BlackKing))) {
				game.SelectedPiece = &engine.Position{Row: row, Col: col}
			}
		} else {
			if row == game.SelectedPiece.Row && col == game.SelectedPiece.Col {
				if !game.JumpActive {
					game.SelectedPiece = nil
				}
			} else {
				moveErr := game.Move(game.SelectedPiece.Row, game.SelectedPiece.Col, row, col)
				if moveErr != nil {
					p := game.Board[row][col]
					if !game.JumpActive && p != engine.Empty {
						game.SelectedPiece = &engine.Position{Row: row, Col: col}
					}
				} else {
					saveState()
				}
			}
		}
		err = templates.CheckersBoard(game).Render(ctx, &buf)
	case "/reset":
		if method == "POST" {
			switch activeGame {
			case "slots":
				slots.Reset()
			case "blackjack":
				blackjack.Reset()
			case "dominoes":
				dominoes.Reset()
				dominoes.Deal(2)
			case "chess":
				chess.Reset()
			default:
				game.Reset()
			}
			saveState()
			switch activeGame {
			case "slots":
				err = templates.SlotMachineGame(slots).Render(ctx, &buf)
			case "blackjack":
				err = templates.BlackjackGame(blackjack).Render(ctx, &buf)
			case "dominoes":
				err = templates.DominoesGame(dominoes).Render(ctx, &buf)
			case "chess":
				err = templates.ChessBoard(chess).Render(ctx, &buf)
			default:
				err = templates.CheckersBoard(game).Render(ctx, &buf)
			}
		}
	default:
		return "404 Not Found"
	}

	if err != nil {
		return fmt.Sprintf("Error rendering: %v", err)
	}

	return buf.String()
}
