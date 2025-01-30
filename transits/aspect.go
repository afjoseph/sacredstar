package transits

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/afjoseph/sacredstar/aspect"
	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/go-playground/errors/v5"
)

type TransitAspect struct {
	Date        time.Time      `json:"date"`
	Aspect      *aspect.Aspect `json:"aspect"`
	Journey     float64        `json:"journey"`
	DaysElapsed int            `json:"daysElapsed"`
	Start       time.Time      `json:"start"`
	End         time.Time      `json:"end"`
}

func (t *TransitAspect) String() string {
	return fmt.Sprintf(
		"TransitAspect{Date: %s, Aspect: %s, Journey: %.2f, DaysElapsed: %d, Start: %s, End: %s}",
		t.Date,
		t.Aspect,
		t.Journey,
		t.DaysElapsed,
		t.Start.Format("2006-01-02"),
		t.End.Format("2006-01-02"),
	)
}

func (t *TransitAspect) Type() TransitType {
	return TransitTypeAspect
}

func (t *TransitAspect) GetJourney() float64 {
	return t.Journey
}

func (t *TransitAspect) GetDuration() int {
	return t.DaysElapsed
}

func (t *TransitAspect) GetStart() time.Time {
	return t.Start
}

func (t *TransitAspect) GetEnd() time.Time {
	return t.End
}

