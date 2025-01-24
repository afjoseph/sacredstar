package transits

import (
	"fmt"
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
	L, R := aspectBigStepSearch(
		targetAspect,
		targetAspectTime,
		time.Duration(-1),
		calculateAspect,
	)
	startingEdge, _ := aspectBinSearch(
		L,
		R,
		targetAspect,
		targetAspectTime,
		calculateAspect,
	)

	// Calculate ending edge
	L, R = aspectBigStepSearch(
		targetAspect,
		targetAspectTime,
		time.Duration(1),
		calculateAspect,
	)
	endingEdge, _ := aspectBinSearch(
		L,
		R,
		targetAspect,
		targetAspectTime,
		calculateAspect,
	)

	duration = endingEdge.Sub(startingEdge)
	journey = float64(targetAspectTime.Sub(startingEdge)) / float64(duration)

	// slog.Info("Aspect journey",
	// 	slog.String("targetAspect", targetAspect.String()),
	// 	slog.Time("targetAspectTime", targetAspectTime),
	// 	slog.String("startingAspect", startingAspect.String()),
	// 	slog.Time("startingEdge", startingEdge),
	// 	slog.String("endingAspect", endingAspect.String()),
	// 	slog.Time("endingEdge", endingEdge),
	// 	slog.Duration("duration", duration),
	// 	slog.Float64("journey", journey),
	// )
	return duration, journey, startingEdge, endingEdge, nil
}

func aspectBigStepSearch(
	targetAspect *aspect.Aspect,
	targetAspectTime time.Time,
	direction time.Duration,
	calculateAspect func(t time.Time) *aspect.Aspect,
) (time.Time, time.Time) {
	var step time.Duration
	p1Step := getStepForPointID(targetAspect.P1)
	p2Step := getStepForPointID(targetAspect.P2)
	if p1Step <= p2Step {
		step = p1Step
	} else {
		step = p2Step
	}
	currentTime := targetAspectTime
	LStepWithinOrb := time.Time{}
	RStepWithinOrb := time.Time{}
	doubleOrb := float64(targetAspect.Orb() * 2)
	orb := float64(targetAspect.Orb())

	var asp *aspect.Aspect
	var nextTime time.Time
	for {
		nextTime = currentTime.Add(step * direction)
		asp = calculateAspect(nextTime)
		normalizedAspectDegree := asp.Degree - float64(
			targetAspect.Type.Degree(),
		)
		// - case 1: within the orb
		if math.Abs(normalizedAspectDegree) <= orb {
			if direction < 0 {
				LStepWithinOrb = nextTime
				RStepWithinOrb = currentTime
			} else {
				LStepWithinOrb = currentTime
				RStepWithinOrb = nextTime
			}
			// slog.Info("case 1",
			// 	slog.Float64("normalizedAspectDegree", normalizedAspectDegree),
			// 	slog.String("aspect", asp.String()),
			// 	slog.Float64("orb", orb),
			// 	slog.Time("currentTime", currentTime),
			// 	slog.Time("prevTime", prevTime),
			// 	slog.Time("LStepWithinOrb", LStepWithinOrb),
			// 	slog.Time("RStepWithinOrb", RStepWithinOrb),
			// )
			currentTime = nextTime
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
			currentTime = nextTime
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

	if LStepWithinOrb.IsZero() || RStepWithinOrb.IsZero() {
		if direction < 0 {
			LStepWithinOrb = targetAspectTime.Add(-step)
			RStepWithinOrb = targetAspectTime
		} else {
			LStepWithinOrb = targetAspectTime
			RStepWithinOrb = targetAspectTime.Add(step)
		}
	} else {
		LStepWithinOrb = LStepWithinOrb.Add(step * direction)
		RStepWithinOrb = RStepWithinOrb.Add(step * direction)
	}
	return LStepWithinOrb, RStepWithinOrb
}

func aspectBinSearch(
	L time.Time,
	R time.Time,
	targetAspect *aspect.Aspect,
	targetAspectTime time.Time,
	calculateAspect func(t time.Time) *aspect.Aspect,
) (time.Time, *aspect.Aspect) {
	orb := float64(targetAspect.Orb())
	aspectTypeDegree := float64(targetAspect.Type.Degree())
	peek := func(L time.Time, R time.Time) (time.Time, *aspect.Aspect, float64) {
		mid := L.Add(R.Sub(L) / 2)
		asp := calculateAspect(mid)
		normalizedAspectDegree := asp.Degree - aspectTypeDegree
		x := math.Abs(math.Abs(normalizedAspectDegree) - orb)
		// slog.Info(
		// 	"endingedgebinsearch: peek",
		// 	slog.Time("L", L),
		// 	slog.Time("R", R),
		// 	slog.Time("mid", mid),
		// 	slog.String("asp", asp.String()),
		// 	slog.Float64("normalizedAspectDegree", normalizedAspectDegree),
		// 	slog.Float64("x", x),
		// 	slog.Float64("orb", orb),
		// 	slog.Float64("aspectTypeDegree", aspectTypeDegree),
		// )
		return mid, asp, x
	}

	for L.Before(R) {
		diff := R.Sub(L)
		if diff < 1*time.Hour {
			return time.Time{}, nil
		}

		mid, asp, x := peek(L, R)
		if x < 0.5 {
			return mid, asp
		}

		_, _, x1 := peek(L, mid)
		_, _, x2 := peek(mid, R)
		if x1 < x2 {
			R = mid
		} else {
			L = mid
		}
	}
	panic(errors.New("could not find ending edge"))
}
