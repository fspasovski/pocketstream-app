package input

import (
	"github.com/veandco/go-sdl2/sdl"
)

type KeyboardMapperStrategy struct{}
type JoyButtonMapperStrategy struct{}
type JoyHatMapperStrategy struct{}

func (s *KeyboardMapperStrategy) ApplicableTo(event sdl.Event) bool {
	_, ok := event.(*sdl.KeyboardEvent)
	return ok
}

func (s *KeyboardMapperStrategy) MapInputToKey(event sdl.Event) Key {
	e, _ := event.(*sdl.KeyboardEvent)

	if e.Type == sdl.KEYDOWN {
		switch e.Keysym.Sym {
		case sdl.K_DOWN:
			return Down
		case sdl.K_UP:
			return Up
		case sdl.K_LEFT:
			return Left
		case sdl.K_RIGHT:
			return Right
		case sdl.K_RETURN:
			return A
		case sdl.K_a:
			return A
		case sdl.K_x:
			return X
		case sdl.K_y:
			return Y
		case sdl.K_b:
			return B
		case sdl.K_ESCAPE:
			return B
		default:
			return Unknown
		}
	}

	return Unknown
}

func (s *JoyButtonMapperStrategy) ApplicableTo(event sdl.Event) bool {
	e, ok := event.(*sdl.JoyButtonEvent)
	if !ok {
		return false
	}

	return e.Type == sdl.JOYBUTTONDOWN
}

func (s *JoyButtonMapperStrategy) MapInputToKey(event sdl.Event) Key {
	e, _ := event.(*sdl.JoyButtonEvent)
	button := e.Button

	switch button {
	case 3:
		return A
	case 4:
		return B
	case 5:
		return Y
	case 6: // X button
		return X
	default:
		return Unknown
	}
}

func (s *JoyHatMapperStrategy) ApplicableTo(event sdl.Event) bool {
	_, ok := event.(*sdl.JoyHatEvent)
	return ok
}

func (s *JoyHatMapperStrategy) MapInputToKey(event sdl.Event) Key {
	e, _ := event.(*sdl.JoyHatEvent)
	value := e.Value

	switch value {
	case sdl.HAT_UP:
		return Up
	case sdl.HAT_DOWN:
		return Down
	case sdl.HAT_LEFT:
		return Left
	case sdl.HAT_RIGHT:
		return Right
	default:
		return Unknown
	}
}
