package chart

import (
	"fmt"
	"time"

	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/wrapper"
)

func SanityCheck(swe *wrapper.SwissEph) error {
	// Calculate a random chart just to see if things work
	bt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := NewChartFromUTC(
		swe,
		bt,
		-0.1278, 51.5074, // London
		TropicalChartType,
		pointid.ClassicalPlanets,
	)
	if err != nil {
		return fmt.Errorf("error calculating chart: %w", err)
	}
	return nil
}