func calculateAspectJourney(
	swe *wrapper.SwissEph,
	targetAspect *aspect.Aspect,
	targetAspectTime time.Time,
) (duration time.Duration, journey float64, start time.Time, end time.Time, err error) {
	if targetAspect.Type == aspect.AspectType_None {
		return 0, 0, time.Time{}, time.Time{}, nil
	}

	calculateAspect := func(t time.Time) *aspect.Aspect {
		chrt, err := chart.NewChartFromJulianDay(
			swe,
			swe.GoTimeToJulianDay(t),
			0, 0,
			chart.TropicalChartType,
			[]pointid.PointID{targetAspect.P1, targetAspect.P2},
		)
		if err != nil {
			panic(errors.Wrapf(
				err,
				"calculating chart for %s [chart: %s]",
				t, chrt,
			))
		}
		p1 := chrt.GetPoint(targetAspect.P1)
		p2 := chrt.GetPoint(targetAspect.P2)
		return p1.GetAspect(p2)
	}

	// Calculate starting edge
	step := 14 * 24 * time.Hour
	currentTime := targetAspectTime
	var asp *aspect.Aspect
	LStepWithinOrb := time.Time{}
	RStepWithinOrb := time.Time{}
	doubleOrb := float64(targetAspect.Orb() * 2)
	orb := float64(targetAspect.Orb())
	// ingressDegreeDoubleOrb := targetAspect.Type.Degree() - doubleOrb
	// ingressDegreeOrb := targetAspect.Type.Degree() - orb
	var prevTime time.Time
	for {
		prevTime = currentTime.Add(-step)
		asp = calculateAspect(prevTime)
		// So if we're working with 64 degrees, and the target aspect is a
		// sextile, the normalizedAspectDegree would be 60
		normalizedAspectDegree := asp.Degree - float64(
			targetAspect.Type.Degree(),
		)
		// - case 1: within the orb
		if math.Abs(normalizedAspectDegree) <= orb {
			// Mark the steps and keep going
			LStepWithinOrb = prevTime
			RStepWithinOrb = currentTime
			// slog.Info("case 1",
			// 	slog.Float64("normalizedAspectDegree", normalizedAspectDegree),
			// 	slog.String("aspect", asp.String()),
			// 	slog.Float64("orb", orb),
			// 	slog.Time("currentTime", currentTime),
			// 	slog.Time("prevTime", prevTime),
			// 	slog.Time("LStepWithinOrb", LStepWithinOrb),
			// 	slog.Time("RStepWithinOrb", RStepWithinOrb),
			// )
			currentTime = prevTime
			continue
		}
		// - case 2: outside orb, within double-orb
		if math.Abs(normalizedAspectDegree) < doubleOrb {
			// slog.Info("case 2",
			// 	slog.Float64("normalizedAspectDegree", normalizedAspectDegree),
			// 	slog.String("aspect", asp.String()),
			// 	slog.Float64("orb", orb),
			// 	slog.Time("currentTime", currentTime),
			// 	slog.Time("prevTime", prevTime),
			// 	slog.Time("LStepWithinOrb", LStepWithinOrb),
			// 	slog.Time("RStepWithinOrb", RStepWithinOrb),
			// )
			// Just keep going
			currentTime = prevTime
			continue
		}
		// case 3: outside double-orb
		// we should break here
		// slog.Info("case 3",
		// 	slog.Float64("normalizedAspectDegree", normalizedAspectDegree),
		// 	slog.String("aspect", asp.String()),
		// 	slog.Float64("orb", orb),
		// 	slog.Time("currentTime", currentTime),
		// 	slog.Time("prevTime", prevTime),
		// 	slog.Time("LStepWithinOrb", LStepWithinOrb),
		// 	slog.Time("RStepWithinOrb", RStepWithinOrb),
		// )
		break
	}

	// Do binary search to find the exact edge
	if LStepWithinOrb.IsZero() || RStepWithinOrb.IsZero() {
		// This means the first step covered double the orb range
		LStepWithinOrb = targetAspectTime.Add(-step)
		RStepWithinOrb = targetAspectTime
	} else {
		LStepWithinOrb = LStepWithinOrb.Add(-step)
		RStepWithinOrb = RStepWithinOrb.Add(-step)
	}
	L := LStepWithinOrb
	R := RStepWithinOrb

	startingEdge, startingAspect := startingEdgeBinSearch(
		L,
		R,
		targetAspect,
		targetAspectTime,
		calculateAspect,
	)

	// Calculate ending edge
	currentTime = targetAspectTime
	LStepWithinOrb = time.Time{}
	RStepWithinOrb = time.Time{}
	nextTime := time.Time{}
	for {
		nextTime = currentTime.Add(step)
		asp = calculateAspect(nextTime)
		normalizedAspectDegree := asp.Degree - float64(
			targetAspect.Type.Degree(),
		)
		if math.Abs(normalizedAspectDegree) < orb {
			LStepWithinOrb = currentTime
			RStepWithinOrb = nextTime
			currentTime = nextTime
			continue
		}
		if math.Abs(normalizedAspectDegree) < doubleOrb {
			currentTime = nextTime
			continue
		}
		break
	}
	if LStepWithinOrb.IsZero() || RStepWithinOrb.IsZero() {
		// This means the first step covered double the orb range
		LStepWithinOrb = targetAspectTime
		RStepWithinOrb = targetAspectTime.Add(step)
	} else {
		LStepWithinOrb = LStepWithinOrb.Add(step)
		RStepWithinOrb = RStepWithinOrb.Add(step)
	}
	L = LStepWithinOrb
	R = RStepWithinOrb

	endingEdge, endingAspect := endingEdgeBinSearch(
		L,
		R,
		targetAspect,
		targetAspectTime,
		calculateAspect,
	)

	// didFind := false
	// edge := math.Abs(float64(targetAspect.Type.Degree()) - orb)

	// Walk from L to R in daily steps to see when we find the edge
	// currentTime = L
	// for L.Before(R) {
	// 	nextTime = currentTime.Add(24 * time.Hour)
	// 	asp = calculateAspect(nextTime)
	// 	edge := math.Abs(float64(targetAspect.Type.Degree()) - orb)
	// 	if asp.Degree > edge {
	// 		// We walked too far: just return the current time
	// 		didFind = true
	// 		break
	// 	}
	// 	// if math.Abs(asp.Degree-edge) < 1.0 {
	// 	// 	didFind = true
	// 	// 	break
	// 	// }
	// 	currentTime = nextTime
	// }
	// if !didFind {
	// 	panic(errors.New("could not find egress edge"))
	// }
	// endingAspect := *asp
	// endingEdge := nextTime

	// for L.Before(R) {
	// 	mid = L.Add(R.Sub(L) / 2)
	// 	asp = calculateAspect(mid)
	// 	// edge := math.Abs(targetAspect.Type.Degree() + orb)
	// 	// normalizedAspectDegree := asp.Degree - float64(
	// 	// 	targetAspect.Type.Degree(),
	// 	// )
	// 	edge := math.Abs((float64(targetAspect.Type.Degree()) - orb))
	// 	normalizedAspectDiff := asp.Degree - edge
	// 	if targetAspect.Type == aspect.AspectType_Conjunction {
	// 		// This is a very special case, but it's required since
	// 		// conjunctions are essentially 0 degrees without a sign
	// 		normalizedAspectDiff = -normalizedAspectDiff
	// 	}
	// 	// var edgeDelta float64
	// 	// if normalizedAspectDegree > 0 {
	// 	// 	edgeDelta = math.Abs(normalizedAspectDegree - orb)
	// 	// 	// edgeDelta = -edgeDelta
	// 	// } else if normalizedAspectDegree < 0 {
	// 	// 	edgeDelta = normalizedAspectDegree + orb
	// 	// } else {
	// 	// 	edgeDelta = 0
	// 	// }
	// 	isWithinEpsilon := math.Abs(normalizedAspectDiff) < 0.1
	// 	// isWithinEpsilon := math.Abs(asp.Degree-edge) < 0.1
	// 	if isWithinEpsilon {
	// 		didFind = true
	// 		break
	// 	}
	// 	// if asp.Degree < edge {
	// 	// 	L = mid
	// 	// } else {
	// 	// 	R = mid
	// 	// }

	// 	if normalizedAspectDiff < 0 {
	// 		R = mid
	// 	} else {
	// 		L = mid
	// 	}
	// }
	// if !didFind {
	// 	panic(errors.New("could not find egress edge"))
	// }
	// endingAspect := *asp
	// endingEdge := mid

	duration = endingEdge.Sub(startingEdge)
	journey = float64(targetAspectTime.Sub(startingEdge)) / float64(duration)

	slog.Info("Aspect journey",
		slog.Duration("step", step),
		slog.String("targetAspect", targetAspect.String()),
		slog.Time("targetAspectTime", targetAspectTime),
		slog.String("startingAspect", startingAspect.String()),
		slog.Time("startingEdge", startingEdge),
		slog.String("endingAspect", endingAspect.String()),
		slog.Time("endingEdge", endingEdge),
		slog.Duration("duration", duration),
		slog.Float64("journey", journey),
	)
	return duration, journey, startingEdge, endingEdge, nil
}

