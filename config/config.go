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
	Width                 int32
	ThumbnailWidth        int32
	ThumbnailHeight       int32
	ProfilePictureSize    int32
	ProfileInfoLeftMargin int32
	ProfileInfoTopMargin  int32
	ProfileNameLeftMargin int32
	TitleLeftMargin       int32
	TitleTopMargin        int32
	MaxTitleLength        int
	Padding               int32
	Height                int32
	LiveBadgeWidth        int32
	LiveBadgeHeight       int32
	LiveBadgeLeftMargin   int32
	LiveBadgeTopMargin    int32
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
	InputBoxTopMargin int32
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

func Load(screenWidth, screenHeight int) *Config {
	thumbnailWidth := int32(float32(screenWidth) * 0.3)
	profilePictureSize := int32(float32(screenHeight) * 0.104)
	headerHeight := int32(float32(screenHeight) * 0.104)
	inputBoxTopMargin := headerHeight + 50
	inputBoxHeight := int32(float32(screenHeight) * 0.075)

	return &Config{
		AppName:    "Pocketstream",
		AppVersion: "v1.0.2",
		Display: DisplayConfig{
			Width:  int32(screenWidth),
			Height: int32(screenHeight),
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
				Width:                 int32(float32(screenWidth) * 0.969),
				Height:                int32(float32(screenHeight) * 0.25),
				ThumbnailWidth:        thumbnailWidth,
				ThumbnailHeight:       int32(float32(screenHeight) * 0.24),
				Padding:               5,
				ProfileInfoLeftMargin: thumbnailWidth + 10,
				ProfileInfoTopMargin:  10,
				ProfileNameLeftMargin: 10,
				ProfilePictureSize:    profilePictureSize,
				LiveBadgeWidth:        int32(float32(screenWidth) * 0.08),
				LiveBadgeHeight:       int32(float32(screenHeight) * 0.063),
				LiveBadgeLeftMargin:   14,
				LiveBadgeTopMargin:    11,
				MaxTitleLength:        50,
				TitleLeftMargin:       thumbnailWidth + 20,
				TitleTopMargin:        profilePictureSize + 20,
			},
			HeaderHeight:      int32(float32(screenHeight) * 0.104),
			FooterHeight:      int32(float32(screenHeight) * 0.083),
			RowHeight:         130,
			FontPath:          "font.ttf",
			FontSize:          int(float32(screenHeight) * 0.040),
			Padding:           5,
			StreamsTopMargin:  20,
			StreamLeftMargin:  10,
			StreamTopMargin:   130,
			InputBoxHeight:    inputBoxHeight,
			InputBoxPadding:   60,
			InputBoxTopMargin: inputBoxTopMargin,
			KeyboardTopMargin: inputBoxTopMargin + inputBoxHeight + 20,
			KeyWidth:          int32(float32(screenWidth) * 0.056),
			KeyHeight:         int32(float32(screenHeight) * 0.056),
			KeySpacingX:       int32(float32(screenWidth) * 0.0125),
			KeySpacingY:       int32(float32(screenHeight) * 0.017),
			VirtualTopPadding: int32(float32(screenHeight) * 0.042),
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
			StreamWidth:  screenWidth,
			StreamHeight: screenHeight,
		},
	}
}
