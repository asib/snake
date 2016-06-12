package menu

import (
	"github.com/asib/snake/system/draw"
	"github.com/asib/snake/system/draw/renderer"
)

const (
	GameTitle = "Snake"
	menuFont  = "Shoguns Clan.ttf"
	titleSize = 120
)

type Menu struct {
}

func Create() *Menu {
	return &Menu{}
}

func (m *Menu) Run(r renderer.Renderer, w, h int) error {
	m.update()
	return m.draw(r, w, h)
}

func (m *Menu) update() {
}

func (m *Menu) draw(r renderer.Renderer, w, h int) (err error) {
	green := draw.CreateColor(0, 0xff, 0, 0xff)
	//white := draw.CreateColor(0xff, 0xff, 0xff, 0xaa)
	black := draw.CreateColor(0, 0, 0, 0xff)

	r.FillRect(draw.CreateRect(0, 0, w, h), black)

	titlew, titleh, err := r.SizeUTF8(menuFont, titleSize, GameTitle)
	if err != nil {
		return
	}
	titlex := w/2 - titlew/2
	titley := 10
	err = r.DrawText("Shoguns Clan.ttf", 120, "Snake", green, titlex, titley)
	if err != nil {
		return
	}

	l1w, l1h, err := r.SizeUTF8("TravelingTypewriter.ttf", 22, "Highscore:")
	if err != nil {
		return
	}
	l1x := w/2 - l1w/2
	l1y := titley + titleh + 40
	err = r.DrawText("TravelingTypewriter.ttf", 22, "Highscore:",
		green, l1x, l1y)
	if err != nil {
		return
	}

	l2w, l2h, err := r.SizeUTF8("TravelingTypewriter.ttf", 22, "Press Enter to start.")
	if err != nil {
		return
	}
	l2x := w/2 - l2w/2
	l2y := l1y + l1h + 40
	err = r.DrawText("TravelingTypewriter.ttf", 22, "Press Enter to start.",
		green, l2x, l2y)
	if err != nil {
		return
	}

	l3w, _, err := r.SizeUTF8("TravelingTypewriter.ttf", 22, "Press ESC to quit.")
	if err != nil {
		return
	}
	l3x := w/2 - l3w/2
	l3y := l2y + l2h + 10
	err = r.DrawText("TravelingTypewriter.ttf", 22, "Press ESC to quit.",
		green, l3x, l3y)
	if err != nil {
		return
	}

	return
}
