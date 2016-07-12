package main

import (
	"image"
	"image/color"
	"math/rand"
	"os"
	"time"

	"github.com/qeedquan/go-media/sdl"
)

func randInt(a, b int) int {
	return a + rand.Intn(b-a+1)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sdlColor(c color.RGBA) sdl.Color {
	return sdl.Color{c.R, c.G, c.B, c.A}
}

func imageRect(r sdl.Rect) image.Rectangle {
	return image.Rect(int(r.X), int(r.Y), int(r.X+r.W), int(r.Y+r.H))
}

func collideWithNonColor(pixArr *Surface, rect sdl.Rect, color color.RGBA) bool {
	rightSide := min(int(rect.X+rect.W), SCR_WIDTH)
	bottomSide := min(int(rect.Y+rect.H), SCR_HEIGHT)

	for x := int(rect.X); x < rightSide; x++ {
		for y := int(rect.Y); y < bottomSide; y++ {
			if pixArr.At(x, y) != color {
				return true
			}
		}
	}
	return false
}

func sleep(d time.Duration) {
	t := time.NewTicker(d)
	defer t.Stop()
loop:
	for {
		select {
		case <-t.C:
			break loop
		default:
		}

		for {
			ev := sdl.PollEvent()
			if ev == nil {
				break
			}
			switch ev := ev.(type) {
			case sdl.QuitEvent:
				os.Exit(0)
			case sdl.KeyDownEvent:
				switch ev.Sym {
				case sdl.K_ESCAPE:
					os.Exit(0)
				}
			}
		}

		drawScreen()
	}
}
