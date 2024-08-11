package tui

import (
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
)

type difficultyOption struct {
	name    string
	disarms int
	mines   int
}

type difficultyOptions []difficultyOption

func (s difficultyOptions) alignedStrings() []string {
	longestNameLength := 0
	longestDisarmNumberLength := 0
	longestMineNumberLength := 0
	for _, option := range s {
		if l := len(option.name); l > longestNameLength {
			longestNameLength = l
		}
		if l := int(math.Log10(float64(option.disarms))) + 1; l > longestDisarmNumberLength {
			longestDisarmNumberLength = l
		}
		if l := int(math.Log10(float64(option.mines))) + 1; l > longestMineNumberLength {
			longestMineNumberLength = l
		}
	}
	optionStrings := make([]string, len(s))
	for i, option := range s {
		optionStrings[i] = fmt.Sprintf("%-*s (%*d disarms, %*d mines)", longestNameLength,
			option.name, longestDisarmNumberLength, option.disarms, longestMineNumberLength,
			option.mines)
	}
	return optionStrings
}

type difficultyMenu struct {
	choices difficultyOptions
	cursor  int
	size    sizeOption
}

func (d difficultyMenu) Init() tea.Cmd {
	return nil
}

func (d difficultyMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return d, tea.Quit

		case "up", "k":
			if d.cursor > 0 {
				d.cursor--
			}
		case "down", "j":
			if d.cursor < len(d.choices)-1 {
				d.cursor++
			}

		case "enter":
			o := d.choices[d.cursor]
			return InitializeGameView(d.size.width, d.size.height, o.disarms, o.mines, true), nil
		}
	}

	return d, nil
}

func (d difficultyMenu) View() string {
	// The header
	str := fmt.Sprintf("Selected size: %s (%dx%d)\n\n", d.size.name, d.size.width, d.size.height)

	str += "Select the difficulty for your game:\n\n"

	// Iterate over our choices
	for i, choice := range d.choices.alignedStrings() {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if d.cursor == i {
			cursor = "■" // cursor! (U+25A0 Black Square)
		}

		// Render the row
		str += fmt.Sprintf("%s [%s]\n", choice, cursor)
	}

	// The footer
	str += "\n(Press q to quit)\n"

	// Send the UI for rendering
	return str
}

func InitializeDifficultyMenu(size sizeOption) difficultyMenu {
	numTiles := size.width * size.height
	return difficultyMenu{
		choices: difficultyOptions{
			difficultyOption{"Easy", numTiles / 36, numTiles / 12},
			difficultyOption{"Medium", numTiles / 24, numTiles / 8},
			difficultyOption{"Hard", numTiles / 15, numTiles / 5},
		},
		size: size,
	}
}
