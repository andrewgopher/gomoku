package main

import (
	"gomoku/game"
	"gomoku/util"
	"image/color"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	white color.RGBA = color.RGBA{255, 255, 255, 255}
	grey  color.RGBA = color.RGBA{127, 127, 127, 255}
	black color.RGBA = color.RGBA{0, 0, 0, 255}
)

func Clear(renderer *sdl.Renderer, color color.RGBA) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.Clear()
}

func Rect(renderer *sdl.Renderer, rect *sdl.Rect, color color.RGBA) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.FillRect(rect)
}

func RectF(renderer *sdl.Renderer, rect *sdl.FRect, color color.RGBA) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.FillRectF(rect)
}

func EllipseF(renderer *sdl.Renderer, rect *sdl.FRect, color color.RGBA) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	gfx.FilledEllipseRGBA(renderer, int32(rect.X), int32(rect.Y), int32(rect.W), int32(rect.H), color.R, color.G, color.B, color.A)
}

func TextF(renderer *sdl.Renderer, text string, posX float32, posY float32, font *ttf.Font, color color.RGBA, center bool) {
	textSurf, err := font.RenderUTF8Solid(text, sdl.Color{R: color.R, G: color.G, B: color.B, A: color.A})
	defer textSurf.Free()
	if err != nil {
		panic(err)
	}
	textTex, err := renderer.CreateTextureFromSurface(textSurf)
	defer textTex.Destroy()

	if err != nil {
		panic(err)
	}
	textW := float32(textSurf.W)
	textH := float32(textSurf.H)
	textRect := &sdl.FRect{X: posX, Y: posY, W: textW, H: textH}
	if center {
		textRect.X -= textW / 2
		textRect.Y -= textH / 2
	}
	renderer.CopyF(textTex, nil, textRect)
}

func RenderState(renderer *sdl.Renderer, state *game.State, rect *sdl.FRect) {
	cellWidth := rect.W / float32(len(state.Board)+1)

}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	err := ttf.Init()
	if err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Gomoku", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 800, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	window.SetResizable(true)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	running := true
	for running {
		Clear(renderer, white)
		w, h := window.GetSize()
		boardRect := &sdl.FRect{}
		if w < h {
			boardRect.X = 0
			boardRect.Y = float32(h-w) / 2
		} else {
			boardRect.X = float32(w-h) / 2
			boardRect.Y = 0
		}
		boardRect.W = float32(util.Min(int(w), int(h)))
		boardRect.H = float32(util.Min(int(w), int(h)))
		renderer.Present()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			default:
			}
		}
	}
}
