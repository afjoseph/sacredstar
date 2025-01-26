package weatherreport

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/afjoseph/sacredstar/aspect"
	"github.com/afjoseph/sacredstar/astropoint"
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

type TransitIngress struct {
	Date        time.Time              `json:"date"`
	P           *astropoint.AstroPoint `json:"p"`
	Journey     float64                `json:"journey"`
	DaysElapsed int                    `json:"daysElapsed"`
	Start       time.Time              `json:"start"`
	End         time.Time              `json:"end"`
}

func (t *TransitIngress) String() string {
	return fmt.Sprintf(
		"TransitIngress{Date: %s, P: %s, Journey: %.2f, DaysElapsed: %d, Start: %s, End: %s}",
		t.Date,
		t.P,
		t.Journey,
		t.DaysElapsed,
		t.Start.Format("2006-01-02"),
		t.End.Format("2006-01-02"),
	)
}

func (t *TransitIngress) Type() TransitType {
	return TransitTypeIngress
}

func (t *TransitIngress) GetJourney() float64 {
	return t.Journey
}

func (t *TransitIngress) GetDuration() int {
	return t.DaysElapsed
}

func (t *TransitIngress) GetStart() time.Time {
	return t.Start
}

func (t *TransitIngress) GetEnd() time.Time {
	return t.End
}

func New(
	swe *wrapper.SwissEph,
	start time.Time,
	end time.Time,
) (map[string][]Transit, error) {
	// For each day from start to end, create a transit
	// and add it to the report
	var err error
	report := make(map[string][]Transit)
	for t := start; t.Before(end); t = t.Add(24 * time.Hour) {
		// fmt.Printf("Calculating transits for %s\n", t)
		report[t.Format("2006-01-02")], err = calculate(swe, t)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"while calculating transits for %s",
				t,
			)
		}
	}
	return report, nil
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

		// fmt.Printf("Calculating aspect journey for %s\n", asp)
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
		// fmt.Printf("Calculating ingress journey for %s\n", p)
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

