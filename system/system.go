package system

import (
	"fmt"

	"github.com/asib/snake/game"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type System struct {
	Window  *sdl.Window
	Screen  *sdl.Surface
	Font    *ttf.Font
	g       *game.Game
	Running bool
	Event   sdl.Event
}

func Create(w, h int) *System {
	return &System{
		Window:  nil,
		Screen:  nil,
		Font:    nil,
		g:       game.Create(w, h),
		Running: true,
		Event:   nil,
	}
}

func (s *System) Run() {
	for s.Running {
		for s.Event = sdl.PollEvent(); s.Event != nil; s.Event = sdl.PollEvent() {
			switch t := s.Event.(type) {
			case *sdl.QuitEvent:
				s.Running = false
			case *sdl.KeyDownEvent:
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					s.Running = false
				case sdl.K_RETURN:
					fmt.Println("Bla")
				}
			}
		}

		s.g.Run()
	}
}

func (s *System) Init() (err error) {
	sdl.Init(sdl.INIT_EVERYTHING)

	win, err := sdl.CreateWindow(s.g.Title, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, s.g.Width, s.g.Height, sdl.WINDOW_SHOWN)
	if err != nil {
		return
	}
	s.Window = win

	scr, err := win.GetSurface()
	if err != nil {
		return
	}
	s.Screen = scr

	if err = ttf.Init(); err != nil {
		return
	}

	font, err := ttf.OpenFont("TimesNewRoman.ttf", 64)
	if err != nil {
		return
	}
	s.Font = font

	/*
	 *tSurf, err := calib.RenderUTF8_Blended("Test text", sdl.Color{0xff, 0, 0, 0xff})
	 *if err != nil {
	 *  return
	 *}
	 */

	err = s.g.Init()

	return
}

func (s *System) Deinit() {
	s.g.Deinit()
	s.Window.Destroy()
	sdl.Quit()
}
