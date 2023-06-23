package wrapper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeConversion(t *testing.T) {
	swe := NewWithBuiltinPath()
	defer swe.Close()
	tm := time.Date(1984, 9, 15, 16, 20, 0, 0, time.UTC)
	jd := swe.GoTimeToJulianDay(tm)
	tm2 := swe.JulianDayToGoTime(jd)
	assert.Equal(t, tm, tm2)
}
