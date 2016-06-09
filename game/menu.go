package game

import (
	keys "github.com/asib/keycodes"
)

type MenuItem uint8

const (
	Play MenuItem = iota
	Quit
)

type Menu struct {
	Highlighted MenuItem
}

func (m *Menu) Reset() {
	m.Highlighted = Play
}

func (m *Menu) Draw() {
}

func (m *Menu) KeyPress(k keys.Keycode) {

}
