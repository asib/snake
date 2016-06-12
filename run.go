package main

import sys "github.com/asib/snake/system"

const DEBUG = true

func initialize(s *sys.System) (err error) {
	if err = s.Init(); err != nil {
		return
	}

	return
}

func main() {
	sys := sys.Create(DEBUG, 640, 480)
	if err := initialize(sys); err != nil {
		panic(err)
	}

	sys.Run()
}
