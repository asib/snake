package play

import (
	"math/rand"
	"time"

	keys "github.com/asib/keycodes"
	"github.com/asib/snake/system/draw"
	"github.com/asib/snake/system/draw/renderer"
)

const (
	TileSize       = 10
	timeStep       = 80
	appleExtension = 3
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
	extend       int
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
		extend:     0,
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

		p.checkCollision()

		p.lastUpdate = time.Now()
	}
}

func (p *Play) checkCollision() {
	// If there's a collision, set p.extend to the appleExtension constant
	// Then, it gets decremented every update tick
	// While it's nonzero, the nextPosition method won't delete the back of the
	// snake each update.
	if p.snake[0].X == p.apple.X && p.snake[0].Y == p.apple.Y {
		p.extend = appleExtension
		p.genApple()
	}
}

func (p *Play) nextPosition(dx, dy int) {
	// Get next position, modulo width/height to allow for teleporting
	front := p.snake[0]
	nx := (front.X + dx) % p.tileW
	ny := (front.Y + dy) % p.tileH
	if nx < 0 {
		nx += p.tileW
	}
	if ny < 0 {
		ny += p.tileH
	}

	if p.extend != 0 {
		// If p.extend is nonzero, don't delete back of snake on update
		p.extend -= 1
		p.snake = append(p.snake, nil)
		copy(p.snake[1:], p.snake[:])
		p.snake[0] = draw.CreatePoint(nx, ny)
	} else {
		// Else, delete back as usual
		p.snake = append(p.snake[:len(p.snake)-1], nil)
		copy(p.snake[1:], p.snake[:])
		p.snake[0] = draw.CreatePoint(nx, ny)
	}
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
	// Make sure the apple isn't on a square that the snake's body occupies
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
	// Make sure the snake isn't on a square that the apple occupies
	for (p.apple.X == x) && (p.apple.Y == y) {
		x, y = rand.Intn(p.tileW), rand.Intn(p.tileH)
	}
	p.snake = append(p.snake, draw.CreatePoint(x, y))
}

func (p *Play) KeyPress(k keys.Keycode) {
	// Don't allow the player to reverse direction when the snake is more than
	// 1 square long
	switch k {
	case keys.K_UP:
		if !(len(p.snake) > 1 && p.direction == down) {
			p.direction = up
		}
	case keys.K_DOWN:
		if !(len(p.snake) > 1 && p.direction == up) {
			p.direction = down
		}
	case keys.K_LEFT:
		if !(len(p.snake) > 1 && p.direction == right) {
			p.direction = left
		}
	case keys.K_RIGHT:
		if !(len(p.snake) > 1 && p.direction == left) {
			p.direction = right
		}
	}
}
