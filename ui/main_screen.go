package ui

import (
	"log"
	"math"

	"github.com/fspasovski/pocketstream-app/app"
	"github.com/fspasovski/pocketstream-app/input"
	"github.com/fspasovski/pocketstream-app/player"
)

type MainScreen struct {
	PageStartIndex int
	PageEndIndex   int
	SelectedStream int
	Player         *player.Player
}

func CreateMainScreen(mediaPlayer *player.Player) *MainScreen {
	return &MainScreen{SelectedStream: 0, PageStartIndex: 0, PageEndIndex: 2, Player: mediaPlayer}
}

func (s *MainScreen) HandleInput(appState *app.App, key input.Key) {
	switch key {
	case input.Up:
		s.handleKeyUp()
	case input.Down:
		s.handleKeyDown(appState)
	case input.A:
		s.handleKeyA(appState)
	case input.B:
		s.handleKeyB(appState)
	case input.X:
		s.handleKeyX(appState)
	}
}

func (s *MainScreen) handleKeyUp() {
	if s.Player.IsPlaying() || s.SelectedStream <= 0 {
		return
	}

	s.SelectedStream--
	if s.SelectedStream < s.PageStartIndex {
		s.PageStartIndex--
		s.PageEndIndex--
	}
}

func (s *MainScreen) handleKeyDown(appState *app.App) {
	if s.Player.IsPlaying() || s.SelectedStream >= len(appState.TopStreams)-1 {
		return
	}

	s.SelectedStream++
	if s.SelectedStream > s.PageEndIndex {
		s.PageStartIndex++
		s.PageEndIndex = int(math.Min(float64(s.PageEndIndex+1), float64(len(appState.TopStreams)-1)))
	}
}

func (s *MainScreen) handleKeyA(app *app.App) {
	if s.Player.IsPlaying() {
		return
	}

	app.StartLoading("Loading " + app.TopStreams[s.SelectedStream].Broadcaster.Login + " stream...")
	go func() {
		err := s.Player.Play(app.TopStreams[s.SelectedStream].Broadcaster.Login)
		if err != nil {
			log.Printf("An error occurred while playing stream: %v", err)
			app.FinishLoading()
		}
	}()
}

func (s *MainScreen) handleKeyB(app *app.App) {
	if s.Player.IsPlaying() {
		s.Player.Stop()
		app.FinishLoading()
		app.RaiseAppWindow()
	} else {
		app.Running = false
	}
}

func (s *MainScreen) handleKeyX(appState *app.App) {
	if !s.Player.IsPlaying() {
		appState.State = CreateSearchScreen(appState, s.Player)
	}
}

func (s *MainScreen) Draw(app *app.App) {
	app.ClearScreen()

	if app.IsLoading {
		app.DrawLoadingScreen()
		return
	}

	DrawStreams(app, app.TopStreams, s.PageStartIndex, s.PageEndIndex, s.SelectedStream)
}
