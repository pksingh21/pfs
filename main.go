package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/pksingh21/pfs/customTypesAndGlobalVariables"
	"github.com/pksingh21/pfs/uiFunctions"
)

func main() {
	err := termbox.Init()
	customTypesAndGlobalVariable.CachedItemsMap = make(map[string][]customTypesAndGlobalVariable.FileItem)
	customTypesAndGlobalVariable.SelectedItemMap = make(map[string]int)
	customTypesAndGlobalVariable.ColumnMap = make(map[int]string)
	customTypesAndGlobalVariable.ColumnStartXMap = make(map[int]int)
	showHiddenFiles := false
	flag.BoolVar(&showHiddenFiles, "hidden", false, "Display hidden directories")
	flag.BoolVar(&showHiddenFiles, "nonHidden", false, "Display non-hidden directories")
	flag.Parse()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rand.New(source)
	termbox.SetInputMode(termbox.InputEsc)

	parentDir := "." // Replace with the desired parent directory path
	parentDir, _ = filepath.Abs(".")
	columnCount := len(customTypesAndGlobalVariable.ColumnColors) + 10
	totalWidth := columnCount*customTypesAndGlobalVariable.ColumnWidth + (columnCount-1)*customTypesAndGlobalVariable.ColumnSpacing
	termWidth, termHeight := termbox.Size()
	startX := 0

	visibleItems := termHeight - 1 // Subtract 1 for scrollbar
	scrollOffset := 0

	pointerLocation := customTypesAndGlobalVariable.Location{X: 2, Y: 1}
	customTypesAndGlobalVariable.SelectedItemMap[parentDir] = 1
	customTypesAndGlobalVariable.ColumnMap[2] = parentDir
	customTypesAndGlobalVariable.ColumnStartXMap[1] = 0
	items, err := getFileItems(getParent(parentDir), showHiddenFiles)
	if err != nil {
		panic(err)
	}
	uiFunctions.RenderCurrentView(items, customTypesAndGlobalVariable.FirstColumnWidth, pointerLocation, 1)
	items, err = getFileItems(parentDir, showHiddenFiles)
	if err != nil {
		panic(err)
	}
	uiFunctions.RenderCurrentView(items, customTypesAndGlobalVariable.SecondColumnWidth, pointerLocation, 2)
	uiFunctions.RenderCurrentView(items, customTypesAndGlobalVariable.ThridColumnWidth, pointerLocation, 3)
	// uiFunctions.RenderColumns(items, columnCount, startX, termWidth-customTypesAndGlobalVariable.ScrollbarWidth, termHeight, visibleItems, scrollOffset, customTypesAndGlobalVariable.Count, pointerLocation)

	termbox.Flush()
	currentColumnPointerAt := 2
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return
			case termbox.KeyArrowDown:
				if customTypesAndGlobalVariable.Count < len(items) {
					if customTypesAndGlobalVariable.Count-scrollOffset >= visibleItems {
						scrollOffset++
					}
					customTypesAndGlobalVariable.Count++
					pointerLocation.Y++
					customTypesAndGlobalVariable.SelectedItemMap[parentDir] = pointerLocation.Y
					// uiFunctions.RenderColumns(items, columnCount, startX, termWidth-customTypesAndGlobalVariable.ScrollbarWidth, termHeight, visibleItems, scrollOffset, customTypesAndGlobalVariable.Count, pointerLocation)
					// uiFunctions.Debugger(customTypesAndGlobalVariable.SelectedItemMap)
					// uiFunctions.Debugger(customTypesAndGlobalVariable.ColumnMap)
					uiFunctions.RenderCurrentView(items, getCurrentWidthForCoulumn(currentColumnPointerAt), pointerLocation, currentColumnPointerAt)
					termbox.Flush()
				}
			case termbox.KeyArrowUp:
				if pointerLocation.Y >= 1 {
					if scrollOffset > 0 {
						scrollOffset--
					}
					customTypesAndGlobalVariable.Count--
					pointerLocation.Y--
					customTypesAndGlobalVariable.SelectedItemMap[parentDir] = pointerLocation.Y
					// uiFunctions.RenderColumns(items, columnCount, startX, termWidth-customTypesAndGlobalVariable.ScrollbarWidth, termHeight, visibleItems, scrollOffset, customTypesAndGlobalVariable.Count, pointerLocation)
					uiFunctions.RenderCurrentView(items, getCurrentWidthForCoulumn(currentColumnPointerAt), pointerLocation, currentColumnPointerAt)
					termbox.Flush()
				}
			case termbox.KeyArrowLeft:
				if customTypesAndGlobalVariable.HorizontalCounter > 1 {
					// customTypesAndGlobalVariable.HorizontalCounter--
					handleLeftClick(&currentColumnPointerAt)
					pointerLocation.X--
					pointerLocation.Y = customTypesAndGlobalVariable.SelectedItemMap[customTypesAndGlobalVariable.ColumnMap[pointerLocation.X]]
					// update the parent directory to the previous directory
					parentDir = customTypesAndGlobalVariable.ColumnMap[pointerLocation.X]
					uiFunctions.RenderCurrentView(items, getCurrentWidthForCoulumn(3), pointerLocation, 3)
					items, _ = getFileItems(parentDir, showHiddenFiles)
					uiFunctions.RenderCurrentView(items, getCurrentWidthForCoulumn(1), pointerLocation, 1)
					uiFunctions.RenderCurrentView(items, getCurrentWidthForCoulumn(2), pointerLocation, 2)
					// startX = customTypesAndGlobalVariable.ColumnStartXMap[pointerLocation.X]
					// uiFunctions.RenderColumns(items, columnCount, startX, termWidth-customTypesAndGlobalVariable.ScrollbarWidth, termHeight, visibleItems, scrollOffset, customTypesAndGlobalVariable.Count, pointerLocation)
					termbox.Flush()
				}
			case termbox.KeyArrowRight:
				if customTypesAndGlobalVariable.HorizontalCounter < columnCount {
					// customTypesAndGlobalVariable.HorizontalCounter++
					handleRightClick(&currentColumnPointerAt)
					pointerLocation.Y = customTypesAndGlobalVariable.SelectedItemMap[customTypesAndGlobalVariable.ColumnMap[pointerLocation.X]]
					// pointerLocation.X++
					parentDir = items[pointerLocation.Y-1].FilePath
					uiFunctions.Debugger(pointerLocation)
					uiFunctions.RenderCurrentView(items, getCurrentWidthForCoulumn(1), pointerLocation, 1)
					items, _ = getFileItems(parentDir, showHiddenFiles)
					uiFunctions.RenderCurrentView(items, getCurrentWidthForCoulumn(2), pointerLocation, 2)
					uiFunctions.RenderCurrentView(items, getCurrentWidthForCoulumn(3), pointerLocation, 3)
					// startX = customTypesAndGlobalVariable.ColumnStartXMap[pointerLocation.X]
					// uiFunctions.RenderColumns(items, columnCount, startX, termWidth-customTypesAndGlobalVariable.ScrollbarWidth, termHeight, visibleItems, scrollOffset, customTypesAndGlobalVariable.Count, pointerLocation)
					termbox.Flush()
				}
			case termbox.KeyEnter:
				if customTypesAndGlobalVariable.Count <= len(items) {
					// Perform the action for the selected item
					// For example, navigate to the selected directory
					selectedItem := items[pointerLocation.Y-1]
					if selectedItem.IsDir {
						parentDir = selectedItem.FilePath
						items, _ = getFileItems(parentDir, showHiddenFiles)
						scrollOffset = 0
						customTypesAndGlobalVariable.SelectedItemMap[parentDir] = customTypesAndGlobalVariable.Count
						pointerLocation.X++
						pointerLocation.Y = 1
						// customTypesAndGlobalVariable.Count = 1
						customTypesAndGlobalVariable.ColumnMap[pointerLocation.X] = parentDir
						// startX = startX + customTypesAndGlobalVariable.ColumnWidth + customTypesAndGlobalVariable.ColumnSpacing
						customTypesAndGlobalVariable.ColumnStartXMap[pointerLocation.X] = startX
						// uiFunctions.RenderColumns(items, columnCount, startX, termWidth-customTypesAndGlobalVariable.ScrollbarWidth, termHeight, visibleItems, scrollOffset, customTypesAndGlobalVariable.Count, pointerLocation)

						uiFunctions.RenderCurrentView(items, customTypesAndGlobalVariable.FirstColumnWidth, pointerLocation, 1)
						termbox.Flush()
					}
				}
			}
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			termWidth, termHeight = termbox.Size()
			startX = (termWidth - totalWidth) / 2
			visibleItems = termHeight - 1 // Subtract 1 for scrollbar
			uiFunctions.RenderColumns(items, columnCount, startX, termWidth-customTypesAndGlobalVariable.ScrollbarWidth, termHeight, visibleItems, scrollOffset, customTypesAndGlobalVariable.Count, pointerLocation)
			termbox.Flush()
		}
	}
}

