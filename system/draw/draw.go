package draw

import "github.com/veandco/go-sdl2/sdl"

type Point struct {
	X int
	Y int
}

func CreatePoint(x, y int) *Point {
	return &Point{
		X: x,
		Y: y,
	}
}

type Rect struct {
	X int
	Y int
	W int
	H int
}

func CreateRect(x, y, w, h int) *Rect {
	return &Rect{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
}

func (r *Rect) ToSDL() *sdl.Rect {
	return &sdl.Rect{
		X: int32(r.X),
		Y: int32(r.Y),
		W: int32(r.W),
		H: int32(r.H),
	}
}

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func CreateColor(r, g, b, a uint8) *Color {
	return &Color{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

func (c *Color) ToSDL() sdl.Color {
	return sdl.Color{
		R: c.R,
		G: c.G,
		B: c.B,
		A: c.A,
	}
}
