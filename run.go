package main

import sys "github.com/asib/snake/system"

func initialize(s *sys.System) (err error) {
	if err = s.Init(); err != nil {
		return
	}

	return
}

func main() {
	sys := sys.Create(640, 480)
	if err := initialize(sys); err != nil {
		panic(err)
	}

	sys.Run()

	/*
	 *  screen.FillRect(&sdl.Rect{0, 0, 640, 480}, 0xffffffff)
	 *
	 *  src := new(sdl.Rect)
	 *  tSurf.GetClipRect(src)
	 *  if err = tSurf.Blit(src, screen, &sdl.Rect{10, 10, src.W, src.H}); err != nil {
	 *    panic(err)
	 *  }
	 *  window.UpdateSurface()
	 */
}