// define a function to get parent for a given directory and return the path to it
func getParent(directory string) string {
	if directory == "." {
		// If the directory is ".", return the absolute path of the current directory
		absPath, err := filepath.Abs(".")
		if err != nil {
			// Handle error if unable to get the absolute path
			fmt.Println("Error:", err)
			return ""
		}

		return getParent(absPath)
	}

	parent := filepath.Dir(directory)
	return parent
}

func getFileItems(parentDir string, showHidden bool) ([]customTypesAndGlobalVariable.FileItem, error) {
	fileInfo, err := os.Stat(parentDir)
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		// Return an empty list if the current path is not a directory
		return []customTypesAndGlobalVariable.FileItem{}, nil
	}

	cacheKey := fmt.Sprintf("%s_%t", parentDir, showHidden)
	cachedItems, exists := customTypesAndGlobalVariable.CachedItemsMap[cacheKey]
	if exists {
		return cachedItems, nil
	}

	files, err := os.ReadDir(parentDir)
	if err != nil {
		return nil, err
	}

	var items []customTypesAndGlobalVariable.FileItem
	nthitem := 0
	for _, file := range files {
		if !showHidden && file.Name()[0] == '.' {
			continue // Skip hidden directories if showHidden is false
		}
		item := customTypesAndGlobalVariable.FileItem{
			Name:         file.Name(),
			IsDir:        file.IsDir(),
			FilePath:     fmt.Sprintf("%s/%s", parentDir, file.Name()),
			RowNumber:    nthitem + 1,
			ColumnNumber: customTypesAndGlobalVariable.HorizontalCounter,
		}
		nthitem++
		if item.IsDir {
			item.Icon = '📁'
			item.Color = customTypesAndGlobalVariable.ColorBlue
		} else {
			item.Icon = '📄'
			item.Color = customTypesAndGlobalVariable.ColorGreen
		}

		items = append(items, item)
	}

	customTypesAndGlobalVariable.HorizontalCounter++
	customTypesAndGlobalVariable.CachedItemsMap[cacheKey] = items
	return items, nil
}
func handleLeftClick(currentColumnPointer *int) {
	if *currentColumnPointer == 1 {
		return
	}
	*currentColumnPointer--
}
func handleRightClick(currentColumnPointer *int) {
	if *currentColumnPointer == 3 {
		return
	}
	*currentColumnPointer++
}
func getCurrentWidthForCoulumn(columnIndex int) int {
	if columnIndex == 1 {
		return customTypesAndGlobalVariable.FirstColumnWidth
	} else if columnIndex == 2 {
		return customTypesAndGlobalVariable.SecondColumnWidth
	}
	return customTypesAndGlobalVariable.ThridColumnWidth
}
