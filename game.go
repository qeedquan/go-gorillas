package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/qeedquan/go-media/sdl"
)

var (
	EXPLOSION_COLOR = color.RGBA{255, 0, 0, 255}
)

func displayScore(screenSurf *Surface, oneScore, twoScore int) {
	text := fmt.Sprintf("%d >Score< %d", oneScore, twoScore)
	drawText(text, screenSurf, 270, 310, WHITE_COLOR, SKY_COLOR, "left")
}

func doExplosion(screenSurf, skylineSurf *Surface, x, y, explosionSize int, speed float64) {
	for r := 1; r < explosionSize; r++ {
		screenSurf.circle(EXPLOSION_COLOR, x, y, r)
		skylineSurf.circle(EXPLOSION_COLOR, x, y, r)
		drawScreen()
		sleep(time.Duration(speed * float64(time.Second)))
	}
	for r := explosionSize; r > 1; r-- {
		screenSurf.circle(SKY_COLOR, x, y, explosionSize)
		skylineSurf.circle(SKY_COLOR, x, y, explosionSize)
		screenSurf.circle(EXPLOSION_COLOR, x, y, r)
		skylineSurf.circle(EXPLOSION_COLOR, x, y, r)
		drawScreen()
		sleep(time.Duration(speed * float64(time.Second)))
	}
	screenSurf.circle(SKY_COLOR, x, y, 2)
	skylineSurf.circle(SKY_COLOR, x, y, 2)
	drawScreen()
}

func plotShot(screenSurf, skylineSurf *Surface, angle, velocity float64, playerNum, wind int, gravity float64, gor1, gor2 image.Point) string {
	angle = angle / 180.0 * math.Pi
	initXVel := math.Cos(angle) * velocity
	initYVel := math.Sin(angle) * velocity
	gorWidth, gorHeight := GOR_DOWN_SURF.width(), GOR_DOWN_SURF.height()
	gor1rect := sdl.Rect{int32(gor1.X), int32(gor1.Y), int32(gorWidth), int32(gorHeight)}
	gor2rect := sdl.Rect{int32(gor2.X), int32(gor2.Y), int32(gorWidth), int32(gorHeight)}

	var gorImg int
	if playerNum == 1 {
		gorImg = LEFT_ARM_UP
	} else {
		gorImg = RIGHT_ARM_UP
	}

	var startx, starty int
	if playerNum == 1 {
		startx = gor1.X
		starty = gor1.Y
	} else if playerNum == 2 {
		startx = gor2.X
		starty = gor2.Y
	}

	drawGorilla(screenSurf, startx, starty, gorImg)
	drawScreen()
	sleep(300 * time.Millisecond)
	drawGorilla(screenSurf, startx, starty, BOTH_ARMS_DOWN)
	drawScreen()

	bananaShape := UP

	if playerNum == 2 {
		startx += GOR_DOWN_SURF.width()
	}

	starty -= int(getBananaRect(0, 0, bananaShape).H) + BAN_UP_SURF.height()

	impact := false
	bananaInPlay := true

	t := 1.0
	sunHit := false

	for !impact && bananaInPlay {
		x := float64(startx) + (float64(initXVel) * t) + (0.5 * (float64(wind) / 5) * t * t)
		y := float64(starty) + ((-1 * (float64(initYVel) * t)) + (0.5 * gravity * t * t))

		if x >= SCR_WIDTH-10 || x <= 3 || y >= SCR_HEIGHT {
			bananaInPlay = false
		}

		var bananaSurf *Surface
		bananaRect := getBananaRect(int(x), int(y), bananaShape)
		switch bananaShape {
		case UP:
			bananaSurf = BAN_UP_SURF
			bananaRect.X -= 2
			bananaRect.Y += 2
		case DOWN:
			bananaSurf = BAN_DOWN_SURF
			bananaRect.X -= 2
			bananaRect.Y += 2
		case LEFT:
			bananaSurf = BAN_LEFT_SURF
		case RIGHT:
			bananaSurf = BAN_RIGHT_SURF
		}

		bananaShape = nextBananaShape(bananaShape)

		if bananaInPlay && y > 0 {
			if (sdl.Point{int32(x), int32(y)}).In(sunRect) {
				sunHit = true
			}

			drawSun(screenSurf, sunHit)

			if bananaRect.Collide(gor1rect) {
				doExplosion(screenSurf, skylineSurf, bananaRect.CenterX(), bananaRect.CenterY(), GOR_EXPLOSION_SIZE*2/3, 0.005)
				doExplosion(screenSurf, skylineSurf, bananaRect.CenterX(), bananaRect.CenterY(), GOR_EXPLOSION_SIZE, 0.005)
				drawSun(screenSurf, false)
				return "gorilla1"
			} else if bananaRect.Collide(gor2rect) {
				doExplosion(screenSurf, skylineSurf, bananaRect.CenterX(), bananaRect.CenterY(), GOR_EXPLOSION_SIZE*2/3, 0.005)
				doExplosion(screenSurf, skylineSurf, bananaRect.CenterX(), bananaRect.CenterY(), GOR_EXPLOSION_SIZE, 0.005)
				screenSurf.fillRect(SKY_COLOR, imageRect(bananaRect))
				drawSun(screenSurf, false)
				return "gorilla2"
			} else if collideWithNonColor(skylineSurf, bananaRect, SKY_COLOR) {
				doExplosion(screenSurf, skylineSurf, bananaRect.CenterX(), bananaRect.CenterY(), BUILD_EXPLOSION_SIZE, 0.05)
				screenSurf.fillRect(SKY_COLOR, imageRect(bananaRect))
				drawSun(screenSurf, false)
				return "building"
			}
		}

		screenSurf.blit(bananaSurf, int(bananaRect.X), int(bananaRect.Y))
		drawScreen()
		sleep(20 * time.Millisecond)

		screenSurf.fillRect(SKY_COLOR, imageRect(bananaRect))

		t += 0.1
	}
	drawSun(screenSurf, false)
	return "miss"
}

