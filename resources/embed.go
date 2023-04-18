package resources

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"

	"github.com/tinne26/edau"
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

	//go:embed sfx/riser10s.ogg
	Riser10s_ogg []byte

	//go:embed sfx/riser5s.ogg
	Riser5s_ogg []byte

	//go:embed sfx/riser2_5s.ogg
	Riser2_5s_ogg []byte

	//go:embed sfx/riser1_25s.ogg
	Riser1_25s_ogg []byte

	//go:embed sfx/riser0_80s.ogg
	Riser0_80s_ogg []byte

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

type audioStream interface {
	io.ReadSeeker
	Length() int64
}

var ACtx *audio.Context

var IntroPlayer, RiserPlayer, VictoryPlayer *audio.Player

var risers map[string]audioStream

func init() {
	ACtx = audio.NewContext(44100)

	foo, err := mp3.DecodeWithSampleRate(44100, bytes.NewReader(Intro_mp3))
	if err != nil {
		log.Fatal(err)
	}
	IntroPlayer, _ = audio.NewPlayer(ACtx, foo)
	//VictoryPlayer = audio.NewPlayerFromBytes(ACtx, Victory_mp3)

	risers = make(map[string]audioStream)
	risers["10s"], err = vorbis.DecodeWithoutResampling(bytes.NewReader(Riser10s_ogg))
	risers["5s"], err = vorbis.DecodeWithoutResampling(bytes.NewReader(Riser5s_ogg))
	risers["2_5s"], err = vorbis.DecodeWithoutResampling(bytes.NewReader(Riser2_5s_ogg))
	risers["1_25s"], err = vorbis.DecodeWithoutResampling(bytes.NewReader(Riser1_25s_ogg))
	risers["0_80s"], err = vorbis.DecodeWithoutResampling(bytes.NewReader(Riser0_80s_ogg))
}

func PlayIntro() {
	IntroPlayer.Rewind()
	IntroPlayer.Play()
}

// picks which riser to play and maybe speeds it up or down via edau?
func PlayRiser(ticksLeft int) {
	var err error
	var shifter *edau.SpeedShifter
	var s audioStream
	var speed float64 = 1
	switch {
	case ticksLeft > 600:
		return
	case ticksLeft > 450:
		s = risers["10s"]
	case ticksLeft > 150:
		s = risers["5s"]
	case ticksLeft > 100:
		s = risers["2_5s"]
	case ticksLeft > 50:
		s = risers["1_25s"]
	default:
		s = risers["0_80s"]
	}

	speed = speed + (rand.NormFloat64() / 20)

	shifter = edau.NewDefaultSpeedShifter(s)
	shifter.SetSpeed(speed)
	fmt.Println("setting speed to: ", speed)
	RiserPlayer, err = ACtx.NewPlayer(shifter)
	if err != nil {
		log.Fatal(err)
	}
	RiserPlayer.Rewind()
	RiserPlayer.Play()
}

func StopRiser() {
	if RiserPlayer != nil {
		RiserPlayer.Close()
	}
}

var fxs = map[string][]wav{
	"hit":     {Hit1_wav, Hit2_wav, Hit3_wav, Hit4_wav, Hit5_wav, Hit6_wav},
	"expl":    {Expl0_wav, Expl1_wav, Expl2_wav, Expl3_wav, Expl4_wav},
	"takeoff": {Takeoff_wav},
	"menu":    {Menu_wav},
	"bonk":    {Bonk_wav},
}

func PlayFX(name string) {
	choices, ok := fxs[name]
	if !ok {
		log.Fatal(name, "Not found")
	}
	x := choices[rand.Intn(len(choices))]
	sePlayer := audio.NewPlayerFromBytes(ACtx, x)
	sePlayer.Play()
}
