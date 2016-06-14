package fonts

import "github.com/veandco/go-sdl2/sdl_ttf"

type fontIdentifier struct {
	path string
	size int
}

type FontManager struct {
	fonts map[fontIdentifier]*ttf.Font
}

func CreateManager() *FontManager {
	return &FontManager{
		fonts: make(map[fontIdentifier]*ttf.Font),
	}
}

func (fm *FontManager) loadFont(path string, size int) (*ttf.Font, error) {
	f, err := ttf.OpenFont(path, size)
	if err != nil {
		return nil, err
	}
	fm.fonts[fontIdentifier{path, size}] = f
	return f, nil
}

func (fm *FontManager) GetFont(path string, size int) (*ttf.Font, error) {
	// try looking for font if already loaded
	f, ok := fm.fonts[fontIdentifier{path, size}]
	if ok {
		return f, nil
	}

	// else load and return
	return fm.loadFont(path, size)
}

func (fm *FontManager) Deinit() {
	for _, f := range fm.fonts {
		f.Close()
	}
}
