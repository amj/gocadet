package main

import (
	"image/color"
	"log"

	resources "github.com/amj/gocadet/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	gofont "golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/opentype"
)

const (
	hugeFontSize  = fontSize * 4
	titleFontSize = fontSize * 1.5
	fontSize      = 24
	smallFontSize = fontSize / 2
)

var (
	hugeArcadeFont  font.Face
	titleArcadeFont font.Face
	arcadeFont      font.Face
	smallArcadeFont font.Face
)

func init() {
	tt, err := opentype.Parse(resources.PressStart_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72

	titleArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    titleFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	smallArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    smallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	tt2, err := opentype.Parse(gofont.TTF)
	if err != nil {
		log.Fatal(err)
	}

	hugeArcadeFont, err = opentype.NewFace(tt2, &opentype.FaceOptions{
		Size:    hugeFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

}

func drawCenteredText(screen *ebiten.Image, txt string, face font.Face, yline int, clr color.Color) {
	var fsize int = face.Metrics().Height.Ceil()
	fsize += int(float32(fsize) * 1.1)
	x := XforCentering(txt, face)
	text.Draw(screen, txt, face, x, fsize*yline, clr)
}

func LFHeight(face font.Face) int {
	var fsize int = face.Metrics().Height.Ceil()
	fsize += int(float32(fsize) * 1.1)
	return fsize
}

func drawText(screen *ebiten.Image, txt string, face font.Face, x, y int, clr color.Color) {

	text.Draw(screen, txt, face, x, y, clr)
}

func XforCentering(txt string, face font.Face) int {
	fsize, _ := face.GlyphAdvance(rune('m'))
	return (screenWidth - len(txt)*fsize.Ceil()) / 2
}

// idx == highlighted character.
func drawTargetWord(screen *ebiten.Image, txt string, idx int, x, y, scale float64) {
	var glyphs []text.Glyph
	glyphs = text.AppendGlyphs(glyphs, hugeArcadeFont, txt) // todo: maybe don't reraster these
	op := &ebiten.DrawImageOptions{}
	// In this example, multiple colors are used to render glyphs.
	for i, gl := range glyphs {
		op.GeoM.Reset()
		op.GeoM.Translate(x, y)
		op.GeoM.Translate(gl.X, gl.Y)
		op.GeoM.Scale(scale, scale)
		op.ColorScale.Reset()
		gb := float32(1)
		if i == idx {
			gb = 0.1
		}
		op.ColorScale.Scale(1, gb, gb, 1)
		screen.DrawImage(gl.Image, op)
	}

}

func lerpf(x1, x2, t float64) float64 {
	return (x1 + (x2-x1)*t)
}
func lerpi(x1, x2 int, t float64) float64 {
	return (float64(x1) + float64(x2-x1)*t)
}

func interpolateColors(c1, c2 color.RGBA64, t float64) color.RGBA64 {
	r := uint16(float64(c1.R) + (float64(c2.R)-float64(c1.R))*t)
	g := uint16(float64(c1.G) + (float64(c2.G)-float64(c1.G))*t)
	b := uint16(float64(c1.B) + (float64(c2.B)-float64(c1.B))*t)
	a := uint16(float64(c1.A) + (float64(c2.A)-float64(c1.A))*t)
	return color.RGBA64{r, g, b, a}
}
