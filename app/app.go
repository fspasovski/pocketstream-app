package app

import (
	"log"
	"time"

	"github.com/fspasovski/pocketstream-app/config"
	"github.com/fspasovski/pocketstream-app/input"
	"github.com/fspasovski/pocketstream-app/model"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type App struct {
	Config      *config.Config
	Running     bool
	State       Screen
	TopStreams  []model.Stream
	Window      *sdl.Window
	Renderer    *sdl.Renderer
	Font        *ttf.Font
	NeedsRedraw bool
	IsLoading   bool
	LoadingText string
}

func (a *App) LoadTopStreams() {
	a.IsLoading = true
	a.LoadingText = "Loading streams..."

	go func() {
		topStreams, err := a.Config.TwitchService.GetTopStreams()
		if err != nil {
			log.Println("Error fetching top streams:", err)
			a.TopStreams = []model.Stream{}
		} else {
			a.TopStreams = topStreams
		}
		a.IsLoading = false
		a.NeedsRedraw = true
	}()
}

func (a *App) StartLoading(text string) {
	a.IsLoading = true
	a.LoadingText = text
}

func (a *App) FinishLoading() {
	a.IsLoading = false
	a.LoadingText = ""
}

func (a *App) RaiseAppWindow() {
	time.Sleep(200 * time.Millisecond)
	a.Window.Hide()
	time.Sleep(50 * time.Millisecond)
	a.Window.Show()
	a.Window.Raise()

	for sdl.PollEvent() != nil {
	}

	a.NeedsRedraw = true
}

type Screen interface {
	HandleInput(appState *App, key input.Key)
	Draw(appState *App)
}

func (a *App) ClearScreen() {
	backgroundColor := a.Config.UI.Colors.BackgroundColor
	a.Renderer.SetDrawColor(backgroundColor.R, backgroundColor.G, backgroundColor.B, backgroundColor.A)
	a.Renderer.Clear()
}

func (a *App) Draw() {
	if a.NeedsRedraw {
		a.redrawUI()
	}

	a.State.Draw(a)
	a.DrawHeader()
	a.DrawFooter()
	a.Renderer.Present()
	sdl.Delay(16)
}

func (a *App) redrawUI() {
	for i := 0; i < 3; i++ {
		a.Renderer.SetDrawColor(0, 0, 0, 255)
		a.Renderer.Clear()
		a.Renderer.Present()
		sdl.Delay(16) // One frame delay between clears
	}

	a.NeedsRedraw = false
}

func (a *App) DrawLoadingScreen() {
	a.ClearScreen()
	a.DrawCenteredText(a.LoadingText, a.Config.UI.Colors.LoadingTextColor)
}

func (a *App) DrawCenteredText(text string, color sdl.Color) {
	centerX := a.Config.Display.Width / 2
	centerY := a.Config.Display.Height / 2
	if text == "" {
		return
	}

	surface, err := a.Font.RenderUTF8Blended(text, sdl.Color{R: color.R, G: color.G, B: color.B, A: color.A})
	if err != nil {
		return
	}
	defer surface.Free()

	texture, err := a.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return
	}
	defer texture.Destroy()

	x := centerX - surface.W/2
	y := centerY - surface.H/2
	a.DrawText(text, color, x, y)
}

func (a *App) DrawCenteredTextInRect(text string, rect *sdl.Rect, color sdl.Color) {
	if text == "" {
		return
	}

	surface, err := a.Font.RenderUTF8Blended(text, color)
	if err != nil {
		return
	}
	defer surface.Free()

	texture, err := a.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return
	}
	defer texture.Destroy()

	x := rect.X + (rect.W-surface.W)/2
	y := rect.Y + (rect.H-surface.H)/2
	dst := sdl.Rect{X: x, Y: y, W: surface.W, H: surface.H}
	a.Renderer.Copy(texture, nil, &dst)
}

func (a *App) DrawText(text string, color sdl.Color, x int32, y int32) {
	if text == "" {
		return
	}

	surface, err := a.Font.RenderUTF8Blended(text, sdl.Color{R: color.R, G: color.G, B: color.B, A: color.A})
	if err != nil {
		return
	}
	defer surface.Free()

	texture, err := a.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return
	}
	defer texture.Destroy()

	dst := sdl.Rect{X: x, Y: y, W: surface.W, H: surface.H}
	a.Renderer.Copy(texture, nil, &dst)
}

