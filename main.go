package main

import (
	"errors"
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Context struct {
	Opts    GameOptions
	Profile UserProfile
	MCfg    MissionConfiguration
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
	Name       string
	Results    map[int]GameResult // per level
	bigramErrs []string           // most recent N mistakes
}

type MissionConfiguration struct {
	level        int
	msPerChar    int
	waves        int
	lives        int
	practiceMode bool
}

type GameResult struct {
	time     int
	Score    int
	Errors   int
	Accuracy float32
	Won      bool
}

var errQuit = errors.New("Quit")

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Keyboard Cadet")
	var m *SceneManager = MakeManager(map[string]Scene{
		"menu": &MenuScene{},
		"game": &KeyScene{},
	})
	if p, ok := ActiveProfile(); ok {
		m.Ctx.Profile = p
	} else {
		fmt.Println("No pilots found")
	}

	m.SwitchTo("menu")
	m.current.OnEnter(m)
	defer SavePilots()

	if err := ebiten.RunGame(m); err != nil && err != errQuit {
		log.Fatal(err)
	}

}
