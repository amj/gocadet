package main

import (
	"fmt"
	"math/rand"
)

const WordsMaxLevel = 11

var missions = map[int]string{
	0:  "home row",
	1:  "heo",
	2:  "ti L-shift",
	3:  "rw.'",
	4:  "ng R-shift",
	5:  "ucp",
	6:  "yx,",
	7:  "mz=",
	8:  "bq?",
	9:  "v\"",
	10: "0-9",
}

var chars = map[int]string{
	0:  "asdfjkl;",
	1:  "heo",
	2:  "ti",
	3:  "rw.'",
	4:  "ng",
	5:  "ucp",
	6:  "yx,",
	7:  "mz=",
	8:  "bq?",
	9:  "v\"",
	10: "0123456789",
}

var words = map[int][]string{
	0:  {"dad", "sad", "lad", "fad"},
	1:  {"he", "she", "led", "head", "fed", "jade", "load", "joss", "hose", "foes", "lose", "does"},
	2:  {"the", "these", "those", "that", "Oath", "side", "tie", "Jot", "fit", "sit", "flit"},
	3:  {"wait", "saw", "was", "wit", "what", "who", "water", "rook", "Look", "took"},
	4:  {"wing", "nog", "go", "wag", "grow", "ring", "rag", "thing", "find"},
	5:  {"Cup", "Puck", "Cog", "Pin", "Pun", "Fun", "Ping", "rung"},
	6:  {"ply", "six", "ox", "fox", "toy", "yip", "fix", "say"},
	7:  {"maze", "max", "yam", "mop", "zap", "map", "book"},
	8:  {"mob", "job", "fob", "slab", "jab", "quit", "boom", "view"},
	9:  {"vibe", "move", "rave", "over", "cave", "save", "pave"},
	10: {"draw", "wire", "flank", "trunk", "fire", "crew", "man", "monk", "wood"},
}

func clampI(x int, max int) int {
	if x > max {
		return max
	}
	if x < 0 {
		return 0
	}
	return x
}

func GetWord(level int) string {
	easier := rand.Float64() > 0.65
	if easier {
		level = clampI(level-rand.Intn(3), level)
	}
	return words[level][rand.Intn(len(words[level]))]
}

func GetBigram(level int) string {
	first := rand.Intn(2) == 1
	off := clampI(level-rand.Intn(3), level)
	c1 := GetChar(level)
	c2 := GetChar(off)

	if first {
		return fmt.Sprintf("%s%s", string(c1), string(c2))
	} else {
		return fmt.Sprintf("%s%s", string(c2), string(c1))
	}
}

func GetChar(level int) byte {
	return chars[level][rand.Intn(len(chars[level]))]
}

func GetTarget(level int, wave int) string {
	var fewWords bool = (level == 0) || (level == 10)

	switch {
	case fewWords && wave%10 == 0:
		return GetWord(level)
	case fewWords && wave%3 == 0:
		return GetBigram(level)
	case fewWords:
		return string(GetChar(level))
	case wave%5 == 0:
		return GetWord(level)
	case wave%3 == 0:
		return GetBigram(level)
	case level > 5:
		if rand.Float64() > 0.5 {
			return GetBigram(level)
		} else {
			return GetWord(level)
		}
	default:
		x := rand.Float64()
		switch {
		case x > 0.75:
			return GetWord(level)
		case x > 0.5:
			easier := rand.Float64() > 0.65
			if easier {
				level = clampI(level-rand.Intn(3), level)
			}
			return string(GetChar(level))
		default:
			return GetBigram(level)
		}

	}
}

var Speeds = map[difficulty]int{
	beginner: 300,
	standard: 120,
	advanced: 60,
	expert:   45,
	master:   30,
}

func WordScore(ticksInState, wordLen int) int {
	return (Speeds[beginner] * wordLen) - ticksInState
}
