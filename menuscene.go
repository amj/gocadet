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
	profile     *UserProfile
	mCfg        *MissionConfiguration
}

const breathTime int = 60
const breathAmt float64 = 1 / float64(breathTime/2)

func (s *MenuScene) OnEnter(sm *SceneManager) error {
	resources.PlayIntro()
	if s.sf == nil { // pick up an alias to the starfield
		s.sf = sm.Ctx.sf
	}
	// TODO prompt for profile name or switch scenes?
	if sm.Ctx.Profile.Name != "" {
		s.profile = &sm.Ctx.Profile
	} else {
		s.profile = nil
	}

	s.mCfg = &sm.Ctx.MCfg
	s.tick = 0
	s.breath = 0
	return nil
}

func (s *MenuScene) OnExit(sm *SceneManager) error {
	s.sf = nil
	resources.IntroPlayer.Pause()
	return nil
}

func (s *MenuScene) Update(sm *SceneManager) error {
	sm.Ctx.sf.Update()
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch s.optSelected {
		case 0: // Do profile switch?
		case 1: // speed
			resources.PlayFX("menu") // TODO -- new sound?
			sm.Ctx.Profile.Speed++
			if sm.Ctx.Profile.Speed > master {
				sm.Ctx.Profile.Speed = beginner
			}
		case 2:
			resources.PlayFX("menu")
			sm.Ctx.MCfg.level++
			if sm.Ctx.MCfg.level > sm.Ctx.Profile.findCurrentLevel() {
				sm.Ctx.MCfg.level = 0
			}
		case 3: // Launch!
			if sm.Ctx.Profile.Name == "" {
				sm.Ctx.Profile = UserProfile{Name: "dad", Results: make(map[int]GameResult)}
			}
			pData.LastUsed = sm.Ctx.Profile.Name
			sm.SwitchTo("game")
		}
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
	s.breath = UpdateBreath(s.tick, s.breath)
	return nil
}

func UpdateBreath(tick int, breath float64) float64 {
	if (tick % breathTime) < (breathTime / 2) {
		breath += breathAmt
	} else {
		breath -= breathAmt
	}
	return clamp(breath)
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

var optText = [...]string{
	"Pilot",
	"Speed",
	"Mission",
	"LAUNCH",
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	s.sf.Draw(screen) // stars.
	drawCenteredText(screen, "Keyboard Cadet", titleArcadeFont, 1, color.White)
	var txt string
	var c color.RGBA64
	c1 := color.RGBA64{0xffff, 0xffff, 0xffff, 0xffff}
	c2 := color.RGBA64{0x0000, 0xeeee, 0x3333, 0xffff}
	for i := range optText {
		if i == s.optSelected {
			txt = fmt.Sprintf(">%s", optText[i])
			c = interpolateColors(c1, c2, s.breath)
		} else {
			txt = optText[i]
			c = c1
		}
		if s.profile != nil {
			switch i {
			case 0:
				txt = fmt.Sprintf("%s: %s", txt, s.profile.Name)
			case 1:
				txt = fmt.Sprintf("%s: %s", txt, s.profile.Speed)
			case 2:
				txt = fmt.Sprintf("%s: %d", txt, s.mCfg.level)
			}
		}

		if i == len(optText)-1 {
			drawCenteredText(screen, txt, titleArcadeFont, 5, c)
		} else {
			drawCenteredText(screen, txt, arcadeFont, 3+i, c)
		}
	}
}
