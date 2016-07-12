package main

import (
	"image/color"

	"github.com/qeedquan/go-media/sdl"
)

const (
	RIGHT = 0
	UP    = 1
	LEFT  = 2
	DOWN  = 3
)

var BAN_COLOR = color.RGBA{255, 255, 82, 255}

var (
	BAN_RIGHT_SURF = makeSurfaceFromASCII(BAN_RIGHT_ASCII, BAN_COLOR, SKY_COLOR)
	BAN_LEFT_SURF  = makeSurfaceFromASCII(BAN_LEFT_ASCII, BAN_COLOR, SKY_COLOR)
	BAN_UP_SURF    = makeSurfaceFromASCII(BAN_UP_ASCII, BAN_COLOR, SKY_COLOR)
	BAN_DOWN_SURF  = makeSurfaceFromASCII(BAN_DOWN_ASCII, BAN_COLOR, SKY_COLOR)
)

func getBananaRect(x, y, shape int) sdl.Rect {
	switch shape {
	case UP:
		return sdl.Rect{int32(x), int32(y), int32(BAN_UP_SURF.width()), int32(BAN_UP_SURF.height())}
	case DOWN:
		return sdl.Rect{int32(x), int32(y), int32(BAN_DOWN_SURF.width()), int32(BAN_DOWN_SURF.height())}
	case LEFT:
		return sdl.Rect{int32(x), int32(y), int32(BAN_LEFT_SURF.width()), int32(BAN_LEFT_SURF.height())}
	case RIGHT:
		return sdl.Rect{int32(x), int32(y), int32(BAN_RIGHT_SURF.width()), int32(BAN_RIGHT_SURF.height())}
	default:
		panic("unreachable")
	}
}

func nextBananaShape(orient int) int {
	return (orient + 1) % 4
}

const BAN_RIGHT_ASCII = `
     XX
    XXX
   XXX
   XXX
   XXX
   XXX
   XXX
    XXX
     XX
`

const BAN_LEFT_ASCII = `
XX
XXX
 XXX
 XXX
 XXX
 XXX
 XXX
XXX
XX
`

const BAN_UP_ASCII = `
XX     XX
XXXXXXXXX
 XXXXXXX
  XXXXX
`

const BAN_DOWN_ASCII = `
  XXXXX
 XXXXXXX
XXXXXXXXX
XX     XX
`
