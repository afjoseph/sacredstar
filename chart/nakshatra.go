package chart

import (
	"fmt"
	"math"
	"time"

	"github.com/afjoseph/sacredstar/astropoint"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/afjoseph/sacredstar/zodiacalpos"
)

type NakshatraType int

const (
	Aswini NakshatraType = iota
	Bharani
	Krittika
	Rohini
	Mrigashirsha
	Ardra
	Punarvasu
	Pushya
	Ashlesha
	Magha
	PurvaPhalguni
	UttaraPhalguni
	Hasta
	Chitra
	Swati
	Vishakha
	Anuradha
	Jyeshtha
	Mula
	PurvaAshadha
	UttaraAshadha
	Shravana
	Dhanishta
	Shatabhisha
	PurvaBhadrapada
	UttaraBhadrapada
	Revati

	NakshatraNone = -1
)

func (nt NakshatraType) String() string {
	switch nt {
	case Aswini:
		return "Aswini"
	case Bharani:
		return "Bharani"
	case Krittika:
		return "Krittika"
	case Rohini:
		return "Rohini"
	case Mrigashirsha:
		return "Mrigashirsha"
	case Ardra:
		return "Ardra"
	case Punarvasu:
		return "Punarvasu"
	case Pushya:
		return "Pushya"
	case Ashlesha:
		return "Ashlesha"
	case Magha:
		return "Magha"
	case PurvaPhalguni:
		return "PurvaPhalguni"
	case UttaraPhalguni:
		return "UttaraPhalguni"
	case Hasta:
		return "Hasta"
	case Chitra:
		return "Chitra"
	case Swati:
		return "Swati"
	case Vishakha:
		return "Vishakha"
	case Anuradha:
		return "Anuradha"
	case Jyeshtha:
		return "Jyeshtha"
	case Mula:
		return "Mula"
	case PurvaAshadha:
		return "PurvaAshadha"
	case UttaraAshadha:
		return "UttaraAshadha"
	case Shravana:
		return "Shravana"
	case Dhanishta:
		return "Dhanishta"
	case Shatabhisha:
		return "Shatabhisha"
	case PurvaBhadrapada:
		return "PurvaBhadrapada"
	case UttaraBhadrapada:
		return "UttaraBhadrapada"
	case Revati:
		return "Revati"
	}
	return "Unknown"
}

func (n NakshatraType) Int() int {
	return int(n) + 1
}

func (n *Nakshatra) String() string {
	return n.Type.String()
}

func NewNakshatraTypeFromInt(i int) (NakshatraType, error) {
	switch i {
	case 0:
		return Aswini, nil
	case 1:
		return Bharani, nil
	case 2:
		return Krittika, nil
	case 3:
		return Rohini, nil
	case 4:
		return Mrigashirsha, nil
	case 5:
		return Ardra, nil
	case 6:
		return Punarvasu, nil
	case 7:
		return Pushya, nil
	case 8:
		return Ashlesha, nil
	case 9:
		return Magha, nil
	case 10:
		return PurvaPhalguni, nil
	case 11:
		return UttaraPhalguni, nil
	case 12:
		return Hasta, nil
	case 13:
		return Chitra, nil
	case 14:
		return Swati, nil
	case 15:
		return Vishakha, nil
	case 16:
		return Anuradha, nil
	case 17:
		return Jyeshtha, nil
	case 18:
		return Mula, nil
	case 19:
		return PurvaAshadha, nil
	case 20:
		return UttaraAshadha, nil
	case 21:
		return Shravana, nil
	case 22:
		return Dhanishta, nil
	case 23:
		return Shatabhisha, nil
	case 24:
		return PurvaBhadrapada, nil
	case 25:
		return UttaraBhadrapada, nil
	case 26:
		return Revati, nil
	}
	return NakshatraNone, fmt.Errorf("invalid nakshatra %d", i)
}

type Nakshatra struct {
	Type NakshatraType
	Pada int
}

func NewNakshatraFromChart(
	swe *wrapper.SwissEph,
	ap *astropoint.AstroPoint,
) (*Nakshatra, error) {
	if ap == nil {
		return nil, fmt.Errorf("astro point is nil")
	}

	// There are 27 nakshatras, each is 13.333333333333334 degrees
	// long. The first nakshatra starts at 0 degrees Aries.
	nakshatraIdx := int(ap.Longitude / 13.333333333333334)
	nid, err := NewNakshatraTypeFromInt(nakshatraIdx)
	if err != nil {
		return nil, fmt.Errorf("NewNakshatraFromInt failed: %v", err)
	}

	pada := int(
		math.Floor(
			math.Mod(ap.Longitude, 13.333333333333334) / 3.3333333333333333,
		),
	)
	return &Nakshatra{nid, pada}, nil
}

