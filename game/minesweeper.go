package game

import "math/rand/v2"

type tile struct {
	adjacentMines int
	mine          bool
	flagged       bool
	revealed      bool
}

type board [][]tile

func newTileGrid(width, height, numMines int) *board {
	grid := make(board, width)
	for i := range grid {
		grid[i] = make([]tile, height)
	}
	grid.propagateMines(numMines)
	grid.propagateMineAdjacency()
	return &grid
}

// Randomly places n mines on the board
func (b *board) propagateMines(n int) {
	minesPlaced := 0
	for minesPlaced < min(n, len(*b)*len((*b)[0])) {
		minePosition := [2]int{rand.IntN(len(*b)), rand.IntN(len((*b)[0]))}
		if (*b)[minePosition[0]][minePosition[1]].mine == false {
			(*b)[minePosition[0]][minePosition[1]].mine = true
			minesPlaced++
		}
	}
}

// Figure out how many mines surround each safe tile and write that value to it
func (b *board) propagateMineAdjacency() {
	for i := range *b {
		for j := range (*b)[i] {
			if (*b)[i][j].mine {
				continue
			}
			(*b)[i][j].adjacentMines = b.minesAdjacentTo(i, j)
		}
	}
}

// Count the number of mines around a given position (x, y) on the board
func (b *board) minesAdjacentTo(x, y int) int {
	mineCount := 0
	for i := max(0, x-1); i < min(len(*b), x+2); i++ {
		for j := max(0, y-1); j < min(len((*b)[i]), y+2); j++ {
			if (*b)[i][j].mine {
				mineCount++
			}
		}
	}
	return mineCount
}

func (b *board) flagsAdjacentTo(x, y int) int {
	adjacentFlags := 0
	for i := max(0, x-1); i < min(len(*b), x+2); i++ {
		for j := max(0, y-1); j < min(len((*b)[i]), y+2); j++ {
			if (*b)[i][j].flagged {
				adjacentFlags++
			}
		}
	}
	return adjacentFlags
}

func (b *board) String() string {
	const ASCII_DIGITS_OFFSET = 48
	cols := len(*b)
	rows := len((*b)[0])
	skew := cols + 1
	boardRunes := make([]rune, skew*rows)
	for i, col := range *b {
		for j, curTile := range col {
			var tileChar rune
			if curTile.flagged {
				tileChar = '⚑' // U+2691 Black Flag
			} else if !curTile.revealed {
				tileChar = '█' // U+2588 Full Block
			} else if curTile.mine {
				tileChar = 'X'
			} else if curTile.adjacentMines == 0 {
				tileChar = ' '
			} else {
				tileChar = rune(curTile.adjacentMines + ASCII_DIGITS_OFFSET)
			}
			boardRunes[i+j*skew] = tileChar
		}
	}
	for i := range rows - 1 {
		boardRunes[i*skew+cols] = '\n'
	}
	return string(boardRunes)
}

type Game struct {
	State struct {
		GameOver bool
		disarms  int
		flags    int
		mines    int
	}
	board   board
	options struct {
		chainRevealing bool
	}
}

func NewGame(boardWidth, boardHeight, numDisarms, numMines int, chainRevealing bool) *Game {
	game := Game{
		State: struct {
			GameOver bool
			disarms  int
			flags    int
			mines    int
		}{
			disarms: numDisarms,
			mines:   numMines,
		},
		board: *newTileGrid(boardWidth, boardHeight, numMines),
		options: struct {
			chainRevealing bool
		}{
			chainRevealing: chainRevealing,
		},
	}
	return &game
}

// Reveals the tile at a given position (x, y)
//
// If chain is true, the block of adjacent safe tiles containing this tile is
// fully revealed
//
// Sets the game over state to true if the revealed tile was a mine
func (g *Game) RevealAt(x, y int) {
	if g.board[x][y].flagged {
		return
	}
	if g.State.disarms > 0 {
		g.disarmAt(x, y)
		g.State.disarms--
	}
	g.board[x][y].revealed = true
	if g.board[x][y].adjacentMines == 0 && g.options.chainRevealing {
		g.revealChain(x, y)
	}
	if g.board[x][y].mine {
		g.State.GameOver = true
	}
}

// Turns a mine into a safe tile
//
// For use as a mercy system at the start of games to prevent instant losses
func (g *Game) disarmAt(x, y int) {
	if !g.board[x][y].mine {
		return
	}
	g.board[x][y].mine = false
	g.State.mines--
	for i := max(0, x-1); i < min(len(g.board), x+2); i++ {
		for j := max(0, y-1); j < min(len(g.board[i]), y+2); j++ {
			if g.board[i][j].adjacentMines > 0 {
				g.board[i][j].adjacentMines--
			}
		}
	}
	g.board[x][y].adjacentMines = g.board.minesAdjacentTo(x, y)
}

// Recursively reveals blocks of safe tiles
func (g *Game) revealChain(x, y int) {
	for i := max(0, x-1); i < (min(len(g.board), x+2)); i++ {
		for j := max(0, y-1); j < (min(len(g.board[0]), y+2)); j++ {
			if !g.board[i][j].revealed {
				g.RevealAt(i, j)
			}
		}
	}
}

func (g *Game) SmartReveal(x, y int) {
	if !g.board[x][y].revealed || g.board[x][y].adjacentMines != g.board.flagsAdjacentTo(x, y) {
		return
	}
	for i := max(0, x-1); i < min(len(g.board), x+2); i++ {
		for j := max(0, y-1); j < min(len(g.board[i]), y+2); j++ {
			if !g.board[i][j].flagged {
				g.RevealAt(i, j)
			}
		}
	}
}

func (g *Game) ToggleFlaggedAt(x, y int) {
	if g.State.flags > g.State.mines {
		return
	}
	if !g.board[x][y].revealed {
		if g.board[x][y].flagged {
			g.State.flags--
		} else {
			g.State.flags++
		}
		g.board[x][y].flagged = !g.board[x][y].flagged
	}
}

func (g *Game) BoardDimensions() (int, int) {
	return len(g.board), len(g.board[0])
}

func (g *Game) BoardString() string {
	return g.board.String()
}

func (g *Game) UnflaggedMines() int {
	return g.State.mines - g.State.flags
}
