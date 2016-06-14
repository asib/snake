package play

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	keys "github.com/asib/keycodes"
	"github.com/asib/snake/system/draw"
	"github.com/asib/snake/system/draw/renderer"
)

const (
	TileSize            = 10
	timeStep            = 80
	appleScore          = 10
	appleExtension      = 3
	directionChanBuffer = 50
	scoreFont           = "TravelingTypewriter.ttf"
	scoreFontSize       = 14
	gameOverFont        = "TravelingTypewriter.ttf"
	gameOverFontSize    = 22
	highscoreFilename   = "highscore"
)

type Direction uint8

const (
	up Direction = iota
	down
	left
	right
	none
)

// Using a directionChange channel so that you can only change direction
// once per update tick (as before you were able to reverse into yourself
// if you pressed the keys quick enough).
type Play struct {
	lastUpdate      time.Time
	paused          bool
	godmode         bool
	gameOver        bool
	highscore       uint
	saveHighscore   func()
	score           uint
	apple           *draw.Point
	snake           []*draw.Point
	direction       Direction
	directionChange chan Direction
	extend          int
	tileW, tileH    int
}

func Create(w, h int) *Play {
	rand.Seed(time.Now().Unix())
	return &Play{
		lastUpdate:      time.Now(),
		paused:          false,
		godmode:         false,
		gameOver:        false,
		highscore:       0,
		score:           0,
		apple:           nil,
		snake:           make([]*draw.Point, 0),
		direction:       none,
		directionChange: make(chan Direction, directionChanBuffer),
		extend:          0,
		tileW:           w / TileSize,
		tileH:           h / TileSize,
	}
}

func (p *Play) Init() error {
	p.genApple()
	p.genSnake()
	return p.loadHighscore()
}

func (p *Play) loadHighscore() (err error) {
	data, err := ioutil.ReadFile(highscoreFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return
	}
	decodedData := make([]byte, hex.DecodedLen(len(data)))
	_, err = hex.Decode(decodedData, data)
	if err != nil {
		return
	}

	hs, err := strconv.Atoi(string(decodedData))
	if err != nil {
		return
	}

	p.highscore = uint(hs)
	return
}

func (p *Play) SaveHighscore() error {
	// Include this check incase of accidental calls
	if p.score > p.highscore {
		s := strconv.Itoa(int(p.score))
		data := make([]byte, hex.EncodedLen(len([]byte(s))))
		hex.Encode(data, []byte(s))
		return ioutil.WriteFile(highscoreFilename, data, 0777)
	}

	return nil
}

func (p *Play) Run(r renderer.Renderer, w, h int) error {
	// We wan't to draw before we check for collision, so we run update, which
	// itself only runs every timestep, and then update returns true if it
	// actually did an update. When it's true, we then also should check for
	// collisions.
	updated := false
	if !p.paused && !p.gameOver {
		updated = p.update()
	}

	if err := p.draw(r); err != nil {
		return err
	}

	if updated && !p.paused && !p.gameOver {
		p.checkCollision()
	}
	return nil
}

func (p *Play) update() bool {
	if time.Since(p.lastUpdate) >= (timeStep * time.Millisecond) {
		// Update direction if necessary
		select {
		case d := <-p.directionChange:
			p.direction = d
		default:
		}

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
		return true
	}

	return false
}