func calculateAspectJourney(
	swe *wrapper.SwissEph,
	targetAspect *aspect.Aspect,
	targetAspectTime time.Time,
) (duration time.Duration, journey float64, start time.Time, end time.Time, err error) {
	if targetAspect.Type == aspect.AspectType_None {
		return 0, 0, time.Time{}, time.Time{}, nil
	}

	calculateAspectDegree := func(t time.Time) (float64, bool) {
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
		isRetrograde := p1.IsRetrograde || p2.IsRetrograde
		asp := p1.GetAspect(p2)
		// For example, we get p1=100 and p2=173 and the target aspect is sextile (60 degreees)
		// so the difference would be 73
		// We have to reduce that by 60 to get the degree difference, which yields 13 and that's what we return

		// Minus the targetAspect
		deg := asp.Degree - float64(targetAspect.Type.Degree())
		slog.Info("calculateAspectDegree",
			slog.String("p1", p1.String()),
			slog.String("p2", p2.String()),
			slog.Float64("deg", deg),
			slog.Time("t", t),
			slog.String("targetAspect", targetAspect.String()),
			slog.String("asp", asp.String()),
		)
		return deg, isRetrograde
		// if deg > float64(orb) {
		// 	return 1
		// }
		// if deg < float64(orb) {
		// 	return -1
		// }
		// if math.Abs(deg-float64(orb)) < epsilon {
		// 	return 0
		// }
		// panic("unreachable")
	}

	// return 0 means the aspect is within epsilon degrees of the orb
	// return 1 means aspect is greater than orb
	// return -1 means aspect is less than orb
	// aspectCmp := func(t time.Time, orb int, epsilon float64) int {
	// 	chrt, err := chart.NewChartFromJulianDay(
	// 		swe,
	// 		swe.GoTimeToJulianDay(t),
	// 		0, 0,
	// 		chart.TropicalChartType,
	// 		[]pointid.PointID{targetAspect.P1, targetAspect.P2},
	// 	)
	// 	if err != nil {
	// 		panic(errors.Wrapf(
	// 			err,
	// 			"calculating chart for %s [chart: %s]",
	// 			t, chrt,
	// 		))
	// 	}

	// 	p1 := chrt.GetPoint(targetAspect.P1)
	// 	p2 := chrt.GetPoint(targetAspect.P2)
	// 	asp := p1.GetAspect(p2)
	// 	// Minus the targetAspect
	// 	deg := asp.Degree - float64(targetAspect.Type.Degree())
	// 	slog.Info("aspectCmp: Aspect",
	// 		slog.String("p1", p1.String()),
	// 		slog.String("p2", p2.String()),
	// 		slog.Float64("deg", deg),
	// 		slog.Time("t", t),
	// 		slog.Int("orb", orb),
	// 		slog.String("targetAspect", targetAspect.String()),
	// 		slog.String("asp", asp.String()),
	// 	)
	// 	if deg > float64(orb) {
	// 		return 1
	// 	}
	// 	if deg < float64(orb) {
	// 		return -1
	// 	}
	// 	if math.Abs(deg-float64(orb)) < epsilon {
	// 		return 0
	// 	}
	// 	panic("unreachable")
	// }

	bigStepSearch := func(direction time.Duration, step time.Duration) (time.Time, time.Time) {
		currentTime := targetAspectTime
		for {
			nextTime := currentTime.Add(step * direction)
			deg, _ := calculateAspectDegree(nextTime)
			// if isEitherRetrograde {
			// 	// Keep going
			// 	slog.Info(
			// 		"bigStepSearch: Retrograde aspect found. Continuing search",
			// 		slog.Time("currentTime", currentTime),
			// 		slog.Time("nextTime", nextTime),
			// 		slog.Duration("step", step),
			// 		slog.Duration("direction", direction),
			// 		slog.Float64("deg", deg),
			// 		slog.Bool("isEitherRetrograde", isEitherRetrograde),
			// 	)
			// 	currentTime = nextTime
			// 	continue
			// }
			isInOrb := math.Abs(deg) < float64(targetAspect.Orb()*2)
			if isInOrb {
				// if aspectCmp(nextTime, targetAspect.Orb()*2, 1.0) < 0 {
				slog.Info("bigStepSearch: Valid aspect found",
					slog.Time("currentTime", currentTime),
					slog.Time("nextTime", nextTime),
					slog.Duration("step", step),
					slog.Duration("direction", direction),
					slog.Float64("deg", deg),
					// slog.Bool("isEitherRetrograde", isEitherRetrograde),
				)
				// if it's still valid, keep going further out
				currentTime = nextTime
				continue
			}
			slog.Info("bigStepSearch: Aspect is no longer valid",
				slog.Time("currentTime", currentTime),
				slog.Time("nextTime", nextTime),
				slog.Duration("step", step),
				slog.Duration("direction", direction),
				slog.Float64("deg", deg),
				// slog.Bool("isEitherRetrograde", isEitherRetrograde),
			)
			return currentTime, nextTime
		}
	}

	binarySearch := func(direction time.Duration, start, end time.Time) time.Time {
		var mid time.Time
		L := start
		R := end
		if direction < 0 {
			L, R = R, L
		}
		slog.Info("binarySearch: Starting",
			slog.Time("L", L),
			slog.Time("R", R),
			slog.Duration("direction", direction),
		)
		for L.Before(R) {
			mid = L.Add(R.Sub(L) / 2)
			deg, _ := calculateAspectDegree(mid)
			// So assuming we return 13 and the orb of a sextile is 3
			orbDiff := deg - float64(targetAspect.Orb())
			// orbDiff would report 13-3 = 10 (or -10 in the other direction),
			// so it is 10 degrees away from the orb's edge, so we know we
			// should set the edge accordingly
			// If direction > 1, L=mid
			isWithinEpsilon := math.Abs(orbDiff) < 1.0
			slog.Info("binarySearch: Checking mid",
				slog.Time("L", L),
				slog.Time("R", R),
				slog.Time("mid", mid),
				slog.Duration("direction", direction),
				slog.Float64("deg", deg),
				slog.Float64("orbDiff", orbDiff),
				// slog.Bool("isEitherRetrograde", isEitherRetrograde),
			)
			if isWithinEpsilon {
				// If the difference is within 1 degree, we're done
				return mid
			}
			// if isWithinEpsilon {
			// if isEitherRetrograde {
			// 	// Keep going
			// 	slog.Info(
			// 		"binarySearch: Retrograde aspect found. Continuing search",
			// 		slog.Time("L", L),
			// 		slog.Time("R", R),
			// 		slog.Time("mid", mid),
			// 		slog.Duration("direction", direction),
			// 		slog.Float64("deg", deg),
			// 		slog.Bool("isEitherRetrograde", isEitherRetrograde),
			// 	)
			// 	if direction > 0 {
			// 		L = mid
			// 	} else {
			// 		R = mid
			// 	}
			// 	continue
			// }
			// return mid
			// }
			if orbDiff < 0 {
				if direction < 0 {
					R = mid
				} else {
					L = mid
				}
			} else if orbDiff > 0 {
				if direction < 0 {
					L = mid
				} else {
					R = mid
				}
			}
		}
		// if we get nothing, return the last mid
		return mid
	}

	findEdge := func(direction time.Duration, step time.Duration, startingPoint time.Time) time.Time {
		// First big-step search: Find the maximum edge
		edgeBeginAt, edgeEndAt := bigStepSearch(direction, step)
		slog.Info("findEdge",
			slog.Time("edgeBeginAt", edgeBeginAt),
			slog.Time("edgeEndAt", edgeEndAt),
		)
		return binarySearch(direction, edgeBeginAt, edgeEndAt)
	}

	// Move at the slowest speed for the search
	speed := getSpeedRate(targetAspect.P1)
	speed2 := getSpeedRate(targetAspect.P2)
	if speed2 < speed {
		speed = speed2
	}

	slog.Info("Discovering aspect journey",
		slog.Duration("speed", speed),
		slog.String("targetAspect", targetAspect.String()),
		slog.Time("targetAspectTime", targetAspectTime),
	)

	start = findEdge(time.Duration(-1), speed, targetAspectTime)
	slog.Info("Found aspect start",
		slog.Time("start", start),
		slog.Duration("speed", speed),
		slog.String("targetAspect", targetAspect.String()),
		slog.Time("targetAspectTime", targetAspectTime),
	)

	end = findEdge(time.Duration(1), speed, targetAspectTime)
	duration = end.Sub(start)
	journey = float64(targetAspectTime.Sub(start)) / float64(duration)
	slog.Info("Aspect journey",
		slog.Duration("speed", speed),
		slog.String("targetAspect", targetAspect.String()),
		slog.Time("targetAspectTime", targetAspectTime),
		slog.Time("start", start),
		slog.Time("end", end),
		slog.Duration("duration", duration),
		slog.Float64("journey", journey),
	)
	return duration, journey, start, end, nil
}

