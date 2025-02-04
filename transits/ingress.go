package transits

import (
	"fmt"
	"math"
	"time"

	"github.com/afjoseph/sacredstar/astropoint"
	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/unixtime"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/afjoseph/sacredstar/zodiacalpos"
	"github.com/go-playground/errors/v5"
)

type TransitIngress struct {
	transitBase
	P *astropoint.AstroPoint `json:"point"`
}

func (t *TransitIngress) String() string {
	return fmt.Sprintf(
		"TransitIngress{Date: %s, P: %s, Journey: %.2f, DaysElapsed: %d, Start: %s, End: %s}",
		t.transitBase.Date.Format("2006-01-02"),
		t.P,
		t.transitBase.Journey,
		t.transitBase.DaysElapsed,
		t.transitBase.Start.Format("2006-01-02"),
		t.transitBase.End.Format("2006-01-02"),
	)
}

func (t *TransitIngress) GetType() TransitType {
	return TransitTypeIngress
}

func (t *TransitIngress) GetJourney() float64 {
	return t.transitBase.Journey
}

func (t *TransitIngress) GetDuration() int {
	return t.transitBase.DaysElapsed
}

func (t *TransitIngress) GetStart() unixtime.UnixTime {
	return t.transitBase.Start
}

func (t *TransitIngress) GetEnd() unixtime.UnixTime {
	return t.transitBase.End
}

func newTransitIngress(
	swe *wrapper.SwissEph,
	targetPoint *astropoint.AstroPoint,
	targetTime time.Time,
) (*TransitIngress, error) {
	duration, journey, start, end, err := calculateIngressJourney(
		swe,
		targetPoint,
		targetTime,
	)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"while calculating journey for %s",
			targetPoint,
		)
	}
	return &TransitIngress{
		transitBase: transitBase{
			Type:        TransitTypeIngress,
			Date:        unixtime.New(targetTime),
			Journey:     journey,
			DaysElapsed: int(duration.Hours() / 24),
			Start:       unixtime.New(start),
			End:         unixtime.New(end),
		},
		P: targetPoint,
	}, nil
}

func calculateIngressJourney(
	swe *wrapper.SwissEph,
	targetPoint *astropoint.AstroPoint,
	targetTime time.Time,
) (duration time.Duration, journey float64, start time.Time, end time.Time, err error) {
	calculateZodiacalPos := func(t time.Time) *zodiacalpos.ZodiacalPos {
		chrt, err := chart.NewChartFromJulianDay(
			swe,
			swe.GoTimeToJulianDay(t),
			0, 0,
			chart.TropicalChartType,
			[]pointid.PointID{targetPoint.ID},
		)
		if err != nil {
			panic(errors.Wrapf(
				err,
				"calculating chart for %s",
				t,
			))
		}
		return chrt.GetPoint(targetPoint.ID).ZodiacalPos
	}

	// Calculate starting edge
	step := getStepForPointID(targetPoint.ID)
	var currZP *zodiacalpos.ZodiacalPos
	targetZP := zodiacalpos.NewZodiacalPos(targetPoint.ZodiacalPos.Sign, 0, 0)
	var nextTime time.Time
	currentTime := targetTime
	for {
		nextTime = currentTime.Add(-step)
		currZP = calculateZodiacalPos(nextTime)
		if currZP.Sign == targetZP.Sign.Previous() {
			break
		}
		currentTime = nextTime
	}

	// Do binary search to find the exact edge
	var mid time.Time
	L := nextTime
	R := currentTime
	didFind := false
	for L.Before(R) {
		mid = L.Add(R.Sub(L) / 2)
		currZP = calculateZodiacalPos(mid)
		diff := currZP.DiffInDirectionalDegrees(targetZP)
		isWithinEpsilon := math.Abs(diff) < 0.1
		if isWithinEpsilon {
			didFind = true
			break
		}
		if diff < 0 {
			R = mid
		} else {
			L = mid
		}
	}
	if !didFind {
		panic(errors.New("could not find ingress edge"))
	}
	// startingZP := *currZP
	startingEdge := mid

	// Calculate ending edge
	targetZP = zodiacalpos.NewZodiacalPos(
		targetPoint.ZodiacalPos.Sign.Next(),
		0,
		0,
	)
	currentTime = targetTime
	for {
		nextTime = currentTime.Add(step)
		currZP := calculateZodiacalPos(nextTime)
		if currZP.Sign == targetZP.Sign {
			break
		}
		currentTime = nextTime
	}
	// Do binary search to find the exact edge
	L = currentTime
	R = nextTime
	didFind = false
	for L.Before(R) {
		mid = L.Add(R.Sub(L) / 2)
		currZP = calculateZodiacalPos(mid)
		diff := currZP.DiffInDirectionalDegrees(targetZP)
		isWithinEpsilon := math.Abs(diff) < 0.1
		if isWithinEpsilon {
			didFind = true
			break
		}
		if diff < 0 {
			R = mid
		} else {
			L = mid
		}
	}
	if !didFind {
		panic(errors.New("could not find egress edge"))
	}
	// endingZP := *currZP
	endingEdge := mid
	// start = findEdge(time.Duration(-1), speed, targetTime)
	// end = findEdge(time.Duration(1), speed, targetTime)
	duration = endingEdge.Sub(startingEdge)
	journey = float64(targetTime.Sub(startingEdge)) / float64(duration)
	// slog.Info("Ingress journey",
	// 	slog.Duration("step", step),
	// 	slog.String("targetPoint", targetPoint.String()),
	// 	slog.Time("targetTime", targetTime),
	// 	slog.String("startingZP", startingZP.String()),
	// 	slog.Time("startingEdge", startingEdge),
	// 	slog.String("endingZP", endingZP.String()),
	// 	slog.Time("endingEdge", endingEdge),
	// 	slog.Duration("duration", duration),
	// 	slog.Float64("journey", journey),
	// )
	return duration, journey, startingEdge, endingEdge, nil
}
