package ui

import (
	"log"
	"time"

	"github.com/fspasovski/pocketstream-app/app"
	"github.com/fspasovski/pocketstream-app/input"
	"github.com/fspasovski/pocketstream-app/player"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	space     = "SPACE"
	enter     = "ENTER"
	backspace = "←"
)

type SearchScreen struct {
	SelectedKeyI int
	SelectedKeyJ int
	Input        string
	Keys         [][]string
	CaretVisible bool
	LastBlink    time.Time
	KeyRects     [][]sdl.Rect
	Shift        bool
	Player       *player.Player
}

func CreateSearchScreen(app *app.App, mediaPlayer *player.Player) *SearchScreen {
	searchState := &SearchScreen{
		SelectedKeyI: 0,
		SelectedKeyJ: 0,
		Input:        "",
		Keys:         getVirtualKeyboardKeys(),
		CaretVisible: false,
		LastBlink:    time.Time{},
		Shift:        false,
		Player:       mediaPlayer,
	}
	searchState.ComputeKeyRects(app)
	return searchState
}

func (s *SearchScreen) HandleInput(app *app.App, key input.Key) {
	switch key {
	case input.Up:
		s.handleKeyUp()
	case input.Down:
		s.handleKeyDown()
	case input.Left:
		s.handleKeyLeft()
	case input.Right:
		s.handleKeyRight()
	case input.A:
		s.handleKeyA(app)
	case input.B, input.X:
		s.handleKeyB(app)
	}
}

func (s *SearchScreen) handleKeyUp() {
	if s.SelectedKeyI-1 >= 0 {
		s.SelectedKeyI--
	}
}

func (s *SearchScreen) handleKeyDown() {
	if s.SelectedKeyI+1 < len(s.Keys) {
		s.SelectedKeyI++
	}
}

func (s *SearchScreen) handleKeyLeft() {
	if s.SelectedKeyJ == 0 {
		s.SelectedKeyJ = len(s.Keys[s.SelectedKeyI]) - 1
	} else {
		s.SelectedKeyJ--
	}
}

func (s *SearchScreen) handleKeyRight() {
	if s.SelectedKeyJ == len(s.Keys[s.SelectedKeyI])-1 {
		s.SelectedKeyJ = 0
	} else {
		s.SelectedKeyJ++
	}
}

func (s *SearchScreen) handleKeyA(app *app.App) {
	keyValue := s.Keys[s.SelectedKeyI][s.SelectedKeyJ]
	if keyValue == space {
		s.Input += " "
	} else if keyValue == enter {
		app.StartLoading("Searching streams...")
		go func() {
			streams, err := app.Config.TwitchService.SearchStreams(s.Input)
			if err != nil {
				log.Printf("An error occurred while fetching streams for: %s, %v", s.Input, err)
			} else {
				app.State = CreateSearchResultsScreen(streams, s.Player)
			}
			app.FinishLoading()
			app.NeedsRedraw = true
		}()
	} else if keyValue == backspace {
		if len(s.Input) > 0 {
			s.Input = s.Input[:len(s.Input)-1]
		}
	} else {
		s.Input += keyValue
	}
}

func (s *SearchScreen) handleKeyB(app *app.App) {
	app.State = CreateMainScreen(s.Player)
}

func (s *SearchScreen) Draw(app *app.App) {
	app.ClearScreen()

	if app.IsLoading {
		app.DrawLoadingScreen()
		return
	}

	drawInputBox(app, s)
	drawVirtualKeyboard(app, s)
}

func drawInputBox(app *app.App, s *SearchScreen) {
	inset := app.Config.UI.InputBoxPadding
	box := sdl.Rect{X: inset, Y: inset, W: app.Config.Display.Width - 2*inset, H: app.Config.UI.InputBoxHeight}

	// border
	app.DrawRect(&box, app.Config.UI.Colors.InputBoxBorderColor)

	// render text
	textX := box.X + 8
	textY := box.Y + (box.H-int32(app.Config.UI.FontSize))/2
	app.DrawText(s.Input, app.Config.UI.Colors.InputTextColor, textX, textY)
}

func drawVirtualKeyboard(app *app.App, s *SearchScreen) {
	for row := 0; row < len(s.Keys); row++ {
		for col := 0; col < len(s.Keys[row]); col++ {
			k := s.Keys[row][col]
			rect := s.KeyRects[row][col]
			selected := row == s.SelectedKeyI && col == s.SelectedKeyJ
			if selected {
				app.FillRect(&rect, app.Config.UI.Colors.SelectedKeyBackgroundColor)
				app.DrawRect(&rect, app.Config.UI.Colors.SelectedKeyBorderColor)
			} else {
				app.FillRect(&rect, app.Config.UI.Colors.KeyBackgroundColor)
				app.DrawRect(&rect, app.Config.UI.Colors.KeyBorderColor)
			}

			drawKey(app, k, &rect, selected)
		}
	}
}

func drawKey(app *app.App, text string, rect *sdl.Rect, selected bool) {
	if text == "" {
		return
	}

	color := app.Config.UI.Colors.KeyColor
	if selected {
		color = app.Config.UI.Colors.SelectedKeyColor
	}

	app.DrawCenteredTextInRect(text, rect, color)
}

func (s *SearchScreen) ComputeKeyRects(app *app.App) {
	s.KeyRects = make([][]sdl.Rect, len(s.Keys))
	y := app.Config.UI.KeyboardTopMargin + app.Config.UI.VirtualTopPadding

	for r := 0; r < len(s.Keys); r++ {
		row := s.Keys[r]
		s.KeyRects[r] = make([]sdl.Rect, len(row))

		// calculate row width to center it
		rowWidth := int32(0)
		for _, k := range row {
			kwidth := app.Config.UI.KeyWidth
			if k == enter || k == space {
				kwidth = app.Config.UI.KeyWidth * 2
			}
			rowWidth += kwidth + app.Config.UI.KeySpacingX
		}
		// remove last spacing
		rowWidth -= app.Config.UI.KeySpacingX

		x := (app.Config.Display.Width - rowWidth) / 2

		for c := 0; c < len(row); c++ {
			k := row[c]
			w := app.Config.UI.KeyWidth
			if k == enter || k == space {
				w = app.Config.UI.KeyWidth * 2
			}
			rect := sdl.Rect{X: x, Y: y, W: w, H: app.Config.UI.KeyHeight}
			s.KeyRects[r][c] = rect
			x += w + app.Config.UI.KeySpacingX
		}

		y += app.Config.UI.KeyHeight + app.Config.UI.KeySpacingY
	}
}

func getVirtualKeyboardKeys() [][]string {
	return [][]string{
		{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
		{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p"},
		{"a", "s", "d", "f", "g", "h", "j", "k", "l", backspace},
		{"z", "x", "c", "v", "b", "n", "m", "_", space, enter},
	}
}
