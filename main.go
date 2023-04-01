package main

import (
	"errors"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Context struct {
	Opts    GameOptions
	Profile UserProfile
	MCfg    MissionConfiguration
	Result  GameResults
	sf      *Starfield
}

type GameOptions struct {
	sounds      bool
	showHands   bool
	screenShake bool
	wordsZoom   bool
	kbLayout    bool // eventually something else?
}

type UserProfile struct {
	name       string
	results    map[int]GameResults // per level
	bigramErrs []string            // most recent N mistakes
}

type MissionConfiguration struct {
	level        int
	msPerChar    int
	waves        int
	lives        int
	practiceMode bool
}

type GameResults struct {
	time     int
	score    int
	errors   int
	accuracy float32
}

var errQuit = errors.New("Quit")

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Keyboard Cadet")
	var m *SceneManager = MakeManager(map[string]Scene{
		"menu": &MenuScene{},
		"game": &KeyScene{},
	})
	m.SwitchTo("menu")
	m.current.OnEnter(m)

	if err := ebiten.RunGame(m); err != nil && err != errQuit {
		log.Fatal(err)
	}
}
