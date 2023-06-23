package rectification

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/sign"
	"github.com/afjoseph/sacredstar/wrapper"
)

// AscendantInterval details the properties of an ascendant during a time
// interval
type AscendantInterval struct {
	ChartType chart.ChartType `json:"chartType"`
	Sign      sign.Sign       `json:"sign"`
	Interval  chart.Interval  `json:"interval"`
}

func analyzeAscendantsForTimeInterval(
	swe *wrapper.SwissEph,
	events []Event,
	inter chart.Interval,
	lon, lat float64,
) ([]*AscendantInterval, error) {
	// - Go through every 10 min of this interval
	m := map[string][]time.Time{}

	chartTypes := []chart.ChartType{
		chart.TropicalChartType,
		chart.D1ChartType,
	}
	for _, e := range events {
		chartTypes = append(chartTypes, e.EventType.VargaChartType())
	}

	// XXX <22-02-2024, afjoseph> We can increase this step since the lowest
	// divisional chart we use is D10 which changes every 12 minutes, but let's
	// keep it lower now for accuracy in our tests
	step := 1 * time.Minute
	for t := inter.Start.Time; t.Before(inter.End.Time) || t.Equal(inter.End.Time); t = t.Add(step) {
		// - Calculate the ascendant for each event type
		for _, ct := range chartTypes {
			timeInJulian := swe.GoTimeToJulianDay(t.UTC())
			asc, err := chart.CalculateAscendant(
				swe,
				timeInJulian,
				lon,
				lat,
				ct,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"while calculating ascendant for chart Type %s at time %v: %w",
					ct,
					t,
					err,
				)
			}

			nodeName := fmt.Sprintf(
				"%s-%s",
				ct,
				asc.ZodiacalPos.Sign.String(),
			)
			m[nodeName] = append(m[nodeName], t)
		}
	}

	ret := []*AscendantInterval{}
	for k, v := range m {
		if len(v) == 0 {
			panic("should not happen")
		}
		split := strings.Split(k, "-")
		chartType := chart.ChartType(split[0])
		sign := sign.Sign(split[1])
		startOfInterval := v[0]
		endOfInterval := v[0]
		for _, t := range v {
			if t.Before(startOfInterval) {
				startOfInterval = t
			}
			if t.After(endOfInterval) {
				endOfInterval = t
			}
		}
		ret = append(ret, &AscendantInterval{
			ChartType: chartType,
			Sign:      sign,
			Interval: chart.NewInterval(
				startOfInterval,
				endOfInterval,
			),
		})
	}

	sortAscendantIntervals(ret)
	return ret, nil
}

// sortAscendantIntervals sorts by chart type and then by start time
func sortAscendantIntervals(in []*AscendantInterval) {
	// Sort by chart type and then by start time
	slices.SortFunc[[]*AscendantInterval](
		in,
		func(a, b *AscendantInterval) int {
			if a.ChartType.Int() < b.ChartType.Int() {
				return -1
			}
			if a.ChartType.Int() > b.ChartType.Int() {
				return 1
			}
			if a.Interval.Start.Time.Before(b.Interval.Start.Time) {
				return -1
			}
			if a.Interval.Start.Time.After(b.Interval.Start.Time) {
				return 1
			}
			return 0
		},
	)
}
