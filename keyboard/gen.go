// Copyright 2015 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build ignore

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	arcadeFontSize = 8
)

var (
	arcadeFont font.Face
)

func init() {
	b, err := os.ReadFile(filepath.Join("..", "resources", "pressstart2p.ttf"))
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(b)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    arcadeFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

var keyboardKeys = [][]string{
	{"Esc", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "=", "Bksp"},
	{"Tab", "Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P", "[", "]", "\\"},
	{"Caps", "A", "S", "D", "F", "G", "H", "J", "K", "L", ";", "'", "Enter"},
	{"Shift", "Z", "X", "C", "V", "B", "N", "M", ",", ".", "/", " Shift"},
	{"Ctrl", "Alt", "Space", "Alt", "Ctrl", "Win"},
}

func keyDisplayNameToKey(name string) ebiten.Key {
	switch name {
	case "Esc":
		return ebiten.KeyEscape
	case "Tab":
		return ebiten.KeyTab
	case "Caps":
		return ebiten.KeyCapsLock
	case "Ctrl":
		return ebiten.KeyControl
	case "Shift":
		return ebiten.KeyShiftLeft
	case " Shift":
		return ebiten.KeyShiftRight
	case "Win":
		return ebiten.KeyMeta
	case "Alt":
		return ebiten.KeyAlt
	case "Space":
		return ebiten.KeySpace
	case "Up":
		return ebiten.KeyUp
	case "Left":
		return ebiten.KeyLeft
	case "Down":
		return ebiten.KeyDown
	case "Right":
		return ebiten.KeyRight
	case "Bksp":
		return ebiten.KeyBackspace
	case "Enter":
		return ebiten.KeyEnter
	case "-":
		return ebiten.KeyMinus
	case "=":
		return ebiten.KeyEqual
	case "\\":
		return ebiten.KeyBackslash
	case "`":
		return ebiten.KeyGraveAccent
	case "[":
		return ebiten.KeyLeftBracket
	case "]":
		return ebiten.KeyRightBracket
	case ";":
		return ebiten.KeySemicolon
	case "'":
		return ebiten.KeyApostrophe
	case ",":
		return ebiten.KeyComma
	case ".":
		return ebiten.KeyPeriod
	case "/":
		return ebiten.KeySlash
	}
	if len(name) != 1 {
		panic("not reached: unknown key " + name)
	}
	c := name[0]
	if '0' <= c && c <= '9' {
		return ebiten.Key0 + ebiten.Key(c-'0')
	}
	if 'A' <= c && c <= 'Z' {
		return ebiten.KeyA + ebiten.Key(c-'A')
	}
	panic("not reached: unknown key " + name)
}

func drawKey(t *ebiten.Image, name string, x, y, width int) {
	const height = 16
	width--
	img := ebiten.NewImage(width, height)
	p := make([]byte, width*height*4)
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			x := (i + j*width) * 4
			switch j {
			case 0, height - 1:
				if 3 <= i && i <= width-4 {
					p[x] = 0xff
					p[x+1] = 0xff
					p[x+2] = 0xff
					p[x+3] = 0xff
				}
			case 1, height - 2:
				if i == 2 || i == width-3 {
					p[x] = 0xff
					p[x+1] = 0xff
					p[x+2] = 0xff
					p[x+3] = 0xff
				}
			case 2, height - 3:
				if i == 1 || i == width-2 {
					p[x] = 0xff
					p[x+1] = 0xff
					p[x+2] = 0xff
					p[x+3] = 0xff
				}
			default:
				if i == 0 || i == width-1 {
					p[x] = 0xff
					p[x+1] = 0xff
					p[x+2] = 0xff
					p[x+3] = 0xff
				}
			}
		}
	}
	img.WritePixels(p)
	const offset = 4
	if name != "Space" {
		text.Draw(img, name, arcadeFont, offset, arcadeFontSize+offset+1, color.White)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	t.DrawImage(img, op)
}

