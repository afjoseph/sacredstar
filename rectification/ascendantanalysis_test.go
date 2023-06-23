package rectification

import (
	"testing"
	"time"

	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/sign"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzeAscendantsForTimeInterval(t *testing.T) {
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()
	type testCase struct {
		name     string
		events   []Event
		inter    chart.Interval
		lon, lat float64
		want     []*AscendantInterval
	}

	tests := []testCase{
		testCase{
			name: "1",
			events: []Event{
				// D4: Moving home
				Event{EventType: EventType_MovingHome},
			},
			inter: chart.NewInterval(
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
			),
			// London coordinates
			lon: -0.1276,
			lat: 51.5072,
			want: []*AscendantInterval{
				&AscendantInterval{
					ChartType: chart.TropicalChartType,
					Sign:      sign.Libra,
					Interval: chart.NewInterval(
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
					),
				},
				&AscendantInterval{
					ChartType: chart.D1ChartType,
					Sign:      sign.Virgo,
					Interval: chart.NewInterval(
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
					),
				},
				&AscendantInterval{
					ChartType: chart.D4ChartType,
					Sign:      sign.Sagittarius,
					Interval: chart.NewInterval(
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 0, 12, 0, 0, time.UTC),
					),
				},
				&AscendantInterval{
					ChartType: chart.D4ChartType,
					Sign:      sign.Pisces,
					Interval: chart.NewInterval(
						time.Date(2024, 1, 1, 0, 13, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 0, 54, 0, 0, time.UTC),
					),
				},
				&AscendantInterval{
					ChartType: chart.D4ChartType,
					Sign:      sign.Gemini,
					Interval: chart.NewInterval(
						time.Date(2024, 1, 1, 0, 55, 0, 0, time.UTC),
						time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
					),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := analyzeAscendantsForTimeInterval(
				swe,
				tt.events,
				tt.inter,
				tt.lon,
				tt.lat,
			)
			assert.NoError(t, err)

			// Sort the slices to make comparison easier
			sortAscendantIntervals(got)
			sortAscendantIntervals(tt.want)
			// Compare the slices
			for i, a := range got {
				b := tt.want[i]
				assert.Equal(t, a.ChartType, b.ChartType)
				assert.Equal(t, a.Sign, b.Sign)
				assert.Equal(t, a.Interval.Start.Time, b.Interval.Start.Time)
				assert.Equal(t, a.Interval.End.Time, b.Interval.End.Time)
			}
		})
	}
}
