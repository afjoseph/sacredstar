package lunation

import (
	"github.com/afjoseph/sacredstar/astropoint"
	"github.com/tidwall/btree"
)

type LunationType string

const (
	LunationTypeFullMoon LunationType = "full-moon"
	LunationTypeNewMoon  LunationType = "new-moon"
)

type Lunation struct {
	Type LunationType `json:"type"`
}

func (l *Lunation) String() string {
	switch l.Type {
	case LunationTypeFullMoon:
		return "FullMoon"
	case LunationTypeNewMoon:
		return "NewMoon"
	default:
		panic("unknown lunation type")
	}
}

func Calculate(moonPoint, sunPoint *astropoint.AstroPoint) *Lunation {
	// This is approximately encompasses 1 day before the lunation and 1 day
	// after
	const orb = 13
	type degreeInterval struct {
		Type LunationType
		Deg  float64
	}

	// A lunation is either a conjunction or opposition between the moon
	// and the sun, with a wide orb capturing 2 days before and after
	diff := moonPoint.ZodiacalPos.DiffInAbsDegrees(sunPoint.ZodiacalPos)
	tree := *btree.NewBTreeG[degreeInterval](func(a, b degreeInterval) bool {
		return (a.Deg - orb) < (b.Deg - orb)
	})
	tree.Set(degreeInterval{Type: LunationTypeFullMoon, Deg: 180})
	tree.Set(degreeInterval{Type: LunationTypeNewMoon, Deg: 0})

	didFind := false
	lunationType := LunationTypeNewMoon
	tree.Scan(func(di degreeInterval) bool {
		if (di.Deg-orb) <= diff && diff <= (di.Deg+orb) {
			didFind = true
			lunationType = di.Type
			return false
		}
		return true
	})
	if !didFind {
		return nil
	}
	return &Lunation{
		Type: lunationType,
	}
}
