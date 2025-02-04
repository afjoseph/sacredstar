package transits

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/lunation"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/stretchr/testify/assert"
)

func TestCalculateLunationJourney(t *testing.T) {
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	type testCase struct {
		name             string
		targetTime       time.Time
		hasLunation      bool
		wantLunationType lunation.LunationType
	}

	tests := []testCase{
		testCase{
			name:             "New Moon",
			targetTime:       time.Date(2025, 1, 29, 12, 36, 0, 0, time.UTC),
			hasLunation:      true,
			wantLunationType: lunation.LunationTypeNewMoon,
		},
		testCase{
			name:        "No lunation",
			targetTime:  time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC),
			hasLunation: false,
		},
		testCase{
			name:             "Full Moon",
			targetTime:       time.Date(2025, 2, 12, 0, 0, 0, 0, time.UTC),
			hasLunation:      true,
			wantLunationType: lunation.LunationTypeFullMoon,
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
			if !tt.hasLunation {
				assert.Nil(t, chrt.Lunation)
				return
			}

			assert.NotNil(t, chrt.Lunation)
			assert.Equal(t, tt.wantLunationType, chrt.Lunation.Type)
			ts := newTransitLunation(
				chrt.Lunation,
				tt.targetTime,
			)
			assert.NotNil(t, ts)
			assert.Equal(t, tt.wantLunationType, ts.Lunation.Type)
		})
	}
}

func TestTransitLunation_MarshalJSON(t *testing.T) {
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	// New moon
	dt := time.Date(2025, 1, 29, 12, 36, 0, 0, time.UTC)
	chrt, err := chart.NewChartFromJulianDay(
		swe,
		swe.GoTimeToJulianDay(dt),
		0, 0, // lon, lat: we're assuming UTC for now
		chart.TropicalChartType,
		pointid.ModernPlanets,
	)
	assert.NoError(t, err)
	assert.NotNil(t, chrt.Lunation)

	ts := newTransitLunation(
		chrt.Lunation,
		dt,
	)
	assert.NotNil(t, ts)

	b, err := json.Marshal(ts)
	assert.NoError(t, err)

	var ts2 TransitLunation
	err = json.Unmarshal(b, &ts2)
	assert.NoError(t, err)
	assert.Equal(t, ts.String(), ts2.String())
}