// func calculateAspectJourney(
// 	swe *wrapper.SwissEph,
// 	targetAspect *aspect.Aspect,
// 	targetAspectTime time.Time,
// ) (duration time.Duration, journey float64, start time.Time, end time.Time, err error) {
// 	if targetAspect.Type == aspect.AspectType_None {
// 		return 0, 0, time.Time{}, time.Time{}, nil
// 	}

// 	findEdge := func(rate time.Duration) (time.Time, error) {
// 		currentTime := targetAspectTime
// 		var lastValidTime time.Time

// 		for {
// 			currentTime = currentTime.Add(rate)
// 			chrt, err := chart.NewChartFromJulianDay(
// 				swe,
// 				swe.GoTimeToJulianDay(currentTime),
// 				0, 0,
// 				chart.TropicalChartType,
// 				[]pointid.PointID{targetAspect.P1, targetAspect.P2},
// 			)
// 			if err != nil {
// 				return time.Time{}, errors.Wrapf(
// 					err,
// 					"calculating chart for %s",
// 					currentTime,
// 				)
// 			}

// 			if chrt.HasAspectIgnoreDegree(targetAspect) {
// 				lastValidTime = currentTime
// 				// fmt.Printf("Found aspect at {%s}\n", currentTime)
// 				continue
// 			}

