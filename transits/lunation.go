package transits

import (
	"fmt"
	"time"

	"github.com/afjoseph/sacredstar/lunation"
	"github.com/afjoseph/sacredstar/unixtime"
)

type LunationType string

type TransitLunation struct {
	transitBase
	Lunation *lunation.Lunation
}

func (t *TransitLunation) String() string {
	return fmt.Sprintf(
		"TransitLunation{Date: %s, Lunation: %s, Journey: %.2f, DaysElapsed: %d, Start: %s, End: %s}",
		t.transitBase.Date.Format("2006-01-02"),
		t.Lunation,
		t.transitBase.Journey,
		t.transitBase.DaysElapsed,
		t.transitBase.Start.Format("2006-01-02"),
		t.transitBase.End.Format("2006-01-02"),
	)
}

func (t *TransitLunation) GetType() TransitType {
	return TransitTypeLunation
}

func (t *TransitLunation) GetJourney() float64 {
	return t.transitBase.Journey
}

func (t *TransitLunation) GetDuration() int {
	return t.transitBase.DaysElapsed
}

func (t *TransitLunation) GetStart() unixtime.UnixTime {
	return t.transitBase.Start
}

func (t *TransitLunation) GetEnd() unixtime.UnixTime {
	return t.transitBase.End
}

func newTransitLunation(
	lunation *lunation.Lunation,
	t time.Time,
) *TransitLunation {
	return &TransitLunation{
		transitBase: transitBase{
			Type:        TransitTypeLunation,
			Date:        unixtime.New(t),
			Journey:     0,
			DaysElapsed: 0,
			Start:       unixtime.New(t),
			End:         unixtime.New(t),
		},
		Lunation: lunation,
	}
}
