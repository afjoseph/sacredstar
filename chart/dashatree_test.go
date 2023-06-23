package chart

import (
	"fmt"
	"testing"
	"time"

	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/stretchr/testify/assert"
)

func TestNewDashaTree(t *testing.T) {
	type testCase struct {
		birthTime              time.Time
		testTime               time.Time
		expectedMahaDashaLord  DashaLord
		expectedAntarDashaLord DashaLord
	}
	tcs := []testCase{
		testCase{
			birthTime:              time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			testTime:               time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			expectedMahaDashaLord:  DashaLordKetu,
			expectedAntarDashaLord: DashaLordMercury,
		},
		testCase{
			birthTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			testTime: time.Date(
				2024,
				10,
				22,
				0,
				0,
				0,
				0,
				time.UTC,
			),
			expectedMahaDashaLord:  DashaLordVenus,
			expectedAntarDashaLord: DashaLordVenus,
		},
	}

	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	for _, tc := range tcs {
		t.Run(
			fmt.Sprintf("%s --> %s", tc.birthTime, tc.testTime),
			func(t *testing.T) {
				moonAstroPoint, err := calculatePlanet(
					swe,
					swe.GoTimeToJulianDay(tc.birthTime),
					pointid.Moon,
					D1ChartType,
					nil,
				)
				assert.NoError(t, err)
				assert.NotNil(t, moonAstroPoint)

				dt, err := NewDashaTree(swe, tc.birthTime, moonAstroPoint)
				assert.NoError(t, err)
				assert.NotNil(t, dt)
				// dt.Dump()

				ret, didFind := dt.GetDashaForTime(tc.testTime)
				assert.True(t, didFind)
				assert.Equal(
					t,
					tc.expectedMahaDashaLord,
					ret.Mahadasha,
					"Expected %s, got %s",
					tc.expectedMahaDashaLord,
					ret.Mahadasha,
				)
				assert.Equal(
					t,
					tc.expectedAntarDashaLord,
					ret.Antardasha,
					"Expected %s, got %s",
					tc.expectedAntarDashaLord,
					ret.Antardasha,
				)
			},
		)
	}
}
