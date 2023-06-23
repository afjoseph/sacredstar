package chart

import (
	"fmt"
	"time"

	"github.com/afjoseph/sacredstar/astropoint"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/tidwall/btree"
)

// In Julian days
const SiderealMoonMonthDelta = 27.326

type DashaTree struct {
	tree      *btree.BTreeG[Dasha]
	startTime time.Time
}

func NewDashaTree(
	swe *wrapper.SwissEph,
	birthTime time.Time,
	moonAstroPoint *astropoint.AstroPoint,
) (*DashaTree, error) {
	chartTimeAsJulian := swe.GoTimeToJulianDay(birthTime.UTC())
	nak, err := NewNakshatraFromChart(swe, moonAstroPoint)
	if err != nil {
		return nil, fmt.Errorf("while calculating Nakshatra: %v", err)
	}
	// Calculate how long elapsed in the dasha
	// 1. Calculate the moon's zodiacal position
	// moonZodPos := zodiacalpos.NewZodiacalPosFromLongitude(moon.Longitude)
	// 2. Calculate the difference between the moon's zodiacal position and
	//    Nakshatra's max/min zodiacal position
	diffBetweenMoonAndMaxNakshatraZodPosInJulian := moonAstroPoint.ZodiacalPos.DiffInAbsDegrees(
		nak.MaxZodiacalPos(),
	)
	diffBetweenMoonAndMinNakshatraZodPosInJulian := moonAstroPoint.ZodiacalPos.DiffInAbsDegrees(
		nak.MinZodiacalPos(),
	)
	// 3. Find approx start/end time of Nakshatra
	approxNakshatraStartTimeInJulianDays := chartTimeAsJulian - calcMoonMovementInJulian(
		diffBetweenMoonAndMinNakshatraZodPosInJulian,
	)
	approxNakshatraEndTimeInJulianDays := chartTimeAsJulian + calcMoonMovementInJulian(
		diffBetweenMoonAndMaxNakshatraZodPosInJulian,
	)
	// 4. Find accurate start/end time of Nakshatra

	// fmt.Printf(
	// 	"moonZodPos: %+v | nak.MaxZodiacalPos: %+v | nak.MinZodiacalPos: %+v | chartTimeAsJulian: %f\n",
	// 	moonAstroPoint.ZodiacalPos,
	// 	nak.MaxZodiacalPos(),
	// 	nak.MinZodiacalPos(),
	// 	chartTimeAsJulian,
	// )
	exactNakshatraStartTimeInJulianDays, err := correctJulianDateToZodiacalPos(
		swe,
		approxNakshatraStartTimeInJulianDays,
		pointid.Moon,
		nak.MinZodiacalPos(),
		nak.Type,
		1.0/60.0/24.0, // Step: 1 minute
		1.0,           // Range: 1 day
	)
	if err != nil {
		return nil, fmt.Errorf(
			"while finding exact start time of Nakshatra: %v",
			err,
		)
	}
	if chartTimeAsJulian < exactNakshatraStartTimeInJulianDays {
		panic(fmt.Sprintf(
			"chartTimeAsJulian (%f) < exactNakshatraStartTimeInJulianDays (%f) for %+v",
			chartTimeAsJulian,
			exactNakshatraStartTimeInJulianDays,
			moonAstroPoint,
		))
	}
	exactNakshatraEndTimeInJulianDays, err := correctJulianDateToZodiacalPos(
		swe,
		approxNakshatraEndTimeInJulianDays,
		pointid.Moon,
		nak.MaxZodiacalPos(),
		nak.Type,
		0.5/60.0/24.0, // Step: 1 minute
		1.0,           // Range: 1 day
	)
	if err != nil {
		return nil, fmt.Errorf(
			"while finding exact start time of Nakshatra: %v",
			err,
		)
	}
	if chartTimeAsJulian > exactNakshatraEndTimeInJulianDays {
		panic(fmt.Sprintf(
			"chartTimeAsJulian (%f) > exactNakshatraEndTimeInJulianDays (%f) for moon at %+v",
			chartTimeAsJulian,
			exactNakshatraEndTimeInJulianDays,
			moonAstroPoint,
		))
	}
	// fmt.Printf(
	// 	"Nakshatra start time: %+v\n",
	// 	swe.JulianDayToGoTime(exactNakshatraStartTimeInJulianDays),
	// )
	// fmt.Printf(
	// 	"Nakshatra end time: %+v\n",
	// 	swe.JulianDayToGoTime(exactNakshatraEndTimeInJulianDays),
	// )
	// 5. Now, we can do accurate measurements on the Nakshatra and have it
	//    reflect to the Vimshottari Dasha calculations.
	//    Find the percentage of the Nakshatra that is remaining
	totalNakshatraDurationInJulianDays := exactNakshatraEndTimeInJulianDays - exactNakshatraStartTimeInJulianDays
	// fmt.Printf(
	// 	"totalNakshatraDuration: %v\n",
	// 	totalNakshatraDurationInJulianDays,
	// )
	// totalNakshatraDurationInMins := totalNakshatraDurationInJulianDays * 24.0 * 60.0
	remainingNakshatraDurationInJulianDays := exactNakshatraEndTimeInJulianDays - chartTimeAsJulian
	// elapsingNakshatraDurationInJulianDays := chart.Time - exactNakshatraStartTimeInJulianDays
	// fmt.Printf("remainingNakshatraDuration: %v\n", swe.JulianDayToGoTime(remainingNakshatraDurationInJulianDays))
	// fmt.Printf("totalNakshatraDuration: %v\n", swe.JulianDayToGoTime(totalNakshatraDurationInJulianDays))
	percNakshatraRemaining := remainingNakshatraDurationInJulianDays / totalNakshatraDurationInJulianDays
	if percNakshatraRemaining > 1 {
		// XXX <08-02-2024, afjoseph> This would mean there's an error in the calculation
		panic(fmt.Sprintf(
			"percNakshatraRemaining > 1: %+v for moon at %+v",
			percNakshatraRemaining,
			moonAstroPoint,
		))
	}
	// percNakshatraElapsing := elapsingNakshatraDurationInJulianDays / totalNakshatraDurationInJulianDays
	// fmt.Printf("percNakshatraRemaining: %+v\n", percNakshatraRemaining)
	// fmt.Printf("percNakshatraElapsing: %+v\n", percNakshatraElapsing)
	mahaDashaLord, antarDashaLord,
		remainingDurationInAntardasha := nak.GetDashaLordPair(birthTime, percNakshatraRemaining)
	// mahadashaDuration := dashaLordPair.MahadashaDuration()
	// fmt.Printf("Dasha: %s | Duration: %+v\n", dl, durafmt.Parse(mahadashaDuration))
	// mahadashaDurationRemaining := time.Duration(
	// 	mahadashaDuration.Minutes()*percNakshatraRemaining,
	// ) * time.Minute
	// and how much of it has passed
	// Finding the antardasha should be easy: just make a function that returns the antardasha
	// given a mahadasha and the remaining duration (or elapsed duration). The same function can be used
	// to find the percentage of the antardasha that has elapsed (and the remaining one).
	// After that's calculated, use it in the dashaEnd calculation below.

	// dashaDurationElapsing := time.Duration(
	// 	mahadashaDuration.Minutes()*percNakshatraElapsing,
	// ) * time.Minute
	// fmt.Printf(
	// 	"dashaDurationRemaining: %+v\n",
	// 	durafmt.Parse(dashaDurationRemaining),
	// )
	// fmt.Printf(
	// 	"dashaDurationElapsing: %+v\n",
	// 	durafmt.Parse(dashaDurationElapsing),
	// )
	// ddashaStart := swe.JulianDayToGoTime(chart.Time).Add(-dashaDurationElapsing)
	// ddashaEnd := swe.JulianDayToGoTime(chart.Time).Add(dashaDurationRemaining)
	// fmt.Printf("ddashaStart: %+v\n", ddashaStart)
	// fmt.Printf("ddashaEnd: %+v\n", ddashaEnd)

	// Now cast a dasha interval tree from this time
	t := btree.NewBTreeG[Dasha](func(a, b Dasha) bool {
		return a.Interval.Start.Before(b.Interval.Start.Time)
	})
	dashaStart := birthTime
	dashaEnd := dashaStart.Add(remainingDurationInAntardasha)
	for i := 0; i < TotalDashaLordPairCombinations; i++ {
		if i != 0 {
			dashaEnd = dashaStart.Add(
				AntardashaDuration(mahaDashaLord, antarDashaLord),
			)
		}
		d := Dasha{
			Mahadasha:  mahaDashaLord,
			Antardasha: antarDashaLord,
			Interval:   NewInterval(dashaStart, dashaEnd),
		}
		// fmt.Printf("Adding dasha: %+v\n", d)
		t.Set(d)
		mahaDashaLord, antarDashaLord = NextAntardasha(
			mahaDashaLord,
			antarDashaLord,
		)
		dashaStart = dashaEnd
	}
	return &DashaTree{
		tree:      t,
		startTime: birthTime,
	}, nil
}

// GetDashaForTime returns the dasha lord and sub dasha lord for the given time
func (dt *DashaTree) GetDashaForTime(
	t time.Time,
) (ret Dasha, didFind bool) {
	dt.tree.Scan(func(d Dasha) bool {
		if d.Interval.IsInBetween(t) {
			ret = d
			didFind = true
			return false
		}
		return true
	})
	if !didFind {
		return Dasha{}, false
	}
	return ret, true
}

func (dt *DashaTree) Dump() {
	for i, d := range dt.tree.Items() {
		fmt.Printf("%d: %+v\n", i, d)
	}
}

// calcMoonMovementInJulian calculates how long does the moon need (in Julian
// days) to move 'deg' degrees
func calcMoonMovementInJulian(deg float64) float64 {
	// 360 degrees in a circle
	// 27.321661 days in a sidereal month
	// 360 -> 27.321661
	// deg -> x
	// x = deg * 27.321661 / 360
	return deg * SiderealMoonMonthDelta / 360.0
}
