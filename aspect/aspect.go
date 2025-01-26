package aspect

import (
	"fmt"
	"math"

	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/zodiacalpos"
	"github.com/tidwall/btree"
)

var orbTable = map[AspectType]float64{
	AspectType_Conjunction: 5,
	AspectType_Opposition:  5,
	AspectType_Trine:       5,
	AspectType_Square:      5,
	AspectType_Sextile:     3,
}

type AspectType int

const (
	AspectType_Conjunction AspectType = iota
	AspectType_Opposition
	AspectType_Trine
	AspectType_Square
	AspectType_Sextile
	// AspectType_SemiSquare
	// AspectType_Sesquiquadrate
	// AspectType_SemiSextile
	AspectType_None
)

func (at AspectType) String() string {
	switch at {
	case AspectType_Conjunction:
		return "Conjunction"
	case AspectType_Opposition:
		return "Opposition"
	case AspectType_Trine:
		return "Trine"
	case AspectType_Square:
		return "Square"
	case AspectType_Sextile:
		return "Sextile"
	}
	return "None"
}

func (a AspectType) Degree() float64 {
	switch a {
	case AspectType_Conjunction:
		return 0
	case AspectType_Opposition:
		return 180
	case AspectType_Trine:
		return 120
	case AspectType_Square:
		return 90
	case AspectType_Sextile:
		return 60
	}
	return 0
}

type Aspect struct {
	P1     pointid.PointID `json:"p1"`
	P2     pointid.PointID `json:"p2"`
	Degree float64         `json:"degree"`
	Type   AspectType      `json:"type"`
}

type degreeInterval struct {
	Type AspectType
	Deg  float64
}

func NewAspect(
	lhsID pointid.PointID,
	lhsZP *zodiacalpos.ZodiacalPos,
	rhsID pointid.PointID,
	rhsZP *zodiacalpos.ZodiacalPos,
) *Aspect {
	diff := lhsZP.DiffInAbsDegrees(rhsZP)

	tree := *btree.NewBTreeG[degreeInterval](func(a, b degreeInterval) bool {
		orbA := orbTable[a.Type]
		orbB := orbTable[b.Type]
		return (a.Deg - orbA) < (b.Deg - orbB)
	})
	tree.Set(degreeInterval{Type: AspectType_Conjunction, Deg: 0})
	tree.Set(degreeInterval{Type: AspectType_Opposition, Deg: 180})
	tree.Set(degreeInterval{Type: AspectType_Trine, Deg: 120})
	tree.Set(degreeInterval{Type: AspectType_Square, Deg: 90})
	tree.Set(degreeInterval{Type: AspectType_Sextile, Deg: 60})

	didFind := false
	aspectType := AspectType_None
	tree.Scan(func(di degreeInterval) bool {
		orb := orbTable[di.Type]
		if (di.Deg-orb) <= diff && diff <= (di.Deg+orb) {
			didFind = true
			aspectType = di.Type
			return false
		}
		return true
	})
	if !didFind {
		aspectType = AspectType_None
	}
	return &Aspect{
		P1:     lhsID,
		P2:     rhsID,
		Degree: diff,
		Type:   aspectType,
	}
}

func (a *Aspect) Orb() int {
	return int(orbTable[a.Type])
}

func (a *Aspect) String() string {
	return fmt.Sprintf("Aspect{P1: %s, P2: %s, Degree: %f, Type: %s}",
		a.P1,
		a.P2,
		a.Degree,
		a.Type,
	)
}

func (a *Aspect) Int() int {
	return int(a.Type)
}

func (a *Aspect) IsHard() bool {
	return a.Type == AspectType_Conjunction ||
		a.Type == AspectType_Opposition ||
		a.Type == AspectType_Square
}

func (a *Aspect) IsSoft() bool {
	return a.Type == AspectType_Trine || a.Type == AspectType_Sextile
}

func (a *Aspect) Equals(other *Aspect, ignoreDegree bool) bool {
	if !ignoreDegree {
		eps := 0.000001
		if math.Abs(a.Degree-other.Degree) > eps {
			return false
		}
	}
	return a.P1 == other.P1 && a.P2 == other.P2 && a.Type == other.Type
}
