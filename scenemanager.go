package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

// Scene is an interface with Draw, Update, On{Enter,Exit}.
type Scene interface {
	OnEnter(sm *SceneManager) error
	OnExit(sm *SceneManager) error
	Update(sm *SceneManager) error
	Draw(screen *ebiten.Image)
}

// SceneManager is a struct that holds the current scene and next scene
// a transitionFlag guards against multiple transitions and ensures
// transitions are "queued up".  (transitioning scenes takes at least one tick)
// Ctx is whatever information needs to go between scenes, and could be quite large.
type SceneManager struct {
	current        Scene
	next           Scene
	transitionFlag bool
	Scenes         map[string]Scene
	Ctx            Context
}

// Pass through "update" to the scene, unless we're mid-transition.
func (s *SceneManager) Update() error {
	if !s.transitionFlag {
		return s.current.Update(s)
	}

	s.transitionFlag = false
	if err := s.current.OnExit(s); err != nil {
		return err
	}
	if err := s.next.OnEnter(s); err != nil {
		return err
	}
	s.current = s.next
	s.next = nil
	return nil
}

// Pass through to current scene.
func (s *SceneManager) Draw(r *ebiten.Image) {
	s.current.Draw(r)
}

// Change current in SceneManager to the scene of the argument.
func (s *SceneManager) SwitchTo(name string) {
	if s.current == nil {
		s.current = s.Scenes[name]
	} else {
		s.next = s.Scenes[name]
		s.transitionFlag = true
	}
}

func MakeManager(scenes map[string]Scene) *SceneManager {
	m := &SceneManager{
		Scenes: scenes,
	}
	m.Ctx.sf = NewStarfield(64)
	m.Ctx.MCfg = MissionConfiguration{
		level: 0,
		waves: 20,
		lives: 5,
	}
	return m
}

func (s *SceneManager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
