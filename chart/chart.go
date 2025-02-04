package chart

import (
	// #cgo pkg-config: swisseph
	// #include <stdio.h>
	// #include <errno.h>
	// #include "swephexp.h"
	"C"
)
import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/afjoseph/sacredstar/aspect"
	"github.com/afjoseph/sacredstar/astropoint"
	"github.com/afjoseph/sacredstar/house"
	"github.com/afjoseph/sacredstar/lunation"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/sign"
	"github.com/afjoseph/sacredstar/timeandzone"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/afjoseph/sacredstar/zodiacalpos"
	"github.com/go-playground/errors/v5"
)

type HouseLordPlacement int

const (
	HouseLordPlacement_Traditional HouseLordPlacement = iota
	HouseLordPlacement_Modern
)

type Chart struct {
	Time      timeandzone.TimeAndZone  `json:"time"`
	ChartType ChartType                `json:"chartType"`
	Points    []*astropoint.AstroPoint `json:"points"`
	Aspects   []*aspect.Aspect         `json:"aspects"`
	Lunation  *lunation.Lunation       `json:"lunations"`
}

func (c *Chart) String() string {
	var sb strings.Builder
	sb.WriteString(
		fmt.Sprintf(
			"Chart{Time: %s, ChartType: %s, Points: [\n",
			c.Time,
			c.ChartType,
		),
	)
	for _, p := range c.Points {
		sb.WriteString("\t" + p.String())
		sb.WriteString("\n")
	}
	sb.WriteString("]")
	sb.WriteString(", Houses: [\n")
	signs, err := c.SignsInOrder()
	if err != nil {
		panic(fmt.Errorf("while getting signs in order: %v", err))
	}
	for i, s := range signs {
		sb.WriteString(fmt.Sprintf("\tHouse %d: %s\n", i+1, s))
	}
	sb.WriteString("]}")
	return sb.String()
}

// SignsInOrder gets the signs in order starting from the ascendant.
// This means the ascendant is in the first house (using whole sign houses)
func (c *Chart) SignsInOrder() ([]sign.Sign, error) {
	asc := c.MustGetPoint(pointid.ASC)
	signs := []sign.Sign{}
	for i := 0; i < 12; i++ {
		s, err := sign.NewSignFromInt(asc.ZodiacalPos.Sign.Int() + i)
		if err != nil {
			return nil, fmt.Errorf("while getting sign for house %d: %v",
				i+1, err)
		}
		signs = append(signs, s)
	}
	return signs, nil
}

func (c *Chart) GetSignOfHouse(h house.House) sign.Sign {
	asc := c.MustGetPoint(pointid.ASC)
	s, err := sign.NewSignFromInt(asc.ZodiacalPos.Sign.Int() + (h.Int() - 1))
	if err != nil {
		panic(fmt.Errorf("while getting sign for house %d: %v", h, err))
	}
	return s
}

func NewChartFromJulianDay(
	swe *wrapper.SwissEph,
	timeInJulian float64,
	lon, lat float64,
	calcType ChartType,
	pointIDs []pointid.PointID,
) (*Chart, error) {
	// Calculate the ascendant always since we use it to calculate
	// the houses for all the other points
	asc, err := CalculateAscendant(
		swe,
		timeInJulian,
		lon,
		lat,
		calcType,
	)
	if err != nil {
		return nil, fmt.Errorf("while calculating ascendant: %v", err)
	}
	points := []*astropoint.AstroPoint{}
	points = append(points, asc)

	// For each point, calculate the longitude, sign and house
	didCalculateRahuKetu := false
	for _, id := range pointIDs {
		if id == pointid.Ketu || id == pointid.Rahu {
			if didCalculateRahuKetu {
				// Already calculated
				continue
			}
			var rahu, ketu *astropoint.AstroPoint
			rahu, ketu, err = calculateRahuKetu(
				swe,
				timeInJulian,
				calcType,
				asc.ZodiacalPos,
			)
			if rahu != nil && ketu != nil {
				didCalculateRahuKetu = true
				points = append(points, rahu)
				points = append(points, ketu)
			}
		} else if id == pointid.ASC {
			// Already calculated
			continue
		} else {
			var p *astropoint.AstroPoint
			p, err = calculatePlanet(
				swe,
				timeInJulian,
				id,
				calcType,
				asc.ZodiacalPos,
			)
			if p != nil {
				points = append(points, p)
			}
		}
		if err != nil {
			return nil, fmt.Errorf("while calculating planet %s: %v",
				id, err)
		}
	}

	// For each point, calculate the aspect with other points
	aspects := []*aspect.Aspect{}
	for _, p1 := range points {
		for _, p2 := range points {
			if p1 == p2 {
				continue
			}
			asp := aspect.NewAspect(
				p1.ID,
				p1.ZodiacalPos,
				p2.ID,
				p2.ZodiacalPos,
			)
			if asp == nil || asp.Type == aspect.AspectType_None {
				continue
			}
			aspects = append(aspects, asp)
		}
	}

	gotime := swe.JulianDayToGoTime(timeInJulian)
	tm := timeandzone.New(gotime)
	chrt := &Chart{
		Time:      tm,
		ChartType: calcType,
		Points:    points,
		Aspects:   aspects,
	}

	// Check if the pointIDs includes both the moon and the sun, else we can't
	// calculate lunations
	hasPID := func(pid pointid.PointID) bool {
		for _, p := range points {
			if p.ID == pid {
				return true
			}
		}
		return false
	}
	if !hasPID(pointid.Moon) || !hasPID(pointid.Sun) {
		return chrt, nil
	}
	chrt.Lunation = lunation.Calculate(
		chrt.MustGetPoint(pointid.Moon),
		chrt.MustGetPoint(pointid.Sun),
	)
	return chrt, nil
}

