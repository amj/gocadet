package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"strings"
	"unicode"

	keyboard "github.com/amj/gocadet/keyboard"
	resources "github.com/amj/gocadet/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var keyboardImage *ebiten.Image
var handsHomeImage *ebiten.Image
var shieldImage *ebiten.Image
var moshipImage *ebiten.Image

const (
	kbdOffsetX = 48
	kbdOffsetY = 280
)

type gameState uint8

const (
	launch     gameState = iota // anim at start
	ready                       // blink at beginning
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
	target       string  // Current word to spell
	tgtX         float64 // current position of target letter, for animating letters.
	tgtY         float64
	ticksInState int
	waveNum      int // number of targets seen
	targetIdx    int // Index of cursor
	ticksLeft    int // Updates until mandated state change
	livesLeft    int
	score        int
	miss         int
	fired        int
	moshipKeys   string
	speed        difficulty
	cfg          MissionConfiguration
	sf           *Starfield
}

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.Keyboard_png))
	if err != nil {
		log.Fatal(err)
	}
	keyboardImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(resources.HandsHome_png))
	if err != nil {
		log.Fatal(err)
	}
	handsHomeImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(resources.Moship_png))
	if err != nil {
		log.Fatal(err)
	}
	moshipImage = ebiten.NewImageFromImage(img)
}

func (g *KeyScene) OnEnter(sm *SceneManager) error {
	fmt.Println("Entered keyscene")
	g.cfg = sm.Ctx.MCfg
	g.speed = sm.Ctx.Profile.Difficulty
	g.livesLeft = g.cfg.lives
	g.waveNum = 0
	g.state = launch
	g.nextState = launch
	g.ticksInState = 0
	g.ticksLeft = 60
	g.score = 0
	g.miss = 0
	g.fired = 0
	g.sf = sm.Ctx.sf // stars!
	g.target = ""
	g.targetIdx = 0
	g.moshipKeys = ""
	return nil
}

