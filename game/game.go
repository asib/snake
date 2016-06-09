package game

const GAME_TITLE = "Snake"

type GameState uint8

const (
	Menu GameState = iota
	Play
	Pause
)

type Game struct {
	Width, Height int
	Title         string
	State         GameState
}

func Create(w, h int) *Game {
	return &Game{
		Width:  w,
		Height: h,
		Title:  GAME_TITLE,
		State:  Menu,
		Menu: &Menu{
			Highlighted: Play,
		},
	}
}

func (g *Game) Run() {

}

func (g *Game) Init() (err error) {
	return
}

func (g *Game) Deinit() {
}