func CalculateAscendant(
	swe *wrapper.SwissEph,
	timeInJulian float64,
	lon, lat float64,
	calcType ChartType,
) (*astropoint.AstroPoint, error) {
	// Calculate houses and ascendant
	cusps := make([]C.double, 13)
	cuspsPtr := &(cusps[0])
	ascmc := make([]C.double, 10)
	ascmcPtr := &(ascmc[0])
	var ascZodiacalPos *zodiacalpos.ZodiacalPos
	if calcType.IsVarga() {
		if ret := C.swe_houses_ex(
			C.double(timeInJulian),
			C.int(C.SEFLG_SIDEREAL),
			C.double(lat),
			C.double(lon),
			// Whole sign
			C.int('W'),
			// Output
			cuspsPtr,
			ascmcPtr,
		); ret < 0 {
			return nil, fmt.Errorf("swe_houses_ex failed")
		}
		// Translate the degree to a sign and degree based on
		// the varga chart type
		var err error
		ascZodiacalPos, err = transformZodiacalPosToVarga(
			pointid.ASC,
			zodiacalpos.NewZodiacalPosFromLongitude(float64(ascmc[0])),
			calcType,
		)
		if err != nil {
			return nil, fmt.Errorf("while transforming ascendant to %s: %v",
				calcType, err)
		}
		// ascDeg = 0
	} else if calcType == TropicalChartType {
		if ret := C.swe_houses(
			C.double(timeInJulian),
			C.double(lat),
			C.double(lon),
			// Whole sign
			C.int('W'),
			// Output
			cuspsPtr,
			ascmcPtr,
		); ret < 0 {
			return nil, fmt.Errorf("swe_houses_ex failed")
		}
		ascZodiacalPos = zodiacalpos.NewZodiacalPosFromLongitude(float64(ascmc[0]))
	}

	return &astropoint.AstroPoint{
		ID:          pointid.ASC,
		Longitude:   float64(ascmc[0]),
		ZodiacalPos: ascZodiacalPos,
		House:       house.House1,
	}, nil
}

func calculateRahuKetu(
	swe *wrapper.SwissEph,
	timeInJulian float64,
	chartCalcType ChartType,
	ascendantZodiacalPos *zodiacalpos.ZodiacalPos,
) (*astropoint.AstroPoint, *astropoint.AstroPoint, error) {
	// XXX <26-01-2024, afjoseph> SwissEph doesn't have a way to
	// calculate Ketu, but it knows Rahu as C.SE_TRUE_NODE.
	// Ketu is basically the opposite of Rahu.
	rahu, err := calculatePlanet(
		swe,
		timeInJulian,
		pointid.Rahu,
		chartCalcType,
		ascendantZodiacalPos,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("while calculating Rahu: %v", err)
	}

	// Ketu as the opposite of Rahu
	ketuZodPos := rahu.ZodiacalPos.Opposite()
	return rahu, &astropoint.AstroPoint{
		ID:          pointid.Ketu,
		Longitude:   ketuZodPos.AbsDegrees(),
		ZodiacalPos: ketuZodPos,
		House:       rahu.House.Opposite(),
	}, nil
}