func (n *Nakshatra) MahadashaLord() DashaLord {
	switch n.Type {
	case Aswini, Magha, Mula:
		return DashaLordKetu
	case Bharani, PurvaPhalguni, UttaraPhalguni:
		return DashaLordVenus
	case Krittika, UttaraAshadha, UttaraBhadrapada:
		return DashaLordSun
	case Rohini, Hasta, Shravana:
		return DashaLordMoon
	case Mrigashirsha, Chitra, Dhanishta:
		return DashaLordMars
	case Ardra, Swati, Shatabhisha:
		return DashaLordRahu
	case Punarvasu, Vishakha, PurvaBhadrapada:
		return DashaLordJupiter
	case Pushya, Anuradha:
		return DashaLordSaturn
	case Ashlesha, Jyeshtha, Revati:
		return DashaLordMercury
	}
	panic(fmt.Sprintf("unknown nakshatra %s", n.Type))
}

func (n *Nakshatra) MinZodiacalPos() *zodiacalpos.ZodiacalPos {
	return zodiacalpos.NewZodiacalPosFromLongitude(
		float64(n.Type.Int()-1) * 13.333333333333334,
	)
}

func (n *Nakshatra) MaxZodiacalPos() *zodiacalpos.ZodiacalPos {
	return zodiacalpos.NewZodiacalPosFromLongitude(
		float64(n.Type.Int()) * 13.333333333333334,
	)
}

func (n *Nakshatra) GetDashaLordPair(
	birthTime time.Time,
	percRemaining float64,
) (
	mahadasha DashaLord,
	antardasha DashaLord,
	remainingDurationInAntardasha time.Duration,
) {
	// Two special cases we can deal with early
	mahadashaLord := n.MahadashaLord()
	if percRemaining == 0 {
		lastAntardasha := mahadashaLord.GetAntardashas()[len(mahadashaLord.GetAntardashas())-1]
		return mahadashaLord, lastAntardasha, 0
	}
	if percRemaining == 1 {
		firstAntardasha := mahadashaLord.GetAntardashas()[0]
		return mahadashaLord, firstAntardasha, mahadashaLord.MahadashaDuration()
	}

	totalMahadashaDuration := n.MahadashaLord().MahadashaDuration()
	mahadashaDurationRemaining := time.Duration(
		totalMahadashaDuration.Minutes()*percRemaining,
	) * time.Minute
	// fmt.Printf(
	// 	"Mahadasha: %s | Duration: %+v\n",
	// 	n.MahadashaLord(),
	// 	durafmt.Parse(totalMahadashaDuration),
	// )
	mahadashaDurationElapsed := totalMahadashaDuration - mahadashaDurationRemaining
	// fmt.Printf(
	// 	"Mahadasha elapsed: %+v\n",
	// 	durafmt.Parse(mahadashaDurationElapsed),
	// )
	mahadashaStartTime := birthTime.Add(-mahadashaDurationElapsed)
	// fmt.Printf("Mahadasha start time: %+v\n", mahadashaStartTime)

	// Walk through the mahadasha from mahadashaStartTime till birthTime and
	// see which Antardasha falls in between
	start := mahadashaStartTime
	end := mahadashaStartTime
	antardashaLord := DashaLordNone
	durationSoFar := time.Duration(0)
	remainingAntardashaDuration := time.Duration(0)
	for _, antardasha := range n.MahadashaLord().GetAntardashas() {
		// fmt.Printf(
		// 	"Antardasha: %s | Duration: %+v\n",
		// 	antardasha,
		// 	durafmt.Parse(antarDuration()),
		// )
		end = end.Add(AntardashaDuration(mahadashaLord, antardasha))
		durationSoFar += end.Sub(start)
		// fmt.Printf(
		// 	"Pair: %s | Start: %+v | End: %+v | Duration: %+v | durationSoFar: %v\n",
		// 	antardasha,
		// 	start,
		// 	end,
		// 	durafmt.Parse(end.Sub(start)),
		// 	durafmt.Parse(durationSoFar),
		// )
		if birthTime.After(start) && birthTime.Before(end) {
			antardashaLord = antardasha
			remainingAntardashaDuration = mahadashaStartTime.Add(
				durationSoFar,
			).Sub(birthTime)
			// fmt.Printf(
			// 	"Found antardasha: %s | Remaining duration: %+v\n",
			// 	targatDashaLordPair,
			// 	durafmt.Parse(remainingAntardashaDuration),
			// )
			// durationSoFar -= end.Sub(birthTime)
			break
		}
		start = end
	}
	if antardashaLord == DashaLordNone {
		// XXX <02-02-2024, afjoseph> This must never happen
		panic(fmt.Sprintf(
			"could not find antardasha for %s at %v", n.MahadashaLord(),
			birthTime,
		))
	}

	return mahadashaLord, antardashaLord, remainingAntardashaDuration
}
