package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/qeedquan/go-media/sdl"
)

func waitForPlayerToPressKey() sdl.Keycode {
	for {
		key := checkForKeyPress()
		if key != 0 {
			return key
		}

		drawScreen()
	}
}

func checkForKeyPress() sdl.Keycode {
	for {
		ev := sdl.PollEvent()
		if ev == nil {
			break
		}
		switch ev := ev.(type) {
		case sdl.QuitEvent:
			os.Exit(0)
		case sdl.KeyUpEvent:
			switch ev.Sym {
			case sdl.K_ESCAPE:
				os.Exit(0)
			case sdl.K_LALT, sdl.K_RALT:
				return 0
			}
			return ev.Sym
		}
	}
	return 0
}

func drawText(text string, surfObj *Surface, x, y int, fgcol, bgcol color.RGBA, pos string) image.Rectangle {
	w, h, _ := GAME_FONT.SizeUTF8(text)
	if pos == "center" {
		x -= w / 2
	}

	col := sdlColor(fgcol)
	bounds, err := GAME_FONT.RenderUTF8BlendedEx(textbuf, text, col)
	ck(err)

	rect := image.Rect(x, y, x+int(bounds.W), y+int(bounds.H))
	draw.Draw(surfObj.RGBA, rect, textbuf, image.ZP, draw.Src)
	for i := rect.Min.Y; i < rect.Max.Y; i++ {
		for j := rect.Min.X; j < rect.Max.X; j++ {
			if surfObj.RGBAAt(j, i) == (color.RGBA{}) {
				surfObj.SetRGBA(j, i, bgcol)
			}
		}
	}

	return image.Rect(x, y, x+w, y+h)
}

func getShot(screenSurf *Surface, p1name, p2name string, playerNum int) (angle, velocity float64) {
	screenSurf.fillRect(SKY_COLOR, image.Rect(0, 0, 200, 50))
	screenSurf.fillRect(SKY_COLOR, image.Rect(550, 0, 550, 50))

	drawText(p1name, screenSurf, 2, 2, WHITE_COLOR, SKY_COLOR, "left")
	drawText(p2name, screenSurf, 500, 2, WHITE_COLOR, SKY_COLOR, "left")

	var x int
	if playerNum == 1 {
		x = 2
	} else {
		x = 500
	}

	allowed := "0123456789"
	for {
		inputAngle, escaped := inputMode("Angle: ", screenSurf, x, 23, WHITE_COLOR, SKY_COLOR, 3, &allowed, "left", "_", false)
		if escaped {
			os.Exit(0)
		}
		if inputAngle == "" {
			continue
		}

		angle, _ = strconv.ParseFloat(inputAngle, 64)
		break
	}

	for {
		inputVelocity, escaped := inputMode("Velocity: ", screenSurf, x, 45, WHITE_COLOR, SKY_COLOR, 3, &allowed, "left", "_", false)
		if escaped {
			os.Exit(0)
		}
		if inputVelocity == "" {
			continue
		}

		velocity, _ = strconv.ParseFloat(inputVelocity, 64)
		break
	}

	drawText(fmt.Sprint("Angle: ", angle), screenSurf, x, 23, SKY_COLOR, SKY_COLOR, "left")
	drawText(fmt.Sprint("Velocity: ", velocity), screenSurf, x, 45, SKY_COLOR, SKY_COLOR, "left")
	drawScreen()

	if playerNum == 2 {
		angle = 180 - angle
	}

	return

}

func inputMode(prompt string, screenSurf *Surface, x, y int, fgcol, bgcol color.RGBA, maxlen int, allowed *string, pos, cursor string, cursorBlink bool) (inputText string, escaped bool) {
	var textrect image.Rectangle

	cursorTimestamp := time.Now()
	cursorShow := cursor
	done := false

	for !done {
		after := time.Now().After(cursorTimestamp.Add(1 * time.Second))
		if cursor != "" && cursorBlink && after {
			if cursorShow == cursor {
				cursorShow = "   "
			} else {
				cursorShow = cursor
			}
			cursorTimestamp = time.Now()
		}

		for {
			ev := sdl.PollEvent()
			if ev == nil {
				break
			}
			switch ev := ev.(type) {
			case sdl.QuitEvent:
				os.Exit(0)

			case sdl.KeyUpEvent:
				switch ev.Sym {
				case sdl.K_ESCAPE:
					for len(inputText) > 0 {
						drawText(prompt+inputText+cursorShow, screenSurf, textrect.Min.X, textrect.Min.Y, fgcol, bgcol, "left")
						inputText = inputText[:len(inputText)-1]
						cursorShow = "   "
						textrect = drawText(prompt+cursorShow, screenSurf, x, y, fgcol, bgcol, pos)
						drawText(prompt+inputText+cursorShow, screenSurf, textrect.Min.X, textrect.Min.Y, fgcol, bgcol, "left")
					}
					drawScreen()
					return "", true
				case sdl.K_RETURN:
					if cursorShow != "" {
						done = true
						cursorShow = "   "
					}
				case sdl.K_BACKSPACE:
					if len(inputText) > 0 {
						drawText(prompt+inputText+cursorShow, screenSurf, textrect.Min.X, textrect.Min.Y, fgcol, bgcol, "left")
						inputText = inputText[:len(inputText)-1]
						cursorShow = "   "
						textrect = drawText(prompt+cursorShow, screenSurf, x, y, fgcol, bgcol, pos)
						drawText(prompt+inputText+cursorShow, screenSurf, textrect.Min.X, textrect.Min.Y, fgcol, bgcol, "left")
					}
				default:
					if allowed != nil && strings.IndexRune(*allowed, rune(ev.Sym)) < 0 {
						continue
					}

					if !(32 <= ev.Sym && ev.Sym < 128) {
						continue
					}

					if len(inputText) >= maxlen {
						continue
					}

					inputText += string(ev.Sym)
				}
			}
		}

		textrect = drawText(prompt+cursorShow, screenSurf, x, y, fgcol, bgcol, pos)
		drawText(prompt+inputText+cursorShow, screenSurf, textrect.Min.X, textrect.Min.Y, fgcol, bgcol, "left")
		drawScreen()
		sdl.Delay(1000 / FPS)
	}
	return
}
