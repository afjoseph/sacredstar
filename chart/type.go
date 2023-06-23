package chart

import (
	"github.com/afjoseph/sacredstar/house"
	"github.com/afjoseph/sacredstar/pointid"
)

type ChartType string

const (
	TropicalChartType = ChartType("tropical")
	D1ChartType       = ChartType("d1")
	D4ChartType       = ChartType("d4")
	D7ChartType       = ChartType("d7")
	D9ChartType       = ChartType("d9")
	D10ChartType      = ChartType("d10")
	// D24ChartType      = ChartType("D24")
)

func (c ChartType) String() string {
	return string(c)
}

// Useful for sorting
func (c ChartType) Int() int {
	switch c {
	case TropicalChartType:
		return 0
	case D1ChartType:
		return 1
	case D4ChartType:
		return 4
	case D7ChartType:
		return 7
	case D9ChartType:
		return 9
	case D10ChartType:
		return 10
		// case D24ChartType:
		// 	return 24
	}
	panic("unreachable")
}

func (c ChartType) Desc() string {
	switch c {
	case D1ChartType:
		return "D1 Rashi"
	case D4ChartType:
		return "D4 Chaturthamsa - Moving home"
	case D7ChartType:
		return "D7 Saptamsa - Children"
	case D9ChartType:
		return "D9 Navamsa - Marriage"
	case D10ChartType:
		return "D10 Dasamsa - Career"
		// case D24ChartType:
		// 	return VargaChartFunctionTypeD24
	}

	panic("unreachable")
}

func (c ChartType) IsVarga() bool {
	return c == D1ChartType ||
		c == D4ChartType ||
		c == D7ChartType ||
		c == D9ChartType ||
		c == D10ChartType
}

func (c ChartType) Karakas() []pointid.PointID {
	switch c {
	case D4ChartType:
		return []pointid.PointID{pointid.Rahu}
	case D7ChartType:
		return []pointid.PointID{pointid.Jupiter}
	case D9ChartType:
		return []pointid.PointID{pointid.Venus}
	case D10ChartType:
		return []pointid.PointID{
			pointid.Sun,
			pointid.Mercury,
			pointid.Jupiter,
			pointid.Saturn,
		}
		// case D24ChartType:
		// 	return []PointID{
		// 		pointid.Mercury,
		// 		pointid.Jupiter,
		// 	}
	}

	return nil
}

func (c ChartType) ImportantHouses() []house.House {
	switch c {
	case D4ChartType:
		return []house.House{house.House7, house.House12}
	case D7ChartType:
		return []house.House{house.House5}
	case D9ChartType:
		return []house.House{house.House7}
	case D10ChartType:
		return []house.House{house.House10}
		// case D24ChartType:
		// 	return []House{House5}
	}

	return nil
}
