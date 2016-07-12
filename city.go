package main

import (
	"image"
	"image/color"
	"math/rand"
)

var (
	BUILDING_COLORS = []color.RGBA{{173, 170, 173, 255}, {0, 170, 173, 255}, {173, 0, 0, 255}}
	LIGHT_WINDOW    = color.RGBA{255, 255, 82, 255}
	DARK_WINDOW     = color.RGBA{82, 85, 82, 255}
)

const BUILD_EXPLOSION_SIZE = SCR_HEIGHT / 50

func getWind() int {
	wind := randInt(5, 15)
	if rand.Intn(2) == 1 {
		wind *= -1
	}
	return wind
}

func drawWind(screenSurf *Surface, wind int) {
	if wind != 0 {
		wind *= 3
		screenSurf.line(EXPLOSION_COLOR, SCR_WIDTH/2, SCR_HEIGHT-7, SCR_WIDTH/2+wind, SCR_HEIGHT-7)

		var arrowDir int
		if wind > 0 {
			arrowDir = -2
		} else {
			arrowDir = 2
		}
		screenSurf.line(EXPLOSION_COLOR, SCR_WIDTH/2+wind, SCR_HEIGHT-7, SCR_WIDTH/2+wind+arrowDir, SCR_HEIGHT-7-2)
		screenSurf.line(EXPLOSION_COLOR, SCR_WIDTH/2+wind, SCR_HEIGHT-7, SCR_WIDTH/2+wind+arrowDir, SCR_HEIGHT-7+2)
	}
}

func makeCityScape() (screenSurf *Surface, buildingCoords []image.Point) {
	screenSurf = newSurface(SCR_WIDTH, SCR_HEIGHT)
	screenSurf.fill(SKY_COLOR)

	var slope string
	var newHeight int
	switch randInt(1, 6) {
	case 1:
		slope = "upward"
		newHeight = 15
	case 2:
		slope = "downward"
		newHeight = 130
	case 3, 4, 5:
		slope = "v"
		newHeight = 15
	default:
		slope = "^"
		newHeight = 130
	}

	bottomLine := 335
	heightInc := 10
	defBuildWidth := 37
	randomHeightDiff := 120
	windowWidth := 4
	windowHeight := 7
	windowSpacingX := 10
	windowSpacingY := 15
	gHeight := 25

	x := 2

	for x < SCR_WIDTH-heightInc {
		switch slope {
		case "upward":
			newHeight += heightInc
		case "downward":
			newHeight -= heightInc
		case "v":
			if x > SCR_WIDTH/2 {
				newHeight -= (2 * heightInc)
			} else {
				newHeight += (2 * heightInc)
			}
		default:
			if x > SCR_WIDTH/2 {
				newHeight += (2 * heightInc)
			} else {
				newHeight -= (2 * heightInc)
			}
		}

		buildWidth := defBuildWidth + randInt(0, defBuildWidth)
		if buildWidth+x > SCR_WIDTH {
			buildWidth = SCR_WIDTH - x - 2
		}

		buildHeight := randInt(heightInc, randomHeightDiff) + newHeight

		if bottomLine-buildHeight <= gHeight {
			buildHeight = gHeight
		}

		buildingColor := BUILDING_COLORS[rand.Intn(len(BUILDING_COLORS))]

		px := x + 1
		py := bottomLine - (buildHeight + 1)
		screenSurf.fillRect(buildingColor, image.Rect(px, py, px+buildWidth-1, py+buildHeight-1))

		buildingCoords = append(buildingCoords, image.Pt(x, bottomLine-buildHeight))

		for winx := 3; winx < buildWidth-windowSpacingX+windowWidth; winx += windowSpacingX {
			for winy := 3; winy < buildHeight-windowSpacingY; winy += windowSpacingY {
				var winColor color.RGBA

				if randInt(1, 4) == 1 {
					winColor = DARK_WINDOW
				} else {
					winColor = LIGHT_WINDOW
				}

				px := x + 1 + winx
				py := (bottomLine - buildHeight) + 1 + winy
				screenSurf.fillRect(winColor, image.Rect(px, py, px+windowWidth, py+windowHeight))
			}
		}

		x += buildWidth
	}

	return
}
