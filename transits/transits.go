package transits

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/go-playground/errors/v5"
)

type TransitType int

const (
	TransitTypeNone TransitType = iota
	TransitTypeAspect
	TransitTypeIngress
)

type Transit interface {
	fmt.Stringer
	Type() TransitType
	GetJourney() float64
	GetDuration() int
	GetStart() time.Time
	GetEnd() time.Time
}

func New(
	swe *wrapper.SwissEph,
	t time.Time,
) ([]Transit, error) {
	tt, err := calculate(swe, t)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"while calculating transits for %s",
			t,
		)
	}
	return tt, err
}

func calculate(swe *wrapper.SwissEph, t time.Time) ([]Transit, error) {
	// Cast a chart
	chrt, err := chart.NewChartFromJulianDay(
		swe,
		swe.GoTimeToJulianDay(t),
		0, 0, // lon, lat: we're assuming UTC for now
		chart.TropicalChartType,
		pointid.ModernPlanets,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "while calculating chart for %s", t)
	}

	transits := []Transit{}
	// Calculate aspect journeys
	for _, asp := range chrt.Aspects {
		if asp.P1 == pointid.ASC || asp.P2 == pointid.ASC {
			// Skip ascendant calculations
			continue
		}

		fmt.Printf("Calculating aspect journey for %s\n", asp)
		duration, journey, start, end, err := calculateAspectJourney(
			swe,
			asp,
			t,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "while calculating journey for %s", t)
		}
		durationInDays := int(duration.Hours() / 24)
		transits = append(transits, &TransitAspect{
			Date:        t,
			Aspect:      asp,
			Journey:     journey,
			DaysElapsed: durationInDays,
			Start:       start,
			End:         end,
		})
	}

	// Calculate ingress journeys
	for _, p := range chrt.Points {
		if p.ID == pointid.ASC {
			// Skip ascendant calculations
			continue
		}
		fmt.Printf("Calculating ingress journey for %s\n", p)
		duration, journey, start, end, err := calculateIngressJourney(swe, p, t)
		if err != nil {
			return nil, errors.Wrapf(err, "while calculating journey for %s", t)
		}
		durationInDays := int(duration.Hours() / 24)
		transits = append(transits, &TransitIngress{
			Date:        t,
			P:           p,
			Journey:     journey,
			DaysElapsed: durationInDays,
			Start:       start,
			End:         end,
		})
	}
	return transits, nil
}

// func calculateAspectJourney(
// 	swe *wrapper.SwissEph,
// 	targetAspect *aspect.Aspect,
// 	targetAspectTime time.Time,
// ) (duration time.Duration, journey float64, start time.Time, end time.Time, err error) {
// 	if targetAspect.Type == aspect.AspectType_None {
// 		return 0, 0, time.Time{}, time.Time{}, nil
// 	}

// 	calculateAspect := func(t time.Time) *aspect.Aspect {
// 		chrt, err := chart.NewChartFromJulianDay(
// 			swe,
// 			swe.GoTimeToJulianDay(t),
// 			0, 0,
// 			chart.TropicalChartType,
// 			[]pointid.PointID{targetAspect.P1, targetAspect.P2},
// 		)
// 		if err != nil {
// 			panic(errors.Wrapf(
// 				err,
// 				"calculating chart for %s [chart: %s]",
// 				t, chrt,
// 			))
// 		}
// 		p1 := chrt.GetPoint(targetAspect.P1)
// 		p2 := chrt.GetPoint(targetAspect.P2)
// 		return p1.GetAspect(p2)
// 		// asp := p1.GetAspect(p2)
// 		// For example, we get p1=100 and p2=173 and the target aspect is
// 		// sextile (60 degreees) so the difference would be 73 We have to
// 		// reduce that by 60 to get the degree difference, which yields 13 and
// 		// that's what we return

// 		// Minus the targetAspect
// 		// aspectType := targetAspect.Type
// 		// deg := asp.Degree - float64(aspectType.Degree())
// 		// slog.Info(
// 		// 	"calculateAspectDegree",
// 		// 	slog.String("p1", p1.String()),
// 		// 	slog.String("p2", p2.String()),
// 		// 	slog.Time("t", t),
// 		// 	slog.String("asp", asp.String()),
// 		// )
// 		// return deg
// 	}

