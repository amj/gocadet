package resources

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

var (
	//go:embed keyboard.png
	Keyboard_png []byte

	//go:embed hands-home.png
	HandsHome_png []byte

	//go:embed sfx/hit1.wav
	Hit1_wav []byte

	//go:embed sfx/hit2.wav
	Hit2_wav []byte

	//go:embed sfx/hit3.wav
	Hit3_wav []byte

	//go:embed sfx/hit4.wav
	Hit4_wav []byte

	//go:embed sfx/hit5.wav
	Hit5_wav []byte

	//go:embed sfx/hit6.wav
	Hit6_wav []byte

	//go:embed sfx/expl0.wav
	Expl0_wav []byte

	//go:embed sfx/expl1.wav
	Expl1_wav []byte

	//go:embed sfx/expl2.wav
	Expl2_wav []byte

	//go:embed sfx/expl3.wav
	Expl3_wav []byte

	//go:embed sfx/expl4.wav
	Expl4_wav []byte

	//go:embed sfx/bonk.wav
	Bonk_wav []byte

	//go:embed sfx/takeoff.wav
	Takeoff_wav []byte

	//go:embed sfx/menu.wav
	Menu_wav []byte

	//go:embed sfx/intro.mp3
	Intro_mp3 []byte

	//go:embed sfx/victory.mp3
	Victory_mp3 []byte

	//go:embed mplus-1p-regular.ttf
	Mplus1p_ttf []byte

	//go:embed pressstart2p.ttf
	PressStart_ttf []byte
)

type wav []byte

var ACtx *audio.Context

var IntroPlayer, VictoryPlayer *audio.Player

func init() {
	ACtx = audio.NewContext(44100)

	foo, err := mp3.DecodeWithSampleRate(44100, bytes.NewReader(Intro_mp3))
	if err != nil {
		log.Fatal(err)
	}
	IntroPlayer, _ = audio.NewPlayer(ACtx, foo)
	//VictoryPlayer = audio.NewPlayerFromBytes(ACtx, Victory_mp3)
}

func PlayIntro() {
	IntroPlayer.Rewind()
	IntroPlayer.Play()
}

var fxs = map[string][]wav{
	"hit":     {Hit1_wav, Hit2_wav, Hit3_wav, Hit4_wav, Hit5_wav, Hit6_wav},
	"expl":    {Expl0_wav, Expl1_wav, Expl2_wav, Expl3_wav, Expl4_wav},
	"takeoff": {Takeoff_wav},
	"menu":    {Menu_wav},
	"bonk":    {Bonk_wav},
}

func PlayFX(name string) {
	fmt.Println("Playing", name)
	choices, ok := fxs[name]
	if !ok {
		log.Fatal(name, "Not found")
	}
	x := choices[rand.Intn(len(choices))]
	sePlayer := audio.NewPlayerFromBytes(ACtx, x)
	sePlayer.Play()
}