func (g *KeyScene) OnExit(sm *SceneManager) error {
	AddProfile(sm.Ctx.Profile)
	g.sf.moveT = zoomin
	g.sf.speed = 1.0 / 64.0
	g.sf = nil
	resources.StopRiser()
	resources.StopSfxPlayer()
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
			g.nextState = ready
			g.sf.moveT = panright
			g.ticksLeft = 30
			g.sf.speed = 1 / 64.0
		}
	case ready:
		if g.ticksInState >= 60 {
			g.nextState = flight
		}
	case flight:
		if g.ticksLeft == 0 {
			g.nextState = targetUp
			g.sf.moveT = stopped
		}
	case targetUp: // target up on screen!  fire!
		if g.target == "" {
			g.SetTargetWord()
		}
		if g.ticksInState == 0 {
			resources.PlayRiser(g.ticksLeft)
		}
		if g.ticksLeft == 0 {
			resources.StopRiser()
			g.nextState = targetMiss
		}
		if len(g.target) < 3 {
			g.tgtX = lerpf(g.tgtX, float64(XforCentering(g.target, hugeArcadeFont)), 1/float64(g.ticksLeft))
			g.tgtY = lerpf(g.tgtY, 200, 1/float64(g.ticksLeft))
		}
		for _, k := range g.keys {
			g.fired++
			if k == rune(g.target[g.targetIdx]) {
				g.targetIdx++
				if g.targetIdx == len(g.target) {
					g.score += WordScore(g.ticksInState, len(g.target))
					g.nextState = targetGot // explosion sfx at 'got' state
					resources.StopRiser()
					break
				} else {
					resources.PlayFX("hit")
				}
			} else {
				g.miss++
			}
		}
	case moship:
		if g.ticksInState == 0 {
			// play the alert sound effect
			resources.PlayFX("moship")
			// reset list of keys to home row.
			g.moshipKeys = "asdfjkl;"
		}
		for _, k := range g.keys {
			// test if 'k' is in the moshipKeys string and remove it if so
			g.moshipKeys = strings.Replace(g.moshipKeys, string(k), "", -1)
			// if moshipKeys is empty, add 100 to the score
			if g.moshipKeys == "" {
				g.score += 100
				//g.score += WordScore(g.ticksInState, len(g.target))
				resources.StopSfxPlayer()
				g.nextState = targetGot // explosion sfx at 'got' state
				g.waveNum++
			}
		}
		if g.ticksInState == 120*(5-int(g.speed)) {
			resources.StopSfxPlayer()
			g.nextState = targetMiss
		}

	case targetGot:
		if g.ticksInState == 0 {
			resources.PlayFX("expl")
		}
		if g.ticksInState == 5 {
			if g.waveNum%5 == 4 {
				g.nextState = moship
			} else {
				g.nextState = targetUp
				g.SetTargetWord()
			}
		}
	case targetMiss:
		if g.ticksInState == 0 {
			resources.PlayFX("bonk")
		}
		if g.ticksInState == 5 {
			g.livesLeft--
			g.nextState = targetUp
			g.SetTargetWord()
		}
	case success:
		if g.ticksInState == 0 {
			resources.PlayFX("takeoff") // TODO: victory fanfare
			g.sf.moveT = zoomin         // speed up and fly off!
			g.sf.speed = 1.0 / 300.0
			// save the game result.  TODO: check for better score
			sm.Ctx.Profile.Results[g.cfg.level] = g.MakeResult(true)
		}
		//acceleration effect
		g.sf.speed = 1.0 / (300.0 - min(285.0, float32(g.ticksInState)))
	case crashed:
		if g.ticksInState == 40 {
			resources.PlayFX("crashed")
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

func (g *KeyScene) MakeResult(won bool) MissionResult {
	return MissionResult{
		Score:    g.score,
		Accuracy: float32(g.fired-g.miss) / float32(g.fired),
		Errors:   g.miss,
		Won:      won,
	}
}

func min(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

func (g *KeyScene) SetTargetWord() {
	g.waveNum++
	var newW string
	// loopcall GetTarget until we get a new word
	for newW = GetTarget(g.cfg.level, g.waveNum); newW == g.target; newW = GetTarget(g.cfg.level, g.waveNum) {
		fmt.Println("duplicate word:", newW, " == ", g.target)
	}
	g.target = newW
	g.targetIdx = 0
	g.ticksLeft = TicksForTarget(len(g.target), g.speed)
	if len(g.target) < 3 { // short words zoom towards the player, long words are still
		g.tgtX = rand.NormFloat64()*float64(Speeds[g.speed]/4) + float64(420)
		g.tgtY = rand.NormFloat64()*float64(Speeds[g.speed]/8) + float64(150)
		fmt.Printf("X: %0.3f, Y: %0.3f", g.tgtX, g.tgtY)
	}
}

func (g *KeyScene) Draw(screen *ebiten.Image) {
	g.sf.Draw(screen)
	if g.state != success && g.state != crashed {
		g.DrawKeyboard(screen)
		g.DrawShields(screen)
	}

	switch g.state {
	case ready:
		if (g.ticksInState/10)%2 == 0 {
			drawCenteredText(screen, "GET READY", titleArcadeFont, 2, color.RGBA{0x22, 0xff, 0x22, 0xff})
		}
	case success:
		drawCenteredText(screen, "MISSION COMPLETE", titleArcadeFont, 2, color.RGBA{0x22, 0xff, 0x22, 0xff})
		drawCenteredText(screen, fmt.Sprintf("Score: %d00", g.score), arcadeFont, 4, color.White)
		drawCenteredText(screen, fmt.Sprintf("Accuracy: %d%%", (100*(g.fired-g.miss))/g.fired), arcadeFont, 5, color.White)
	case crashed:
		drawCenteredText(screen, "CRASHED!", titleArcadeFont, 2, color.RGBA{0xff, 0x22, 0x22, 0xff})
		drawCenteredText(screen, fmt.Sprintf("Score: %d00", g.score), arcadeFont, 4, color.White)
		drawCenteredText(screen, "Don't give up!", arcadeFont, 5, color.White)
	case moship:
		g.DrawMoship(screen)
		drawCenteredText(screen, "asdfjkl;", hugeArcadeFont, 1, color.RGBA{0x22, 0xff, 0x22, 0xff})
		g.DrawHighlightKeys(screen)

	case targetUp:
		if len(g.target) < 3 {
			scale := 0.7 + lerpf(0, 0.3, float64(g.ticksInState)/float64(TicksForTarget(len(g.target), g.speed)))
			drawTargetWord(screen, g.target, g.targetIdx,
				float64(g.tgtX), float64(g.tgtY), scale)
		} else {
			drawTargetWord(screen, g.target, g.targetIdx, float64(XforCentering(g.target, hugeArcadeFont)), 200, 1.0)
		}
		g.DrawHighlightKeys(screen)
		fallthrough

	default:
		g.drawStatusText(screen)
		text.Draw(screen, g.state.String(), smallArcadeFont, 530, 20, color.White)
	}
}

func (g *KeyScene) drawStatusText(screen *ebiten.Image) {
	var tLeft float32 = 0.0
	if g.state == targetUp {
		tLeft = float32(g.ticksLeft) / 60.0
	}

	text.Draw(screen, fmt.Sprintf("%d00", g.score), smallArcadeFont, 30, 20, color.White)
	text.Draw(screen, fmt.Sprintf("%.3f", tLeft), smallArcadeFont, 30, 40, color.White)
	text.Draw(screen, fmt.Sprintf("%03d/%03d", g.waveNum, g.cfg.waves), smallArcadeFont, 530, 40, color.White)
}

func (g *KeyScene) DrawHighlightedKey(screen *ebiten.Image, rn rune) {
	tgt := unicode.ToUpper(rn)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	r, ok := keyboard.RuneRect(tgt)
	if !ok {
		return
	}
	op.ColorScale.Scale(0.9, 0.5, 0.0, 1)
	op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(kbdOffsetX, kbdOffsetY)
	screen.DrawImage(keyboardImage.SubImage(r).(*ebiten.Image), op)
}

func (g *KeyScene) DrawHighlightKeys(screen *ebiten.Image) {
	if (g.ticksInState % 30) > 15 { // blink it
		return
	}
	if g.state == targetUp {
		if g.targetIdx >= len(g.target) {
			return
		}
		g.DrawHighlightedKey(screen, rune(g.target[g.targetIdx]))
	} else if g.state == moship {
		for _, c := range "asdfjkl;" {
			g.DrawHighlightedKey(screen, c)
		}
	}
}

func (g *KeyScene) DrawMoship(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, 1)
	op.GeoM.Translate(kbdOffsetX, 75)
	screen.DrawImage(moshipImage, op)
}

func (g *KeyScene) DrawKeyboard(screen *ebiten.Image) {
	// Draw the base (grayed) keyboard image.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(kbdOffsetX, kbdOffsetY)
	op.ColorScale.Scale(0.8, 0.8, 0.8, 1)
	screen.DrawImage(keyboardImage, op)
	op.ColorScale.Scale(0.8, 0.8, 0.8, 0.01) // where's my alpha? >:(
	screen.DrawImage(handsHomeImage, op)
}

func (g *KeyScene) DrawShields(screen *ebiten.Image) {
	const (
		margin int = 10
		xOff       = 15
		yOff       = 470
		width      = 20
	)
	height := (180 - (margin * (g.cfg.lives - 1))) / g.cfg.lives
	shieldImage = ebiten.NewImage(width, height)
	shieldImage.Fill(color.RGBA{0x00, 0x99, 0xff, 0xff})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(xOff, yOff)
	rop := &ebiten.DrawImageOptions{}
	rop.GeoM.Translate(screenWidth-width-xOff, yOff)
	for i := 0; i < g.livesLeft; i++ {
		op.GeoM.Translate(0, float64(-1*(height+margin)))
		screen.DrawImage(shieldImage, op)
		rop.GeoM.Translate(0, float64(-1*(height+margin)))
		screen.DrawImage(shieldImage, rop)
	}

}
