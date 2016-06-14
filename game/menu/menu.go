package menu

import (
	"fmt"

	"github.com/asib/snake/system/draw"
	"github.com/asib/snake/system/draw/renderer"
)

const (
	GameTitle            = "Snake"
	titleFont            = "Shoguns Clan.ttf"
	titleFontSize        = 120
	instructionsFont     = "TravelingTypewriter.ttf"
	instructionsFontSize = 22
)

type Menu struct {
	highscore uint
}

func Create() *Menu {
	return &Menu{
		highscore: 0,
	}
}

func (m *Menu) Init(hs uint) {
	m.highscore = hs
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

	titlew, titleh, err := r.SizeUTF8(titleFont, titleFontSize, GameTitle)
	if err != nil {
		return
	}
	titlex := w/2 - titlew/2
	titley := 10
	err = r.DrawText(titleFont, titleFontSize, "Snake", green, titlex, titley)
	if err != nil {
		return
	}

	highscoreString := fmt.Sprintf("Highscore: %d", m.highscore)
	l1w, l1h, err := r.SizeUTF8(instructionsFont, instructionsFontSize, highscoreString)
	if err != nil {
		return
	}
	l1x := w/2 - l1w/2
	l1y := titley + titleh + 40
	err = r.DrawText(instructionsFont, instructionsFontSize, highscoreString,
		green, l1x, l1y)
	if err != nil {
		return
	}

	l2w, l2h, err := r.SizeUTF8(instructionsFont, instructionsFontSize, "Press Enter to start.")
	if err != nil {
		return
	}
	l2x := w/2 - l2w/2
	l2y := l1y + l1h + 40
	err = r.DrawText(instructionsFont, instructionsFontSize, "Press Enter to start.",
		green, l2x, l2y)
	if err != nil {
		return
	}

	l3w, _, err := r.SizeUTF8(instructionsFont, instructionsFontSize, "Press ESC to quit.")
	if err != nil {
		return
	}
	l3x := w/2 - l3w/2
	l3y := l2y + l2h + 10
	err = r.DrawText(instructionsFont, instructionsFontSize, "Press ESC to quit.",
		green, l3x, l3y)
	if err != nil {
		return
	}

	return
}
