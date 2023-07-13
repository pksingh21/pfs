package customTypesAndGlobalVariable

import "github.com/nsf/termbox-go"

type FileItem struct {
	Name               string
	IsDir              bool
	Icon               rune
	Color              termbox.Attribute
	FilePath           string
	IsCurrentSelection bool
	RowNumber          int
	ColumnNumber       int
}
type Location struct {
	X int
	Y int
}

const (
	// Terminal colors
	ColorDefault = termbox.ColorDefault
	ColorRed     = termbox.ColorRed
	ColorGreen   = termbox.ColorGreen
	ColorYellow  = termbox.ColorYellow
	ColorBlue    = termbox.ColorBlue
	ColorMagenta = termbox.ColorMagenta
	ColorCyan    = termbox.ColorCyan
	ColorWhite   = termbox.ColorWhite

	// Column widths
	ColumnWidth = 25
	Ellipsis    = "..."
	// Spacing between columns
	ColumnSpacing = 4

	// Scrollbar
	ScrollbarWidth              = 1
	FirstColumnWidth            = 20 // represents the width of current parent directory's column
	SecondColumnWidth           = 40 // represent the current parent directory being rendered
	ThridColumnWidth            = 40 // represents the child of current directory under selection
)

var (
	ColumnColors = []termbox.Attribute{
		ColorRed,
		ColorGreen,
		ColorYellow,
		ColorBlue,
		ColorMagenta,
		ColorCyan,
		ColorWhite,
	}

	Counter           = 1
	HorizontalCounter = 1
	// maintains cache of the current directories which have been accessed
	CachedItemsMap map[string][]FileItem
	// maintain a hashmap of the current directory and the index of the selected item
	// so that when the user navigates back to the directory, the selected item is
	// highlighted
	SelectedItemMap map[string]int
	// create a hashmap to maintain the column number and the directory path
	ColumnMap map[int]string
	// create a hashmap to maintain the column number and the startX position
	ColumnStartXMap map[int]int
	Count           = 0
)
