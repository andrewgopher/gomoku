package main

import (
	"fmt"
	"gomoku/game"
	"gomoku/util"
	"image/color"
	"math"
	"os"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	boardSize = 19
	numConsec = 5
)

const (
	assetsFolder string = "assets"
)

var (
	white  color.RGBA = color.RGBA{255, 255, 255, 255}
	grey   color.RGBA = color.RGBA{127, 127, 127, 255}
	black  color.RGBA = color.RGBA{0, 0, 0, 255}
	yellow color.RGBA = color.RGBA{242, 202, 92, 255}
)

var openSans *ttf.Font

var playerToColor map[game.Player]color.RGBA = map[game.Player]color.RGBA{
	game.BlackPlayer: black,
	game.WhitePlayer: white,
}

var (
	smallCircleRadius float32 = 0.0035
	bigCircleRadius   float32 = 0.02
)

type uiState struct {
	gameState *game.State
	hover     *game.Point[int]
	enabled   bool
}

func LoadFonts(renderer *sdl.Renderer) error {
	var err error
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	openSans, err = ttf.OpenFont(fmt.Sprintf("%v/%v/opensans.ttf", wd, assetsFolder), 60)
	return err
}

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

func EllipseF(renderer *sdl.Renderer, center *sdl.FPoint, radiusX float32, radiusY float32, color color.RGBA) {
	gfx.FilledEllipseRGBA(renderer, int32(center.X), int32(center.Y), int32(radiusX), int32(radiusY), color.R, color.G, color.B, color.A)
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

func LineF(renderer *sdl.Renderer, rect *sdl.FRect, color color.RGBA) {
	gfx.AALineRGBA(renderer, int32(rect.X), int32(rect.Y), int32(rect.X+rect.W), int32(rect.Y+rect.H), color.R, color.G, color.B, color.A)
}

func distF(a *sdl.FPoint, b *sdl.FPoint) float32 {
	return float32(math.Sqrt(float64((a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y))))
}

func RenderUIState(renderer *sdl.Renderer, state *uiState, rect *sdl.FRect, mouseX int32, mouseY int32) {
	smallCircleRadiusConv := smallCircleRadius * rect.W
	bigCircleRadiusConv := bigCircleRadius * rect.W
	cellWidth := rect.W / float32(len(state.gameState.Board)+1)
	for i := 0; i < len(state.gameState.Board); i++ {
		LineF(renderer, &sdl.FRect{X: rect.X + cellWidth, Y: rect.Y + cellWidth*float32(i+1), W: cellWidth * float32(len(state.gameState.Board)-1), H: 0}, grey)
		LineF(renderer, &sdl.FRect{X: rect.X + cellWidth*float32(i+1), Y: rect.Y + cellWidth, W: 0, H: cellWidth * float32(len(state.gameState.Board)-1)}, grey)
	}
	foundHover := false
	for i := 0; i < len(state.gameState.Board); i++ {
		for j := 0; j < len(state.gameState.Board); j++ {
			center := &sdl.FPoint{X: rect.X + cellWidth*float32(i+1), Y: rect.Y + cellWidth*float32(j+1)}
			EllipseF(renderer, center, smallCircleRadiusConv, smallCircleRadiusConv, grey)
			if state.enabled && distF(center, &sdl.FPoint{X: float32(mouseX), Y: float32(mouseY)}) <= bigCircleRadiusConv && state.gameState.Board[i][j] == game.NilPiece {
				foundHover = true
				EllipseF(renderer, center, bigCircleRadiusConv, bigCircleRadiusConv, playerToColor[state.gameState.Turn])
				state.hover = &game.Point[int]{X: i, Y: j}
			}
			if state.gameState.Board[i][j] != game.NilPiece {
				EllipseF(renderer, center, bigCircleRadiusConv, bigCircleRadiusConv, playerToColor[game.PieceToPlayer[state.gameState.Board[i][j]]])
			}
		}
	}
	if !foundHover {
		state.hover = nil
	}
	if state.gameState.Winner != game.NilPlayer {
		TextF(renderer, fmt.Sprintf("%v won", game.PlayerToString[state.gameState.Winner]), rect.X+rect.W/2, rect.Y+rect.H/2, openSans, black, true)
	}
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

	err = LoadFonts(renderer)
	if err != nil {
		panic(err)
	}

	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	running := true
	state := game.NewState(boardSize, numConsec, game.BlackPlayer)
	currUIState := &uiState{
		gameState: state,
		hover:     nil,
		enabled:   true,
	}
	for running {
		Clear(renderer, yellow)
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
		mouseX, mouseY, _ := sdl.GetMouseState()
		RenderUIState(renderer, currUIState, boardRect, mouseX, mouseY)
		renderer.Present()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				if e.Button == sdl.BUTTON_LEFT && e.Type == sdl.MOUSEBUTTONDOWN {
					if currUIState.hover != nil {
						state.MakeMove(currUIState.hover)
						if state.Winner != game.NilPlayer {
							currUIState.enabled = false
						}
					}
				}
			}
		}
	}
}
