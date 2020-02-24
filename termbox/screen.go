package termbox

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"os"
	"strings"
)

type Termbox struct {
	Width                int
	Height               int
	SearchFormHeight     int
	CommandHistoryHeight int
	SelectionFile        string
	SelectionIndex       int
	FileList             []string
	FilteredFileList     []string
	Filter
}

type Filter struct {
	SearchQuery []rune
}

func NewFilter() Filter {
	return Filter{}
}

func (f *Filter) Append(c rune) {
	f.SearchQuery = append(f.SearchQuery, c)
}

func NewTermbox(fileList []string) *Termbox {
	return &Termbox{
		Filter:   NewFilter(),
		FileList: fileList,
	}
}

func (t *Termbox) Init() error {
	err := termbox.Init()
	if err != nil {
		return fmt.Errorf("init termbox error")
	}

	return nil
}

func (t *Termbox) Print(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		width := runewidth.RuneWidth(c)
		x = x + width
	}
}

func (t *Termbox) PrintSearchQuery(x, y int, fg, bg termbox.Attribute) {
	for _, c := range t.Filter.SearchQuery {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (f *Filter) FilterResult(commandHistory []string) []string {
	var filterResult []string
	for _, command := range commandHistory {
		filterIndex := f.getIndex(command)
		if filterIndex == -1 {
			continue
		}
		filterResult = append(filterResult, command)
	}
	return filterResult
}

func (f *Filter) getIndex(command string) int {
	SearchQueryToString := string(f.SearchQuery)
	index := strings.Index(command, SearchQueryToString)
	return index
}

func (t *Termbox) Draw() {
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)
	displaySearchForm := "QUERY>"
	t.Print(0, 0, termbox.ColorWhite, termbox.ColorDefault, displaySearchForm)
	t.PrintSearchQuery(len(displaySearchForm), 0, termbox.ColorWhite, termbox.ColorDefault)
	t.FilteredFileList = t.FileList

	if len(t.Filter.SearchQuery) != 0 {
		t.FilteredFileList = t.Filter.FilterResult(t.FilteredFileList)
		if len(t.FilteredFileList) > 0 {
			//t.SelectionFile = commandHistory[0]
			t.SelectionFile = t.FilteredFileList[t.SelectionIndex]
		} else {
			t.SelectionFile = ""
		}
	} else {
		t.SelectionFile = t.FileList[t.SelectionIndex]
	}

	// display file name in order
	for i := 0; i < len(t.FilteredFileList); i++ {
		//cursor := "\U0001F449 "
		cursor := "> "
		command := t.FilteredFileList[i]
		// if selection file then display cursor
		if i == t.SelectionIndex {
			command = cursor + command
		}
		t.Print(0, i+1, termbox.ColorWhite, termbox.ColorDefault, command)
	}

	termbox.Flush()
}

func (t *Termbox) Display() {
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)
	t.Draw()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if t.SelectionIndex != 0 {
					t.SelectionIndex--
				}
			case termbox.KeyArrowDown:
				if t.SelectionIndex < len(t.FilteredFileList)-1 {
					t.SelectionIndex++
				}
			case termbox.KeyCtrlC:
				os.Exit(0)
			case termbox.KeyEsc:
				break loop
			case termbox.KeySpace:
			case termbox.KeyCtrlD:
			case termbox.KeyBackspace2:
				if len(t.Filter.SearchQuery) > 0 {
					t.Filter.SearchQuery = t.Filter.SearchQuery[:len(t.Filter.SearchQuery)-1]
				}
			case termbox.KeyEnter:
				if len(t.SelectionFile) == 0 {
					continue
				}
				break loop
			default:
				if ev.Ch != 0 {
					t.SelectionIndex = 0
					t.Filter.Append(ev.Ch)
				}
			}
		}
		t.Draw()
	}
}
