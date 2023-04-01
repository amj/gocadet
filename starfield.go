package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	width      = 640
	height     = 480
	starsCount = 1024
)

type Star struct {
	fromx, fromy, tox, toy, brightness float32
}

type moveType uint8

const (
	zoomin moveType = iota
	panright
	panleft
	stopped
)

func (me moveType) String() string {
	return [...]string{"zoomin", "panR", "panL", "stop"}[me]
}

type Starfield struct {
	lightspeed bool
	scale      float32
	speed      float32
	moveT      moveType
	s          [starsCount]Star
}

func NewStarfield(scale float32) (sf *Starfield) {
	sf = &Starfield{scale: scale, moveT: zoomin, speed: 1.0 / 64.0}
	for i := 0; i < starsCount; i++ {
		sf.s[i].Init(sf.scale)
	}
	return
}

func (sf *Starfield) Draw(screen *ebiten.Image) {
	for i := 0; i < starsCount; i++ {
		sf.s[i].Draw(screen, sf.scale)
	}
}

func (sf *Starfield) Update() {
	if sf.moveT == stopped {
		for i := range sf.s {
			sf.s[i].tox = sf.s[i].fromx + sf.scale
			sf.s[i].toy = sf.s[i].fromy + sf.scale
		}
		return
	}
	switch sf.moveT {
	case panright:
		for i := range sf.s {
			sf.s[i].Update(sf.s[i].fromx+6400.0, sf.s[i].fromy, sf.scale, sf.speed, sf.moveT)
		}
	case panleft:
		for i := range sf.s {
			sf.s[i].Update(sf.s[i].fromx-6400.0, sf.s[i].fromy, sf.scale, sf.speed, sf.moveT)
		}
	case zoomin:
		x, y := ebiten.CursorPosition()
		for i := range sf.s {
			sf.s[i].Update(float32(x)*sf.scale, float32(y)*sf.scale, sf.scale, sf.speed, sf.moveT)
		}
	}
}

func (s *Star) Init(scale float32) {
	s.tox = rand.Float32() * width * scale
	s.fromx = s.tox
	s.toy = rand.Float32() * height * scale
	s.fromy = s.toy
	s.brightness = rand.Float32() * 0xff
}

// Init on either the left or right edge
func (s *Star) InitLR(x, scale float32) {
	s.tox = x
	s.fromx = s.tox
	s.toy = rand.Float32() * height * scale
	s.fromy = s.toy
	s.brightness = rand.Float32()*0xbb + 0x33
}

func (s *Star) Update(x, y, scale, speed float32, mt moveType) {
	s.fromx = s.tox
	s.fromy = s.toy
	s.tox += (s.tox - x) * speed
	s.toy += (s.toy - y) * speed
	switch mt {
	case zoomin:
		s.brightness += 1
		if 0xff < s.brightness {
			s.brightness = 0xff
		}
		if s.fromx < 0 || width*scale < s.fromx || s.fromy < 0 || height*scale < s.fromy {
			s.Init(scale)
		}
	case panright:
		if s.fromx < 0 {
			s.InitLR(width*scale, scale)
		}
	case panleft:
		if s.fromx >= width*scale {
			s.InitLR(0, scale)
		}
	}
}

func (s *Star) Draw(screen *ebiten.Image, scale float32) {
	c := color.RGBA{
		R: uint8(0xff * s.brightness / 0xff),
		G: uint8(0xee * s.brightness / 0xff),
		B: uint8(0xcc * s.brightness / 0xff),
		A: 0xff}
	vector.StrokeLine(screen, s.fromx/scale, s.fromy/scale, s.tox/scale, s.toy/scale, 1, c, true)
}
