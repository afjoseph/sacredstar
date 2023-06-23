package chart

import (
	"fmt"

	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/afjoseph/sacredstar/zodiacalpos"
)

// correctToDegree finds the exact time a 'pid' was in 'zodPos' and corrects
// the chart to that time
func correctJulianDateToZodiacalPos(
	swe *wrapper.SwissEph,
	initialJulDay float64,
	pid pointid.PointID,
	targetZodPos *zodiacalpos.ZodiacalPos,
	targetNakshatraType NakshatraType,
	step float64,
	_range float64,
) (float64, error) {
	foo := func(id string, jd float64) (ret float64, didFind bool, zodPos *zodiacalpos.ZodiacalPos, err error) {
		ap, err := calculatePlanet(
			swe,
			jd,
			pointid.Moon,
			D1ChartType,
			// ascendantZodiacalPos: we don't need to provide it since we're
			// not calculating houses here
			nil,
		)
		if err != nil {
			return 0, false, nil, fmt.Errorf(
				"while casting chart at timeInJulian %f: %v",
				jd,
				err,
			)
		}
		diff := targetZodPos.DiffInAbsDegrees(ap.ZodiacalPos)
		// nak, err := NewNakshatraFromChart(swe, ap)
		// if err != nil {
		// 	return 0, false, ap.ZodiacalPos, fmt.Errorf(
		// 		"while calculating Nakshatra: %v", err)
		// }
		// fmt.Printf(
		// 	"id: %s | initalJulDay: %f | targetZodPos: %+v | targetZodPos (abs): %v | targetNak: %s | jd: %f | zodPos: %+v | zodPos (abs): %v | diff: %v | Nak: %s\n",
		// 	id,
		// 	initialJulDay,
		// 	targetZodPos,
		// 	targetZodPos.AbsDegrees(),
		// 	targetNakshatraType,
		// 	jd,
		// 	ap.ZodiacalPos,
		// 	ap.ZodiacalPos.AbsDegrees(),
		// 	diff,
		// 	nak.Type,
		// )
		if diff > 0.00000001 {
			// Didn't find it
			return 0, false, ap.ZodiacalPos, nil
		}
		// if nak.Type != targetNakshatraType {
		// 	// Try the next one.
		// 	// fmt.Printf(
		// 	// 	"Nakshatra mismatch: %s != %s. Will try the next one\n",
		// 	// 	nak.Type,
		// 	// 	targetNakshatraType,
		// 	// )
		// 	return 0, false, zodPos, nil
		// }
		// fmt.Printf(
		// 	"Found a good correction: jd: %f | zodPos: %+v\n",
		// 	jd,
		// 	zodPos,
		// )
		return jd, true, ap.ZodiacalPos, nil
	}

	// We'll need to use a binary search algorithm to find the exact time
	// in the most efficient way (i.e., O(log n))
	low := initialJulDay - _range
	high := initialJulDay + _range
	var mid float64
	for low <= high {
		// foo("low", low)
		// foo("high", high)
		mid = (low + high) / 2.0
		// fmt.Printf("low = %+v, high = %+v, mid = %+v\n", low, high, mid)
		ret, didFind, zodPos, err := foo("mid", mid)
		if err != nil {
			return 0, err
		}
		if didFind {
			return ret, nil
		}

		if zodPos.LessThan(targetZodPos) {
			low = mid + step
		} else {
			high = mid - step
		}
		// fmt.Println("---------------")
		// fmt.Println("---------------")
	}
	// If we failed to find the correct one, just return the closest one we could find
	// fmt.Printf(
	// 	"Failed to find the exact time. Returning the closest one: %f\n",
	// 	mid,
	// )
	return mid, nil
}
