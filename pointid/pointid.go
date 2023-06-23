package pointid

import (
	"fmt"
	"strings"
)

type PointID string

func (p PointID) String() string {
	return string(p)
}

func (p PointID) SwissEphID() int {
	switch p {
	case Sun:
		return 0
	case Moon:
		return 1
	case Mercury:
		return 2
	case Venus:
		return 3
	case Mars:
		return 4
	case Jupiter:
		return 5
	case Saturn:
		return 6
	case Uranus:
		return 7
	case Neptune:
		return 8
	case Pluto:
		return 9
	case Rahu:
		return 11 // North Node: equal to C.SE_TRUE_NODE
	case Ketu:
		return -1 // Doesn't exist in SwissEph
	case ASC:
		return -1 // Calculated differently
	default:
		return -1
	}
}

func NewPointID(s string) (PointID, error) {
	switch strings.ToLower(s) {
	case "mercury":
		return Mercury, nil
	case "venus":
		return Venus, nil
	case "mars":
		return Mars, nil
	case "jupiter":
		return Jupiter, nil
	case "saturn":
		return Saturn, nil
	case "uranus":
		return Uranus, nil
	case "neptune":
		return Neptune, nil
	case "pluto":
		return Pluto, nil
	case "moon":
		return Moon, nil
	case "sun":
		return Sun, nil
	case "asc":
		return ASC, nil
	case "ketu":
		return Ketu, nil
	case "rahu":
		return Rahu, nil
	default:
		return None, fmt.Errorf("Unknown planet: %s", s)
	}
}

var TraditionalPlanets = []PointID{
	Sun,
	Moon,
	Mercury,
	Venus,
	Mars,
	Jupiter,
	Saturn,
}

var ClassicalPlanets = []PointID{
	Sun,
	Moon,
	Mercury,
	Venus,
	Mars,
	Jupiter,
	Saturn,
	Uranus,
	Neptune,
	Pluto,
}

var VedicPlanets = []PointID{
	Sun,
	Moon,
	Mercury,
	Venus,
	Mars,
	Jupiter,
	Saturn,
	Rahu,
	Ketu,
}

var (
	ASC     = PointID("asc")
	Sun     = PointID("sun")
	Moon    = PointID("moon")
	Mercury = PointID("mercury")
	Venus   = PointID("venus")
	Mars    = PointID("mars")
	Jupiter = PointID("jupiter")
	Saturn  = PointID("saturn")
	Uranus  = PointID("uranus")
	Neptune = PointID("neptune")
	Pluto   = PointID("pluto")
	Rahu    = PointID("rahu") // North Node: equal to C.SE_TRUE_NODE
	Ketu    = PointID("ketu")

	None = PointID("")
)
