package main

import (
	"image"
	"image/color"
	"time"
)

const (
	BOTH_ARMS_DOWN = 0
	LEFT_ARM_UP    = 1
	RIGHT_ARM_UP   = 2
)

var (
	GOR_COLOR = color.RGBA{255, 170, 82, 255}
)

var (
	GOR_DOWN_SURF  = makeSurfaceFromASCII(GOR_DOWN_ASCII, GOR_COLOR, SKY_COLOR)
	GOR_LEFT_SURF  = makeSurfaceFromASCII(GOR_LEFT_ASCII, GOR_COLOR, SKY_COLOR)
	GOR_RIGHT_SURF = makeSurfaceFromASCII(GOR_RIGHT_ASCII, GOR_COLOR, SKY_COLOR)
)

const GOR_EXPLOSION_SIZE = 30

func placeGorillas(buildCoords []image.Point) (gorPos [2]image.Point) {
	xAdj := GOR_DOWN_SURF.width() / 2
	yAdj := GOR_DOWN_SURF.height()

	for i := range gorPos {
		var buildNum int
		if i == 0 {
			buildNum = randInt(1, 2)
		} else {
			buildNum = randInt(len(buildCoords)-3, len(buildCoords)-2)
		}

		buildWidth := buildCoords[buildNum+1].X - buildCoords[buildNum].X
		gorPos[i] = image.Pt(buildCoords[buildNum].X+int(buildWidth/2)-xAdj, buildCoords[buildNum].Y-yAdj-1)
	}
	return
}

func drawGorilla(screenSurf *Surface, x, y, arms int) {
	var gorSurf *Surface
	switch arms {
	case BOTH_ARMS_DOWN:
		gorSurf = GOR_DOWN_SURF
	case LEFT_ARM_UP:
		gorSurf = GOR_LEFT_SURF
	case RIGHT_ARM_UP:
		gorSurf = GOR_RIGHT_SURF
	}
	screenSurf.blit(gorSurf, x, y)
}

func victoryDance(screenSurf *Surface, x, y int) {
	for i := 0; i < 4; i++ {
		screenSurf.blit(GOR_LEFT_SURF, x, y)
		sleep(300 * time.Millisecond)
		screenSurf.blit(GOR_RIGHT_SURF, x, y)
		sleep(300 * time.Millisecond)
	}
}

const GOR_DOWN_ASCII = `

          XXXXXXXX
          XXXXXXXX
         XX      XX
         XXXXXXXXXX
         XXX  X  XX
          XXXXXXXX
          XXXXXXXX
           XXXXXX
      XXXXXXXXXXXXXXXX
   XXXXXXXXXXXXXXXXXXXXXX
  XXXXXXXXXXXX XXXXXXXXXXX
 XXXXXXXXXXXXX XXXXXXXXXXXX
 XXXXXXXXXXXX X XXXXXXXXXXX
XXXXX XXXXXX XXX XXXXX XXXXX
XXXXX XXX   XXXXX   XX XXXXX
XXXXX   XXXXXXXXXXXX   XXXXX
 XXXXX  XXXXXXXXXXXX  XXXXX
 XXXXX  XXXXXXXXXXXX  XXXXX
  XXXXX XXXXXXXXXXXX XXXXX
   XXXXXXXXXXXXXXXXXXXXXX
       XXXXXXXXXXXXX
     XXXXXX     XXXXXX
     XXXXX       XXXXX
    XXXXX         XXXXX
    XXXXX         XXXXX
    XXXXX         XXXXX
    XXXXX         XXXXX
    XXXXX         XXXXX
     XXXXX       XXXXX
`

const GOR_LEFT_ASCII = `
   XXXXX
  XXXXX   XXXXXXXX
 XXXXX    XXXXXXXX
 XXXXX   XX      XX
XXXXX    XXXXXXXXXX
XXXXX    XXX  X  XX
XXXXX     XXXXXXXX
 XXXXX    XXXXXXXX
 XXXXX     XXXXXX
  XXXXXXXXXXXXXXXXXXXX
   XXXXXXXXXXXXXXXXXXXXXX
      XXXXXXXX XXXXXXXXXXX
      XXXXXXXX XXXXXXXXXXXX
      XXXXXXX X XXXXXXXXXXX
      XXXXXX XXX XXXXX XXXXX
      XXX   XXXXX   XX XXXXX
        XXXXXXXXXXXX   XXXXX
        XXXXXXXXXXXX  XXXXX
        XXXXXXXXXXXX  XXXXX
        XXXXXXXXXXXX XXXXX
       XXXXXXXXXXXXXXXXXX
       XXXXXXXXXXXXX
     XXXXXX     XXXXXX
     XXXXX       XXXXX
    XXXXX         XXXXX
    XXXXX         XXXXX
    XXXXX         XXXXX
    XXXXX         XXXXX
    XXXXX         XXXXX
     XXXXX       XXXXX
`

const GOR_RIGHT_ASCII = `
                    XXXXX
          XXXXXXXX   XXXXX
          XXXXXXXX    XXXXX
         XX      XX   XXXXX
         XXXXXXXXXX    XXXXX
         XXX  X  XX    XXXXX
          XXXXXXXX     XXXXX
          XXXXXXXX    XXXXX
           XXXXXX     XXXXX
      XXXXXXXXXXXXXXXXXXXX
   XXXXXXXXXXXXXXXXXXXXXX
  XXXXXXXXXXXX XXXXXXX
 XXXXXXXXXXXXX XXXXXXX
 XXXXXXXXXXXX X XXXXXX
XXXXX XXXXXX XXX XXXXX
XXXXX XXX   XXXXX   XX
XXXXX   XXXXXXXXXXXX
 XXXXX  XXXXXXXXXXXX
 XXXXX  XXXXXXXXXXXX
  XXXXX XXXXXXXXXXXX
   XXXXXXXXXXXXXXXXX
       XXXXXXXXXXXXX
     XXXXXX     XXXXXX
     XXXXX       XXXXX
    XXXXX         XXXXX
    XXXXX         XXXXX
    XXXXX         XXXXX
    XXXXX         XXXXX
    XXXXX         XXXXX
     XXXXX       XXXXX
`
