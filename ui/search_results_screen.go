package ui

import (
	"log"
	"math"

	"github.com/fspasovski/pocketstream-app/app"
	"github.com/fspasovski/pocketstream-app/input"
	"github.com/fspasovski/pocketstream-app/model"
	"github.com/fspasovski/pocketstream-app/player"
)

type SearchResultsScreen struct {
	SelectedStream int
	PageStartIndex int
	PageEndIndex   int
	Streams        []model.Stream
	Player         *player.Player
}

func CreateSearchResultsScreen(streams []model.Stream, mediaPlayer *player.Player) *SearchResultsScreen {
	return &SearchResultsScreen{
		Streams:        streams,
		PageStartIndex: 0,
		PageEndIndex:   int(math.Min(float64(2), float64(len(streams)-1))),
		Player:         mediaPlayer,
	}
}

func (s *SearchResultsScreen) HandleInput(appState *app.App, key input.Key) {
	switch key {
	case input.Up:
		s.handleKeyUp()
	case input.Down:
		s.handleKeyDown()
	case input.A:
		s.handleKeyA(appState)
	case input.B:
		s.handleKeyB(appState)
	case input.X:
		appState.State = s.handleKeyX(appState)
	}
}

func (s *SearchResultsScreen) handleKeyX(appState *app.App) *SearchScreen {
	return CreateSearchScreen(appState, s.Player)
}

func (s *SearchResultsScreen) handleKeyB(app *app.App) {
	if s.Player.IsPlaying() {
		s.Player.Stop()
		app.FinishLoading()
		app.RaiseAppWindow()
	} else {
		app.State = CreateMainScreen(s.Player)
	}
}

func (s *SearchResultsScreen) handleKeyA(app *app.App) {
	if s.Player.IsPlaying() {
		return
	}

	app.StartLoading("Loading " + s.Streams[s.SelectedStream].Broadcaster.Login + " stream...")
	go func() {
		err := s.Player.Play(s.Streams[s.SelectedStream].Broadcaster.Login)
		if err != nil {
			log.Printf("An error occurred while playing stream: %v", err)
			app.FinishLoading()
		}
	}()
}

func (s *SearchResultsScreen) handleKeyDown() {
	if s.Player.IsPlaying() || s.SelectedStream >= len(s.Streams)-1 {
		return
	}

	s.SelectedStream++
	if s.SelectedStream > s.PageEndIndex {
		s.PageStartIndex++
		s.PageEndIndex = int(math.Min(float64(s.PageEndIndex+1), float64(len(s.Streams)-1)))
	}
}

func (s *SearchResultsScreen) handleKeyUp() {
	if s.Player.IsPlaying() || s.SelectedStream <= 0 {
		return
	}

	s.SelectedStream--
	if s.SelectedStream < s.PageStartIndex {
		s.PageStartIndex--
		s.PageEndIndex--
	}
}

func (s *SearchResultsScreen) Draw(app *app.App) {
	app.ClearScreen()

	if app.IsLoading {
		app.DrawLoadingScreen()
		return
	}

	DrawStreams(app, s.Streams, s.PageStartIndex, s.PageEndIndex, s.SelectedStream)
}
