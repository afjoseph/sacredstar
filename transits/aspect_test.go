package transits

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/afjoseph/sacredstar/aspect"
	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/stretchr/testify/assert"
)

func TestCalculateAspectJourney(t *testing.T) {
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	type testCase struct {
		name         string
		targetP1     pointid.PointID
		targetP2     pointid.PointID
		targetTime   time.Time
		wantDuration time.Duration
		wantJourney  float64
	}

	tests := []testCase{
		testCase{
			name:         "fast conjunction in Aquarius",
			targetP1:     pointid.Pluto,
			targetP2:     pointid.Sun,
			targetTime:   time.Date(2025, 1, 24, 0, 0, 0, 0, time.UTC),
			wantDuration: 252 * time.Hour,
			wantJourney:  0.75,
		},
		testCase{
			name:       "conjunction-venus-retrograde",
			targetP1:   pointid.Venus,
			targetP2:   pointid.Neptune,
			targetTime: time.Date(2025, 1, 29, 0, 0, 0, 0, time.UTC),
			// This is a retrograde of Venus in Pisces/Aries where Neptune is
			// copresent
			wantDuration: 315 * time.Hour,
			wantJourney:  0.20,
		},
		testCase{
			name:         "slow moving sextile",
			targetP1:     pointid.Neptune,
			targetP2:     pointid.Pluto,
			targetTime:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
			wantDuration: 122040 * time.Hour,
			wantJourney:  0.358,
		},
		testCase{
			name:         "saturn-neptune",
			targetP1:     pointid.Saturn,
			targetP2:     pointid.Neptune,
			targetTime:   time.Date(2025, 4, 30, 0, 0, 0, 0, time.UTC),
			wantDuration: 9000 * time.Hour,
			wantJourney:  0.060,
		},
		testCase{
			name:         "fast moving trine",
			targetP1:     pointid.Moon,
			targetP2:     pointid.Uranus,
			targetTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantDuration: 17 * time.Hour,
			wantJourney:  0.538,
		},
		testCase{
			name:         "mars-neptune-retrograde",
			targetP1:     pointid.Mars,
			targetP2:     pointid.Neptune,
			targetTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantDuration: 2394 * time.Hour,
			wantJourney:  0.761,
		},
		testCase{
			name:         "mars-pluto",
			targetP1:     pointid.Mars,
			targetP2:     pointid.Pluto,
			targetTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantDuration: 2079 * time.Hour,
			wantJourney:  0.835,
		},
		testCase{
			name:         "jupiter-saturn",
			targetP1:     pointid.Jupiter,
			targetP2:     pointid.Saturn,
			targetTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantDuration: 8370 * time.Hour,
			wantJourney:  0.451,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrt, err := chart.NewChartFromJulianDay(
				swe,
				swe.GoTimeToJulianDay(tt.targetTime),
				0, 0,
				chart.TropicalChartType,
				pointid.ModernPlanets,
			)
			assert.NoError(t, err)
			asp := chrt.GetPoint(tt.targetP1).
				GetAspect(chrt.GetPoint(tt.targetP2))

			gotDuration, gotJourney, _, _, err := calculateAspectJourney(
				swe,
				asp,
				tt.targetTime,
			)
			assert.NoError(t, err)
			durationDiff := math.Abs(
				gotDuration.Hours() - tt.wantDuration.Hours(),
			)
			assert.True(t, durationDiff < 1, "actual: %s, want: %s",
				gotDuration,
				tt.wantDuration,
			)
			assert.InEpsilon(
				t,
				tt.wantJourney,
				gotJourney,
				0.01,
				"actual: %f, want: %f",
				gotJourney,
				tt.wantJourney,
			)
		})
	}
}

func TestTransitAspect_MarshalJSON(t *testing.T) {
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	dt := time.Now()
	chrt, err := chart.NewChartFromJulianDay(
		swe,
		swe.GoTimeToJulianDay(dt),
		0, 0, // lon, lat: we're assuming UTC for now
		chart.TropicalChartType,
		pointid.ModernPlanets,
	)
	assert.NoError(t, err)

	var asp *aspect.Aspect
	for _, a := range chrt.Aspects {
		// XXX <04-02-2025,afjoseph> Get an aspect that is not the ascendant:
		// we don't do ascendant calculations now
		if a.P1 == pointid.ASC || a.P2 == pointid.ASC {
			continue
		}
		asp = a
		break
	}

	ts, err := newTransitAspect(
		swe,
		// Doesn't matter which aspect
		asp,
		dt,
	)
	assert.NoError(t, err)

	b, err := json.Marshal(ts)
	assert.NoError(t, err)

	var ts2 TransitAspect
	err = json.Unmarshal(b, &ts2)
	assert.NoError(t, err)
	assert.Equal(t, ts.String(), ts2.String())
}
