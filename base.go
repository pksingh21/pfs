package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	// Terminal colors
	colorDefault = termbox.ColorDefault
	colorRed     = termbox.ColorRed
	colorGreen   = termbox.ColorGreen
	colorYellow  = termbox.ColorYellow
	colorBlue    = termbox.ColorBlue
	colorMagenta = termbox.ColorMagenta
	colorCyan    = termbox.ColorCyan
	colorWhite   = termbox.ColorWhite

	// Column widths
	columnWidth = 20

	// Spacing between columns
	columnSpacing = 4

	// Scrollbar
	scrollbarWidth = 1
)

var (
	columnColors = []termbox.Attribute{
		colorRed,
		colorGreen,
		colorYellow,
		colorBlue,
		colorMagenta,
		colorCyan,
		colorWhite,
	}

	columnPositions = []int{}
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// rand.Seed(time.Now().UnixNano())
	rand.NewSource(time.Now().UnixNano())
	termbox.SetInputMode(termbox.InputEsc)

	parentDir := "." // Replace with the desired parent directory path

	directories, err := getDirectories(parentDir)
	if err != nil {
		panic(err)
	}

	columnCount := len(columnColors)
	totalWidth := columnCount*columnWidth + (columnCount-1)*columnSpacing
	termWidth, termHeight := termbox.Size()
	startX := (termWidth - totalWidth) / 2

	visibleDirectories := termHeight - 1 // Subtract 1 for scrollbar
	scrollOffset := 0

	renderColumns(directories, columnCount, startX, termWidth-scrollbarWidth, termHeight, visibleDirectories, scrollOffset)

	renderScrollbar(scrollbarWidth, termHeight, len(directories), visibleDirectories, scrollOffset)

	termbox.Flush()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return
			case termbox.KeyArrowDown:
				if scrollOffset < len(directories)-visibleDirectories {
					scrollOffset++
					renderColumns(directories, columnCount, startX, termWidth-scrollbarWidth, termHeight, visibleDirectories, scrollOffset)
					renderScrollbar(scrollbarWidth, termHeight, len(directories), visibleDirectories, scrollOffset)
					termbox.Flush()
				}
			case termbox.KeyArrowUp:
				if scrollOffset > 0 {
					scrollOffset--
					renderColumns(directories, columnCount, startX, termWidth-scrollbarWidth, termHeight, visibleDirectories, scrollOffset)
					renderScrollbar(scrollbarWidth, termHeight, len(directories), visibleDirectories, scrollOffset)
					termbox.Flush()
				}
			}
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			termWidth, termHeight = termbox.Size()
			startX = (termWidth - totalWidth) / 2
			visibleDirectories = termHeight - 1 // Subtract 1 for scrollbar
			renderColumns(directories, columnCount, startX, termWidth-scrollbarWidth, termHeight, visibleDirectories, scrollOffset)
			renderScrollbar(scrollbarWidth, termHeight, len(directories), visibleDirectories, scrollOffset)
			termbox.Flush()
		}
	}
}

func renderColumns(directories []string, columnCount, startX, termWidth, termHeight, visibleDirectories, scrollOffset int) {
	columnPositions = []int{}
	for i := 0; i < columnCount; i++ {
		pos := startX + i*(columnWidth+columnSpacing)
		columnPositions = append(columnPositions, pos)
	}

	for i := 0; i < visibleDirectories; i++ {
		index := i + scrollOffset
		if index >= len(directories) {
			break
		}

		dir := directories[index]
		color := columnColors[i%len(columnColors)]

		// Calculate row position based on index
		row := i + 1

		// Calculate column position based on index
		columnIndex := i % columnCount
		columnPosition := columnPositions[columnIndex]

		renderColumn(dir, columnPosition, row, color)
	}
}

func renderColumn(dir string, x, y int, color termbox.Attribute) {
	runes := []rune(dir)

	for i := 0; i < columnWidth; i++ {
		if i < len(runes) {
			termbox.SetCell(i, y, runes[i], color, termbox.ColorDefault)
		} else {
			termbox.SetCell(i, y, ' ', color, termbox.ColorDefault)
		}
	}
}

func renderScrollbar(scrollbarWidth, termHeight, totalItems, visibleItems, scrollOffset int) {
	// Calculate scrollbar dimensions
	scrollbarHeight := termHeight - 1
	scrollbarPos := (scrollbarHeight * scrollOffset) / totalItems
	scrollbarSize := (scrollbarHeight * visibleItems) / totalItems

	// Render the scrollbar
	for row := 0; row < termHeight; row++ {
		width, _ := termbox.Size()
		termbox.SetCell(width-scrollbarWidth, row, ' ', colorDefault, colorDefault)
		if row >= scrollbarPos && row < scrollbarPos+scrollbarSize {
		} else {
			termbox.SetCell(width-scrollbarWidth, row, ' ', colorDefault, colorDefault)
		}
	}
}

func getDirectories(parentDir string) ([]string, error) {
	files, err := os.ReadDir(parentDir)
	if err != nil {
		return nil, err
	}

	var directories []string
	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file.Name())
		}
	}

	return directories, nil
}
