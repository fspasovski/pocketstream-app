package ui

import (
	"log"
	"math"

	"github.com/fspasovski/pocketstream-app/app"
	"github.com/fspasovski/pocketstream-app/input"
	"github.com/fspasovski/pocketstream-app/model"
	"github.com/fspasovski/pocketstream-app/player"
)

type FavoriteBroadcastersScreen struct {
	SelectedStream int
	PageStartIndex int
	PageEndIndex   int
	Streams        []model.Stream
	Player         *player.Player
}

func CreateFavoriteBroadcastersScreen(app *app.App, mediaPlayer *player.Player) *FavoriteBroadcastersScreen {
	if app.UserDataManager.NoFavoriteBroadcasters() {
		return &FavoriteBroadcastersScreen{
			Streams:        make([]model.Stream, 0),
			PageStartIndex: 0,
			PageEndIndex:   0,
			Player:         mediaPlayer,
		}
	}

	favoriteBroadcasterLogins := make([]string, 0)

	for login := range app.UserDataManager.Data.FavoriteBroadcasters {
		favoriteBroadcasterLogins = append(favoriteBroadcasterLogins, login)
	}

	favoriteStreams := app.PocketstreamService.GetStreams(favoriteBroadcasterLogins)

	imageUrls := make([]string, 0)

	for _, stream := range favoriteStreams {
		imageUrls = append(imageUrls, stream.PreviewImageURL)
		imageUrls = append(imageUrls, app.UserDataManager.GetBroadcasterImageUrl(stream.Broadcaster.Login))
	}

	imageData := app.ImageDataService.GetImageData(imageUrls)
	result := make([]model.Stream, len(favoriteStreams))
	for i := 0; i < len(favoriteStreams); i++ {
		result[i] = favoriteStreams[i].WithImageData(
			imageData[favoriteStreams[i].PreviewImageURL],
			imageData[app.UserDataManager.GetBroadcasterImageUrl(favoriteStreams[i].Broadcaster.Login)],
		)
	}

	return &FavoriteBroadcastersScreen{
		Streams:        result,
		PageStartIndex: 0,
		PageEndIndex:   int(math.Min(float64(2), float64(len(favoriteStreams)-1))),
		Player:         mediaPlayer,
	}
}

func (s *FavoriteBroadcastersScreen) HandleInput(appState *app.App, key input.Key) {
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
		s.handleKeyX(appState)
	case input.Y:
		s.handleKeyY(appState)
	case input.Left:
		s.handleKeyRight(appState)
	}
}

func (s *FavoriteBroadcastersScreen) handleKeyY(appState *app.App) {
	if len(s.Streams) > 0 {
		appState.UserDataManager.ToggleFavoriteBroadcaster(s.Streams[s.SelectedStream].Broadcaster)
	}
}

func (s *FavoriteBroadcastersScreen) handleKeyB(app *app.App) {
	if s.Player.IsPlaying() {
		s.Player.Stop()
		app.FinishLoading()
		app.RaiseAppWindow()
	} else {
		app.State = CreateMainScreen(s.Player)
	}
}

func (s *FavoriteBroadcastersScreen) handleKeyA(app *app.App) {
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

func (s *FavoriteBroadcastersScreen) handleKeyDown() {
	if s.Player.IsPlaying() || s.SelectedStream >= len(s.Streams)-1 {
		return
	}

	s.SelectedStream++
	if s.SelectedStream > s.PageEndIndex {
		s.PageStartIndex++
		s.PageEndIndex = int(math.Min(float64(s.PageEndIndex+1), float64(len(s.Streams)-1)))
	}
}

func (s *FavoriteBroadcastersScreen) handleKeyUp() {
	if s.Player.IsPlaying() || s.SelectedStream <= 0 {
		return
	}

	s.SelectedStream--
	if s.SelectedStream < s.PageStartIndex {
		s.PageStartIndex--
		s.PageEndIndex--
	}
}

func (s *FavoriteBroadcastersScreen) Draw(app *app.App) {
	app.ClearScreen()

	if app.IsLoading {
		app.DrawLoadingScreen()
		return
	}

	DrawStreams(app, s.Streams, s.PageStartIndex, s.PageEndIndex, s.SelectedStream)
}

func (s *FavoriteBroadcastersScreen) handleKeyRight(app *app.App) {
	app.State = CreateMainScreen(s.Player)
}

func (s *FavoriteBroadcastersScreen) handleKeyX(appState *app.App) {
	if !s.Player.IsPlaying() {
		appState.State = CreateSearchScreen(appState, s.Player)
	}
}
