package game

import (
	"log"

	"github.com/asib/keycodes"
	"github.com/asib/snake/game/menu"
	"github.com/asib/snake/game/play"
	"github.com/asib/snake/system/draw"
	"github.com/asib/snake/system/draw/renderer"
)

type GameState uint8

const (
	Menu GameState = iota
	Play
	Pause
)

type Game struct {
	Debug         bool
	Width, Height int
	Running       bool
	Title         string
	State         GameState
	Menu          *menu.Menu
	Play          *play.Play
}

func Create(debug bool, w, h int) *Game {
	return &Game{
		Debug:   debug,
		Width:   w,
		Height:  h,
		Running: true,
		Title:   menu.GameTitle,
		State:   Menu,
		Menu:    menu.Create(),
		Play:    nil,
	}
}

func (g *Game) Run(r renderer.Renderer) {
	err := r.Clear(draw.CreateColor(0, 0, 0, 0))
	if err != nil {
		log.Println("Clear failure: " + err.Error())
	}

	switch g.State {
	case Menu:
		if err = g.Menu.Run(r, g.Width, g.Height); err != nil {
			panic(err)
		}
	case Play:
		if err = g.Play.Run(r, g.Width, g.Height); err != nil {
			panic(err)
		}
	}

	r.Present()
}

func (g *Game) KeyPress(k keycodes.Keycode) {
	switch g.State {
	case Menu:
		if k == keycodes.K_RETURN {
			g.State = Play
			g.Play = play.Create(g.Width, g.Height)
		} else if k == keycodes.K_ESCAPE {
			g.Running = false
		}
	case Play:
		g.Play.KeyPress(k)
	}
}

func (g *Game) Init() (err error) {
	return
}

func (g *Game) Deinit() {
}