// 			// If there's no aspect:
// 			// - there's a chance the planets form the same aspect in the
// 			//   future, either because one goes retrograde or just because they
// 			//   move away a bit from each other and rejoin in the future
// 			//   - Astrologically, we still wanna include that as part of the
// 			//   journey
// 			// - So keep going until the aspect between those two planets is
// 			//   triple the orb of the original aspect
// 			//   - if the aspect reaches triple the degree of the orb, we're
// 			//   pretty clear these two planets have concluded their journey
// 			//   - So for a conjunction of orb 5 degrees, keep going until the
// 			//     planets reach a degree of 15 between them
// 			//   - If they do and never form an orb again, just use 'lastTime'
// 			//   above
// 			//   - But if they do form an orb again, calculate that one instead

// 			p1 := chrt.GetPoint(targetAspect.P1)
// 			p2 := chrt.GetPoint(targetAspect.P2)
// 			asp := p1.GetAspect(p2)
// 			if int(asp.Degree) > targetAspect.Orb()*3 {
// 				// fmt.Printf(
// 				// 	"Reached orb of %d degrees at {%s}\n",
// 				// 	targetAspect.Orb(),
// 				// 	currentTime,
// 				// )
// 				if lastValidTime.IsZero() {
// 					// fmt.Printf("Returning current time {%s}\n", currentTime)
// 					return currentTime, nil
// 				}
// 				// fmt.Printf("Returning last valid time {%s}\n", lastValidTime)
// 				return lastValidTime, nil
// 			}
// 		}
// 	}

// 	// Move at the slowest speed
// 	speed := getSpeedRate(targetAspect.P1)
// 	speed2 := getSpeedRate(targetAspect.P2)
// 	if speed2 < speed {
// 		speed = speed2
// 	}

// 	start, err = findEdge(-speed)
// 	if err != nil {
// 		return 0, 0, time.Time{}, time.Time{}, errors.Wrap(
// 			err,
// 			"finding aspect start",
// 		)
// 	}

// 	end, err = findEdge(speed)
// 	if err != nil {
// 		return 0, 0, time.Time{}, time.Time{}, errors.Wrap(
// 			err,
// 			"finding aspect end",
// 		)
// 	}

// 	duration = end.Sub(start)
// 	journey = float64(targetAspectTime.Sub(start)) / float64(duration)
// 	return duration, journey, start, end, nil
// }

// func calculateIngressJourney(
// 	swe *wrapper.SwissEph,
// 	targetPoint *astropoint.AstroPoint,
// 	targetTime time.Time,
// ) (duration time.Duration, journey float64, start time.Time, end time.Time, err error) {
// 	// Check cache
// 	// cacheKey := getCacheKey(targetTime, targetPoint.ID, 0)
// 	// if cached, ok := ingressCache[cacheKey]; ok {
// 	//     return time.Duration(cached.DaysElapsed) * 24 * time.Hour,
// 	//         cached.Journey,
// 	//         cached.Start,
// 	//         cached.End,
// 	//         nil
// 	// }

// 	findEdge := func(rate time.Duration) (time.Time, error) {
// 		left := targetTime
// 		right := targetTime.Add(rate)
// 		var lastValid time.Time

// 		for left.Before(right) {
// 			mid := left.Add(right.Sub(left) / 2)

// 			chrt, err := chart.NewChartFromJulianDay(
// 				swe,
// 				swe.GoTimeToJulianDay(mid),
// 				0, 0,
// 				chart.TropicalChartType,
// 				[]pointid.PointID{targetPoint.ID},
// 			)
// 			if err != nil {
// 				return time.Time{}, err
// 			}

// 			currPoint := chrt.GetPoint(targetPoint.ID)
// 			if currPoint.ZodiacalPos.Sign != targetPoint.ZodiacalPos.Sign {
// 				lastValid = mid
// 				if rate < 0 {
// 					right = mid
// 				} else {
// 					left = mid
// 				}
// 			} else {
// 				if rate < 0 {
// 					left = mid
// 				} else {
// 					right = mid
// 				}
// 			}
// 		}
// 		return lastValid, nil
// 	}