func getDirection(t time.Time, targetAspect *aspect.Aspect,
	calculateAspect func(t time.Time) *aspect.Aspect,
) time.Duration {
	// Go forward 1 hour
	step := 24 * time.Hour
	nextTime := t.Add(step)
	asp := calculateAspect(nextTime)
	if asp.Degree > targetAspect.Degree {
		return time.Duration(1)
	}
	return time.Duration(-1)
}

func startingEdgeBinSearch(
	L time.Time,
	R time.Time,
	targetAspect *aspect.Aspect,
	targetAspectTime time.Time,
	calculateAspect func(t time.Time) *aspect.Aspect,
) (time.Time, *aspect.Aspect) {
	orb := float64(targetAspect.Orb())
	aspectTypeDegree := float64(targetAspect.Type.Degree())
	direction := getDirection(L, targetAspect, calculateAspect)

	// if targetAspect.Type == aspect.AspectType_Conjunction {
	// 	var asp *aspect.Aspect
	// 	for L.Before(R) {
	// 		mid := L.Add(R.Sub(L) / 2)
	// 		asp = calculateAspect(mid)
	// 		normalizedAspectDegree := asp.Degree - aspectTypeDegree
	// 		x := math.Abs(normalizedAspectDegree - orb)
	// 		// ingressEdge := orb + aspectTypeDegree
	// 		slog.Info(
	// 			"startingedgebinsearch",
	// 			slog.Duration("direction", direction),
	// 			slog.Time("L", L),
	// 			slog.Time("R", R),
	// 			slog.Time("mid", mid),
	// 			slog.String("asp", asp.String()),
	// 			slog.Float64("normalizedAspectDegree", normalizedAspectDegree),
	// 			slog.Float64("x", x),
	// 			slog.Float64("orb", orb),
	// 			slog.Float64("aspectTypeDegree", aspectTypeDegree),
	// 		)
	// 		if x < 1.0 {
	// 			return mid, asp
	// 		}
	// 		if normalizedAspectDegree < orb {
	// 			L = mid
	// 		} else {
	// 			R = mid
	// 		}
	// 	}
	// 	panic(errors.New("could not find starting edge"))
	// }

	do := func(direction time.Duration, L time.Time, R time.Time) (time.Time, *aspect.Aspect) {
		var asp *aspect.Aspect
		for L.Before(R) {
			diff := R.Sub(L)
			if diff < 1*time.Hour {
				return time.Time{}, nil
			}

			mid := L.Add(R.Sub(L) / 2)
			asp = calculateAspect(mid)
			normalizedAspectDegree := asp.Degree - aspectTypeDegree
			x := math.Abs(math.Abs(normalizedAspectDegree) - orb)
			// ingressEdge := orb + aspectTypeDegree
			slog.Info(
				"startingedgebinsearch",
				slog.Time("L", L),
				slog.Time("R", R),
				slog.Time("mid", mid),
				slog.Duration("direction", direction),
				slog.String("asp", asp.String()),
				slog.Float64("normalizedAspectDegree", normalizedAspectDegree),
				slog.Float64("x", x),
				slog.Float64("orb", orb),
				slog.Float64("aspectTypeDegree", aspectTypeDegree),
			)
			if x < 1.0 {
				return mid, asp
			}
			if normalizedAspectDegree < orb {
				if direction < 0 {
					R = mid
				} else {
					L = mid
				}
			} else {
				if direction < 0 {
					L = mid
				} else {
					R = mid
				}
			}

			// 			if direction < 0 {
			// 				R = mid
			// 			} else {
			// 				L = mid
			// 			}
		}
		return time.Time{}, nil
	}

	t, asp := do(direction, L, R)
	if t.IsZero() {
		t, asp = do(-direction, L, R)
	}
	if t.IsZero() {
		panic(errors.New("could not find starting edge"))
	}
	return t, asp
}

