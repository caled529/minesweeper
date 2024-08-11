package tui

import (
	"fmt"
	"strings"

	"github.com/caled529/minesweeper/game"
	tea "github.com/charmbracelet/bubbletea"
)

type gameView struct {
	cursorX     int
	cursorY     int
	game        *game.Game
	initDisarms int
	initMines   int
	width       int
	height      int
}

func (g gameView) Init() tea.Cmd {
	return nil
}

func (g gameView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		msgString := msg.String()
		if g.game.State.GameOver {
			switch msgString {
			case "r":
				return InitializeGameView(g.width, g.height, g.initDisarms, g.initMines, true), nil
			case "R":
				return InitializeSizeMenu(), nil
			}
		} else {
			switch msgString {
			case "left", "h":
				if g.cursorX > 0 {
					g.cursorX--
				}
			case "down", "j":
				if g.cursorY < g.height-1 {
					g.cursorY++
				}
			case "up", "k":
				if g.cursorY > 0 {
					g.cursorY--
				}
			case "right", "l":
				if g.cursorX < g.width-1 {
					g.cursorX++
				}
			case "enter", ";":
				g.game.RevealAt(g.cursorX, g.cursorY)
			case "'":
				g.game.SmartReveal(g.cursorX, g.cursorY)
			case " ":
				g.game.ToggleFlaggedAt(g.cursorX, g.cursorY)
			}
		}
		switch msgString {
		case "ctrl+c", "q":
			return g, tea.Quit
		}
	}
	return g, nil
}

func (g gameView) View() string {
	s := "Minesweeper\n\n"
	s += fmt.Sprintf("Mines remaining: %d\n", g.game.UnflaggedMines())
	for i, line := range strings.Split(g.game.BoardString(), "\n") {
		runeline := []rune(line)
		if g.game.State.GameOver {
			if i == g.height/2 {
				l := len(runeline)
				// Puts the text "GAME OVER" in the middle of the line
				runeline = append(runeline[:max(0, (l-9)/2)],
					append([]rune("GAME OVER"),
						runeline[min(l, l-(l-9)/2)-1+l%2:]...,
					)...)
			}
		} else if i == g.cursorY {
			runeline[g.cursorX] = '■'
		}
		s += fmt.Sprintf("%s\n", string(runeline))
	}
	if g.game.State.GameOver {
		s += "\n(Press r to restart)\t (Press R to return to the menu)"
	}
	s += "\n(Press q to quit)\n"
	return s
}

func InitializeGameView(boardWidth, boardHeight, numDisarms, numMines int, chainRevealing bool) gameView {
	return gameView{
		initDisarms: numDisarms,
		initMines:   numMines,
		game:        game.NewGame(boardWidth, boardHeight, numDisarms, numMines, chainRevealing),
		width:       boardWidth,
		height:      boardHeight,
	}
}
