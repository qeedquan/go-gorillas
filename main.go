package main

import (
	"flag"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

const (
	SCR_WIDTH  = 640
	SCR_HEIGHT = 350
	FPS        = 30
)

var (
	SKY_COLOR      = color.RGBA{0, 0, 173, 255}
	DARK_RED_COLOR = color.RGBA{173, 0, 0, 255}
	BLACK_COLOR    = color.RGBA{0, 0, 0, 255}
	WHITE_COLOR    = color.RGBA{255, 255, 255, 255}
	GRAY_COLOR     = color.RGBA{173, 170, 173, 255}
)

var (
	assets     = flag.String("assets", filepath.Join(sdl.GetBasePath(), "assets"), "assets directory")
	fullscreen = flag.Bool("fullscreen", false, "fullscreen")

	GAME_FONT  *sdlttf.Font
	textbuf    *sdl.Surface
	winSurface *Surface
	screen     *Display
	texture    *sdl.Texture
)

func main() {
	runtime.LockOSThread()
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)
	flag.Parse()
	initSDL()

	winSurface = newSurface(SCR_WIDTH, SCR_HEIGHT)
	showStartScreen(winSurface)

	var wind int
	var skylineSurf *Surface
	var buildCoords []image.Point
	var gorPos [2]image.Point
	for {
		p1name, p2name, winPoints, gravity, nextScreen := showSettingsScreen(winSurface)
		if nextScreen == 'v' {
			showIntroScreen(winSurface, p1name, p2name)
		}

		p1score := 0
		p2score := 0
		turn := 1

		newRound := true
		for p1score < winPoints && p2score < winPoints {
			if newRound {
				skylineSurf, buildCoords = makeCityScape()
				gorPos = placeGorillas(buildCoords)
				wind = getWind()
				newRound = false
			}

			winSurface.blit(skylineSurf, 0, 0)
			drawGorilla(winSurface, gorPos[0].X, gorPos[0].Y, 0)
			drawGorilla(winSurface, gorPos[1].X, gorPos[1].Y, 0)
			drawWind(winSurface, wind)
			drawSun(winSurface, false)
			displayScore(winSurface, p1score, p2score)
			drawScreen()

			angle, velocity := getShot(winSurface, p1name, p2name, turn)
			result := plotShot(winSurface, skylineSurf, angle, velocity, turn, wind, gravity, gorPos[0], gorPos[1])

			if result == "gorilla1" {
				victoryDance(winSurface, gorPos[1].X, gorPos[1].Y)
				p2score += 1
				newRound = true
			} else if result == "gorilla2" {
				victoryDance(winSurface, gorPos[0].X, gorPos[0].Y)
				p1score += 1
				newRound = true
			}

			turn = (turn % 2) + 1
		}
		showGameOverScreen(winSurface, p1name, p1score, p2name, p2score)
	}
}

type Display struct {
	*sdl.Window
	*sdl.Renderer
}

func newDisplay(w, h int, flags sdl.WindowFlags) *Display {
	window, renderer, err := sdl.CreateWindowAndRenderer(w, h, flags)
	ck(err)
	return &Display{window, renderer}
}

func initSDL() {
	err := sdl.Init(sdl.INIT_EVERYTHING &^ sdl.INIT_AUDIO)
	ck(err)

	err = sdlttf.Init()
	ck(err)

	wflag := sdl.WINDOW_RESIZABLE
	if *fullscreen {
		wflag |= sdl.WINDOW_FULLSCREEN_DESKTOP
	}
	screen = newDisplay(SCR_WIDTH, SCR_HEIGHT, wflag)
	screen.SetTitle("Gorillas")
	screen.SetLogicalSize(SCR_WIDTH, SCR_HEIGHT)

	name := filepath.Join(*assets, "vga437.ttf")
	GAME_FONT, err = sdlttf.OpenFont(name, 16)
	ck(err)

	texture, err = screen.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, SCR_WIDTH, SCR_HEIGHT)
	ck(err)

	textbuf, err = sdl.CreateRGBSurface(sdl.SWSURFACE, SCR_WIDTH, SCR_HEIGHT, 32, 0x00FF0000, 0x0000FF00, 0x000000FF, 0xFF000000)
	ck(err)

	sdl.ShowCursor(0)
}

func drawScreen() {
	screen.SetDrawColor(sdlcolor.Black)
	screen.Clear()
	texture.Update(nil, winSurface.Pix, winSurface.Stride)
	screen.Copy(texture, nil, nil)
	screen.Present()
}

func ck(err error) {
	if err != nil {
		sdl.LogCritical(sdl.LOG_CATEGORY_APPLICATION, "%v", err)
		if screen != nil {
			sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", err.Error(), screen.Window)
		}
		os.Exit(1)
	}
}
