package transits

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/unixtime"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/go-playground/errors/v5"
)

type TransitType string

const (
	TransitTypeAspect  TransitType = "aspect"
	TransitTypeIngress TransitType = "ingress"
)

type Transit interface {
	fmt.Stringer
	GetType() TransitType
	GetJourney() float64
	GetDuration() int
	GetStart() unixtime.UnixTime
	GetEnd() unixtime.UnixTime
}

type transitBase struct {
	Type        TransitType       `json:"type"`
	Date        unixtime.UnixTime `json:"date"`
	Journey     float64           `json:"journey"`
	DaysElapsed int               `json:"daysElapsed"`
	Start       unixtime.UnixTime `json:"start"`
	End         unixtime.UnixTime `json:"end"`
}

type Transits []Transit

func (tss *Transits) UnmarshalJSON(data []byte) error {
	// Unmarshal into []json.RawMessage
	var rawList []json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return errors.Wrapf(err, "while unmarshalling transits")
	}

	// Process the type
	for _, raw := range rawList {
		var base transitBase
		if err := json.Unmarshal(raw, &base); err != nil {
			return errors.Wrapf(err, "while unmarshalling transit base")
		}

		var t Transit
		switch base.Type {
		case TransitTypeAspect:
			var ta TransitAspect
			if err := json.Unmarshal(raw, &ta); err != nil {
				return errors.Wrapf(err, "while unmarshalling transit aspect")
			}
			t = &ta
		case TransitTypeIngress:
			var ti TransitIngress
			if err := json.Unmarshal(raw, &ti); err != nil {
				return errors.Wrapf(err, "while unmarshalling transit ingress")
			}
			t = &ti
		default:
			return errors.Newf("unknown transit type: %s", base.Type)
		}
		*tss = append(*tss, t)
	}
	return nil
}

func New(
	swe *wrapper.SwissEph,
	t time.Time,
) (Transits, error) {
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

func calculate(swe *wrapper.SwissEph, t time.Time) (Transits, error) {
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
		ts, err := newTransitIngress(swe, p, t)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"while calculating ingress journey for %s",
				t,
			)
		}
		transits = append(transits, ts)
	}

	// Calculate aspect journeys
	for _, asp := range chrt.Aspects {
		if asp.P1 == pointid.ASC || asp.P2 == pointid.ASC {
			// Skip ascendant calculations
			continue
		}

		// fmt.Printf("Calculating aspect journey for %s\n", asp)
		ts, err := newTransitAspect(swe, asp, t)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"while calculating aspect journey for %s",
				t,
			)
		}
		transits = append(transits, ts)
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
