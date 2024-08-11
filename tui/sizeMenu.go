package tui

import (
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
)

type sizeOption struct {
	name   string
	width  int
	height int
}

type sizeOptions []sizeOption

func (s sizeOptions) alignedStrings() []string {
	longestNameLength := 0
	longestDimensionLength := 0
	for _, option := range s {
		if l := len(option.name); l > longestNameLength {
			longestNameLength = l
		}
		if l := int(math.Floor(math.Log10(float64(option.width)))) + 1; l > longestDimensionLength {
			longestDimensionLength = l
		}
		if l := int(math.Floor(math.Log10(float64(option.height)))) + 1; l > longestDimensionLength {
			longestDimensionLength = l
		}
	}
	optionStrings := make([]string, len(s))
	for i, option := range s {
		optionStrings[i] = fmt.Sprintf("%-*s (%0*dx%0*d)", longestNameLength, option.name,
			longestDimensionLength, option.width, longestDimensionLength, option.height)
	}
	return optionStrings
}

type sizeMenu struct {
	choices sizeOptions
	cursor  int
}

func (s sizeMenu) Init() tea.Cmd {
	return nil
}

func (s sizeMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return s, tea.Quit

		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
		case "down", "j":
			if s.cursor < len(s.choices)-1 {
				s.cursor++
			}

		case "enter":
			return InitializeDifficultyMenu(s.choices[s.cursor]), nil
		}
	}

	return s, nil
}

func (s sizeMenu) View() string {
	// The header
	str := "Select a size for your game:\n\n"

	// Iterate over our choices
	for i, choice := range s.choices.alignedStrings() {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if s.cursor == i {
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

func InitializeSizeMenu() sizeMenu {
	return sizeMenu{
		choices: sizeOptions{
			sizeOption{
				name:   "Small",
				width:  15,
				height: 5,
			},
			sizeOption{
				name:   "Medium",
				width:  30,
				height: 10,
			},
			sizeOption{
				name:   "Large",
				width:  60,
				height: 20,
			},
		},
	}
}
