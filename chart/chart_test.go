package chart

import (
	"testing"
	"time"

	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/sign"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/afjoseph/sacredstar/zodiacalpos"
	"github.com/stretchr/testify/assert"
)

func TestNewChart(t *testing.T) {
	type testcase struct {
		title                     string
		inputDate                 time.Time
		lon                       float64
		lat                       float64
		ChartType                 ChartType
		expectedZodiacalPositions map[pointid.PointID]*zodiacalpos.ZodiacalPos
	}

	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	for _, tc := range []testcase{
		{
			title:     "2024-01-01 Tropical",
			inputDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			ChartType: TropicalChartType,
			// London coordinates
			lon: -0.1278,
			lat: 51.5074,
			expectedZodiacalPositions: map[pointid.PointID]*zodiacalpos.ZodiacalPos{
				pointid.ASC:     zodiacalpos.NewZodiacalPos(sign.Libra, 7, 03),
				pointid.Sun:     zodiacalpos.NewZodiacalPos(sign.Capricorn, 10, 02),
				pointid.Moon:    zodiacalpos.NewZodiacalPos(sign.Virgo, 5, 59),
				pointid.Mercury: zodiacalpos.NewZodiacalPos(sign.Sagittarius, 22, 16),
				pointid.Venus:   zodiacalpos.NewZodiacalPos(sign.Sagittarius, 2, 36),
				pointid.Mars:    zodiacalpos.NewZodiacalPos(sign.Sagittarius, 27, 18),
				pointid.Jupiter: zodiacalpos.NewZodiacalPos(sign.Taurus, 5, 34),
				pointid.Saturn:  zodiacalpos.NewZodiacalPos(sign.Pisces, 3, 14),
				pointid.Rahu:    zodiacalpos.NewZodiacalPos(sign.Aries, 21, 04),
				pointid.Ketu:    zodiacalpos.NewZodiacalPos(sign.Libra, 21, 3),
			},
		},
		{
			title:     "2024-01-01 D1",
			inputDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			ChartType: D1ChartType,
			// London coordinates
			lon: -0.1278,
			lat: 51.5074,
			expectedZodiacalPositions: map[pointid.PointID]*zodiacalpos.ZodiacalPos{
				pointid.ASC:     zodiacalpos.NewZodiacalPos(sign.Virgo, 12, 52),
				pointid.Sun:     zodiacalpos.NewZodiacalPos(sign.Sagittarius, 15, 50),
				pointid.Moon:    zodiacalpos.NewZodiacalPos(sign.Leo, 11, 48),
				pointid.Mars:    zodiacalpos.NewZodiacalPos(sign.Sagittarius, 3, 07),
				pointid.Mercury: zodiacalpos.NewZodiacalPos(sign.Scorpio, 28, 05),
				pointid.Jupiter: zodiacalpos.NewZodiacalPos(sign.Aries, 11, 23),
				pointid.Venus:   zodiacalpos.NewZodiacalPos(sign.Scorpio, 8, 25),
				pointid.Saturn:  zodiacalpos.NewZodiacalPos(sign.Aquarius, 9, 03),
				pointid.Rahu:    zodiacalpos.NewZodiacalPos(sign.Pisces, 26, 53),
				pointid.Ketu:    zodiacalpos.NewZodiacalPos(sign.Virgo, 26, 53),
			},
		},
		{
			title:     "1992-06-13 D1",
			inputDate: time.Date(1992, 6, 13, 4, 40, 0, 0, time.UTC),
			ChartType: D1ChartType,
			// Syria coordinates
			lon: 36.3,
			lat: 33.5,
			expectedZodiacalPositions: map[pointid.PointID]*zodiacalpos.ZodiacalPos{
				pointid.ASC:     zodiacalpos.NewZodiacalPos(sign.Gemini, 27, 49),
				pointid.Sun:     zodiacalpos.NewZodiacalPos(sign.Taurus, 28, 39),
				pointid.Moon:    zodiacalpos.NewZodiacalPos(sign.Scorpio, 5, 16),
				pointid.Mercury: zodiacalpos.NewZodiacalPos(sign.Gemini, 13, 03),
				pointid.Venus:   zodiacalpos.NewZodiacalPos(sign.Taurus, 28, 31),
				pointid.Mars:    zodiacalpos.NewZodiacalPos(sign.Aries, 5, 9),
				pointid.Jupiter: zodiacalpos.NewZodiacalPos(sign.Leo, 13, 34),
				pointid.Saturn:  zodiacalpos.NewZodiacalPos(sign.Capricorn, 24, 32),
				pointid.Rahu:    zodiacalpos.NewZodiacalPos(sign.Sagittarius, 6, 54),
				pointid.Ketu:    zodiacalpos.NewZodiacalPos(sign.Gemini, 6, 53),
			},
		},
		{
			title:     "2024-01-01 D4",
			inputDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			ChartType: D4ChartType,
			// London coordinates
			lon: -0.1278,
			lat: 51.5074,
			expectedZodiacalPositions: map[pointid.PointID]*zodiacalpos.ZodiacalPos{
				pointid.ASC:     zodiacalpos.NewZodiacalPos(sign.Sagittarius, 21, 29),
				pointid.Sun:     zodiacalpos.NewZodiacalPos(sign.Gemini, 3, 23),
				pointid.Moon:    zodiacalpos.NewZodiacalPos(sign.Scorpio, 17, 12),
				pointid.Mercury: zodiacalpos.NewZodiacalPos(sign.Leo, 22, 21),
				pointid.Venus:   zodiacalpos.NewZodiacalPos(sign.Aquarius, 3, 41),
				pointid.Mars:    zodiacalpos.NewZodiacalPos(sign.Sagittarius, 12, 28),
				pointid.Jupiter: zodiacalpos.NewZodiacalPos(sign.Cancer, 15, 34),
				pointid.Saturn:  zodiacalpos.NewZodiacalPos(sign.Taurus, 6, 12),
				pointid.Rahu:    zodiacalpos.NewZodiacalPos(sign.Sagittarius, 17, 32),
				pointid.Ketu:    zodiacalpos.NewZodiacalPos(sign.Gemini, 17, 30),
			},
		},
		{
			title:     "2024-01-01 D7",
			inputDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			ChartType: D7ChartType,
			// London coordinates
			lon: -0.1278,
			lat: 51.5074,
			expectedZodiacalPositions: map[pointid.PointID]*zodiacalpos.ZodiacalPos{
				pointid.ASC:     zodiacalpos.NewZodiacalPos(sign.Gemini, 0, 06),
				pointid.Sun:     zodiacalpos.NewZodiacalPos(sign.Pisces, 20, 56),
				pointid.Moon:    zodiacalpos.NewZodiacalPos(sign.Libra, 22, 36),
				pointid.Mercury: zodiacalpos.NewZodiacalPos(sign.Scorpio, 16, 38),
				pointid.Venus:   zodiacalpos.NewZodiacalPos(sign.Gemini, 28, 57),
				pointid.Mars:    zodiacalpos.NewZodiacalPos(sign.Sagittarius, 21, 49),
				pointid.Jupiter: zodiacalpos.NewZodiacalPos(sign.Gemini, 19, 44),
				pointid.Saturn:  zodiacalpos.NewZodiacalPos(sign.Aries, 3, 22),
				pointid.Rahu:    zodiacalpos.NewZodiacalPos(sign.Pisces, 8, 12),
				pointid.Ketu:    zodiacalpos.NewZodiacalPos(sign.Virgo, 8, 11),
			},
		},
		{
			title:     "2024-01-01 D9",
			inputDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			ChartType: D9ChartType,
			// London coordinates
			lon: -0.1278,
			lat: 51.5074,
			expectedZodiacalPositions: map[pointid.PointID]*zodiacalpos.ZodiacalPos{
				pointid.ASC:     zodiacalpos.NewZodiacalPos(sign.Aries, 25, 51),
				pointid.Sun:     zodiacalpos.NewZodiacalPos(sign.Leo, 22, 38),
				pointid.Moon:    zodiacalpos.NewZodiacalPos(sign.Cancer, 16, 13),
				pointid.Mercury: zodiacalpos.NewZodiacalPos(sign.Pisces, 12, 49),
				pointid.Venus:   zodiacalpos.NewZodiacalPos(sign.Virgo, 15, 47),
				pointid.Mars:    zodiacalpos.NewZodiacalPos(sign.Aries, 28, 03),
				pointid.Jupiter: zodiacalpos.NewZodiacalPos(sign.Cancer, 12, 31),
				pointid.Saturn:  zodiacalpos.NewZodiacalPos(sign.Sagittarius, 21, 28),
				pointid.Rahu:    zodiacalpos.NewZodiacalPos(sign.Pisces, 1, 58),
				pointid.Ketu:    zodiacalpos.NewZodiacalPos(sign.Virgo, 1, 56),
			},
		},
		{
			title:     "2024-01-01 D10",
			inputDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			ChartType: D10ChartType,
			// London coordinates
			lon: -0.1278,
			lat: 51.5074,
			expectedZodiacalPositions: map[pointid.PointID]*zodiacalpos.ZodiacalPos{
				pointid.ASC:     zodiacalpos.NewZodiacalPos(sign.Virgo, 8, 44),
				pointid.Sun:     zodiacalpos.NewZodiacalPos(sign.Taurus, 8, 29),
				pointid.Moon:    zodiacalpos.NewZodiacalPos(sign.Scorpio, 28, 01),
				pointid.Mercury: zodiacalpos.NewZodiacalPos(sign.Aries, 10, 54),
				pointid.Venus:   zodiacalpos.NewZodiacalPos(sign.Virgo, 24, 13),
				pointid.Mars:    zodiacalpos.NewZodiacalPos(sign.Capricorn, 1, 10),
				pointid.Jupiter: zodiacalpos.NewZodiacalPos(sign.Cancer, 23, 55),
				pointid.Saturn:  zodiacalpos.NewZodiacalPos(sign.Taurus, 0, 31),
				pointid.Rahu:    zodiacalpos.NewZodiacalPos(sign.Cancer, 28, 52),
				pointid.Ketu:    zodiacalpos.NewZodiacalPos(sign.Capricorn, 28, 48),
			},
		},
		// {
		// 	title:                "2024-01-01 D24",
		// 	inputDate:            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		// 	ChartType: D24ChartType,
		// 	// London coordinates
		// 	lon: -0.1278,
		// 	lat: 51.5074,
		// 	expectedZodiacalPositions: map[pointid.PointID]*zodiacalpos.ZodiacalPos{
		// 		pointid.ASC:     zodiacalpos.NewZodiacalPos(Taurus, 8, 57),
		// 		pointid.Sun:     zodiacalpos.NewZodiacalPos(Leo, 20, 22),
		// 		pointid.Moon:    zodiacalpos.NewZodiacalPos(Taurus, 13, 15),
		// 		pointid.Mercury: zodiacalpos.NewZodiacalPos(Taurus, 14, 11),
		// 		pointid.Venus:   zodiacalpos.NewZodiacalPos(Capricorn, 22, 07),
		// 		pointid.Mars:    zodiacalpos.NewZodiacalPos(Libra, 14, 49),
		// 		pointid.Jupiter: zodiacalpos.NewZodiacalPos(Taurus, 3, 24),
		// 		pointid.Saturn:  zodiacalpos.NewZodiacalPos(Pisces, 7, 16),
		// 		pointid.Rahu:    zodiacalpos.NewZodiacalPos(Aries, 15, 17),
		// 		pointid.Ketu:    zodiacalpos.NewZodiacalPos(Aries, 1, 1),
		// 	},
		// },
	} {
		t.Run(tc.title, func(t *testing.T) {
			c, err := NewChartFromUTC(
				swe,
				tc.inputDate,
				tc.lon, tc.lat,
				tc.ChartType,
				pointid.VedicPlanets,
			)
			assert.NoError(t, err)
			assert.NotNil(t, c)
			// fmt.Printf("Chart: %s\n", c)
			for pid, expectedZodPos := range tc.expectedZodiacalPositions {
				p := c.GetPoint(pid)
				assert.NotNil(t, p, "expected point %s to exist in chart", pid)
				zodPos := p.ZodiacalPos
				assert.True(
					t,
					zodPos.DiffInAbsDegrees(expectedZodPos) < 0.2,
					"For %s: expected generated zodPos {%s} to be equal to expected zodPos: {%s}",
					pid,
					zodPos,
					expectedZodPos,
				)
			}
		})
	}
}

