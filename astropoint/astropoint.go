package astropoint

import (
	"fmt"

	"github.com/afjoseph/sacredstar/aspect"
	"github.com/afjoseph/sacredstar/house"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/sign"
	"github.com/afjoseph/sacredstar/zodiacalpos"
)

type AstroPoint struct {
	ID           pointid.PointID          `json:"id"`
	Longitude    float64                  `json:"longitude"`
	ZodiacalPos  *zodiacalpos.ZodiacalPos `json:"zodiacalPos"`
	House        house.House              `json:"house"`
	IsRetrograde bool                     `json:"isRetrograde"`
}

func (p *AstroPoint) String() string {
	return fmt.Sprintf(
		"AstroPoint{ID: %s, Longitude: %f, ZodiacalPos: %s, House: %s, IsRetrograde: %t}",
		p.ID,
		p.Longitude,
		p.ZodiacalPos,
		p.House,
		p.IsRetrograde,
	)
}

func (p *AstroPoint) SimpleDescription() string {
	if p.IsRetrograde {
		return fmt.Sprintf(
			"%s in %s in %s house (rx)",
			p.ID,
			p.ZodiacalPos.Sign,
			p.House,
		)
	} else {
		return fmt.Sprintf(
			"%s in %s in %s house",
			p.ID,
			p.ZodiacalPos.Sign,
			p.House,
		)
	}
}

func (p *AstroPoint) IsDomicile() bool {
	s := p.ZodiacalPos.Sign
	switch p.ID {
	case pointid.Sun:
		return s == sign.Leo
	case pointid.Moon:
		return s == sign.Cancer
	case pointid.Mercury:
		return s == sign.Gemini || s == sign.Virgo
	case pointid.Venus:
		return s == sign.Taurus || s == sign.Libra
	case pointid.Mars:
		return s == sign.Aries || s == sign.Scorpio
	case pointid.Jupiter:
		return s == sign.Sagittarius || s == sign.Pisces
	case pointid.Saturn:
		return s == sign.Capricorn || s == sign.Aquarius
	}
	return false
}

func (p *AstroPoint) IsDetriment() bool {
	s := p.ZodiacalPos.Sign
	switch p.ID {
	case pointid.Sun:
		return s == sign.Aquarius
	case pointid.Moon:
		return s == sign.Capricorn
	case pointid.Mercury:
		return s == sign.Sagittarius || s == sign.Pisces
	case pointid.Venus:
		return s == sign.Aries || s == sign.Scorpio
	case pointid.Mars:
		return s == sign.Taurus || s == sign.Libra
	case pointid.Jupiter:
		return s == sign.Gemini || s == sign.Virgo
	case pointid.Saturn:
		return s == sign.Cancer || s == sign.Leo
	}
	return false
}

func (p *AstroPoint) IsExalted() bool {
	s := p.ZodiacalPos.Sign
	switch p.ID {
	case pointid.Sun:
		return s == sign.Aries
	case pointid.Moon:
		return s == sign.Taurus
	case pointid.Mercury:
		return s == sign.Virgo
	case pointid.Venus:
		return s == sign.Pisces
	case pointid.Mars:
		return s == sign.Capricorn
	case pointid.Jupiter:
		return s == sign.Cancer
	case pointid.Saturn:
		return s == sign.Libra
	}
	return false
}

func (p *AstroPoint) IsFall() bool {
	s := p.ZodiacalPos.Sign
	switch p.ID {
	case pointid.Sun:
		return s == sign.Libra
	case pointid.Moon:
		return s == sign.Scorpio
	case pointid.Mercury:
		return s == sign.Pisces
	case pointid.Venus:
		return s == sign.Virgo
	case pointid.Mars:
		return s == sign.Cancer
	case pointid.Jupiter:
		return s == sign.Capricorn
	case pointid.Saturn:
		return s == sign.Aries
	}
	return false
}

func (ap *AstroPoint) GetAspect(rhs *AstroPoint) *aspect.Aspect {
	return aspect.NewAspect(ap.ID, ap.ZodiacalPos, rhs.ID, rhs.ZodiacalPos)
}
