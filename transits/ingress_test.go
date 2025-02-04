package transits

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/afjoseph/sacredstar/astropoint"
	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/stretchr/testify/assert"
)

func TestCalculateIngressJourney(t *testing.T) {
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	type testCase struct {
		name          string
		targetPointID pointid.PointID
		targetTime    time.Time
		wantDuration  time.Duration
		wantJourney   float64
	}

	tests := []testCase{
		testCase{
			name:          "1",
			targetPointID: pointid.Venus,
			targetTime:    time.Date(2025, 1, 24, 0, 0, 0, 0, time.UTC),
			wantDuration:  771 * time.Hour,
			wantJourney:   0.67,
		},
		testCase{
			name:          "2",
			targetPointID: pointid.Mercury,
			targetTime:    time.Date(2025, 1, 24, 0, 0, 0, 0, time.UTC),
			wantDuration:  472 * time.Hour,
			wantJourney:   0.79,
		},
		testCase{
			name:          "3",
			targetPointID: pointid.Mars,
			targetTime:    time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC),
			// This includes a retrograde journey
			wantDuration: 5407 * time.Hour,
			wantJourney:  0.55,
		},
		testCase{
			name:          "4",
			targetPointID: pointid.Mars,
			targetTime:    time.Date(2025, 2, 25, 0, 0, 0, 0, time.UTC),
			// This includes a retrograde journey
			wantDuration: 5407 * time.Hour,
			wantJourney:  0.77,
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
			p := chrt.GetPoint(tt.targetPointID)
			assert.NotNil(t, p)

			gotDuration, gotJourney, _, _, err := calculateIngressJourney(
				swe,
				p,
				tt.targetTime,
			)
			assert.NoError(t, err)
			durationDiff := math.Abs(
				gotDuration.Hours() - tt.wantDuration.Hours(),
			)
			assert.True(t, durationDiff < 1.0, "actual: %s, want: %s",
				gotDuration,
				tt.wantDuration,
			)
			assert.InEpsilon(
				t,
				tt.wantJourney,
				gotJourney,
				0.1,
				"actual: %.2f, want: %.2f",
				gotJourney,
				tt.wantJourney,
			)
		})
	}
}

func TestTransitIngress_MarshalJSON(t *testing.T) {
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

	// XXX <04-02-2025,afjoseph> Get an ingress that is not the ascendant:
	// we don't do ascendant calculations now
	var p *astropoint.AstroPoint
	for _, pp := range chrt.Points {
		if pp.ID == pointid.ASC {
			continue
		}
		p = pp
		break
	}
	ts, err := newTransitIngress(
		swe,
		// Doesn't matter which ingress
		p,
		dt,
	)
	assert.NoError(t, err)

	b, err := json.Marshal(ts)
	assert.NoError(t, err)

	var ts2 TransitIngress
	err = json.Unmarshal(b, &ts2)
	assert.NoError(t, err)
	assert.Equal(t, ts.String(), ts2.String())
}