// calculatePlanet calculates an astropoint.AstroPoint for a given time and a
// point. If ascendantZodiacalPos is not nil, it will calculate the house for
// the point as well.
func calculatePlanet(
	swe *wrapper.SwissEph,
	timeInJulian float64,
	pid pointid.PointID,
	chartType ChartType,
	ascendantZodiacalPos *zodiacalpos.ZodiacalPos,
) (*astropoint.AstroPoint, error) {
	flag := C.int(0)
	if chartType.IsVarga() {
		flag = C.int(C.SEFLG_SIDEREAL)
	}
	// Add SEFLG_SPEED to flag
	flag |= C.int(C.SEFLG_SPEED)
	errBytes := make([]byte, C.AS_MAXCH)
	errPtr := (*C.char)(C.CBytes(errBytes))
	defer C.free(unsafe.Pointer(errPtr))
	// - lon
	// - lat
	// - dist
	// - speed in lon
	// - speed in lat
	// - speed in dist
	xx := make([]C.double, 6)
	xxPtr := &(xx[0])
	ret := C.swe_calc_ut(
		// Julian day
		C.double(timeInJulian),
		// Planet ID
		// XXX <19-01-2024, afjoseph> PointID are organized
		// in the same way swisseph accepts them, so Sun is 0, Moon
		// is 1, etc.
		C.int(pid.SwissEphID()),
		// iflag
		flag,
		// C.int(0),
		// Output
		xxPtr,
		// Error
		errPtr,
	)
	if ret < 0 {
		return nil, fmt.Errorf("swe_calc_ut failed: %s",
			C.GoString(errPtr))
	}

	var err error
	zp := zodiacalpos.NewZodiacalPosFromLongitude(float64(xx[0]))
	if chartType.IsVarga() {
		zp, err = transformZodiacalPosToVarga(pid, zp, chartType)
		if err != nil {
			return nil, fmt.Errorf("while transforming zodiacal pos to %s: %v",
				chartType, err)
		}
	}
	// isRetrograde :=
	h := house.HouseNone
	if ascendantZodiacalPos != nil {
		h = house.NewHouseFromSign(zp.Sign, ascendantZodiacalPos.Sign)
	}
	p := &astropoint.AstroPoint{
		ID:           pid,
		Longitude:    float64(xx[0]),
		ZodiacalPos:  zp,
		House:        h,
		IsRetrograde: float64(xx[3]) < 0,
	}

	return p, nil
}

// https://www.astro.com/swisseph/swephprg.htm#_Toc112948950
func NewChartFromUTC(
	swe *wrapper.SwissEph,
	date time.Time,
	lon, lat float64,
	calcType ChartType,
	pointIDs []pointid.PointID,
) (*Chart, error) {
	return NewChartFromJulianDay(
		swe,
		swe.GoTimeToJulianDay(date),
		lon,
		lat,
		calcType,
		pointIDs,
	)
}

func (c *Chart) GetHouseLordFor(
	h house.House,
	placementType HouseLordPlacement,
) (pointid.PointID, error) {
	houseAsInt := h.Int()
	asc := c.MustGetPoint(pointid.ASC)
	// Get the sign that corresponds to the house
	s, err := sign.NewSignFromInt(asc.ZodiacalPos.Sign.Int() + (houseAsInt - 1))
	if err != nil {
		return pointid.None, fmt.Errorf("while getting sign for house %d: %v",
			houseAsInt, err)
	}
	switch placementType {
	case HouseLordPlacement_Traditional:
		return s.TraditionalRuler(), nil
	case HouseLordPlacement_Modern:
		return s.ModernRuler(), nil
	default:
		return pointid.None, fmt.Errorf("invalid placement type: %d",
			placementType)
	}
}

func (c *Chart) GetZodiacalPosFor(id pointid.PointID) *zodiacalpos.ZodiacalPos {
	for _, p := range c.Points {
		if p.ID == id {
			return p.ZodiacalPos
		}
	}
	return nil
}

func (c *Chart) GetHouseFor(id pointid.PointID) house.House {
	p := c.GetPoint(id)
	if p == nil {
		return house.HouseNone
	}
	return p.House
}

func (c *Chart) GetPoint(id pointid.PointID) *astropoint.AstroPoint {
	for _, p := range c.Points {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func (c *Chart) MustGetPoint(id pointid.PointID) *astropoint.AstroPoint {
	p := c.GetPoint(id)
	if p == nil {
		panic(errors.Newf("point %s not found in chart: %s", id, c.String()))
	}
	return p
}

func (c *Chart) HasAspectIgnoreDegree(asp *aspect.Aspect) bool {
	for _, a := range c.Aspects {
		if asp.Equals(a, true) {
			return true
		}
	}
	return false
}
