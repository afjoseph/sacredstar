package transits

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

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
				"TransitIngress{Date: 2025-01-01, P: AstroPoint{ID: sun, Longitude: 280.813611, ZodiacalPos: capricorn 10 48, House: 4th, IsRetrograde: false}, Journey: 0.36, DaysElapsed: 29, Start: 2024-12-21, End: 2025-01-19}",
				"TransitIngress{Date: 2025-01-01, P: AstroPoint{ID: moon, Longitude: 293.913579, ZodiacalPos: capricorn 23 54, House: 4th, IsRetrograde: false}, Journey: 0.75, DaysElapsed: 111, Start: 2024-10-09, End: 2025-01-28}",
				"TransitIngress{Date: 2025-01-01, P: AstroPoint{ID: mercury, Longitude: 259.869978, ZodiacalPos: sagittarius 19 52, House: 3rd, IsRetrograde: false}, Journey: 0.89, DaysElapsed: 66, Start: 2024-11-02, End: 2025-01-08}",
				"TransitIngress{Date: 2025-01-01, P: AstroPoint{ID: venus, Longitude: 327.712110, ZodiacalPos: aquarius 27 42, House: 5th, IsRetrograde: false}, Journey: 0.92, DaysElapsed: 26, Start: 2024-12-07, End: 2025-01-03}",
				"TransitIngress{Date: 2025-01-01, P: AstroPoint{ID: mars, Longitude: 121.917923, ZodiacalPos: leo 1 55, House: 11th, IsRetrograde: true}, Journey: 0.26, DaysElapsed: 225, Start: 2024-11-04, End: 2025-06-17}",
				"TransitIngress{Date: 2025-01-01, P: AstroPoint{ID: jupiter, Longitude: 73.215463, ZodiacalPos: gemini 13 12, House: 9th, IsRetrograde: true}, Journey: 0.58, DaysElapsed: 380, Start: 2024-05-25, End: 2025-06-10}",
				"TransitIngress{Date: 2025-01-01, P: AstroPoint{ID: saturn, Longitude: 344.524071, ZodiacalPos: pisces 14 31, House: 6th, IsRetrograde: false}, Journey: 0.82, DaysElapsed: 810, Start: 2023-03-07, End: 2025-05-25}",
				"TransitIngress{Date: 2025-01-01, P: AstroPoint{ID: uranus, Longitude: 53.635827, ZodiacalPos: taurus 23 38, House: 8th, IsRetrograde: true}, Journey: 0.82, DaysElapsed: 2607, Start: 2019-03-07, End: 2026-04-26}",
				"TransitIngress{Date: 2025-01-01, P: AstroPoint{ID: neptune, Longitude: 357.297825, ZodiacalPos: pisces 27 17, House: 6th, IsRetrograde: false}, Journey: 0.98, DaysElapsed: 4803, Start: 2012-02-05, End: 2025-04-01}",
				"TransitIngress{Date: 2025-01-01, P: AstroPoint{ID: pluto, Longitude: 301.064802, ZodiacalPos: aquarius 1 3, House: 5th, IsRetrograde: false}, Journey: 0.05, DaysElapsed: 6986, Start: 2024-01-23, End: 2043-03-11}",
				"TransitAspect{Date: 2025-01-01, Aspect: Aspect{P1: moon, P2: uranus, Degree: 119.733333, Type: Trine}, Journey: 0.54, DaysElapsed: 0, Start: 2024-12-31, End: 2025-01-01}",
				"TransitAspect{Date: 2025-01-01, Aspect: Aspect{P1: venus, P2: uranus, Degree: 85.933333, Type: Square}, Journey: 0.90, DaysElapsed: 8, Start: 2024-12-24, End: 2025-01-01}",
				"TransitAspect{Date: 2025-01-01, Aspect: Aspect{P1: mars, P2: neptune, Degree: 124.633333, Type: Trine}, Journey: 0.75, DaysElapsed: 99, Start: 2024-10-17, End: 2025-01-25}",
				"TransitAspect{Date: 2025-01-01, Aspect: Aspect{P1: mars, P2: pluto, Degree: 179.133333, Type: Opposition}, Journey: 0.83, DaysElapsed: 86, Start: 2024-10-21, End: 2025-01-15}",
				"TransitAspect{Date: 2025-01-01, Aspect: Aspect{P1: jupiter, P2: saturn, Degree: 88.683333, Type: Square}, Journey: 0.45, DaysElapsed: 348, Start: 2024-07-27, End: 2025-07-11}",
				"TransitAspect{Date: 2025-01-01, Aspect: Aspect{P1: saturn, P2: jupiter, Degree: 88.683333, Type: Square}, Journey: 0.45, DaysElapsed: 348, Start: 2024-07-27, End: 2025-07-11}",
				"TransitAspect{Date: 2025-01-01, Aspect: Aspect{P1: uranus, P2: moon, Degree: 119.733333, Type: Trine}, Journey: 0.54, DaysElapsed: 0, Start: 2024-12-31, End: 2025-01-01}",
				"TransitAspect{Date: 2025-01-01, Aspect: Aspect{P1: uranus, P2: venus, Degree: 85.933333, Type: Square}, Journey: 0.90, DaysElapsed: 8, Start: 2024-12-24, End: 2025-01-01}",
				"TransitAspect{Date: 2025-01-01, Aspect: Aspect{P1: neptune, P2: mars, Degree: 124.633333, Type: Trine}, Journey: 0.75, DaysElapsed: 99, Start: 2024-10-17, End: 2025-01-25}",
				"TransitAspect{Date: 2025-01-01, Aspect: Aspect{P1: pluto, P2: mars, Degree: 179.133333, Type: Opposition}, Journey: 0.83, DaysElapsed: 86, Start: 2024-10-21, End: 2025-01-15}",
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

func TestMarshalJSON_Transits(t *testing.T) {
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	tss, err := New(swe, time.Now())
	assert.NoError(t, err)

	b, err := json.Marshal(tss)
	assert.NoError(t, err)

	var tss2 Transits
	err = json.Unmarshal(b, &tss2)
	assert.NoError(t, err)
}