// 	aspectValidator := func(t time.Time, isBigStep bool) bool {
// 		asp := calculateAspect(t)
// 		// So if we're working with 64 degrees, and the target aspect is a
// 		// sextile, the normalizedAspectDegree would be 60
// 		normalizedAspectDegree := asp.Degree - float64(
// 			targetAspect.Type.Degree(),
// 		)
// 		var orbSize float64
// 		if isBigStep {
// 			orbSize = float64(targetAspect.Orb() * 2)
// 		} else {
// 			orbSize = float64(targetAspect.Orb())
// 		}
// 		isAspectWithinOrb := math.Abs(normalizedAspectDegree) < orbSize
// 		slog.Info(
// 			"aspectvalidator",
// 			slog.Time("t", t),
// 			slog.String("asp", asp.String()),
// 			slog.Float64("normalizedAspectDegree", normalizedAspectDegree),
// 			slog.Float64("orbSize", orbSize),
// 			slog.Bool("isAspectWithinOrb", isAspectWithinOrb),
// 		)
// 		return isAspectWithinOrb
// 	}
// 	aspectComparer := func(direction time.Duration, t time.Time) float64 {
// 		asp := calculateAspect(t)
// 		targetAspectDegree := float64(targetAspect.Type.Degree())
// 		normalizedAspectDegree := asp.Degree - targetAspectDegree
// 		if targetAspect.Type == aspect.AspectType_Conjunction && direction > 0 {
// 			// This is a very special case, but it's required since
// 			// conjunctions are essentially 0 degrees without a sign
// 			normalizedAspectDegree = -normalizedAspectDegree
// 		}
// 		orbSize := float64(targetAspect.Orb())
// 		// How far from the edge of the orb are we?
// 		// So if we're working with 64 degrees, and the target aspect is a
// 		// sextile, the normalizedAspectDegree would be 64-60=4
// 		// and the edgeDelta would be 4-3=1 (for the L-edge)
// 		// similarly, if we're at 43 degrees, and the target aspect is a
// 		// sextile, the normalizedAspectDegree would be 43-60=-17
// 		// and the edgeDelta would be -17+3=-14 (for the R-edge)
// 		var edgeDelta float64
// 		if normalizedAspectDegree > 0 {
// 			edgeDelta = normalizedAspectDegree - orbSize
// 		} else if normalizedAspectDegree < 0 {
// 			edgeDelta = normalizedAspectDegree + orbSize
// 		} else {
// 			edgeDelta = 0
// 		}
// 		slog.Info(
// 			"aspectcomparer",
// 			slog.Time("t", t),
// 			slog.String("asp", asp.String()),
// 			slog.Float64("targetAspectDegree", targetAspectDegree),
// 			slog.Float64("normalizedAspectDegree", normalizedAspectDegree),
// 			slog.Float64("orbSize", orbSize),
// 			slog.Duration("direction", direction),
// 			slog.Float64("edgeDelta", edgeDelta),
// 		)
// 		return edgeDelta
// 	}
// 	findEdge := func(direction time.Duration, step time.Duration, startingPoint time.Time) time.Time {
// 		edgeBeginAt, edgeEndAt := bigStepSearch(
// 			direction,
// 			step,
// 			startingPoint,
// 			aspectValidator,
// 			true,
// 		)
// 		// If the difference between the edgeBeginAt and edgeEndAt is more than
// 		// a year, step by weeks
// 		if edgeEndAt.Sub(edgeBeginAt) > time.Hour*24*365 {
// 			edgeEndAt, edgeBeginAt = bigStepSearch(
// 				direction,
// 				time.Hour*24*14,
// 				edgeEndAt,
// 				aspectValidator,
// 				false,
// 			)
// 		}
// 		return smallStepSearch(
// 			direction,
// 			edgeBeginAt,
// 			edgeEndAt,
// 			aspectComparer,
// 			0.5,
// 			false,
// 		)
// 	}

// 	// Move at the slowest speed for the search
// 	speed := getSpeedRate(targetAspect.P1)
// 	speed2 := getSpeedRate(targetAspect.P2)
// 	if speed2 < speed {
// 		speed = speed2
// 	}
// 	start = findEdge(time.Duration(-1), speed, targetAspectTime)
// 	end = findEdge(time.Duration(1), speed, targetAspectTime)
// 	duration = end.Sub(start)
// 	journey = float64(targetAspectTime.Sub(start)) / float64(duration)
// 	slog.Info("Aspect journey",
// 		slog.Duration("speed", speed),
// 		slog.String("targetAspect", targetAspect.String()),
// 		slog.Time("targetAspectTime", targetAspectTime),
// 		slog.Time("start", start),
// 		slog.Time("end", end),
// 		slog.Duration("duration", duration),
// 		slog.Float64("journey", journey),
// 	)
// 	return duration, journey, start, end, nil
// }

// func calculateIngressJourney(
// 	swe *wrapper.SwissEph,
// 	targetPoint *astropoint.AstroPoint,
// 	targetTime time.Time,
// ) (duration time.Duration, journey float64, start time.Time, end time.Time, err error) {
// 	calculateZodiacalPos := func(t time.Time) *zodiacalpos.ZodiacalPos {
// 		chrt, err := chart.NewChartFromJulianDay(
// 			swe,
// 			swe.GoTimeToJulianDay(t),
// 			0, 0,
// 			chart.TropicalChartType,
// 			[]pointid.PointID{targetPoint.ID},
// 		)
// 		if err != nil {
// 			panic(errors.Wrapf(
// 				err,
// 				"calculating chart for %s",
// 				t,
// 			))
// 		}
// 		return chrt.GetPoint(targetPoint.ID).ZodiacalPos
// 	}