func TestNakshatra(t *testing.T) {
	// Calculate different nakshatras for different birth times
	// and compare with expected results
	swe := wrapper.NewWithBuiltinPath()
	defer swe.Close()

	for _, tc := range []struct {
		inputTime             time.Time
		expectedNakshatraType NakshatraType
	}{
		{
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Magha,
		},
		{
			time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			Magha + 1,
		},
		{
			time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			Magha + 2,
		},
		{
			time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
			Magha + 3,
		},
		{
			time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			Magha + 4,
		},
		{
			time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
			Magha + 5,
		},
		{
			time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
			Magha + 6,
		},
		{
			time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			Magha + 7,
		},
		{
			time.Date(2024, 1, 9, 0, 0, 0, 0, time.UTC),
			Magha + 8,
		},
		{
			time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Magha + 9,
		},
		{
			time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			Magha + 10,
		},
		{
			time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			Magha + 11,
		},
		{
			time.Date(2024, 1, 13, 0, 0, 0, 0, time.UTC),
			Magha + 12,
		},
		{
			time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC),
			Magha + 13,
		},
		{
			time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Magha + 14,
		},
		{
			time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
			Magha + 15,
		},
		{
			time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC),
			Magha + 17,
		},
		{
			time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC),
			Aswini,
		},
		{
			time.Date(2024, 1, 19, 0, 0, 0, 0, time.UTC),
			Aswini + 1,
		},
		{
			time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			Aswini + 2,
		},
		{
			time.Date(2024, 1, 21, 0, 0, 0, 0, time.UTC),
			Aswini + 3,
		},
		{
			time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC),
			Aswini + 4,
		},
		{
			time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC),
			Aswini + 5,
		},
		{
			time.Date(2024, 1, 24, 0, 0, 0, 0, time.UTC),
			Aswini + 5,
		},
		{
			time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			Aswini + 6,
		},
		{
			time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC),
			Aswini + 7,
		},
		{
			time.Date(2024, 1, 27, 0, 0, 0, 0, time.UTC),
			Aswini + 8,
		},
		{
			time.Date(2024, 1, 28, 0, 0, 0, 0, time.UTC),
			Magha,
		},
	} {
		ap, err := calculatePlanet(
			swe,
			swe.GoTimeToJulianDay(tc.inputTime),
			pointid.Moon,
			D1ChartType,
			nil,
		)
		assert.NoError(t, err)
		assert.NotNil(t, ap)

		n, err := NewNakshatraFromChart(swe, ap)
		assert.NoError(t, err)
		// fmt.Printf("For %s, nakshatra = %s, pada = %s\n", tc.inputTime, n, p)
		assert.Equal(
			t,
			tc.expectedNakshatraType,
			n.Type,
			"for %s",
			tc.inputTime,
		)
	}
}
