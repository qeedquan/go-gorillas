package main

var (
	STAR_SURF = makeSurfaceFromASCII(STAR_ASCII, DARK_RED_COLOR, BLACK_COLOR)
)

func drawStars(screenSurf *Surface, vertAdj, horAdj int) {
	for i := 0; i < 16; i++ {
		// draw top row of stars
		screenSurf.blit(STAR_SURF, 2+(((3-vertAdj)+i*4)*STAR_SURF.width()), 3)

		// draw bottom row of stars
		screenSurf.blit(STAR_SURF, 2+((vertAdj+i*4)*STAR_SURF.width()), SCR_HEIGHT-7-STAR_SURF.height())
	}
	for i := 0; i < 4; i++ {
		// draw left column of stars going down
		screenSurf.blit(STAR_SURF, 5, 6+STAR_SURF.height()+(horAdj+i*84))

		// draw right column of stars going up
		screenSurf.blit(STAR_SURF, SCR_WIDTH-5-STAR_SURF.width(), (SCR_HEIGHT - (6 + STAR_SURF.height() + (horAdj + i*84))))
	}
}

const STAR_ASCII = `


   XX  XX
    XXXX
  XXXXXXXX
    XXXX
   XX  XX
`
