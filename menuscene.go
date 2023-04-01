package main

import (
	"fmt"
	"image/color"

	resources "github.com/amj/gocadet/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "github.com/hajimehoshi/ebiten/v2/text"
)

type MenuScene struct {
	sf          *Starfield
	tick        int
	breath      float64
	optSelected int
}

var optText = [...]string{
	"Profile",
	"Speed  ",
	"Mission",
	"LAUNCH",
}

const breathTime int = 120

func (s *MenuScene) OnEnter(sm *SceneManager) error {
	if s.sf == nil { // pick up an alias to the starfield
		s.sf = sm.Ctx.sf
	}
	s.tick = 0
	s.breath = 0
	return nil
}

func (s *MenuScene) OnExit(sm *SceneManager) error {
	s.sf = nil
	return nil
}

func (s *MenuScene) Update(sm *SceneManager) error {
	sm.Ctx.sf.Update()
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		var c = MissionConfiguration{
			level:     9,
			msPerChar: 1000,
			waves:     20,
			lives:     5,
		}
		sm.Ctx.MCfg = c
		sm.SwitchTo("game")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) && s.optSelected < (len(optText)-1) {
		s.optSelected++
		resources.PlayFX("menu")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) && s.optSelected > 0 {
		s.optSelected--
		resources.PlayFX("menu")
	}
	s.tick++
	if (s.tick % breathTime) < (breathTime / 2) {
		s.breath += 0.01
	} else {
		s.breath -= 0.01
	}
	s.breath = clamp(s.breath)
	return nil
}

func clamp(x float64) float64 {
	if x > 1 {
		return 1
	}
	if x < 0 {
		return 0
	}
	return x
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	s.sf.Draw(screen)
	drawCenteredText(screen, "Keyboard Cadet", titleArcadeFont, 1, color.White)
	var txt string
	var c color.RGBA64
	c1 := color.RGBA64{0xffff, 0xffff, 0xffff, 0xffff}
	c2 := color.RGBA64{0xffff, 0x4444, 0x4444, 0xffff}
	for i := range optText {
		if i == s.optSelected {
			txt = fmt.Sprintf(">%s", optText[i])
			c = interpolateColors(c1, c2, s.breath)
		} else {
			txt = optText[i]
			c = c1
		}
		if i == len(optText)-1 {
			drawCenteredText(screen, txt, titleArcadeFont, 5, c)
		} else {
			drawCenteredText(screen, txt, arcadeFont, 3+i, c)
		}
	}
}