func endingEdgeBinSearch(
	L time.Time,
	R time.Time,
	targetAspect *aspect.Aspect,
	targetAspectTime time.Time,
	calculateAspect func(t time.Time) *aspect.Aspect,
) (time.Time, *aspect.Aspect) {
	orb := float64(targetAspect.Orb())
	aspectTypeDegree := float64(targetAspect.Type.Degree())

	var asp *aspect.Aspect
	for L.Before(R) {
		mid := L.Add(R.Sub(L) / 2)
		asp = calculateAspect(mid)
		normalizedAspectDegree := asp.Degree - aspectTypeDegree
		if targetAspect.Type == aspect.AspectType_Conjunction {
			// This is a very special case, but it's required since
			// conjunctions are essentially 0 degrees without a sign
			normalizedAspectDegree = -normalizedAspectDegree
		}
		x := math.Abs(normalizedAspectDegree) - orb
		slog.Info(
			"endingedgebinsearch",
			slog.Time("L", L),
			slog.Time("R", R),
			slog.Time("mid", mid),
			slog.String("asp", asp.String()),
			slog.Float64("normalizedAspectDegree", normalizedAspectDegree),
			slog.Float64("x", x),
			slog.Float64("orb", orb),
			slog.Float64("aspectTypeDegree", aspectTypeDegree),
		)
		if x < 1.0 {
			return mid, asp
		}
		if normalizedAspectDegree < -orb {
			R = mid
		} else {
			L = mid
		}
		// var edgeDelta float64
		// if normalizedAspectDegree > 0 {
		// 	edgeDelta = normalizedAspectDegree - float64(targetAspect.Orb())
		// } else if normalizedAspectDegree < 0 {
		// 	edgeDelta = math.Abs(normalizedAspectDegree + float64(targetAspect.Orb()))
		// } else {
		// 	edgeDelta = 0
		// }
		// isWithinEpsilon := math.Abs(edgeDelta) < 1.0
		// if isWithinEpsilon {
		// 	return mid, asp
		// }
		// if edgeDelta < 0 {
		// 	R = mid
		// } else {
		// 	L = mid
		// }
	}
	// if L.Before(R) {
	panic(errors.New("could not find starting edge"))
}
