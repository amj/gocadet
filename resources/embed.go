package resources

import _ "embed"
import "log"
import "math/rand"
import "github.com/hajimehoshi/ebiten/v2/audio"

var (
	//go:embed keyboard.png
	Keyboard_png []byte

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

	//go:embed sfx/takeoff.wav
	Takeoff_wav []byte

	//go:embed sfx/menu.wav
	Menu_wav []byte
)

type wav []byte

var ACtx *audio.Context

func init() {
	ACtx = audio.NewContext(44100)
}

var fxs = map[string][]wav{
	"hit": {Hit1_wav, Hit2_wav, Hit3_wav, Hit4_wav, Hit5_wav, Hit6_wav},
	"expl": {Expl0_wav, Expl1_wav, Expl2_wav, Expl3_wav, Expl4_wav},
	"takeoff": {Takeoff_wav},
	"menu": {Menu_wav},
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
