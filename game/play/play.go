package play

import (
	"math/rand"
	"time"

	keys "github.com/asib/keycodes"
	"github.com/asib/snake/system/draw"
	"github.com/asib/snake/system/draw/renderer"
)

const (
	TileSize = 10
	timeStep = 100
)

type Direction uint8

const (
	up Direction = iota
	down
	left
	right
	none
)

type Play struct {
	lastUpdate   time.Time
	score        uint
	apple        *draw.Point
	snake        []*draw.Point
	direction    Direction
	tileW, tileH int
}

func Create(w, h int) *Play {
	tw := w / TileSize
	th := h / TileSize
	p := &Play{
		lastUpdate: time.Now(),
		score:      0,
		apple:      nil,
		snake:      make([]*draw.Point, 0),
		direction:  none,
		tileW:      tw,
		tileH:      th,
	}

	p.genApple()
	p.genSnake()

	return p
}

func (p *Play) Run(r renderer.Renderer, w, h int) error {
	// Don't forget to check for full board
	p.update()
	return p.draw(r)
}

func (p *Play) update() {
	if time.Since(p.lastUpdate) >= (timeStep * time.Millisecond) {
		switch p.direction {
		case up:
			p.nextPosition(0, -1)
		case down:
			p.nextPosition(0, 1)
		case left:
			p.nextPosition(-1, 0)
		case right:
			p.nextPosition(1, 0)
		}

		p.lastUpdate = time.Now()
	}
}

func (p *Play) nextPosition(dx, dy int) {
	front := p.snake[0]
	nx := (front.X + dx) % p.tileW
	ny := (front.Y + dy) % p.tileH
	if nx < 0 {
		nx += p.tileW
	}
	if ny < 0 {
		ny += p.tileH
	}

	p.snake = append(p.snake[:len(p.snake)-1], nil)
	copy(p.snake[1:], p.snake[:])
	p.snake[0] = draw.CreatePoint(nx, ny)
}

func (p *Play) draw(r renderer.Renderer) (err error) {
	if err = p.drawSnake(r); err != nil {
		return
	}
	return p.drawApple(r)
}

func (p *Play) drawSnake(r renderer.Renderer) (err error) {
	snakeColor := draw.CreateColor(0, 0xff, 0, 0xff)

	for i := range p.snake {
		err = r.FillRect(draw.CreateRect(p.snake[i].X*TileSize,
			p.snake[i].Y*TileSize, TileSize, TileSize), snakeColor)
		if err != nil {
			return
		}
	}

	return
}

func (p *Play) drawApple(r renderer.Renderer) error {
	appleColor := draw.CreateColor(0xff, 0, 0, 0xff)
	return r.FillRect(draw.CreateRect(p.apple.X*TileSize, p.apple.Y*TileSize, TileSize, TileSize),
		appleColor)
}

func (p *Play) genApple() {
	var x, y int
	overlapping := true
	for overlapping {
		overlapping = false
		x, y = rand.Intn(p.tileW), rand.Intn(p.tileH)
		for i := range p.snake {
			if (x == p.snake[i].X) && (y == p.snake[i].Y) {
				overlapping = true
				break
			}
		}
	}

	p.apple = draw.CreatePoint(x, y)
}

func (p *Play) genSnake() {
	x, y := rand.Intn(p.tileW), rand.Intn(p.tileH)
	for (p.apple.X == x) && (p.apple.Y == y) {
		x, y = rand.Intn(p.tileW), rand.Intn(p.tileH)
	}
	p.snake = append(p.snake, draw.CreatePoint(x, y))
}

func (p *Play) KeyPress(k keys.Keycode) {
	switch k {
	case keys.K_UP:
		p.direction = up
	case keys.K_DOWN:
		p.direction = down
	case keys.K_LEFT:
		p.direction = left
	case keys.K_RIGHT:
		p.direction = right
	}
}
