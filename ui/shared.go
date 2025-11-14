package ui

import (
	"fmt"
	"math"

	"github.com/fspasovski/pocketstream-app/app"
	"github.com/fspasovski/pocketstream-app/model"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func DrawStreams(app *app.App, streams []model.Stream, startIndex int, endIndex int, selectedIndex int) {
	app.ClearScreen()

	if len(streams) == 0 {
		app.DrawCenteredText("No results.", app.Config.UI.Colors.NoResultsTextColor)
		return
	}

	y := app.Config.UI.HeaderHeight + app.Config.UI.StreamsTopMargin
	for i := startIndex; i <= endIndex; i++ {
		selected := i == selectedIndex
		err := drawStream(&streams[i], app, app.Config.UI.StreamLeftMargin, y, selected)
		if err == nil {
			y += app.Config.UI.StreamsUiConfig.Height
		}
	}
}

func drawStream(stream *model.Stream, app *app.App, x int32, y int32, selected bool) error {
	thumbnailBg := sdl.Rect{
		X: x,
		Y: y,
		W: app.Config.UI.StreamsUiConfig.ThumbnailWidth,
		H: app.Config.UI.StreamsUiConfig.ThumbnailHeight,
	}
	app.FillRect(&thumbnailBg, app.Config.UI.Colors.StreamThumbnailBackgroundColor)

	if len(stream.PreviewImageData) > 0 {
		rw, err := sdl.RWFromMem(stream.PreviewImageData)
		if err == nil {
			surface, err := img.LoadRW(rw, true)
			if err == nil {
				defer surface.Free()

				tex, err := app.CreateTextureFromSurface(surface)
				if err == nil {
					defer tex.Destroy()

					previewDst := sdl.Rect{
						X: x + app.Config.UI.StreamsUiConfig.Padding,
						Y: y + app.Config.UI.StreamsUiConfig.Padding,
						W: app.Config.UI.StreamsUiConfig.ThumbnailWidth - 2*app.Config.UI.StreamsUiConfig.Padding,
						H: app.Config.UI.StreamsUiConfig.ThumbnailHeight - 2*app.Config.UI.StreamsUiConfig.Padding,
					}
					app.CopyTexture(tex, nil, &previewDst)
				}
			}
		}
	}

	liveBadge := sdl.Rect{
		X: x + app.Config.UI.StreamsUiConfig.LiveBadgeLeftMargin,
		Y: y + app.Config.UI.StreamsUiConfig.LiveBadgeTopMargin,
		W: app.Config.UI.StreamsUiConfig.LiveBadgeWidth,
		H: app.Config.UI.StreamsUiConfig.LiveBadgeHeight,
	}
	app.FillRect(&liveBadge, app.Config.UI.Colors.StreamLiveBadgeBackgroundColor)

	// Draw LIVE text
	app.Font.SetStyle(ttf.STYLE_BOLD)
	app.DrawCenteredTextInRect("LIVE", &liveBadge, app.Config.UI.Colors.LiveTextColor)

	// Draw viewer count badge (bottom right of thumbnail)
	viewerText := fmt.Sprintf("%s", formatViewerCount(stream.ViewersCount))
	app.Font.SetStyle(ttf.STYLE_NORMAL)
	viewerSurface, err := app.Font.RenderUTF8Blended(viewerText, app.Config.UI.Colors.ViewersCountTextColor)
	if err == nil {
		defer viewerSurface.Free()
		viewerTexture, err := app.Renderer.CreateTextureFromSurface(viewerSurface)
		if err == nil {
			defer viewerTexture.Destroy()
			_, _, vw, vh, _ := viewerTexture.Query()

			//Draw semi-transparent background for viewer count
			viewerBg := sdl.Rect{
				X: x + app.Config.UI.StreamsUiConfig.ThumbnailWidth - vw - 15,
				Y: y + app.Config.UI.StreamsUiConfig.ThumbnailHeight - vh - 10,
				W: vw + 10,
				H: vh + 4,
			}
			app.Renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
			app.FillRect(&viewerBg, app.Config.UI.Colors.ViewersCountBackgroundColor)

			viewerDst := sdl.Rect{
				X: x + app.Config.UI.StreamsUiConfig.ThumbnailWidth - vw - 10,
				Y: y + app.Config.UI.StreamsUiConfig.ThumbnailHeight - vh - 8,
				W: vw,
				H: vh,
			}
			app.CopyTexture(viewerTexture, nil, &viewerDst)
		}
	}

	// Draw profile picture background
	profileBg := sdl.Rect{
		X: x + app.Config.UI.StreamsUiConfig.ProfileInfoLeftMargin,
		Y: y + app.Config.UI.StreamsUiConfig.ProfileInfoTopMargin,
		W: app.Config.UI.StreamsUiConfig.ProfilePictureSize,
		H: app.Config.UI.StreamsUiConfig.ProfilePictureSize,
	}
	app.FillRect(&profileBg, app.Config.UI.Colors.ProfilePictureBackgroundColor)

	//Draw profile picture if available
	profileX := x + app.Config.UI.StreamsUiConfig.ProfileInfoLeftMargin
	if len(stream.Broadcaster.ProfileImageData) > 0 {
		rw, err := sdl.RWFromMem(stream.Broadcaster.ProfileImageData)
		if err == nil {
			surface, err := img.LoadRW(rw, true)
			if err == nil {
				defer surface.Free()

				tex, err := app.CreateTextureFromSurface(surface)
				if err == nil {
					defer tex.Destroy()

					profileDst := sdl.Rect{
						X: profileX,
						Y: y + 10,
						W: app.Config.UI.StreamsUiConfig.ProfilePictureSize,
						H: app.Config.UI.StreamsUiConfig.ProfilePictureSize,
					}
					app.CopyTexture(tex, nil, &profileDst)
				}
			}
		}
	}

	// Draw streamer name (bold) - offset by profile picture width + spacing
	nameX := profileX + app.Config.UI.StreamsUiConfig.ProfilePictureSize + app.Config.UI.StreamsUiConfig.ProfileNameLeftMargin
	app.Font.SetStyle(ttf.STYLE_BOLD)
	nameSurface, err := app.Font.RenderUTF8Blended(stream.Broadcaster.Login, app.Config.UI.Colors.StreamerNameTextColor)
	if err == nil {
		defer nameSurface.Free()
		nameTexture, err := app.CreateTextureFromSurface(nameSurface)
		if err == nil {
			defer nameTexture.Destroy()
			_, _, nw, nh, _ := nameTexture.Query()
			// Center the name vertically with the profile picture
			nameY := y + 10 + (app.Config.UI.StreamsUiConfig.ProfilePictureSize-nh)/2
			nameDst := sdl.Rect{X: nameX, Y: nameY, W: nw, H: nh}
			app.CopyTexture(nameTexture, nil, &nameDst)
		}
	}

	if app.UserDataManager.IsFavoriteBroadcaster(stream.Broadcaster) {
		drawFilledStar(
			app,
			app.Config.Display.Width-app.Config.UI.StreamsUiConfig.FavoriteIconRightMargin,
			y+app.Config.UI.StreamsUiConfig.FavoriteIconTopMargin,
			app.Config.UI.StreamsUiConfig.FavoriteIconSize,
			app.Config.UI.Colors.FavoriteIconColor,
		)
	}

	// Draw stream title (gray, truncated if needed)
	app.Font.SetStyle(ttf.STYLE_NORMAL)
	truncatedTitle := truncateText(stream.Title, app.Config.UI.StreamsUiConfig.MaxTitleLength)
	titleSurface, err := app.Font.RenderUTF8Blended(truncatedTitle, app.Config.UI.Colors.StreamTitleColor)
	if err == nil {
		defer titleSurface.Free()
		titleTexture, err := app.CreateTextureFromSurface(titleSurface)
		if err == nil {
			defer titleTexture.Destroy()
			_, _, tw, th, _ := titleTexture.Query()
			titleDst := sdl.Rect{X: app.Config.UI.StreamsUiConfig.TitleLeftMargin, Y: y + app.Config.UI.StreamsUiConfig.TitleTopMargin, W: tw, H: th}
			app.CopyTexture(titleTexture, nil, &titleDst)
		}
	}

	if selected {
		selectionRect := sdl.Rect{X: x, Y: y, W: app.Config.UI.StreamsUiConfig.Width, H: app.Config.UI.StreamsUiConfig.Height}

		for i := int32(0); i < 3; i++ {
			borderRect := sdl.Rect{X: selectionRect.X + i, Y: selectionRect.Y + i, W: selectionRect.W - 2*i, H: selectionRect.H - 2*i}
			app.DrawRect(&borderRect, app.Config.UI.Colors.SelectedStreamBorderColor)
		}
	}

	return nil
}

func formatViewerCount(count int) string {
	if count >= 1000 {
		return fmt.Sprintf("%.1fK", float64(count)/1000.0)
	}
	return fmt.Sprintf("%d", count)
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

func drawFilledStar(app *app.App, centerX, centerY, size int32, color sdl.Color) {
	// Calculate star points
	points := make([]sdl.Point, 11) // 10 points + 1 to close the shape

	outerRadius := float64(size)
	innerRadius := outerRadius * 0.4 // Inner radius is 40% of outer radius

	for i := 0; i < 10; i++ {
		angle := float64(i)*math.Pi/5.0 - math.Pi/2.0 // Start from top

		var radius float64
		if i%2 == 0 {
			radius = outerRadius // Outer point
		} else {
			radius = innerRadius // Inner point
		}

		x := centerX + int32(radius*math.Cos(angle))
		y := centerY + int32(radius*math.Sin(angle))

		points[i] = sdl.Point{X: x, Y: y}
	}
	points[10] = points[0] // Close the shape

	// Fill the star using scanline algorithm
	// First, set the draw color

	// Draw filled triangles from center to create a filled star
	for i := 0; i < 10; i++ {
		drawFilledTriangle(
			app,
			centerX, centerY,
			points[i].X, points[i].Y,
			points[i+1].X, points[i+1].Y,
			color,
		)
	}
}

// drawFilledTriangle draws a filled triangle using scanline algorithm
func drawFilledTriangle(app *app.App, x1, y1, x2, y2, x3, y3 int32, color sdl.Color) {
	// Sort vertices by Y coordinate (y1 <= y2 <= y3)
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	if y1 > y3 {
		x1, x3 = x3, x1
		y1, y3 = y3, y1
	}
	if y2 > y3 {
		x2, x3 = x3, x2
		y2, y3 = y3, y2
	}

	// Handle degenerate cases
	if y1 == y3 {
		return
	}

	// Scanline fill
	for y := y1; y <= y3; y++ {
		var xStart, xEnd int32

		if y < y2 {
			// Upper part of triangle
			if y2 != y1 {
				xStart = x1 + (x2-x1)*(y-y1)/(y2-y1)
			} else {
				xStart = x1
			}
			if y3 != y1 {
				xEnd = x1 + (x3-x1)*(y-y1)/(y3-y1)
			} else {
				xEnd = x1
			}
		} else {
			// Lower part of triangle
			if y3 != y2 {
				xStart = x2 + (x3-x2)*(y-y2)/(y3-y2)
			} else {
				xStart = x2
			}
			if y3 != y1 {
				xEnd = x1 + (x3-x1)*(y-y1)/(y3-y1)
			} else {
				xEnd = x1
			}
		}

		if xStart > xEnd {
			xStart, xEnd = xEnd, xStart
		}

		app.DrawLine(xStart, y, xEnd, y, color)
	}
}

// drawStarOutline draws just the outline of a star (lighter weight alternative)
func drawStarOutline(app *app.App, centerX, centerY, size int32, color sdl.Color, thickness int32) {
	points := make([]sdl.Point, 11)

	outerRadius := float64(size)
	innerRadius := outerRadius * 0.4

	for i := 0; i < 10; i++ {
		angle := float64(i)*math.Pi/5.0 - math.Pi/2.0

		var radius float64
		if i%2 == 0 {
			radius = outerRadius
		} else {
			radius = innerRadius
		}

		x := centerX + int32(radius*math.Cos(angle))
		y := centerY + int32(radius*math.Sin(angle))

		points[i] = sdl.Point{X: x, Y: y}
	}
	points[10] = points[0]

	// Draw multiple lines for thickness
	for t := int32(0); t < thickness; t++ {
		for i := 0; i < 10; i++ {
			app.DrawLine(points[i].X, points[i].Y+t, points[i+1].X, points[i+1].Y+t, color)
			app.DrawLine(points[i].X+t, points[i].Y, points[i+1].X+t, points[i+1].Y, color)
		}
	}
}

// Example usage in your drawStream function:
func exampleUsageInDrawStream() {
	/*
		// Add this to your drawStream function, after drawing the thumbnail:

		// Check if broadcaster is a favorite
		if app.Config.IsFavorite(stream.Broadcaster.Login) {
			// Position star in top-right corner of thumbnail
			starX := x + app.Config.UI.StreamsUiConfig.ThumbnailWidth - 32
			starY := y + 20
			starSize := int32(12) // Size of the star

			// Gold color for the star
			goldColor := sdl.Color{R: 255, G: 215, B: 0, A: 255}

			// Optional: Draw a dark background circle for better visibility
			drawFilledCircle(app.Renderer, starX, starY, starSize+4,
				sdl.Color{R: 0, G: 0, B: 0, A: 180})

			// Draw the star
			drawFilledStar(app.Renderer, starX, starY, starSize, goldColor)
		}
	*/
}

// Bonus: Helper function to draw a filled circle (for star background)
func drawFilledCircle(renderer *sdl.Renderer, centerX, centerY, radius int32, color sdl.Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)

	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				renderer.DrawPoint(centerX+x, centerY+y)
			}
		}
	}
}
