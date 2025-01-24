package transits

import (
	"fmt"
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

	return transits, nil
}

func getStepForPointID(p pointid.PointID) time.Duration {
	switch p {
	case pointid.Moon:
		return 7 * 24 * time.Hour
	case pointid.Sun:
		return 14 * 24 * time.Hour
	case pointid.Mercury:
		return 14 * 24 * time.Hour
	case pointid.Venus:
		return 14 * 24 * time.Hour
	case pointid.Mars:
		return 14 * 24 * time.Hour
	case pointid.Jupiter:
		return 1 * 30 * 24 * time.Hour
	case pointid.Saturn:
		return 1 * 30 * 24 * time.Hour
	case pointid.Uranus:
		return 6 * 30 * 24 * time.Hour
	case pointid.Neptune:
		return 6 * 30 * 24 * time.Hour
	case pointid.Pluto:
		return 6 * 30 * 24 * time.Hour
	default:
		panic(errors.Newf("Unknown planet: %s", p))
	}
}
