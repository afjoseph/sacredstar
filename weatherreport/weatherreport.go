package weatherreport

import (
	"fmt"
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

	findEdge := func(rate time.Duration) (time.Time, error) {
		currentTime := targetAspectTime
		var lastValidTime time.Time

		for {
			currentTime = currentTime.Add(rate)
			chrt, err := chart.NewChartFromJulianDay(
				swe,
				swe.GoTimeToJulianDay(currentTime),
				0, 0,
				chart.TropicalChartType,
				[]pointid.PointID{targetAspect.P1, targetAspect.P2},
			)
			if err != nil {
				return time.Time{}, errors.Wrapf(
					err,
					"calculating chart for %s",
					currentTime,
				)
			}

			if chrt.HasAspectIgnoreDegree(targetAspect) {
				lastValidTime = currentTime
				// fmt.Printf("Found aspect at {%s}\n", currentTime)
				continue
			}

			// If there's no aspect:
			// - there's a chance the planets form the same aspect in the
			//   future, either because one goes retrograde or just because they
			//   move away a bit from each other and rejoin in the future
			//   - Astrologically, we still wanna include that as part of the
			//   journey
			// - So keep going until the aspect between those two planets is
			//   triple the orb of the original aspect
			//   - if the aspect reaches triple the degree of the orb, we're
			//   pretty clear these two planets have concluded their journey
			//   - So for a conjunction of orb 5 degrees, keep going until the
			//     planets reach a degree of 15 between them
			//   - If they do and never form an orb again, just use 'lastTime'
			//   above
			//   - But if they do form an orb again, calculate that one instead

			p1 := chrt.GetPoint(targetAspect.P1)
			p2 := chrt.GetPoint(targetAspect.P2)
			asp := p1.GetAspect(p2)
			if int(asp.Degree) > targetAspect.Orb()*3 {
				// fmt.Printf(
				// 	"Reached orb of %d degrees at {%s}\n",
				// 	targetAspect.Orb(),
				// 	currentTime,
				// )
				if lastValidTime.IsZero() {
					// fmt.Printf("Returning current time {%s}\n", currentTime)
					return currentTime, nil
				}
				// fmt.Printf("Returning last valid time {%s}\n", lastValidTime)
				return lastValidTime, nil
			}
		}
	}

	// Move at the slowest speed
	speed := getSpeedRate(targetAspect.P1)
	speed2 := getSpeedRate(targetAspect.P2)
	if speed2 < speed {
		speed = speed2
	}

	start, err = findEdge(-speed)
	if err != nil {
		return 0, 0, time.Time{}, time.Time{}, errors.Wrap(
			err,
			"finding aspect start",
		)
	}

	end, err = findEdge(speed)
	if err != nil {
		return 0, 0, time.Time{}, time.Time{}, errors.Wrap(
			err,
			"finding aspect end",
		)
	}

	duration = end.Sub(start)
	journey = float64(targetAspectTime.Sub(start)) / float64(duration)
	return duration, journey, start, end, nil
}

func calculateIngressJourney(
	swe *wrapper.SwissEph,
	targetPoint *astropoint.AstroPoint,
	targetTime time.Time,
) (duration time.Duration, journey float64, start time.Time, end time.Time, err error) {
	findEdge := func(rate time.Duration) (time.Time, error) {
		currentTime := targetTime

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
		fallthrough
	case pointid.Moon:
		fallthrough
	case pointid.Mercury:
		fallthrough
	case pointid.Venus:
		fallthrough
	case pointid.Mars:
		// 6 hours for fast planets
		return 6 * time.Hour
	case pointid.Jupiter:
		fallthrough
	case pointid.Saturn:
		// 24 hours for traditional big planets
		return 24 * time.Hour
	case pointid.Uranus:
		fallthrough
	case pointid.Neptune:
		fallthrough
	case pointid.Pluto:
		// 4 days for outer planets
		return 4 * (24 * time.Hour)
	default:
		panic(errors.Newf("Unknown planet: %s", p))
	}
}