func outputKeyboardImage() (map[ebiten.Key]image.Rectangle, map[rune]image.Rectangle, error) {
	keyMap := map[ebiten.Key]image.Rectangle{}
	runeMap := map[rune]image.Rectangle{}
	img := ebiten.NewImage(320, 120)
	x, y := 0, 0
	for j, line := range keyboardKeys {
		x = 0
		const height = 18
		for i, keyDisplayName := range line {
			width := 16
			switch j {
			default:
				switch i {
				case 0:
					width = 16 + 8*(j+2)
				case len(line) - 1:
					if j > 0 {
						width = 16 + 8*(j+2)
					} else {
						width = 32 + 8*(j+2)
					}

				}
			case 4: // last row.
				switch i {
				case 0:
					width = 40
				case 1:
					width = 32
				case 2: // spacebar
					width = (16 * 6)
				case 3:
					width = 32
				case 4:
					width = 40
				case 5:
					width = 32
				}
			case 6, 7:
				width = 16 * 3
			}
			if keyDisplayName != "" {
				drawKey(img, keyDisplayName, x, y, width)
				if keyDisplayName != " " {
					key := keyDisplayNameToKey(keyDisplayName)
					keyMap[key] = image.Rect(x, y, x+width, y+height)
					if len(keyDisplayName) == 1 {
						runeMap[rune(keyDisplayName[0])] = image.Rect(x, y, x+width, y+height)
						fmt.Println("rune: ", rune(keyDisplayName[0]))
					} else {
						fmt.Println("no rune for display name: ", keyDisplayName)
					}
				}
			}
			x += width
		}
		y += height
	}

	out, err := os.Create(filepath.Join("..", "resources", "keyboard.png"))
	if err != nil {
		return nil, nil, err
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		return nil, nil, err
	}

	return keyMap, runeMap, nil
}

const license = `// Copyright 2013 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
`

const keyRectTmpl = `{{.License}}

// Code generated by gen.go using 'go generate'. DO NOT EDIT.

package keyboard

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var keyboardKeyRects = map[ebiten.Key]image.Rectangle{}
var keyboardRuneRects = map[rune]image.Rectangle{}

func init() {
{{range $key, $rect := .KeyRectsMap}}	keyboardKeyRects[ebiten.Key{{$key}}] = image.Rect({{$rect.Min.X}}, {{$rect.Min.Y}}, {{$rect.Max.X}}, {{$rect.Max.Y}})
{{end}}
{{range $rune, $rect := .RuneRectsMap}}	keyboardRuneRects[{{$rune}}] = image.Rect({{$rect.Min.X}}, {{$rect.Min.Y}}, {{$rect.Max.X}}, {{$rect.Max.Y}})
{{end}}
}

func RuneRect(r rune) (image.Rectangle, bool) {
	rect, ok := keyboardRuneRects[r]
	return rect, ok
}

func KeyRect(key ebiten.Key) (image.Rectangle, bool) {
	r, ok := keyboardKeyRects[key]
	return r, ok
}`

func outputKeyRectsGo(k map[ebiten.Key]image.Rectangle, r map[rune]image.Rectangle) error {
	path := "keyrects.go"

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := template.New(path).Parse(keyRectTmpl)
	if err != nil {
		return err
	}

	fmt.Println("Runemap:")
	for key, val := range r {
		fmt.Println(string(key), val.String())
	}
	fmt.Println("========")
	return tmpl.Execute(f, map[string]any{
		"License":      license,
		"KeyRectsMap":  k,
		"RuneRectsMap": r,
	})
}

type game struct {
	rects map[ebiten.Key]image.Rectangle
	runes map[rune]image.Rectangle
}

func (g *game) Update() error {
	var err error
	g.rects, g.runes, err = outputKeyboardImage()
	if err != nil {
		return err
	}
	return ebiten.Termination
}

func (g *game) Draw(_ *ebiten.Image) {
}

func (g *game) Layout(outw, outh int) (int, int) {
	return 256, 256
}

func main() {
	g := &game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
	if err := outputKeyRectsGo(g.rects, g.runes); err != nil {
		log.Fatal(err)
	}
}
