package uiFunctions

import (
	"fmt"
	"strconv"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/pksingh21/pfs/customTypesAndGlobalVariables"
)

func RenderCurrentView(items []customTypesAndGlobalVariable.FileItem, currentColumnWidthPercentage int, currentPointerLocation customTypesAndGlobalVariable.Location, currentColumnIndex int) {
	// get current window width and height
	termWidth, _ := termbox.Size()
	//
	currentColumnWidth := termWidth * currentColumnWidthPercentage / 100
	// clear the current column width
	for i := 0; i < currentColumnWidth; i++ {
		termbox.SetCell(i, 0, ' ', customTypesAndGlobalVariable.ColorRed, termbox.ColorDefault)
	}
	// get indices for items array in a for loop

	for i := 0; i < len(items); i++ {
		currentItem := items[i]
		if currentItem.RowNumber == currentPointerLocation.Y && currentItem.ColumnNumber == currentPointerLocation.X {
			currentItem.IsCurrentSelection = true
		} else {
			currentItem.IsCurrentSelection = false
		}
		distanceFromLeft := 0
		if currentColumnIndex == 1 {
			distanceFromLeft = 0
		} else if currentColumnIndex == 2 {
			distanceFromLeft = termWidth * 20 / 100
		} else if currentColumnIndex == 3 {
			distanceFromLeft = termWidth * 60 / 100
		}
		RenderCoordinates := customTypesAndGlobalVariable.Location{X: distanceFromLeft, Y: i}
		RenderItemv2(currentItem, currentColumnWidth, RenderCoordinates)
	}

}

func RenderItemv2(item customTypesAndGlobalVariable.FileItem, currentColumnWidth int, CoorindateToRender customTypesAndGlobalVariable.Location) {
	runes := convertLongRune(item.Name, currentColumnWidth)
	for i := 0; i < currentColumnWidth; i++ {
		termbox.SetCell(CoorindateToRender.X+i, CoorindateToRender.Y, ' ', item.Color, termbox.ColorDefault)
	}
	for i := 0; i < currentColumnWidth; i++ {
		if i < len(runes) {
			termbox.SetCell(CoorindateToRender.X+i, CoorindateToRender.Y, runes[i], item.Color, termbox.ColorDefault)
		}
	}
	if item.IsCurrentSelection {
		for i := 0; i < currentColumnWidth; i++ {
			if i < len(runes) {
				termbox.SetBg(CoorindateToRender.X+i, CoorindateToRender.Y, item.Color)
				termbox.SetFg(CoorindateToRender.X+i, CoorindateToRender.Y, termbox.ColorDefault)
			}
		}
	}
}
func RenderItem(item customTypesAndGlobalVariable.FileItem, x, y int, counter int, pointerLocation customTypesAndGlobalVariable.Location) {
	runes := convertLongRune(item.Name, customTypesAndGlobalVariable.ColumnWidth)
	for i := 0; i < customTypesAndGlobalVariable.ColumnWidth; i++ {
		termbox.SetCell(i+x, y, ' ', item.Color, termbox.ColorDefault)
	}
	// render counter variable used only for debugging purposes
	counterRune := []rune("")
	for i := 0; i < customTypesAndGlobalVariable.ColumnWidth-1; i++ {
		if i < len(counterRune) {
			termbox.SetCell(x+i, y, counterRune[i], item.Color, termbox.ColorDefault)
		}
	}
	// Display the row number
	if item.IsCurrentSelection {
		for i := 0; i < customTypesAndGlobalVariable.ColumnWidth-1; i++ {
			if i < len(runes) {
				termbox.SetCell(i+len(counterRune)+x, y, runes[i], item.Color, termbox.ColorDefault)
			} else if i == customTypesAndGlobalVariable.ColumnWidth-2 {
				termbox.SetCell(i+x, y, item.Icon, item.Color, termbox.ColorDefault)
			}
		}

		for i := 0; i < customTypesAndGlobalVariable.ColumnWidth; i++ {
			termbox.SetBg(x+i, y, item.Color)
			termbox.SetFg(x+i, y, termbox.ColorDefault)
		}
	} else {
		for i := 0; i < customTypesAndGlobalVariable.ColumnWidth-1; i++ {
			if i < len(runes) {
				termbox.SetCell(i+len(counterRune)+x, y, runes[i], item.Color, termbox.ColorDefault)
			} else if i == customTypesAndGlobalVariable.ColumnWidth-2 {
				termbox.SetCell(i+x, y, item.Icon, item.Color, termbox.ColorDefault)
			}
		}
	}
}

func Debugger(input interface{}) {
	// Convert the input to a string
	output := fmt.Sprintf("%v", input)
	println(output, "output")
	// Set the position for rendering the output
	x := 1
	y := 0

	// Render the output on the screen
	for _, char := range output {
		termbox.SetCell(x, y, char, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}

	termbox.Flush()
}

func RenderColumns(items []customTypesAndGlobalVariable.FileItem, columnCount, startX, termWidth, termHeight, visibleItems, scrollOffset int, counter int, pointerLocation customTypesAndGlobalVariable.Location) {
	if counter-scrollOffset >= visibleItems {
		scrollOffset = counter - visibleItems + 1
	}

	// Calculate the number of visible items based on the actual item count
	// calculate lenght of items and render it on terminal
	counterRune := []rune(strconv.Itoa(customTypesAndGlobalVariable.Counter))
	customTypesAndGlobalVariable.Counter++
	for i := 0; i < customTypesAndGlobalVariable.ColumnWidth; i++ {
		termbox.SetCell(i, 0, ' ', customTypesAndGlobalVariable.ColorRed, termbox.ColorDefault)
	}
	for i := 0; i < customTypesAndGlobalVariable.ColumnWidth; i++ {
		if i < len(counterRune) {
			termbox.SetCell(i, 0, counterRune[i], customTypesAndGlobalVariable.ColorRed, termbox.ColorDefault)
		}
	}

	actualVisibleItems := len(items) - scrollOffset
	if actualVisibleItems > visibleItems {
		actualVisibleItems = visibleItems
	}

	// Render the visible items
	for i := 0; i < actualVisibleItems; i++ {
		index := i + scrollOffset
		item := items[index]

		// Calculate row position based on index
		row := i + 1

		// Calculate column position based on index
		columnPosition := startX

		if item.RowNumber == pointerLocation.Y && item.ColumnNumber == pointerLocation.X {
			item.IsCurrentSelection = true
		} else {
			item.IsCurrentSelection = false
		}

		RenderItem(item, columnPosition, row, counter, pointerLocation)
	}

	// Erase the remaining vertical space
	for i := actualVisibleItems + 1; i <= visibleItems; i++ {
		for j := 0; j < termWidth; j++ {
			termbox.SetCell(j+startX, i, ' ', customTypesAndGlobalVariable.ColorDefault, customTypesAndGlobalVariable.ColorDefault)
		}
	}
}
func convertLongRune(longRune string, maxWidth int) []rune {
	runes := []rune(longRune)
	runeWidth := runewidth.StringWidth(longRune)

	if runeWidth <= maxWidth {
		return runes
	}

	truncatedWidth := maxWidth - runewidth.StringWidth(customTypesAndGlobalVariable.Ellipsis)
	truncatedRunes := []rune(runewidth.Truncate(string(runes), truncatedWidth, customTypesAndGlobalVariable.Ellipsis))
	return truncatedRunes
}
