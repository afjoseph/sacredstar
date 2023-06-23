package chart

import (
	"fmt"

	"github.com/hako/durafmt"
)

type Dasha struct {
	Mahadasha  DashaLord `json:"mahadasha"`
	Antardasha DashaLord `json:"antardasha"`
	Interval   Interval  `json:"interval"`
}

func (vd Dasha) String() string {
	return fmt.Sprintf(
		"Dasha{Mahadasha: %s, Antardasha: %s, Interval: %+v, Duration: %+v}",
		vd.Mahadasha,
		vd.Antardasha,
		vd.Interval,
		durafmt.Parse(vd.Interval.End.Sub(vd.Interval.Start.Time)),
	)
}