func showStartScreen(screenSurf *Surface) {
	vertAdj := 0
	horAdj := 0
	for checkForKeyPress() == 0 {
		screenSurf.fill(BLACK_COLOR)
		drawStars(screenSurf, vertAdj, horAdj)

		vertAdj = (vertAdj + 1) % 4
		horAdj = (horAdj + 12) % 84

		drawText("G  O  R  I  L  L  A  S", screenSurf, SCR_WIDTH/2, 50, WHITE_COLOR, BLACK_COLOR, "center")
		drawText("Your mission is to hit your opponent with the exploding", screenSurf, SCR_WIDTH/2, 110, GRAY_COLOR, BLACK_COLOR, "center")
		drawText("banana by varying the angle and power of your throw, taking", screenSurf, SCR_WIDTH/2, 130, GRAY_COLOR, BLACK_COLOR, "center")
		drawText("into account wind speed, gravity, and the city skyline.", screenSurf, SCR_WIDTH/2, 150, GRAY_COLOR, BLACK_COLOR, "center")
		drawText("The wind speed is shown by a directional arrow at the bottom", screenSurf, SCR_WIDTH/2, 170, GRAY_COLOR, BLACK_COLOR, "center")
		drawText("of the playing field, its length relative to its strength.", screenSurf, SCR_WIDTH/2, 190, GRAY_COLOR, BLACK_COLOR, "center")
		drawText("Press any key to continue", screenSurf, SCR_WIDTH/2, 300, GRAY_COLOR, BLACK_COLOR, "center")

		drawScreen()
		sdl.Delay(1000 / FPS)
	}
}

func showSettingsScreen(screenSurf *Surface) (p1name, p2name string, points int, gravity float64, nextScreen rune) {
	screenSurf.fill(BLACK_COLOR)

	for {
		text := "Name of Player 1 (Default = 'Player 1'): "
		width, _, _ := GAME_FONT.SizeUTF8(text)
		name, escaped := inputMode(text, screenSurf, (SCR_WIDTH-width)/2, 50, GRAY_COLOR, BLACK_COLOR, 10, nil, "left", "_", true)
		if !escaped {
			p1name = name
			if strings.TrimSpace(p1name) == "" {
				p1name = "Player 1"
			}
			break
		}
	}

	for {
		text := "Name of Player 2 (Default = 'Player 2'): "
		width, _, _ := GAME_FONT.SizeUTF8(text)
		name, escaped := inputMode(text, screenSurf, (SCR_WIDTH-width)/2, 80, GRAY_COLOR, BLACK_COLOR, 10, nil, "left", "_", true)
		if !escaped {
			p2name = name
			if strings.TrimSpace(p2name) == "" {
				p2name = "Player 2"
			}
			break
		}
	}

	for {
		text := "Play to how many total points (Default = 3)? "
		width, _, _ := GAME_FONT.SizeUTF8(text)
		allowed := "0123456789"
		input, escaped := inputMode(text, screenSurf, (SCR_WIDTH-width)/2, 110, GRAY_COLOR, BLACK_COLOR, 6, &allowed, "left", "_", true)
		if !escaped {
			points, _ = strconv.Atoi(input)
			if input == "" {
				points = 3
			}
			break
		}
	}

	for {
		text := "Gravity in Meters/Sec (Earth = 9.8)? "
		width, _, _ := GAME_FONT.SizeUTF8(text)
		allowed := "0123456789."
		input, escaped := inputMode(text, screenSurf, (SCR_WIDTH-width)/2, 140, GRAY_COLOR, BLACK_COLOR, 6, &allowed, "left", "_", true)
		if !escaped {
			gravity, _ = strconv.ParseFloat(input, 64)
			if input == "" {
				gravity = 9.8
			}
			break
		}
	}

	drawText("--------------", screenSurf, SCR_WIDTH/2-10, 170, GRAY_COLOR, BLACK_COLOR, "center")
	drawText("V = View Intro", screenSurf, SCR_WIDTH/2-10, 200, GRAY_COLOR, BLACK_COLOR, "center")
	drawText("P = Play Game", screenSurf, SCR_WIDTH/2-10, 230, GRAY_COLOR, BLACK_COLOR, "center")
	drawText("Your Choice?", screenSurf, SCR_WIDTH/2-10, 260, GRAY_COLOR, BLACK_COLOR, "center")
	drawScreen()

	for {
		k := waitForPlayerToPressKey()
		if k == sdl.K_v {
			nextScreen = 'v'
			break
		}

		if k == sdl.K_p {
			nextScreen = 'p'
			break
		}
	}

	return
}

