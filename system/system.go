package system

import (
	"log"

	"github.com/asib/keycodes"
	"github.com/asib/snake/game"
	"github.com/asib/snake/system/draw/renderer"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type System struct {
	Debug    bool
	Window   *sdl.Window
	Screen   *sdl.Surface
	Renderer *renderer.SDLRenderer
	g        *game.Game
	Event    sdl.Event
}

func Create(debug bool, w, h int) *System {
	return &System{
		Debug:    debug,
		Window:   nil,
		Screen:   nil,
		Renderer: nil,
		g:        game.Create(debug, w, h),
		Event:    nil,
	}
}

func (s *System) Run() {
	for s.g.Running {
		for s.Event = sdl.PollEvent(); s.Event != nil; s.Event = sdl.PollEvent() {
			switch t := s.Event.(type) {
			case *sdl.QuitEvent:
				s.g.Running = false
				if err := s.g.Play.SaveHighscore(); err != nil {
					log.Println(err)
				}
			case *sdl.KeyDownEvent:
				s.g.KeyPress(keycodes.FromSDL(t.Keysym.Sym))
			}
		}

		s.g.Run(s.Renderer)
	}
}

func (s *System) Init() (err error) {
	// Init SDL, create renderer, load fonts
	sdl.Init(sdl.INIT_EVERYTHING)

	win, err := sdl.CreateWindow(s.g.Title, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, s.g.Width, s.g.Height, sdl.WINDOW_SHOWN)
	if err != nil {
		return
	}
	s.Window = win

	r, err := sdl.CreateRenderer(win, -1, sdl.RENDERER_SOFTWARE|sdl.RENDERER_TARGETTEXTURE)
	if err != nil {
		return
	}
	s.Renderer = renderer.CreateSDLRenderer(r)

	scr, err := win.GetSurface()
	if err != nil {
		return
	}
	s.Screen = scr

	if err = ttf.Init(); err != nil {
		return
	}

	// Allow game to init
	err = s.g.Init()

	return
}

func (s *System) Deinit() {
	// Release SDL resources
	s.g.Deinit()
	s.Renderer.Deinit()
	s.Window.Destroy()
	sdl.Quit()
}
