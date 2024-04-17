package main

import (
	"bytes"
	"fmt"
	"image"
	"log"

	resources "github.com/amj/gocadet/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var helpImage *ebiten.Image
var moshipHelpImage *ebiten.Image

type HelpScene struct {
	sf *Starfield
}

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.Instr_png))
	if err != nil {
		log.Fatal(err)
	}
	helpImage = ebiten.NewImageFromImage(img)
}

func (s *HelpScene) OnEnter(sm *SceneManager) error {
	fmt.Println("Entered helpscene")
	s.sf = sm.Ctx.sf // stars!
	s.sf.moveT = stopped
	return nil
}

func (s *HelpScene) OnExit(sm *SceneManager) error {
	fmt.Println("Exiting helpscene")
	s.sf = nil
	return nil
}

func (s *HelpScene) OnUpdate(sm *SceneManager) error {

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		sm.SwitchTo("menu")
	}
	return nil
}

func (s *HelpScene) Draw(screen *ebiten.Image) {
	s.sf.Draw(screen)
}