func showIntroScreen(screenSurf *Surface, p1name, p2name string) {
	screenSurf.fill(SKY_COLOR)
	drawText("P  y  t  h  o  n     G  O  R  I  L  L  A  S", screenSurf, SCR_WIDTH/2, 15, WHITE_COLOR, SKY_COLOR, "center")
	drawText("STARRING:", screenSurf, SCR_WIDTH/2, 55, WHITE_COLOR, SKY_COLOR, "center")
	drawText(fmt.Sprintf("%s AND %s", p1name, p2name), screenSurf, SCR_WIDTH/2, 115, WHITE_COLOR, SKY_COLOR, "center")

	x := 278
	y := 175

	for i := 0; i < 2; i++ {
		drawGorilla(screenSurf, x-13, y, RIGHT_ARM_UP)
		drawGorilla(screenSurf, x+47, y, LEFT_ARM_UP)
		drawScreen()

		sleep(2 * time.Second)

		drawGorilla(screenSurf, x-13, y, LEFT_ARM_UP)
		drawGorilla(screenSurf, x+47, y, RIGHT_ARM_UP)
		drawScreen()

		sleep(2 * time.Second)
	}

	for i := 0; i < 4; i++ {
		drawGorilla(screenSurf, x-13, y, LEFT_ARM_UP)
		drawGorilla(screenSurf, x+47, y, RIGHT_ARM_UP)
		drawScreen()

		sleep(300 * time.Millisecond)

		drawGorilla(screenSurf, x-13, y, RIGHT_ARM_UP)
		drawGorilla(screenSurf, x+47, y, LEFT_ARM_UP)
		drawScreen()

		sleep(300 * time.Millisecond)
	}
}

func showGameOverScreen(screenSurf *Surface, p1name string, p1score int, p2name string, p2score int) {
	vertAdj := 0
	horAdj := 0
	for checkForKeyPress() == 0 {
		screenSurf.fill(BLACK_COLOR)

		drawStars(screenSurf, vertAdj, horAdj)
		vertAdj = (vertAdj + 1) % 4
		horAdj = (horAdj + 12) % 84

		drawText("GAME OVER!", screenSurf, SCR_WIDTH/2, 120, GRAY_COLOR, BLACK_COLOR, "center")
		drawText("Score:", screenSurf, SCR_WIDTH/2, 155, GRAY_COLOR, BLACK_COLOR, "center")
		drawText(p1name, screenSurf, 225, 175, GRAY_COLOR, BLACK_COLOR, "left")
		drawText(fmt.Sprint(p1score), screenSurf, 395, 175, GRAY_COLOR, BLACK_COLOR, "left")
		drawText(p2name, screenSurf, 225, 195, GRAY_COLOR, BLACK_COLOR, "left")
		drawText(fmt.Sprint(p2score), screenSurf, 395, 195, GRAY_COLOR, BLACK_COLOR, "left")
		drawText("Press any key to continue", screenSurf, SCR_WIDTH/2, 298, GRAY_COLOR, BLACK_COLOR, "center")

		drawScreen()
		sdl.Delay(1000 / FPS)
	}
}
