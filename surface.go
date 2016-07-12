package main

import (
	"image"
	"image/color"
	"image/draw"
	"strings"
)

func makeSurfaceFromASCII(ascii string, fgColor, bgColor color.RGBA) *Surface {
	lines := strings.Split(ascii, "\n")
	height := len(lines)
	width := 0
	for _, line := range lines {
		width = max(width, len(line))
	}

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), image.NewUniform(bgColor), image.ZP, draw.Src)
	for y := range lines {
		for x := range lines[y] {
			if lines[y][x] == 'X' {
				rgba.Set(x, y, fgColor)
			}
		}
	}
	return &Surface{rgba}
}

type Surface struct {
	*image.RGBA
}

func newSurface(width, height int) *Surface {
	return &Surface{image.NewRGBA(image.Rect(0, 0, width, height))}
}

func (s *Surface) width() int {
	return s.Bounds().Dx()
}

func (s *Surface) height() int {
	return s.Bounds().Dy()
}

func (s *Surface) fill(c color.RGBA) {
	draw.Draw(s.RGBA, s.Bounds(), image.NewUniform(c), image.ZP, draw.Src)
}

func (s *Surface) fillRect(c color.RGBA, r image.Rectangle) {
	draw.Draw(s.RGBA, r, image.NewUniform(c), image.ZP, draw.Src)
}

func (s *Surface) blit(p *Surface, x, y int) {
	r := image.Rect(x, y, x+p.width(), y+p.height())
	draw.Draw(s.RGBA, r, p.RGBA, image.ZP, draw.Over)
}

func (s *Surface) circle(c color.RGBA, x, y, r int) {
	rect := image.Rect(x-r, y-r, x+r, y+r)
	sp := image.Pt(x-r, y-r)
	circle := &circle{image.Pt(x, y), r, c}
	draw.Draw(s.RGBA, rect, circle, sp, draw.Over)
}

func (s *Surface) line(c color.RGBA, x0, y0, x1, y1 int) {
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}
	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	for {
		s.SetRGBA(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

type circle struct {
	p image.Point
	r int
	c color.RGBA
}

func (c *circle) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return c.c
	}
	return color.Transparent
}
