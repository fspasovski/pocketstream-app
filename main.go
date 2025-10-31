package main

import (
	"fmt"
	"log"

	"github.com/fspasovski/pocketstream-app/app"
	"github.com/fspasovski/pocketstream-app/config"
	"github.com/fspasovski/pocketstream-app/input"
	"github.com/fspasovski/pocketstream-app/model"
	"github.com/fspasovski/pocketstream-app/player"
	"github.com/fspasovski/pocketstream-app/ui"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	cfg := config.Load()

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS); err != nil {
		log.Fatalf("Failed to init SDL: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		log.Fatalf("TTF init failed: %s\n", err)
	}
	defer ttf.Quit()

	sdl.InitSubSystem(sdl.INIT_JOYSTICK)
	joystick := initJoystick()
	if joystick != nil {
		defer joystick.Close()
	}

	if sdl.NumJoysticks() > 0 {
		if sdl.IsGameController(0) {
			controller := sdl.GameControllerOpen(0)
			defer controller.Close()
		}
	}

	font, err := ttf.OpenFont(cfg.UI.FontPath, cfg.UI.FontSize)
	if err != nil {
		log.Fatal(err)
	}
	defer font.Close()

	window, err := sdl.CreateWindow("Pocketstream", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, cfg.Display.Width, cfg.Display.Height, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	if err := img.Init(img.INIT_PNG); err != nil {
		log.Fatalf("could not initialize SDL_image: %v", err)
	}
	defer img.Quit()

	if err != nil {
		log.Fatalf("could not create texture: %v", err)
	}

	mediaPlayer := &player.Player{Cfg: cfg, BroadcasterStreamingUrls: make(map[string]string)}
	app := &app.App{
		Running:     true,
		State:       ui.CreateMainScreen(mediaPlayer),
		TopStreams:  make([]model.Stream, 0),
		Window:      window,
		Renderer:    renderer,
		Font:        font,
		NeedsRedraw: true,
		IsLoading:   false,
		LoadingText: "",
		Config:      cfg,
	}

	app.LoadTopStreams()

	for app.Running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				app.Running = false
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_SHOWN ||
					e.Event == sdl.WINDOWEVENT_EXPOSED ||
					e.Event == sdl.WINDOWEVENT_FOCUS_GAINED {
					app.NeedsRedraw = true
				}
			default:
				keyMapperStrategy := input.GetKeyMapperStrategy(e)
				if keyMapperStrategy != nil {
					key := keyMapperStrategy.MapInputToKey(e)
					if key != input.Unknown {
						app.State.HandleInput(app, keyMapperStrategy.MapInputToKey(e))
					}
				}
			}
		}

		app.Draw()
	}
}

func initJoystick() *sdl.Joystick {
	if sdl.NumJoysticks() > 0 {
		joystick := sdl.JoystickOpen(0)
		if joystick != nil {
			fmt.Printf("Joystick initialized: %s\n", sdl.JoystickNameForIndex(0))
			return joystick
		}
	}
	return nil
}
