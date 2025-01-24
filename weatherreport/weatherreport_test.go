package weatherreport

import (
	"testing"
	"time"

	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	type testCase struct {
		name    string
		start   time.Time
		end     time.Time
		wantRet map[string][]string
	}

	tests := []testCase{
		testCase{
			name:  "1",
			start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2025, 1, 2, 1, 0, 0, 0, time.UTC),
			wantRet: map[string][]string{
				"2025-01-01": []string{
					"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: moon, P2: uranus, Degree: 119.733333, Type: Trine}, Journey: 0.50, DaysElapsed: 0, Start: 2024-12-31, End: 2025-01-01}",
					"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: venus, P2: uranus, Degree: 85.933333, Type: Square}, Journey: 0.91, DaysElapsed: 8, Start: 2024-12-24, End: 2025-01-01}",
					"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: mars, P2: neptune, Degree: 124.633333, Type: Trine}, Journey: 0.04, DaysElapsed: 25, Start: 2024-12-31, End: 2025-01-25}",
					"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: mars, P2: pluto, Degree: 179.133333, Type: Opposition}, Journey: 0.53, DaysElapsed: 30, Start: 2024-12-16, End: 2025-01-15}",
					"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: jupiter, P2: saturn, Degree: 88.683333, Type: Square}, Journey: 0.61, DaysElapsed: 57, Start: 2024-11-27, End: 2025-01-23}",
					"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: saturn, P2: jupiter, Degree: 88.683333, Type: Square}, Journey: 0.61, DaysElapsed: 57, Start: 2024-11-27, End: 2025-01-23}",
					"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: uranus, P2: moon, Degree: 119.733333, Type: Trine}, Journey: 0.50, DaysElapsed: 0, Start: 2024-12-31, End: 2025-01-01}",
					"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: uranus, P2: venus, Degree: 85.933333, Type: Square}, Journey: 0.91, DaysElapsed: 8, Start: 2024-12-24, End: 2025-01-01}",
					"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: neptune, P2: mars, Degree: 124.633333, Type: Trine}, Journey: 0.04, DaysElapsed: 25, Start: 2024-12-31, End: 2025-01-25}",
					"TransitAspect{Date: 2025-01-01 00:00:00 +0000 UTC, Aspect: Aspect{P1: pluto, P2: mars, Degree: 179.133333, Type: Opposition}, Journey: 0.53, DaysElapsed: 30, Start: 2024-12-16, End: 2025-01-15}",
					"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: sun: longitude = 280.813611, zodpos = {capricorn 10 48}, house = 4th, Journey: 0.36, DaysElapsed: 29, Start: 2024-12-21, End: 2025-01-19}",
					"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: moon: longitude = 293.913579, zodpos = {capricorn 23 54}, house = 4th, Journey: 0.88, DaysElapsed: 2, Start: 2024-12-30, End: 2025-01-01}",
					"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: mercury: longitude = 259.869978, zodpos = {sagittarius 19 52}, house = 3rd, Journey: 0.89, DaysElapsed: 66, Start: 2024-11-03, End: 2025-01-08}",
					"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: venus: longitude = 327.712110, zodpos = {aquarius 27 42}, house = 5th, Journey: 0.92, DaysElapsed: 26, Start: 2024-12-07, End: 2025-01-03}",
					"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: mars: longitude = 121.917923, zodpos = {leo 1 55}, house = 11th, Journey: 0.92, DaysElapsed: 63, Start: 2024-11-04, End: 2025-01-06}",
					"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: jupiter: longitude = 73.215463, zodpos = {gemini 13 12}, house = 9th, Journey: 0.58, DaysElapsed: 379, Start: 2024-05-26, End: 2025-06-09}",
					"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: saturn: longitude = 344.524071, zodpos = {pisces 14 31}, house = 6th, Journey: 0.82, DaysElapsed: 809, Start: 2023-03-08, End: 2025-05-25}",
					"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: uranus: longitude = 53.635827, zodpos = {taurus 23 38}, house = 8th, Journey: 0.92, DaysElapsed: 2308, Start: 2019-03-10, End: 2025-07-04}",
					"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: neptune: longitude = 357.297825, zodpos = {pisces 27 17}, house = 6th, Journey: 0.98, DaysElapsed: 4800, Start: 2012-02-07, End: 2025-03-30}",
					"TransitIngress{Date: 2025-01-01 00:00:00 +0000 UTC, P: pluto: longitude = 301.064802, zodpos = {aquarius 1 3}, house = 5th, Journey: 0.01, DaysElapsed: 6680, Start: 2024-11-22, End: 2043-03-08}",
				},
				"2025-01-02": []string{
					"TransitAspect{Date: 2025-01-02 00:00:00 +0000 UTC, Aspect: Aspect{P1: sun, P2: saturn, Degree: 62.783333, Type: Sextile}, Journey: 0.04, DaysElapsed: 6, Start: 2025-01-01, End: 2025-01-08}",
					"TransitAspect{Date: 2025-01-02 00:00:00 +0000 UTC, Aspect: Aspect{P1: mars, P2: neptune, Degree: 124.283333, Type: Trine}, Journey: 0.08, DaysElapsed: 25, Start: 2024-12-31, End: 2025-01-25}",
					"TransitAspect{Date: 2025-01-02 00:00:00 +0000 UTC, Aspect: Aspect{P1: mars, P2: pluto, Degree: 179.500000, Type: Opposition}, Journey: 0.56, DaysElapsed: 30, Start: 2024-12-16, End: 2025-01-15}",
					"TransitAspect{Date: 2025-01-02 00:00:00 +0000 UTC, Aspect: Aspect{P1: jupiter, P2: saturn, Degree: 88.500000, Type: Square}, Journey: 0.63, DaysElapsed: 57, Start: 2024-11-27, End: 2025-01-23}",
					"TransitAspect{Date: 2025-01-02 00:00:00 +0000 UTC, Aspect: Aspect{P1: saturn, P2: sun, Degree: 62.783333, Type: Sextile}, Journey: 0.04, DaysElapsed: 6, Start: 2025-01-01, End: 2025-01-08}",
					"TransitAspect{Date: 2025-01-02 00:00:00 +0000 UTC, Aspect: Aspect{P1: saturn, P2: jupiter, Degree: 88.500000, Type: Square}, Journey: 0.63, DaysElapsed: 57, Start: 2024-11-27, End: 2025-01-23}",
					"TransitAspect{Date: 2025-01-02 00:00:00 +0000 UTC, Aspect: Aspect{P1: neptune, P2: mars, Degree: 124.283333, Type: Trine}, Journey: 0.08, DaysElapsed: 25, Start: 2024-12-31, End: 2025-01-25}",
					"TransitAspect{Date: 2025-01-02 00:00:00 +0000 UTC, Aspect: Aspect{P1: pluto, P2: mars, Degree: 179.500000, Type: Opposition}, Journey: 0.56, DaysElapsed: 30, Start: 2024-12-16, End: 2025-01-15}",
					"TransitIngress{Date: 2025-01-02 00:00:00 +0000 UTC, P: sun: longitude = 281.833252, zodpos = {capricorn 11 49}, house = 4th, Journey: 0.39, DaysElapsed: 29, Start: 2024-12-21, End: 2025-01-19}",
					"TransitIngress{Date: 2025-01-02 00:00:00 +0000 UTC, P: moon: longitude = 307.455797, zodpos = {aquarius 7 27}, house = 5th, Journey: 0.25, DaysElapsed: 2, Start: 2025-01-01, End: 2025-01-03}",
					"TransitIngress{Date: 2025-01-02 00:00:00 +0000 UTC, P: mercury: longitude = 261.168927, zodpos = {sagittarius 21 10}, house = 3rd, Journey: 0.91, DaysElapsed: 66, Start: 2024-11-03, End: 2025-01-08}",
					"TransitIngress{Date: 2025-01-02 00:00:00 +0000 UTC, P: venus: longitude = 328.783260, zodpos = {aquarius 28 46}, house = 5th, Journey: 0.96, DaysElapsed: 26, Start: 2024-12-07, End: 2025-01-03}",
					"TransitIngress{Date: 2025-01-02 00:00:00 +0000 UTC, P: mars: longitude = 121.585043, zodpos = {leo 1 35}, house = 11th, Journey: 0.93, DaysElapsed: 63, Start: 2024-11-04, End: 2025-01-06}",
					"TransitIngress{Date: 2025-01-02 00:00:00 +0000 UTC, P: jupiter: longitude = 73.110170, zodpos = {gemini 13 6}, house = 9th, Journey: 0.58, DaysElapsed: 379, Start: 2024-05-26, End: 2025-06-09}",
					"TransitIngress{Date: 2025-01-02 00:00:00 +0000 UTC, P: saturn: longitude = 344.600359, zodpos = {pisces 14 36}, house = 6th, Journey: 0.82, DaysElapsed: 809, Start: 2023-03-08, End: 2025-05-25}",
					"TransitIngress{Date: 2025-01-02 00:00:00 +0000 UTC, P: uranus: longitude = 53.611804, zodpos = {taurus 23 36}, house = 8th, Journey: 0.92, DaysElapsed: 2312, Start: 2019-03-07, End: 2025-07-05}",
					"TransitIngress{Date: 2025-01-02 00:00:00 +0000 UTC, P: neptune: longitude = 357.311901, zodpos = {pisces 27 18}, house = 6th, Journey: 0.98, DaysElapsed: 4800, Start: 2012-02-04, End: 2025-03-27}",
					"TransitIngress{Date: 2025-01-02 00:00:00 +0000 UTC, P: pluto: longitude = 301.095569, zodpos = {aquarius 1 5}, house = 5th, Journey: 0.01, DaysElapsed: 6680, Start: 2024-11-23, End: 2043-03-09}",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(swe, tt.start, tt.end)
			// spew.Dump(got)
			assert.NoError(t, err)
			m := make(map[string][]string)
			for k, v := range got {
				m[k] = []string{}
				for _, s := range v {
					m[k] = append(m[k], s.String())
				}
			}
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
			name:         "1",
			targetP1:     pointid.Pluto,
			targetP2:     pointid.Sun,
			targetTime:   time.Date(2025, 1, 24, 16, 49, 0, 0, time.UTC),
			wantDuration: 240 * time.Hour,
			wantJourney:  0.825,
		},
		testCase{
			name:       "2",
			targetP1:   pointid.Venus,
			targetP2:   pointid.Neptune,
			targetTime: time.Date(2025, 1, 26, 16, 49, 0, 0, time.UTC),
			// This is a retrograde of Venus in Pisces/Aries where Neptune is
			// copresent
			wantDuration: 2496 * time.Hour,
			wantJourney:  0.002404,
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
			assert.Equal(t, tt.wantDuration, gotDuration)
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
			targetTime:    time.Date(2025, 1, 24, 16, 49, 0, 0, time.UTC),
			wantDuration:  768 * time.Hour,
			wantJourney:   0.67,
		},
		testCase{
			name:          "2",
			targetPointID: pointid.Mercury,
			targetTime:    time.Date(2025, 1, 24, 16, 49, 0, 0, time.UTC),
			wantDuration:  468 * time.Hour,
			wantJourney:   0.83,
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
			assert.Equal(t, tt.wantDuration, gotDuration)
			assert.InEpsilon(
				t,
				tt.wantJourney,
				gotJourney,
				0.01,
				"actual: %.2f, want: %.2f",
				gotJourney,
				tt.wantJourney,
			)
		})
	}
}
