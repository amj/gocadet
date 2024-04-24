package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	resources "github.com/amj/gocadet/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var helpImage *ebiten.Image
var moshipHelpImage *ebiten.Image

type HelpScene struct {
	sf         *Starfield
	decelTicks int
}

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.Instr_png))
	if err != nil {
		log.Fatal(err)
	}
	helpImage = ebiten.NewImageFromImage(img)
	moimg, _, err := image.Decode(bytes.NewReader(resources.MoshipInstr_png))
	if err != nil {
		log.Fatal(err)
	}
	moshipHelpImage = ebiten.NewImageFromImage(moimg)
}

func (s *HelpScene) OnEnter(sm *SceneManager) error {
	fmt.Println("Entered helpscene")
	s.sf = sm.Ctx.sf // stars!
	s.decelTicks = 15
	return nil
}

func (s *HelpScene) OnExit(sm *SceneManager) error {
	fmt.Println("Exiting helpscene")
	s.sf = nil
	return nil
}

func (s *HelpScene) Update(sm *SceneManager) error {
	if s.decelTicks > 0 {
		s.decelTicks--
		s.sf.speed /= 1.1
	} else {
		s.sf.moveT = stopped
	}
	s.sf.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		sm.SwitchTo("menu")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		sm.SwitchTo("menu")
	}
	return nil
}

func (s *HelpScene) Draw(screen *ebiten.Image) {
	s.sf.Draw(screen)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.25, 0.25)
	op.GeoM.Translate(0, 0)
	screen.DrawImage(helpImage, op)
	op.GeoM.Translate(screenWidth*0.5, screenHeight*0.5)
	screen.DrawImage(moshipHelpImage, op)

	x := int(math.Round(screenWidth * 0.52))
	y := 2 * LFHeight(smallArcadeFont)
	text.Draw(screen, "Zap letters by pressing", smallArcadeFont, x, y, color.White)
	y += LFHeight(smallArcadeFont)
	text.Draw(screen, "keys as letters appear!", smallArcadeFont, x, y, color.White)
	y += 2 * LFHeight(smallArcadeFont)
	text.Draw(screen, "Accuracy counts! Get a", smallArcadeFont, x, y, color.White)
	y += LFHeight(smallArcadeFont)
	text.Draw(screen, "high enough score to ", smallArcadeFont, x, y, color.White)
	y += LFHeight(smallArcadeFont)
	text.Draw(screen, "unlock the next mission!", smallArcadeFont, x, y, color.White)

	x = int(math.Round(screenWidth * 0.02))
	y = 2*LFHeight(smallArcadeFont) + int(math.Round(screenHeight*0.5))
	text.Draw(screen, "Keep your hands on the", smallArcadeFont, x, y, color.White)
	y += LFHeight(smallArcadeFont)
	text.Draw(screen, "home row!  When the ", smallArcadeFont, x, y, color.White)
	y += LFHeight(smallArcadeFont)
	text.Draw(screen, "mothership appears, hit", smallArcadeFont, x, y, color.White)
	y += LFHeight(smallArcadeFont)
	text.Draw(screen, "all the home keys at once!", smallArcadeFont, x, y, color.White)
	y += 2 * LFHeight(smallArcadeFont)
	text.Draw(screen, "YOU CAN DO IT!", smallArcadeFont, x, y, color.White)
}
