package config

import (
	"net/http"

	"github.com/fspasovski/pocketstream-app/twitch"
	"github.com/veandco/go-sdl2/sdl"
)

type Config struct {
	AppName       string
	AppVersion    string
	Display       DisplayConfig
	TwitchService *twitch.TwitchService
	UI            UIConfig
	Player        PlayerConfig
}

type DisplayConfig struct {
	Width  int32
	Height int32
}

type Colors struct {
	BackgroundColor                sdl.Color
	HeaderBackgroundColor          sdl.Color
	HeaderTextColor                sdl.Color
	FooterBackgroundColor          sdl.Color
	FooterTextColor                sdl.Color
	InputBoxBorderColor            sdl.Color
	InputBoxBackgroundColor        sdl.Color
	InputTextColor                 sdl.Color
	LoadingTextColor               sdl.Color
	NoResultsTextColor             sdl.Color
	KeyBackgroundColor             sdl.Color
	KeyBorderColor                 sdl.Color
	SelectedKeyBackgroundColor     sdl.Color
	SelectedKeyBorderColor         sdl.Color
	KeyColor                       sdl.Color
	SelectedKeyColor               sdl.Color
	StreamThumbnailBackgroundColor sdl.Color
	StreamLiveBadgeBackgroundColor sdl.Color
	LiveTextColor                  sdl.Color
	ViewersCountBackgroundColor    sdl.Color
	ViewersCountTextColor          sdl.Color
	ProfilePictureBackgroundColor  sdl.Color
	StreamerNameTextColor          sdl.Color
	StreamTitleColor               sdl.Color
	SelectedStreamBorderColor      sdl.Color
}

type StreamsUiConfig struct {
	ThumbnailWidth      int32
	ThumbnailHeight     int32
	ProfilePictureSize  int32
	InfoTextX           int32
	Padding             int32
	RowHeight           int32
	LiveBadgeWidth      int32
	LiveBadgeHeight     int32
	LiveBadgeLeftMargin int32
	LiveBadgeTopMargin  int32
	LiveTextLeftMargin  int32
	LiveTextTopMargin   int32
}

type UIConfig struct {
	HeaderHeight      int32
	FooterHeight      int32
	RowHeight         int32
	FontPath          string
	FontSize          int
	Padding           int32
	StreamsTopMargin  int32
	StreamLeftMargin  int32
	StreamTopMargin   int32
	InputBoxHeight    int32
	InputBoxPadding   int32
	KeyboardTopMargin int32
	KeyWidth          int32
	KeyHeight         int32
	KeySpacingX       int32
	KeySpacingY       int32
	VirtualTopPadding int32
	Colors            Colors
	StreamsUiConfig   StreamsUiConfig
}

type PlayerConfig struct {
	StreamWidth  int
	StreamHeight int
}

func Load() *Config {
	return &Config{
		AppName:    "Pocketstream",
		AppVersion: "v1.0.0",
		Display: DisplayConfig{
			Width:  640,
			Height: 480,
		},
		TwitchService: &twitch.TwitchService{
			Config: twitch.TwitchConfig{
				ClientId:               "kimne78kx3ncx6brgo4mv6wki5h1ko",
				GqlUrl:                 "https://gql.twitch.tv/gql",
				UsherUrl:               "https://usher.ttvnw.net/api/channel/hls",
				TopStreamsLimit:        10,
				HttpClient:             &http.Client{},
				StreamResolution:       "RESOLUTION=852x480",
				BrowsPagePopularSha256: "75a4899f0a765cc08576125512f710e157b147897c06f96325de72d4c5a64890",
				SearchResultsSha256:    "845698a3efbde3c2d1cc31e77ca1160cde6a21c556ad808106910ff63e727b98",
			},
		},
		UI: UIConfig{
			StreamsUiConfig: StreamsUiConfig{
				ThumbnailWidth:      200,
				ThumbnailHeight:     112,
				Padding:             5,
				RowHeight:           116,
				InfoTextX:           210,
				ProfilePictureSize:  32,
				LiveBadgeWidth:      43,
				LiveBadgeHeight:     22,
				LiveBadgeLeftMargin: 14,
				LiveBadgeTopMargin:  11,
				LiveTextLeftMargin:  17,
				LiveTextTopMargin:   13,
			},
			HeaderHeight:      40,
			FooterHeight:      22,
			RowHeight:         130,
			FontPath:          "font.ttf",
			FontSize:          16,
			Padding:           5,
			StreamsTopMargin:  64,
			StreamLeftMargin:  10,
			StreamTopMargin:   130,
			InputBoxHeight:    36,
			InputBoxPadding:   60,
			KeyboardTopMargin: 140,
			KeyWidth:          36,
			KeyHeight:         27,
			KeySpacingX:       8,
			KeySpacingY:       8,
			VirtualTopPadding: 20,
			Colors: Colors{
				BackgroundColor:                sdl.Color{15, 23, 42, 255},
				HeaderBackgroundColor:          sdl.Color{30, 58, 138, 255},
				HeaderTextColor:                sdl.Color{240, 249, 255, 255},
				FooterBackgroundColor:          sdl.Color{30, 41, 59, 255},
				FooterTextColor:                sdl.Color{148, 163, 184, 255},
				InputBoxBorderColor:            sdl.Color{59, 130, 246, 255},
				InputBoxBackgroundColor:        sdl.Color{15, 23, 42, 255},
				InputTextColor:                 sdl.Color{240, 249, 255, 255},
				LoadingTextColor:               sdl.Color{240, 249, 255, 255},
				NoResultsTextColor:             sdl.Color{240, 249, 255, 255},
				KeyBackgroundColor:             sdl.Color{51, 65, 85, 255},
				KeyBorderColor:                 sdl.Color{100, 116, 139, 255},
				SelectedKeyBackgroundColor:     sdl.Color{59, 130, 246, 255},
				SelectedKeyBorderColor:         sdl.Color{37, 99, 235, 255},
				KeyColor:                       sdl.Color{226, 232, 240, 255},
				SelectedKeyColor:               sdl.Color{255, 255, 255, 255},
				StreamThumbnailBackgroundColor: sdl.Color{30, 41, 59, 255},
				StreamLiveBadgeBackgroundColor: sdl.Color{239, 68, 68, 255},
				LiveTextColor:                  sdl.Color{255, 255, 255, 255},
				ViewersCountBackgroundColor:    sdl.Color{15, 23, 42, 200},
				ViewersCountTextColor:          sdl.Color{240, 249, 255, 255},
				ProfilePictureBackgroundColor:  sdl.Color{30, 41, 59, 255},
				StreamerNameTextColor:          sdl.Color{240, 249, 255, 255},
				StreamTitleColor:               sdl.Color{148, 163, 184, 255},
				SelectedStreamBorderColor:      sdl.Color{59, 130, 246, 255},
			},
		},
		Player: PlayerConfig{
			StreamWidth:  640,
			StreamHeight: 480,
		},
	}
}
