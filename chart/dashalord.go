package chart

import (
	"fmt"
	"strconv"
	"time"

	"github.com/afjoseph/sacredstar/pointid"
)

type DashaLord int

const (
	DashaLordKetu DashaLord = iota
	DashaLordVenus
	DashaLordSun
	DashaLordMoon
	DashaLordMars
	DashaLordRahu
	DashaLordJupiter
	DashaLordSaturn
	DashaLordMercury
	DashaLordNone DashaLord = -1

	TotalDashaLords                = 9
	TotalDashaLordPairCombinations = 9 * 9
)

var MahadashaList = []DashaLord{
	DashaLordKetu,
	DashaLordVenus,
	DashaLordSun,
	DashaLordMoon,
	DashaLordMars,
	DashaLordRahu,
	DashaLordJupiter,
	DashaLordSaturn,
	DashaLordMercury,
}

func NewDashaLordFromPointID(id pointid.PointID) (DashaLord, error) {
	switch id {
	case pointid.Ketu:
		return DashaLordKetu, nil
	case pointid.Venus:
		return DashaLordVenus, nil
	case pointid.Sun:
		return DashaLordSun, nil
	case pointid.Moon:
		return DashaLordMoon, nil
	case pointid.Mars:
		return DashaLordMars, nil
	case pointid.Rahu:
		return DashaLordRahu, nil
	case pointid.Jupiter:
		return DashaLordJupiter, nil
	case pointid.Saturn:
		return DashaLordSaturn, nil
	case pointid.Mercury:
		return DashaLordMercury, nil
	}
	return DashaLordNone, fmt.Errorf("Invalid PointID: %s", id)
}

func (dl DashaLord) MarshalJSON() ([]byte, error) {
	return []byte(`"` + dl.ToPointID().String() + `"`), nil
}

func (dl *DashaLord) UnmarshalJSON(b []byte) error {
	var err error
	s := string(b)
	s, err = strconv.Unquote(s)
	if err != nil {
		return fmt.Errorf("Invalid DashaLord: %s", string(b))
	}
	*dl, err = NewDashaLordFromPointID(pointid.PointID(s))
	if err != nil {
		return fmt.Errorf("Invalid DashaLord: %s", string(b))
	}
	return nil
}

func (dl DashaLord) Int() int {
	return int(dl)
}

func (dl DashaLord) String() string {
	switch dl {
	case DashaLordKetu:
		return "Ketu"
	case DashaLordVenus:
		return "Venus"
	case DashaLordSun:
		return "Sun"
	case DashaLordMoon:
		return "Moon"
	case DashaLordMars:
		return "Mars"
	case DashaLordRahu:
		return "Rahu"
	case DashaLordJupiter:
		return "Jupiter"
	case DashaLordSaturn:
		return "Saturn"
	case DashaLordMercury:
		return "Mercury"
	}
	panic("Invalid DashaLord")
}

func (dl DashaLord) ToPointID() pointid.PointID {
	switch dl {
	case DashaLordKetu:
		return pointid.Ketu
	case DashaLordVenus:
		return pointid.Venus
	case DashaLordSun:
		return pointid.Sun
	case DashaLordMoon:
		return pointid.Moon
	case DashaLordMars:
		return pointid.Mars
	case DashaLordRahu:
		return pointid.Rahu
	case DashaLordJupiter:
		return pointid.Jupiter
	case DashaLordSaturn:
		return pointid.Saturn
	case DashaLordMercury:
		return pointid.Mercury
	}
	panic("Invalid DashaLord")
}

func (dl DashaLord) Next() DashaLord {
	i := dl.Int()
	i++
	if i == TotalDashaLords {
		i = 0
	}
	return DashaLord(i)
}

func parseDuration(year, mon, day int) time.Duration {
	d := year * 8760 // Hours in a calendar year
	d += mon * 730   // Hours in a calendar month
	d += day * 24    // Hours in a day
	return time.Duration(d) * time.Hour
}

func (dl DashaLord) MahadashaDuration() time.Duration {
	switch dl {
	case DashaLordKetu:
		return parseDuration(7, 0, 0)
	case DashaLordVenus:
		return parseDuration(20, 0, 0)
	case DashaLordSun:
		return parseDuration(6, 0, 0)
	case DashaLordMoon:
		return parseDuration(10, 0, 0)
	case DashaLordMars:
		return parseDuration(7, 0, 0)
	case DashaLordRahu:
		return parseDuration(18, 0, 0)
	case DashaLordJupiter:
		return parseDuration(16, 0, 0)
	case DashaLordSaturn:
		return parseDuration(19, 0, 0)
	case DashaLordMercury:
		return parseDuration(17, 0, 0)
	}
	panic("Invalid DashaLord")
}

