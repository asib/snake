package renderer

import (
	"github.com/asib/snake/system/draw"
	"github.com/asib/snake/system/draw/fonts"
	"github.com/veandco/go-sdl2/sdl"
)

type Renderer interface {
	Present()
	Clear(c *draw.Color) error
	FillRect(rect *draw.Rect, color *draw.Color) error
	SizeUTF8(path string, size int, text string) (int, int, error)
	DrawText(path string, size int, text string, color *draw.Color,
		x, y int) error
}

type SDLRenderer struct {
	sdlRenderer *sdl.Renderer
	fontManager *fonts.FontManager
}

func CreateSDLRenderer(r *sdl.Renderer) *SDLRenderer {
	return &SDLRenderer{
		sdlRenderer: r,
		fontManager: fonts.CreateManager(),
	}
}

func (r *SDLRenderer) Present() {
	r.sdlRenderer.Present()
}

func (r *SDLRenderer) Clear(c *draw.Color) error {
	err := r.sdlRenderer.SetDrawColor(c.R, c.G, c.B, c.A)
	if err != nil {
		return err
	}
	return r.sdlRenderer.Clear()
}

func (r *SDLRenderer) SizeUTF8(path string, size int, text string) (int, int,
	error) {
	font, err := r.fontManager.GetFont(path, size)
	if err != nil {
		return 0, 0, err
	}
	return font.SizeUTF8(text)
}

/*
 *func (r *SDLRenderer) FillRect(rect *draw.Rect, color *draw.Color) (err error) {
 *  src := rect.ToSDL()
 *  src.X = 0
 *  src.Y = 0
 *
 *  t, err := r.sdlRenderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
 *    sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
 *  if err != nil {
 *    return
 *  }
 *  defer t.Destroy()
 *
 *  err = t.SetBlendMode(sdl.BLENDMODE_BLEND)
 *  if err != nil {
 *    return
 *  }
 *
 *  err = r.sdlRenderer.SetRenderTarget(t)
 *  if err != nil {
 *    return
 *  }
 *  err = r.sdlRenderer.SetDrawColor(color.R, color.G, color.B, color.A)
 *  if err != nil {
 *    return
 *  }
 *  err = r.sdlRenderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
 *  if err != nil {
 *    return err
 *  }
 *
 *  r.sdlRenderer.FillRect(nil)
 *  err = r.sdlRenderer.SetRenderTarget(nil)
 *  if err != nil {
 *    return err
 *  }
 *  //log.Println(src, rect)
 *  r.sdlRenderer.Copy(t, src, rect.ToSDL())
 *  return
 *}
 */

func (r *SDLRenderer) FillRect(rect *draw.Rect, color *draw.Color) error {
	err := r.sdlRenderer.SetDrawColor(color.R, color.G, color.B, color.A)
	if err != nil {
		return err
	}

	return r.sdlRenderer.FillRect(rect.ToSDL())
}

func (r *SDLRenderer) DrawText(path string, size int, text string,
	color *draw.Color, x, y int) error {
	font, err := r.fontManager.GetFont(path, size)
	if err != nil {
		return err
	}

	surf, err := font.RenderUTF8_Blended(text, color.ToSDL())
	if err != nil {
		return err
	}
	defer surf.Free()
	src := new(sdl.Rect)
	surf.GetClipRect(src)
	dst := &sdl.Rect{
		X: int32(x),
		Y: int32(y),
		W: src.W,
		H: src.H,
	}

	texture, err := r.sdlRenderer.CreateTextureFromSurface(surf)
	if err != nil {
		return err
	}
	defer texture.Destroy()
	err = texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		return err
	}
	err = r.sdlRenderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		return err
	}

	r.sdlRenderer.Copy(texture, src, dst)
	return nil
}

func (r *SDLRenderer) Deinit() {
	r.fontManager.Deinit()
	r.sdlRenderer.Destroy()
}