// 	findEdge := func(direction time.Duration, step time.Duration, startingPoint time.Time) time.Time {
// 		var s sign.Sign
// 		if direction < 0 {
// 			s = targetPoint.ZodiacalPos.Sign
// 		} else {
// 			s = targetPoint.ZodiacalPos.Sign.Next()
// 		}
// 		targetZP := zodiacalpos.NewZodiacalPos(s, 0, 0)

// 		ingressValidator := func(t time.Time, isBigStep bool) bool {
// 			currZP := calculateZodiacalPos(t)
// 			return currZP.Sign == targetZP.Sign
// 		}
// 		ingressComparer := func(direction time.Duration, t time.Time) float64 {
// 			currZP := calculateZodiacalPos(t)
// 			return currZP.DiffInDirectionalDegrees(targetZP)
// 		}
// 		edgeBeginAt, edgeEndAt := bigStepSearch(
// 			direction,
// 			step,
// 			startingPoint,
// 			ingressValidator,
// 			true,
// 		)
// 		return smallStepSearch(
// 			direction,
// 			edgeBeginAt,
// 			edgeEndAt,
// 			ingressComparer,
// 			0.1,
// 			false,
// 		)
// 	}

// 	speed := getSpeedRate(targetPoint.ID)
// 	start = findEdge(time.Duration(-1), speed, targetTime)
// 	end = findEdge(time.Duration(1), speed, targetTime)
// 	duration = end.Sub(start)
// 	journey = float64(targetTime.Sub(start)) / float64(duration)
// 	slog.Info("Ingress journey",
// 		slog.Duration("speed", speed),
// 		slog.String("targetPoint", targetPoint.String()),
// 		slog.Time("targetTime", targetTime),
// 		slog.Time("start", start),
// 		slog.Time("end", end),
// 		slog.Duration("duration", duration),
// 		slog.Float64("journey", journey),
// 	)
// 	return duration, journey, start, end, nil
// }

func getSpeedRate(p pointid.PointID) time.Duration {
	switch p {
	case pointid.Sun:
		return 30 * 24 * time.Hour
	case pointid.Moon:
		return 7 * 24 * time.Hour
	case pointid.Mercury:
		return 40 * 24 * time.Hour
	case pointid.Venus:
		return 65 * 24 * time.Hour
	case pointid.Mars:
		return 90 * 24 * time.Hour
	case pointid.Jupiter:
		return 5 * 30 * 24 * time.Hour
	case pointid.Saturn:
		return 5 * 365 * 24 * time.Hour
	case pointid.Uranus:
		return 5 * 365 * 24 * time.Hour
	case pointid.Neptune:
		return 5 * 365 * 24 * time.Hour
	case pointid.Pluto:
		return 34 * 365 * 24 * time.Hour
	default:
		panic(errors.Newf("Unknown planet: %s", p))
	}
}

// SearchValidator is a function type that determines if a given time point is
// valid for the search criteria
type SearchValidator func(time.Time, bool) bool

// SearchComparer is a function type that compares two time points and returns
// a float64 indicating their relative position/difference
type SearchComparer func(time.Duration, time.Time) float64

// bigStepSearch performs a broad search to find approximate boundaries
func bigStepSearch(
	direction time.Duration,
	step time.Duration,
	startTime time.Time,
	validator SearchValidator,
	isBigStep bool,
) (time.Time, time.Time) {
	currentTime := startTime
	for {
		nextTime := currentTime.Add(step * direction)
		isValid := validator(nextTime, isBigStep)
		slog.Info("bigstepsearch",
			slog.Time("currentTime", currentTime),
			slog.Time("nextTime", nextTime),
			slog.Duration("step", step),
			slog.Duration("direction", direction),
			slog.Bool("isValid", isValid),
		)
		if !isValid {
			return nextTime, currentTime
		}
		currentTime = nextTime
	}
}

func smallStepSearch(
	direction time.Duration,
	start time.Time,
	end time.Time,
	comparer SearchComparer,
	epsilon float64,
	accountForDirection bool,
) time.Time {
	var mid time.Time
	L := start
	R := end
	if direction < 0 {
		L, R = R, L
	}

	for L.Before(R) {
		mid = L.Add(R.Sub(L) / 2)
		diff := comparer(direction, mid)
		isWithinEpsilon := math.Abs(diff) < epsilon
		slog.Info("smallstepsearch",
			slog.Time("L", L),
			slog.Time("R", R),
			slog.Time("mid", mid),
			slog.Duration("direction", direction),
			slog.Float64("diff", diff),
		)
		if isWithinEpsilon {
			return mid
		}

		if diff < 0 {
			if accountForDirection {
				if direction < 0 {
					R = mid
				} else {
					L = mid
				}
			} else {
				R = mid
			}
		} else {
			if accountForDirection {
				if direction < 0 {
					L = mid
				} else {
					R = mid
				}
			} else {
				L = mid
			}
		}
	}
	return mid
}