func AntardashaDuration(mdl DashaLord, adl DashaLord) time.Duration {
	switch {
	case mdl == DashaLordKetu && adl == DashaLordKetu:
		return parseDuration(0, 4, 27)
	case mdl == DashaLordKetu && adl == DashaLordVenus:
		return parseDuration(1, 2, 0)
	case mdl == DashaLordKetu && adl == DashaLordSun:
		return parseDuration(0, 4, 6)
	case mdl == DashaLordKetu && adl == DashaLordMoon:
		return parseDuration(0, 7, 0)
	case mdl == DashaLordKetu && adl == DashaLordMars:
		return parseDuration(0, 4, 27)
	case mdl == DashaLordKetu && adl == DashaLordRahu:
		return parseDuration(1, 0, 18)
	case mdl == DashaLordKetu && adl == DashaLordJupiter:
		return parseDuration(0, 11, 6)
	case mdl == DashaLordKetu && adl == DashaLordSaturn:
		return parseDuration(1, 1, 9)
	case mdl == DashaLordKetu && adl == DashaLordMercury:
		return parseDuration(0, 11, 27)

	case mdl == DashaLordVenus && adl == DashaLordVenus:
		return parseDuration(3, 4, 0)
	case mdl == DashaLordVenus && adl == DashaLordSun:
		return parseDuration(1, 0, 0)
	case mdl == DashaLordVenus && adl == DashaLordMoon:
		return parseDuration(1, 8, 0)
	case mdl == DashaLordVenus && adl == DashaLordMars:
		return parseDuration(1, 2, 0)
	case mdl == DashaLordVenus && adl == DashaLordRahu:
		return parseDuration(3, 0, 0)
	case mdl == DashaLordVenus && adl == DashaLordJupiter:
		return parseDuration(2, 8, 0)
	case mdl == DashaLordVenus && adl == DashaLordSaturn:
		return parseDuration(3, 2, 0)
	case mdl == DashaLordVenus && adl == DashaLordMercury:
		return parseDuration(2, 10, 0)
	case mdl == DashaLordVenus && adl == DashaLordKetu:
		return parseDuration(1, 2, 0)

	case mdl == DashaLordSun && adl == DashaLordSun:
		return parseDuration(0, 3, 18)
	case mdl == DashaLordSun && adl == DashaLordMoon:
		return parseDuration(0, 6, 0)
	case mdl == DashaLordSun && adl == DashaLordMars:
		return parseDuration(0, 4, 6)
	case mdl == DashaLordSun && adl == DashaLordRahu:
		return parseDuration(0, 10, 24)
	case mdl == DashaLordSun && adl == DashaLordJupiter:
		return parseDuration(0, 9, 18)
	case mdl == DashaLordSun && adl == DashaLordSaturn:
		return parseDuration(0, 11, 12)
	case mdl == DashaLordSun && adl == DashaLordMercury:
		return parseDuration(0, 10, 6)
	case mdl == DashaLordSun && adl == DashaLordKetu:
		return parseDuration(0, 4, 6)
	case mdl == DashaLordSun && adl == DashaLordVenus:
		return parseDuration(1, 0, 0)

	case mdl == DashaLordMoon && adl == DashaLordMoon:
		return parseDuration(0, 10, 0)
	case mdl == DashaLordMoon && adl == DashaLordMars:
		return parseDuration(0, 7, 0)
	case mdl == DashaLordMoon && adl == DashaLordRahu:
		return parseDuration(1, 6, 0)
	case mdl == DashaLordMoon && adl == DashaLordJupiter:
		return parseDuration(1, 4, 0)
	case mdl == DashaLordMoon && adl == DashaLordSaturn:
		return parseDuration(1, 7, 0)
	case mdl == DashaLordMoon && adl == DashaLordMercury:
		return parseDuration(1, 5, 0)
	case mdl == DashaLordMoon && adl == DashaLordKetu:
		return parseDuration(0, 7, 0)
	case mdl == DashaLordMoon && adl == DashaLordVenus:
		return parseDuration(1, 8, 0)
	case mdl == DashaLordMoon && adl == DashaLordSun:
		return parseDuration(0, 6, 0)

	case mdl == DashaLordMars && adl == DashaLordMars:
		return parseDuration(0, 4, 27)
	case mdl == DashaLordMars && adl == DashaLordRahu:
		return parseDuration(1, 0, 18)
	case mdl == DashaLordMars && adl == DashaLordJupiter:
		return parseDuration(0, 11, 6)
	case mdl == DashaLordMars && adl == DashaLordSaturn:
		return parseDuration(1, 1, 9)
	case mdl == DashaLordMars && adl == DashaLordMercury:
		return parseDuration(0, 11, 27)
	case mdl == DashaLordMars && adl == DashaLordKetu:
		return parseDuration(0, 4, 27)
	case mdl == DashaLordMars && adl == DashaLordVenus:
		return parseDuration(1, 2, 0)
	case mdl == DashaLordMars && adl == DashaLordSun:
		return parseDuration(0, 4, 6)
	case mdl == DashaLordMars && adl == DashaLordMoon:
		return parseDuration(0, 7, 0)

	case mdl == DashaLordRahu && adl == DashaLordRahu:
		return parseDuration(2, 8, 12)
	case mdl == DashaLordRahu && adl == DashaLordJupiter:
		return parseDuration(2, 4, 24)
	case mdl == DashaLordRahu && adl == DashaLordSaturn:
		return parseDuration(2, 10, 6)
	case mdl == DashaLordRahu && adl == DashaLordMercury:
		return parseDuration(2, 6, 18)
	case mdl == DashaLordRahu && adl == DashaLordKetu:
		return parseDuration(1, 0, 18)
	case mdl == DashaLordRahu && adl == DashaLordVenus:
		return parseDuration(3, 0, 0)
	case mdl == DashaLordRahu && adl == DashaLordSun:
		return parseDuration(0, 10, 24)
	case mdl == DashaLordRahu && adl == DashaLordMoon:
		return parseDuration(1, 6, 0)
	case mdl == DashaLordRahu && adl == DashaLordMars:
		return parseDuration(1, 0, 18)

	case mdl == DashaLordJupiter && adl == DashaLordJupiter:
		return parseDuration(2, 1, 18)
	case mdl == DashaLordJupiter && adl == DashaLordSaturn:
		return parseDuration(2, 6, 12)
	case mdl == DashaLordJupiter && adl == DashaLordMercury:
		return parseDuration(2, 3, 6)
	case mdl == DashaLordJupiter && adl == DashaLordKetu:
		return parseDuration(0, 11, 6)
	case mdl == DashaLordJupiter && adl == DashaLordVenus:
		return parseDuration(2, 8, 0)
	case mdl == DashaLordJupiter && adl == DashaLordSun:
		return parseDuration(0, 9, 18)
	case mdl == DashaLordJupiter && adl == DashaLordMoon:
		return parseDuration(1, 4, 0)
	case mdl == DashaLordJupiter && adl == DashaLordMars:
		return parseDuration(0, 11, 6)
	case mdl == DashaLordJupiter && adl == DashaLordRahu:
		return parseDuration(2, 4, 24)

	case mdl == DashaLordSaturn && adl == DashaLordSaturn:
		return parseDuration(3, 0, 3)
	case mdl == DashaLordSaturn && adl == DashaLordMercury:
		return parseDuration(2, 8, 9)
	case mdl == DashaLordSaturn && adl == DashaLordKetu:
		return parseDuration(1, 1, 9)
	case mdl == DashaLordSaturn && adl == DashaLordVenus:
		return parseDuration(3, 2, 0)
	case mdl == DashaLordSaturn && adl == DashaLordSun:
		return parseDuration(0, 11, 12)
	case mdl == DashaLordSaturn && adl == DashaLordMoon:
		return parseDuration(1, 7, 0)
	case mdl == DashaLordSaturn && adl == DashaLordMars:
		return parseDuration(1, 1, 9)
	case mdl == DashaLordSaturn && adl == DashaLordRahu:
		return parseDuration(2, 10, 6)
	case mdl == DashaLordSaturn && adl == DashaLordJupiter:
		return parseDuration(2, 6, 12)

	case mdl == DashaLordMercury && adl == DashaLordMercury:
		return parseDuration(2, 4, 27)
	case mdl == DashaLordMercury && adl == DashaLordKetu:
		return parseDuration(0, 11, 27)
	case mdl == DashaLordMercury && adl == DashaLordVenus:
		return parseDuration(2, 10, 0)
	case mdl == DashaLordMercury && adl == DashaLordSun:
		return parseDuration(0, 10, 6)
	case mdl == DashaLordMercury && adl == DashaLordMoon:
		return parseDuration(1, 5, 0)
	case mdl == DashaLordMercury && adl == DashaLordMars:
		return parseDuration(0, 11, 27)
	case mdl == DashaLordMercury && adl == DashaLordRahu:
		return parseDuration(2, 6, 18)
	case mdl == DashaLordMercury && adl == DashaLordJupiter:
		return parseDuration(2, 3, 6)
	case mdl == DashaLordMercury && adl == DashaLordSaturn:
		return parseDuration(2, 8, 9)
	}

	panic("Invalid DashaLord combination")
}