func (p *Play) checkCollision() {
	if !p.godmode {
		head := p.snake[0]
		for i := range p.snake {
			if i != 0 && p.snake[i].X == head.X && p.snake[i].Y == head.Y {
				p.gameOver = true
				if err := p.SaveHighscore(); err != nil {
					log.Println(err)
				}
			}
		}
	}

	// If there's a collision, set p.extend to the appleExtension constant
	// Then, it gets decremented every update tick
	// While it's nonzero, the nextPosition method won't delete the back of the
	// snake each update.
	if p.snake[0].X == p.apple.X && p.snake[0].Y == p.apple.Y {
		p.score += appleScore
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
	if !p.gameOver {
		if err = p.drawScore(r); err != nil {
			return
		}
	}
	if err = p.drawApple(r); err != nil {
		return
	}
	if err = p.drawSnake(r); err != nil {
		return
	}

	if p.gameOver {
		return p.drawGameOver(r)
	}

	return
}

func (p *Play) drawScore(r renderer.Renderer) (err error) {
	scoreColor := draw.CreateColor(0, 0xff, 0, 0xff)

	x, y := 5, 5
	return r.DrawText(scoreFont, scoreFontSize, fmt.Sprintf("Score: %d", p.score),
		scoreColor, x, y)
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

func (p *Play) drawGameOver(r renderer.Renderer) (err error) {
	white := draw.CreateColor(0xff, 0xff, 0xff, 0xff)
	black := draw.CreateColor(0, 0, 0, 0xff)

	l1w, l1h, err := r.SizeUTF8(gameOverFont, gameOverFontSize, "Game over!")
	if err != nil {
		return
	}
	l1x := p.pixelWidth()/2 - l1w/2
	l1y := 10
	err = r.DrawTextBg(gameOverFont, gameOverFontSize, "Game over!", white, black, l1x, l1y)
	if err != nil {
		return
	}

	l2String := fmt.Sprintf("You scored %d", p.score)
	l2w, l2h, err := r.SizeUTF8(gameOverFont, gameOverFontSize, l2String)
	if err != nil {
		return
	}
	l2x := p.pixelWidth()/2 - l2w/2
	l2y := l1y + l1h + 10
	err = r.DrawTextBg(gameOverFont, gameOverFontSize, l2String, white, black, l2x, l2y)
	if err != nil {
		return
	}

	l3w, l3h, err := r.SizeUTF8(gameOverFont, gameOverFontSize, "New highscore!")
	if err != nil {
		return
	}
	l3x := p.pixelWidth()/2 - l3w/2
	l3y := l2y + l2h + 10
	if p.score > p.highscore {
		err = r.DrawTextBg(gameOverFont, gameOverFontSize, "New highscore!", white, black, l3x, l3y)
		if err != nil {
			return
		}
	}

	l4w, _, err := r.SizeUTF8(gameOverFont, gameOverFontSize, "Press Enter to start a new game.")
	if err != nil {
		return
	}
	l4x := p.pixelWidth()/2 - l4w/2
	l4y := l3y + l3h + 10
	err = r.DrawTextBg(gameOverFont, gameOverFontSize, "Press Enter to start a new game.",
		white, black, l4x, l4y)
	if err != nil {
		return
	}

	return
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

func (p *Play) restart() {
	rand.Seed(time.Now().Unix())
	p.lastUpdate = time.Now()
	p.paused = false
	p.godmode = false
	p.gameOver = false
	p.highscore = 0
	p.score = 0
	p.snake = make([]*draw.Point, 0)
	p.direction = none
	p.directionChange = make(chan Direction, directionChanBuffer)
	p.extend = 0
	p.Init()
}

func (p *Play) KeyPress(k keys.Keycode) {
	// Don't allow the player to reverse direction when the snake is more than
	// 1 square long
	switch k {
	case keys.K_RETURN:
		if p.gameOver {
			p.restart()
		}
	case keys.K_p:
		p.pause()
	case keys.K_g:
		p.toggleGod()
	case keys.K_UP:
		if !(len(p.snake) > 1 && p.direction == down) {
			p.directionChange <- up
		}
	case keys.K_DOWN:
		if !(len(p.snake) > 1 && p.direction == up) {
			p.directionChange <- down
		}
	case keys.K_LEFT:
		if !(len(p.snake) > 1 && p.direction == right) {
			p.directionChange <- left
		}
	case keys.K_RIGHT:
		if !(len(p.snake) > 1 && p.direction == left) {
			p.directionChange <- right
		}
	}
}

func (p *Play) Score() uint {
	return p.score
}

func (p *Play) Highscore() uint {
	return p.highscore
}

func (p *Play) pause() {
	p.paused = !p.paused
}

func (p *Play) toggleGod() {
	p.godmode = !p.godmode
}

func (p *Play) pixelWidth() int {
	return p.tileW * TileSize
}

func (p *Play) pixelHeight() int {
	return p.tileH * TileSize
}
