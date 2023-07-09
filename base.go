package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mattn/go-runewidth"
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
	ellipsis    = "..."
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
	counter         = 1
)

type FileItem struct {
	Name               string
	IsDir              bool
	Icon               rune
	Color              termbox.Attribute
	FilePath           string
	IsCurrentSelection bool
	RowNumber          int
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	rand.Seed(time.Now().UnixNano())

	termbox.SetInputMode(termbox.InputEsc)

	parentDir := "." // Replace with the desired parent directory path

	items, err := getFileItems(parentDir)
	if err != nil {
		panic(err)
	}

	columnCount := len(columnColors)
	totalWidth := columnCount*columnWidth + (columnCount-1)*columnSpacing
	termWidth, termHeight := termbox.Size()
	startX := (termWidth - totalWidth) / 2

	visibleItems := termHeight - 1 // Subtract 1 for scrollbar
	scrollOffset := 0
	currentSelectionRow := 1

	renderColumns(items, columnCount, startX, termWidth-scrollbarWidth, termHeight, visibleItems, scrollOffset, counter)

	renderScrollbar(scrollbarWidth, termHeight, len(items), visibleItems, scrollOffset)

	termbox.Flush()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return
			case termbox.KeyArrowDown:
				if counter < len(items) {
					if counter-scrollOffset >= visibleItems {
						scrollOffset++
					}
					counter++
					renderColumns(items, columnCount, startX, termWidth-scrollbarWidth, termHeight, visibleItems, scrollOffset, counter)
					renderScrollbar(scrollbarWidth, termHeight, len(items), visibleItems, scrollOffset)
					termbox.Flush()
				}
			case termbox.KeyArrowUp:
				if counter > 1 {
					if scrollOffset > 0 {
						scrollOffset--
					}
					counter--
					renderColumns(items, columnCount, startX, termWidth-scrollbarWidth, termHeight, visibleItems, scrollOffset, counter)
					termbox.Flush()
				}
			case termbox.KeyEnter:
				if currentSelectionRow <= len(items) {
					// Perform the action for the selected item
					// For example, navigate to the selected directory
					selectedItem := items[currentSelectionRow-1]
					if selectedItem.IsDir {
						parentDir = selectedItem.FilePath
						items, _ = getFileItems(parentDir)
						currentSelectionRow = 1
						scrollOffset = 0
						counter = 1
						renderColumns(items, columnCount, startX, termWidth-scrollbarWidth, termHeight, visibleItems, scrollOffset, counter)
						renderScrollbar(scrollbarWidth, termHeight, len(items), visibleItems, scrollOffset)
						termbox.Flush()
					}
				}
			}
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			termWidth, termHeight = termbox.Size()
			startX = (termWidth - totalWidth) / 2
			visibleItems = termHeight - 1 // Subtract 1 for scrollbar
			renderColumns(items, columnCount, startX, termWidth-scrollbarWidth, termHeight, visibleItems, scrollOffset, counter)
			renderScrollbar(scrollbarWidth, termHeight, len(items), visibleItems, scrollOffset)
			termbox.Flush()
		}
	}
}

func convertLongRune(longRune string, maxWidth int) []rune {
	runes := []rune(longRune)
	runeWidth := runewidth.StringWidth(longRune)

	if runeWidth <= maxWidth {
		return runes
	}

	truncatedWidth := maxWidth - runewidth.StringWidth(ellipsis)
	truncatedRunes := []rune(runewidth.Truncate(string(runes), truncatedWidth, ellipsis))
	return truncatedRunes
}

func renderColumns(items []FileItem, columnCount, startX, termWidth, termHeight, visibleItems, scrollOffset int, counter int) {
	columnPositions = []int{}
	for i := 0; i < columnCount; i++ {
		pos := startX + i*(columnWidth+columnSpacing)
		columnPositions = append(columnPositions, pos)
	}

	if counter-scrollOffset >= visibleItems {
		scrollOffset = counter - visibleItems + 1
	}

	for i := 0; i < visibleItems; i++ {
		index := i + scrollOffset
		if index >= len(items) {
			break
		}

		item := items[index]

		// Calculate row position based on index
		row := i + 1

		// Calculate column position based on index
		columnIndex := i % columnCount
		columnPosition := columnPositions[columnIndex]

		if item.RowNumber == counter {
			item.IsCurrentSelection = true
		} else {
			item.IsCurrentSelection = false
		}

		renderItem(item, columnPosition, row, counter)
	}
}

func renderItem(item FileItem, x, y int, counter int) {
	runes := convertLongRune(item.Name, columnWidth)
	for i := 0; i < columnWidth; i++ {
		termbox.SetCell(i, y, ' ', item.Color, termbox.ColorDefault)
	}
	// render counter variable used only for debugging purposes
	// counterRune := []rune(strconv.Itoa(item.RowNumber) + "." + strconv.Itoa(counter))
	// for i := 0; i < columnWidth-1; i++ {
	// 	if i < len(counterRune) {
	// 		termbox.SetCell(i, y, counterRune[i], item.Color, termbox.ColorDefault)
	// 	}
	// }
	// Display the row number
	if item.IsCurrentSelection {
		for i := 0; i < columnWidth-1; i++ {
			if i < len(runes) {
				termbox.SetCell(i, y, runes[i], item.Color, termbox.ColorDefault)
			} else if i == columnWidth-2 {
				termbox.SetCell(i, y, item.Icon, item.Color, termbox.ColorDefault)
			}
		}

		for i := 0; i < columnWidth; i++ {
			termbox.SetBg(i, y, item.Color)
			termbox.SetFg(i, y, termbox.ColorDefault)
		}
	} else {
		for i := 0; i < columnWidth-1; i++ {
			if i < len(runes) {
				termbox.SetCell(i, y, runes[i], item.Color, termbox.ColorDefault)
			} else if i == columnWidth-2 {
				termbox.SetCell(i, y, item.Icon, item.Color, termbox.ColorDefault)
			}
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
			termbox.SetCell(width-scrollbarWidth, row, '▌', colorDefault, colorDefault)
		} else {
			termbox.SetCell(width-scrollbarWidth, row, ' ', colorDefault, colorDefault)
		}
	}
}

func getFileItems(parentDir string) ([]FileItem, error) {
	files, err := os.ReadDir(parentDir)
	if err != nil {
		return nil, err
	}

	var items []FileItem
	for i, file := range files {
		item := FileItem{
			Name:      file.Name(),
			IsDir:     file.IsDir(),
			FilePath:  fmt.Sprintf("%s/%s", parentDir, file.Name()),
			RowNumber: i + 1,
		}

		if item.IsDir {
			item.Icon = '📁'
			item.Color = colorBlue
		} else {
			item.Icon = '📄'
			item.Color = colorGreen
		}

		items = append(items, item)
	}

	return items, nil
}
