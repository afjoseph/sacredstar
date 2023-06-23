package chart

import (
	"testing"

	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/sign"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/afjoseph/sacredstar/zodiacalpos"
	"github.com/stretchr/testify/assert"
)

func TestCorrectJulianDateToZodiacalPos(t *testing.T) {
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	type TestCase struct {
		title          string
		initialJulDay  float64
		pid            pointid.PointID
		targetZodPos   *zodiacalpos.ZodiacalPos
		targetNakType  NakshatraType
		step           float64
		_range         float64
		expectedReturn float64
	}

	for _, tc := range []TestCase{
		TestCase{
			title:          "test1",
			initialJulDay:  2455550.9396453705,
			pid:            pointid.Moon,
			targetZodPos:   zodiacalpos.NewZodiacalPos(sign.Taurus, 23, 20),
			targetNakType:  Mrigashirsha,
			step:           (1.0 / 60.0 / 24.0), // 1 minute
			_range:         1.0,                 // 1 day
			expectedReturn: 2455550.9611446653,
		},
		TestCase{
			title:          "test2",
			initialJulDay:  2455551.9517194447,
			pid:            pointid.Moon,
			targetZodPos:   zodiacalpos.NewZodiacalPos(sign.Gemini, 6, 40),
			targetNakType:  Mrigashirsha,
			step:           (1.0 / 60.0 / 24.0), // 1 minute
			_range:         1.0,                 // 1 day
			expectedReturn: 2455551.939015,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			actual, err := correctJulianDateToZodiacalPos(
				swe,
				tc.initialJulDay,
				tc.pid,
				tc.targetZodPos,
				tc.targetNakType,
				tc.step,
				tc._range,
			)
			assert.NoError(t, err)
			assert.InDelta(t, tc.expectedReturn, actual, 0.001)
		})
	}
}
