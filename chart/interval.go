package chart

import (
	"time"

	"github.com/afjoseph/sacredstar/timeandzone"
)

type Interval struct {
	Start timeandzone.TimeAndZone `json:"start"`
	End   timeandzone.TimeAndZone `json:"end"`
}

func (it Interval) Checksum() []byte {
	return []byte(it.Start.String() + it.End.String())
}

func NewInterval(start time.Time, end time.Time) Interval {
	return Interval{
		Start: timeandzone.New(start),
		End:   timeandzone.New(end),
	}
}

func (i Interval) IsInBetween(t time.Time) bool {
	return t.After(i.Start.Time) && t.Before(i.End.Time)
}

func (it Interval) String() string {
	return it.Start.String() + " - " + it.End.String()
}