func NextAntardasha(mdl DashaLord, adl DashaLord) (DashaLord, DashaLord) {
	// Get the sister antardashas for this mahadasha
	sisters := mdl.GetAntardashas()
	// If this is the last antardasha, return the next mahadasha and the first
	// antardasha with it
	if adl == sisters[len(sisters)-1] {
		mdl = mdl.Next()
		sisters = mdl.GetAntardashas()
		return mdl, sisters[0]
	}
	// Else, return the same mahadasha and the next antardasha
	return mdl, adl.Next()
}

func (dl DashaLord) GetAntardashas() []DashaLord {
	switch dl {
	case DashaLordKetu:
		return []DashaLord{
			DashaLordKetu,
			DashaLordVenus,
			DashaLordSun,
			DashaLordMoon,
			DashaLordMars,
			DashaLordRahu,
			DashaLordJupiter,
			DashaLordSaturn,
			DashaLordMercury,
		}

	case DashaLordVenus:
		return []DashaLord{
			DashaLordVenus,
			DashaLordSun,
			DashaLordMoon,
			DashaLordMars,
			DashaLordRahu,
			DashaLordJupiter,
			DashaLordSaturn,
			DashaLordMercury,
			DashaLordKetu,
		}

	case DashaLordSun:
		return []DashaLord{
			DashaLordSun,
			DashaLordMoon,
			DashaLordMars,
			DashaLordRahu,
			DashaLordJupiter,
			DashaLordSaturn,
			DashaLordMercury,
			DashaLordKetu,
			DashaLordVenus,
		}

	case DashaLordMoon:
		return []DashaLord{
			DashaLordMoon,
			DashaLordMars,
			DashaLordRahu,
			DashaLordJupiter,
			DashaLordSaturn,
			DashaLordMercury,
			DashaLordKetu,
			DashaLordVenus,
			DashaLordSun,
		}

	case DashaLordMars:
		return []DashaLord{
			DashaLordMars,
			DashaLordRahu,
			DashaLordJupiter,
			DashaLordSaturn,
			DashaLordMercury,
			DashaLordKetu,
			DashaLordVenus,
			DashaLordSun,
			DashaLordMoon,
		}

	case DashaLordRahu:
		return []DashaLord{
			DashaLordRahu,
			DashaLordJupiter,
			DashaLordSaturn,
			DashaLordMercury,
			DashaLordKetu,
			DashaLordVenus,
			DashaLordSun,
			DashaLordMoon,
			DashaLordMars,
		}

	case DashaLordJupiter:
		return []DashaLord{
			DashaLordJupiter,
			DashaLordSaturn,
			DashaLordMercury,
			DashaLordKetu,
			DashaLordVenus,
			DashaLordSun,
			DashaLordMoon,
			DashaLordMars,
			DashaLordRahu,
		}

	case DashaLordSaturn:
		return []DashaLord{
			DashaLordSaturn,
			DashaLordMercury,
			DashaLordKetu,
			DashaLordVenus,
			DashaLordSun,
			DashaLordMoon,
			DashaLordMars,
			DashaLordRahu,
			DashaLordJupiter,
		}

	case DashaLordMercury:
		return []DashaLord{
			DashaLordMercury,
			DashaLordKetu,
			DashaLordVenus,
			DashaLordSun,
			DashaLordMoon,
			DashaLordMars,
			DashaLordRahu,
			DashaLordJupiter,
			DashaLordSaturn,
		}
	}
	panic("unreachable")
}
