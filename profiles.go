package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

const profPath = "pilots.gob"

type difficulty uint8

const (
	beginner difficulty = iota
	standard
	advanced
	expert
	master
)

func (me difficulty) String() string {
	return [...]string{"slow", "fast", "super", "Hyper!", "ULTRA"}[me]
}

type UserProfile struct {
	Name       string
	Speed      difficulty         // difficulty setting.
	Results    map[int]GameResult // per level
	bigramErrs []string           // most recent N mistakes
}

type PilotData struct {
	LastUsed string
	Profiles map[string]UserProfile
}

var pData PilotData

func init() {
	// load pilots?
	fd, err := os.Open(profPath)
	defer fd.Close()

	if err != nil {
		fmt.Println("err loading pilots")
		log.Println(err)
		return
	}

	dec := gob.NewDecoder(fd)
	if err := dec.Decode(&pData); err != nil {
		fmt.Println("err decoding pilot data")
		log.Fatal(err)
	}

	fmt.Println("loaded: ", pData)
	fmt.Println("last: ", pData.LastUsed)
}

func ActiveProfile() (p UserProfile, ok bool) {
	if pData.LastUsed == "" {
		fmt.Println("No lastused saved")
		return UserProfile{}, false
	}
	p, ok = pData.Profiles[pData.LastUsed]
	if !ok {
		fmt.Printf("Last used profile not found: %s", pData.LastUsed)
	}
	return p, ok
}

func AddProfile(p UserProfile) error {
	if pData.Profiles == nil {
		pData.Profiles = make(map[string]UserProfile)
	}
	pData.Profiles[p.Name] = p
	pData.LastUsed = p.Name
	return nil
}

func SavePilots() {
	fd, err := os.Create(profPath)
	defer fd.Close()

	if err != nil {
		fmt.Println("err writing pilots")
		log.Println(err)
		return
	}

	enc := gob.NewEncoder(fd)
	if err := enc.Encode(pData); err != nil {
		fmt.Println("err encoding pilot data")
		log.Fatal(err)
	}
}

// Finds the lowest number mission for a profile
// Either the first mission they haven't finished,
// or the first mission where the accuracy was less than 90%
func (u UserProfile) findCurrentLevel() int {
	for i := 0; i < WordsMaxLevel; i++ {
		res, ok := u.Results[i]
		if !ok {
			return i
		}
		if !res.Won || res.Accuracy < .9 {
			return i
		}
	}
	return WordsMaxLevel
}
