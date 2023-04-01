package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"unicode"

	keyboard "github.com/amj/gocadet/keyboard"
	resources "github.com/amj/gocadet/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var keyboardImage *ebiten.Image

const (
	kbdOffsetX = 24
	kbdOffsetY = 125
)

type gameState uint8

const (
	launch     gameState = iota // anim at start
	flight                      // anim between targets
	targetUp                    // target present, timer ticking
	targetGot                   // success
	targetMiss                  // missed
	moship                      // mothership event
	crashed                     // out of chances
	success                     // mission completed!
)

func (me gameState) String() string {
	return [...]string{"launch", "flight", "tUp", "tGot", "tMiss", "mship", "crashed", "success"}[me]
}

type KeyScene struct {
	keys         []rune
	state        gameState
	nextState    gameState
	target       string // Current word to spell
	ticksInState int
	waveNum      int // number of targets seen
	targetIdx    int // Index of cursor
	ticksLeft    int // Updates until mandated state change
	livesLeft    int
	score        int
	cfg          MissionConfiguration
	sf           *Starfield
}

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.Keyboard_png))
	if err != nil {
		log.Fatal(err)
	}

	keyboardImage = ebiten.NewImageFromImage(img)
}

func (g *KeyScene) OnEnter(sm *SceneManager) error {
	fmt.Println("Entered keyscene")
	g.cfg = sm.Ctx.MCfg
	g.livesLeft = g.cfg.lives
	g.waveNum = 0
	g.state = launch
	g.nextState = launch
	g.ticksInState = 0
	g.ticksLeft = 60
	g.sf = sm.Ctx.sf
	return nil
}

func (g *KeyScene) OnExit(sm *SceneManager) error {
	g.sf.moveT = zoomin
	g.sf.speed = 1.0 / 64.0
	g.sf = nil
	return nil
}

func (g *KeyScene) Update(sm *SceneManager) error {
	g.sf.Update()
	if g.nextState != g.state {
		g.state = g.nextState
		g.ticksInState = 0
	}
	g.keys = ebiten.AppendInputChars(g.keys[:0])
	//g.keys = inpututil.AppendJustPressedKeys(g.keys[:0])
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		sm.SwitchTo("menu")
	}

	// update game logic.
	switch g.state {
	case launch:
		if g.ticksInState == 1 {
			resources.PlayFX("takeoff")
		}
		if g.ticksLeft >= 40 {
			g.sf.speed *= 1.15
		} else {
			g.sf.speed /= 1.1
		}
		if g.ticksLeft == 0 {
			g.nextState = flight
			g.sf.moveT = panright
			g.ticksLeft = 30
			g.sf.speed = 1 / 64.0
		}
	case flight:
		if g.ticksLeft == 0 {
			g.nextState = targetUp
			g.sf.moveT = stopped
		}
	case targetUp:
		if g.target == "" {
			g.SetTargetWord()
		}
		if g.ticksLeft == 0 {
			g.nextState = targetMiss
		}
		for _, k := range g.keys {
			if k == rune(g.target[g.targetIdx]) {
				g.targetIdx++
				if g.targetIdx == len(g.target) {
					g.nextState = targetGot
					g.score += g.ticksLeft
				} else {
					resources.PlayFX("hit")
				}
			} else {
				fmt.Println("k   : ", k)
				fmt.Println("tgt : ", rune(g.target[g.targetIdx]))
			}
		}
	case targetGot:
		if g.ticksInState == 0 {
			resources.PlayFX("expl")
		}
		if g.ticksInState == 5 {
			g.nextState = targetUp
			g.SetTargetWord()
		}
	case targetMiss:
		//play a sound?
		if g.ticksInState == 5 {
			g.livesLeft--
			g.nextState = targetUp
			g.SetTargetWord()
		}
	}

	// End-scene checks.
	if g.livesLeft == 0 {
		g.nextState = crashed
	}
	if g.waveNum >= g.cfg.waves {
		g.nextState = success
	}
	// Timer cleanup
	if g.ticksLeft > 0 {
		g.ticksLeft--
	}
	g.ticksInState++
	return nil
}

func (g *KeyScene) SetTargetWord() {
	g.waveNum++
	var newW string
	for newW = GetTarget(g.cfg.level, g.waveNum); newW == g.target; {
	} // do it until we get a new word
	g.target = newW
	g.targetIdx = 0
	g.ticksLeft = len(g.target) * 120
}

func (g *KeyScene) Draw(screen *ebiten.Image) {
	g.sf.Draw(screen)
	g.DrawKeyboard(screen)

	switch g.state {
	case targetUp:
		drawTargetWord(screen, g.target, g.targetIdx, float64(XforCentering(g.target, hugeArcadeFont)), 200)
		g.DrawHighlightKeys(screen)
	}
	g.drawStatusText(screen)
	text.Draw(screen, g.state.String(), smallArcadeFont, 30, 20, color.White)
}

func (g *KeyScene) drawStatusText(screen *ebiten.Image) {
	var tLeft float32 = 0.0
	if g.state == targetUp {
		tLeft = float32(g.ticksLeft) / 60.0
	}

	text.Draw(screen, fmt.Sprintf("%.3f", tLeft), smallArcadeFont, 30, 40, color.White)
	text.Draw(screen, fmt.Sprintf("%03d/%03d", g.waveNum, g.cfg.waves), smallArcadeFont, 530, 40, color.White)
}

func (g *KeyScene) DrawHighlightKeys(screen *ebiten.Image) {
	if g.targetIdx >= len(g.target) {
		return
	}
	if (g.ticksInState % 30) > 15 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	tgt := unicode.ToUpper(rune(g.target[g.targetIdx]))
	r, ok := keyboard.RuneRect(tgt)
	if !ok {
		return
	}
	op.ColorScale.Scale(0.9, 0.5, 0.0, 1)
	op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
	op.GeoM.Translate(kbdOffsetX, kbdOffsetY)
	op.GeoM.Scale(2, 2)
	screen.DrawImage(keyboardImage.SubImage(r).(*ebiten.Image), op)
}

func (g *KeyScene) DrawKeyboard(screen *ebiten.Image) {
	// Draw the base (grayed) keyboard image.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(kbdOffsetX, kbdOffsetY)
	op.GeoM.Scale(2, 2)
	op.ColorScale.Scale(0.8, 0.8, 0.8, 1)
	screen.DrawImage(keyboardImage, op)
}
