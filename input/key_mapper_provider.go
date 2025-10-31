package input

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Key int

const (
	Up Key = iota
	Down
	Left
	Right
	A
	B
	X
	Y
	Unknown
)

type KeyMapperStrategy interface {
	ApplicableTo(event sdl.Event) bool
	MapInputToKey(event sdl.Event) Key
}

var strategies = []KeyMapperStrategy{
	&KeyboardMapperStrategy{},
	&JoyButtonMapperStrategy{},
	&JoyHatMapperStrategy{},
}

func GetKeyMapperStrategy(event sdl.Event) KeyMapperStrategy {
	for i := 0; i < len(strategies); i++ {
		if strategies[i].ApplicableTo(event) {
			return strategies[i]
		}
	}
	return nil
}