// 	speed := getSpeedRate(targetPoint.ID)
// 	start, err = findEdge(-speed)
// 	if err != nil {
// 		return 0, 0, time.Time{}, time.Time{}, err
// 	}

// 	end, err = findEdge(speed)
// 	if err != nil {
// 		return 0, 0, time.Time{}, time.Time{}, err
// 	}

// 	duration = end.Sub(start)
// 	journey = float64(targetTime.Sub(start)) / float64(duration)

// 	// Cache result
// 	// ingressCache[cacheKey] = &TransitIngress{
// 	// 	Date:        targetTime,
// 	// 	P:           targetPoint,
// 	// 	Journey:     journey,
// 	// 	DaysElapsed: int(duration.Hours() / 24),
// 	// 	Start:       start,
// 	// 	End:         end,
// 	// }

// return duration, journey, start, end, nil
// }

func calculateIngressJourney(
	swe *wrapper.SwissEph,
	targetPoint *astropoint.AstroPoint,
	targetTime time.Time,
) (duration time.Duration, journey float64, start time.Time, end time.Time, err error) {
	findEdge := func(rate time.Duration) (time.Time, error) {
		currentTime := targetTime

		// var targetZP *zodiacalpos.ZodiacalPos
		// if rate < 0 {
		// 	targetZP = zodiacalpos.NewZodiacalPos(
		// 		targetPoint.ZodiacalPos.Sign,
		// 		0,
		// 		0,
		// 	)
		// } else {
		// 	targetZP = zodiacalpos.NewZodiacalPos(targetPoint.ZodiacalPos.Sign.Next(), 0, 0)
		// }
		// fmt.Printf("Target ZP: %s\n", targetZP)

		for {
			currentTime = currentTime.Add(rate)
			chrt, err := chart.NewChartFromJulianDay(
				swe,
				swe.GoTimeToJulianDay(currentTime),
				0, 0,
				chart.TropicalChartType,
				[]pointid.PointID{targetPoint.ID},
			)
			if err != nil {
				return time.Time{}, errors.Wrapf(
					err,
					"calculating chart for %s",
					currentTime,
				)
			}

			currPoint := chrt.GetPoint(targetPoint.ID)
			if currPoint.ZodiacalPos.Sign != targetPoint.ZodiacalPos.Sign {
				// If the point got out of the sign, break
				return currentTime.Add(-rate), nil
			}
		}
	}

	start, err = findEdge(-getSpeedRate(targetPoint.ID))
	if err != nil {
		return 0, 0, time.Time{}, time.Time{}, errors.Wrap(
			err,
			"finding ingress start",
		)
	}
	end, err = findEdge(getSpeedRate(targetPoint.ID))
	if err != nil {
		return 0, 0, time.Time{}, time.Time{}, errors.Wrap(
			err,
			"finding ingress end",
		)
	}

	duration = end.Sub(start)
	journey = float64(targetTime.Sub(start)) / float64(duration)
	return duration, journey, start, end, nil
}

func getSpeedRate(p pointid.PointID) time.Duration {
	switch p {
	case pointid.Sun:
		return 30 * 24 * time.Hour
	case pointid.Moon:
		return 2.5 * 24 * time.Hour
	case pointid.Mercury:
		return 40 * 24 * time.Hour
	case pointid.Venus:
		return 65 * 24 * time.Hour
	case pointid.Mars:
		return 90 * 24 * time.Hour
	case pointid.Jupiter:
		return 14 * 30 * 24 * time.Hour
	case pointid.Saturn:
		return 4 * 365 * 24 * time.Hour
	case pointid.Uranus:
		return 9 * 365 * 24 * time.Hour
	case pointid.Neptune:
		return 15 * 365 * 24 * time.Hour
	case pointid.Pluto:
		return 34 * 365 * 24 * time.Hour
	default:
		panic(errors.Newf("Unknown planet: %s", p))
	}
}