func (a *App) DrawHeader() error {
	// Draw header background
	headerBg := sdl.Rect{X: 0, Y: 0, W: a.Config.Display.Width, H: a.Config.UI.HeaderHeight}
	a.FillRect(&headerBg, a.Config.UI.Colors.HeaderBackgroundColor)

	// Draw app name (left side)
	a.Font.SetStyle(ttf.STYLE_BOLD)
	nameSurface, err := a.Font.RenderUTF8Blended(a.Config.AppName, a.Config.UI.Colors.HeaderTextColor)
	if err != nil {
		return err
	}
	defer nameSurface.Free()

	nameTexture, err := a.Renderer.CreateTextureFromSurface(nameSurface)
	if err != nil {
		return err
	}
	defer nameTexture.Destroy()

	_, _, nw, nh, err := nameTexture.Query()
	if err != nil {
		return err
	}

	// Center vertically in header
	nameY := (a.Config.UI.HeaderHeight - nh) / 2
	nameDst := sdl.Rect{X: 20, Y: nameY, W: nw, H: nh}
	a.Renderer.Copy(nameTexture, nil, &nameDst)

	// Draw version (right side)
	a.Font.SetStyle(ttf.STYLE_NORMAL)
	versionSurface, err := a.Font.RenderUTF8Blended(a.Config.AppVersion, a.Config.UI.Colors.HeaderTextColor)
	if err != nil {
		return err
	}
	defer versionSurface.Free()

	versionTexture, err := a.Renderer.CreateTextureFromSurface(versionSurface)
	if err != nil {
		return err
	}
	defer versionTexture.Destroy()

	_, _, vw, vh, err := versionTexture.Query()
	if err != nil {
		return err
	}

	// Center vertically, align to right with padding
	versionY := (a.Config.UI.HeaderHeight - vh) / 2
	versionDst := sdl.Rect{X: a.Config.Display.Width - vw - 20, Y: versionY, W: vw, H: vh}
	a.Renderer.Copy(versionTexture, nil, &versionDst)

	return nil
}

func (a *App) DrawFooter() error {
	footerY := a.Config.Display.Height - a.Config.UI.FooterHeight

	// Draw footer background
	footerBg := sdl.Rect{X: 0, Y: footerY, W: a.Config.Display.Width, H: a.Config.UI.FooterHeight}
	a.FillRect(&footerBg, a.Config.UI.Colors.FooterBackgroundColor)

	// Draw button hint text
	a.Font.SetStyle(ttf.STYLE_NORMAL)
	hintText := "D-Pad: Navigate  |  A: Select  |  B: Back  |  X: Search"

	hintSurface, err := a.Font.RenderUTF8Blended(hintText, a.Config.UI.Colors.FooterTextColor)
	if err != nil {
		return err
	}
	defer hintSurface.Free()

	hintTexture, err := a.Renderer.CreateTextureFromSurface(hintSurface)
	if err != nil {
		return err
	}
	defer hintTexture.Destroy()

	_, _, w, h, err := hintTexture.Query()
	if err != nil {
		return err
	}

	// Center the text vertically in the footer, with left padding
	textY := footerY + (a.Config.UI.FooterHeight-h)/2
	hintDst := sdl.Rect{X: 20, Y: textY, W: w, H: h}

	return a.Renderer.Copy(hintTexture, nil, &hintDst)
}

func (a *App) DrawRect(rect *sdl.Rect, color sdl.Color) {
	a.Renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	a.Renderer.DrawRect(rect)
}

func (a *App) FillRect(rect *sdl.Rect, color sdl.Color) {
	a.Renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	a.Renderer.FillRect(rect)
}

func (a *App) CreateTextureFromSurface(surface *sdl.Surface) (*sdl.Texture, error) {
	return a.Renderer.CreateTextureFromSurface(surface)
}

func (a *App) CopyTexture(texture *sdl.Texture, src, dst *sdl.Rect) error {
	return a.Renderer.Copy(texture, src, dst)
}
