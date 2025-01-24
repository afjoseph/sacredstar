package transits

import (
	"fmt"
	"log/slog"
	"math"
	"os"
	"testing"
	"time"

	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/prettyslog"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/stretchr/testify/assert"
)

func init() {
	slog.SetDefault(slog.New(prettyslog.NewPrettyJSONHandler(os.Stdout, nil)))
}

// TestNew_NoPanic is a test to ensure that New() does not panic after X months
// of calculation
func TestNew_NoPanic(t *testing.T) {
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	// Loop for 6 months
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for d := start; d.Before(start.Add(6 * 30 * 24 * time.Hour)); d = d.Add(24 * time.Hour) {
		t.Run(d.Format("2006-01-02"), func(t *testing.T) {
			// t.Parallel()
			startTime := time.Now()
			_, err := New(swe, d)
			duration := time.Since(startTime)
			assert.NoError(t, err)
			fmt.Printf("%s: took %s\n", d.Format("2006-01-02"), duration)
		})
	}
}

func TestNew_Measure(t *testing.T) {
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	type testCase struct {
		name    string
		day     time.Time
		wantRet []string
	}

	tests := []testCase{
		testCase{
			name: "1",
			day:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantRet: []string{
				"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: AstroPoint{ID: sun, Longitude: 280.813611, ZodiacalPos: capricorn 10 48, House: 4th, IsRetrograde: false}, Journey: 0.36, DaysElapsed: 29, Start: 2024-12-21, End: 2025-01-19}",
				"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: AstroPoint{ID: moon, Longitude: 293.913579, ZodiacalPos: capricorn 23 54, House: 4th, IsRetrograde: false}, Journey: 0.75, DaysElapsed: 111, Start: 2024-10-09, End: 2025-01-28}",
				"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: AstroPoint{ID: mercury, Longitude: 259.869978, ZodiacalPos: sagittarius 19 52, House: 3rd, IsRetrograde: false}, Journey: 0.89, DaysElapsed: 66, Start: 2024-11-02, End: 2025-01-08}",
				"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: AstroPoint{ID: venus, Longitude: 327.712110, ZodiacalPos: aquarius 27 42, House: 5th, IsRetrograde: false}, Journey: 0.92, DaysElapsed: 26, Start: 2024-12-07, End: 2025-01-03}",
				"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: AstroPoint{ID: mars, Longitude: 121.917923, ZodiacalPos: leo 1 55, House: 11th, IsRetrograde: true}, Journey: 0.26, DaysElapsed: 225, Start: 2024-11-04, End: 2025-06-17}",
				"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: AstroPoint{ID: jupiter, Longitude: 73.215463, ZodiacalPos: gemini 13 12, House: 9th, IsRetrograde: true}, Journey: 0.58, DaysElapsed: 380, Start: 2024-05-25, End: 2025-06-10}",
				"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: AstroPoint{ID: saturn, Longitude: 344.524071, ZodiacalPos: pisces 14 31, House: 6th, IsRetrograde: false}, Journey: 0.82, DaysElapsed: 810, Start: 2023-03-07, End: 2025-05-25}",
				"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: AstroPoint{ID: uranus, Longitude: 53.635827, ZodiacalPos: taurus 23 38, House: 8th, IsRetrograde: true}, Journey: 0.82, DaysElapsed: 2607, Start: 2019-03-07, End: 2026-04-26}",
				"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: AstroPoint{ID: neptune, Longitude: 357.297825, ZodiacalPos: pisces 27 17, House: 6th, IsRetrograde: false}, Journey: 0.98, DaysElapsed: 4803, Start: 2012-02-05, End: 2025-04-01}",
				"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: AstroPoint{ID: pluto, Longitude: 301.064802, ZodiacalPos: aquarius 1 3, House: 5th, IsRetrograde: false}, Journey: 0.05, DaysElapsed: 6986, Start: 2024-01-23, End: 2043-03-11}",
				"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: moon, P2: uranus, Degree: 119.733333, Type: Trine}, Journey: 0.54, DaysElapsed: 0, Start: 2024-12-31, End: 2025-01-01}",
				"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: venus, P2: uranus, Degree: 85.933333, Type: Square}, Journey: 0.90, DaysElapsed: 8, Start: 2024-12-24, End: 2025-01-01}",
				"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: mars, P2: neptune, Degree: 124.633333, Type: Trine}, Journey: 0.75, DaysElapsed: 99, Start: 2024-10-17, End: 2025-01-25}",
				"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: mars, P2: pluto, Degree: 179.133333, Type: Opposition}, Journey: 0.83, DaysElapsed: 86, Start: 2024-10-21, End: 2025-01-15}",
				"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: jupiter, P2: saturn, Degree: 88.683333, Type: Square}, Journey: 0.45, DaysElapsed: 348, Start: 2024-07-27, End: 2025-07-11}",
				"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: saturn, P2: jupiter, Degree: 88.683333, Type: Square}, Journey: 0.45, DaysElapsed: 348, Start: 2024-07-27, End: 2025-07-11}",
				"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: uranus, P2: moon, Degree: 119.733333, Type: Trine}, Journey: 0.54, DaysElapsed: 0, Start: 2024-12-31, End: 2025-01-01}",
				"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: uranus, P2: venus, Degree: 85.933333, Type: Square}, Journey: 0.90, DaysElapsed: 8, Start: 2024-12-24, End: 2025-01-01}",
				"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: neptune, P2: mars, Degree: 124.633333, Type: Trine}, Journey: 0.75, DaysElapsed: 99, Start: 2024-10-17, End: 2025-01-25}",
				"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: pluto, P2: mars, Degree: 179.133333, Type: Opposition}, Journey: 0.83, DaysElapsed: 86, Start: 2024-10-21, End: 2025-01-15}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(swe, tt.day)
			assert.NoError(t, err)

			m := []string{}
			for _, v := range got {
				m = append(m, v.String())
			}
			// spew.Dump(m)
			assert.Equal(t, tt.wantRet, m)
		})
	}
}

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
