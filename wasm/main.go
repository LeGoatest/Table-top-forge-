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
	game      *engine.Checkers
	slots     *engine.SlotMachine
	database  *db.Database
	activeGame string
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
			if activeGame == "slots" {
				slots.Reset()
			} else {
				game.Reset()
			}
			saveState()
			if activeGame == "slots" {
				err = templates.SlotMachineGame(slots).Render(ctx, &buf)
			} else {
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
