package timeandzone

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/stretchr/testify/assert"
)

func TestTimeAndZoneJSONSerialization(t *testing.T) {
	tm := New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	b, err := json.Marshal(tm)
	assert.NoError(t, err)
	assert.Equal(t, `{"t":"2020-01-01T00:00","z":"UTC"}`, string(b))
	var stm2 TimeAndZone
	err = json.Unmarshal(b, &stm2)
	assert.NoError(t, err)
	assert.Equal(t, tm, stm2)
}

func TestTimeConversion(t *testing.T) {
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	tm := time.Date(1984, 9, 15, 16, 20, 0, 0, time.UTC)
	jd := swe.GoTimeToJulianDay(tm)
	tm2 := swe.JulianDayToGoTime(jd)
	assert.Equal(t, tm, tm2)
	tm3 := New(tm)
	assert.Equal(t, tm, tm3.Time)
}
