package zodiacalpos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAspect(t *testing.T) {
	type testCase struct {
		title              string
		expectedAspectType AspectType
		expectedDegree     float64
		expectedShouldFind bool
		inputLHS           *ZodiacalPos
		inputRHS           *ZodiacalPos
	}

	for _, tc := range []testCase{
		testCase{
			title:              "conjunction",
			inputLHS:           NewZodiacalPosFromLongitude(0),
			inputRHS:           NewZodiacalPosFromLongitude(0),
			expectedAspectType: AspectType_Conjunction,
			expectedDegree:     0,
			expectedShouldFind: true,
		},
		// Test for opposition (180 degrees apart)
		testCase{
			title:              "opposition",
			inputLHS:           NewZodiacalPosFromLongitude(0),
			inputRHS:           NewZodiacalPosFromLongitude(180),
			expectedAspectType: AspectType_Opposition,
			expectedDegree:     180,
			expectedShouldFind: true,
		},
		// Test for trine (120 degrees apart)
		testCase{
			title:              "trine",
			inputLHS:           NewZodiacalPosFromLongitude(0),
			inputRHS:           NewZodiacalPosFromLongitude(120),
			expectedAspectType: AspectType_Trine,
			expectedDegree:     120,
			expectedShouldFind: true,
		},
		// Test for square (90 degrees apart)
		testCase{
			title:              "square",
			inputLHS:           NewZodiacalPosFromLongitude(0),
			inputRHS:           NewZodiacalPosFromLongitude(90),
			expectedAspectType: AspectType_Square,
			expectedDegree:     90,
			expectedShouldFind: true,
		},
		// Test for sextile (60 degrees apart)
		testCase{
			title:              "sextile",
			inputLHS:           NewZodiacalPosFromLongitude(0),
			inputRHS:           NewZodiacalPosFromLongitude(60),
			expectedAspectType: AspectType_Sextile,
			expectedDegree:     60,
			expectedShouldFind: true,
		},
		// Test for no aspect (not an exact aspect degree)
		testCase{
			title:              "no aspect",
			inputLHS:           NewZodiacalPosFromLongitude(0),
			inputRHS:           NewZodiacalPosFromLongitude(50),
			expectedAspectType: AspectType_None,
			expectedDegree:     0, // or the actual degree difference if that's what NewAspect returns
			expectedShouldFind: false,
		},
		testCase{
			title:              "conjunction at the edge of the zodiac",
			inputLHS:           NewZodiacalPosFromLongitude(359),
			inputRHS:           NewZodiacalPosFromLongitude(1),
			expectedAspectType: AspectType_Conjunction,
			expectedDegree:     2,
			expectedShouldFind: true,
		},
		testCase{
			title:              "opposition across the zodiac edge",
			inputLHS:           NewZodiacalPosFromLongitude(359),
			inputRHS:           NewZodiacalPosFromLongitude(179),
			expectedAspectType: AspectType_Opposition,
			expectedDegree:     180,
			expectedShouldFind: true,
		},
		testCase{
			title:              "square across the zodiac edge",
			inputLHS:           NewZodiacalPosFromLongitude(359),
			inputRHS:           NewZodiacalPosFromLongitude(89),
			expectedAspectType: AspectType_Square,
			expectedDegree:     90,
			expectedShouldFind: true,
		},
		testCase{
			title:              "trine at 0 degrees",
			inputLHS:           NewZodiacalPosFromLongitude(0),
			inputRHS:           NewZodiacalPosFromLongitude(120),
			expectedAspectType: AspectType_Trine,
			expectedDegree:     120,
			expectedShouldFind: true,
		},
		testCase{
			title:              "conjunction at the same edge degree",
			inputLHS:           NewZodiacalPosFromLongitude(360),
			inputRHS:           NewZodiacalPosFromLongitude(0),
			expectedAspectType: AspectType_Conjunction,
			expectedDegree:     0,
			expectedShouldFind: true,
		},
		testCase{
			title:              "test normalization of large degree numbers",
			inputLHS:           NewZodiacalPosFromLongitude(480), // 480 should normalize to 120
			inputRHS:           NewZodiacalPosFromLongitude(600), // 600 should normalize to 240
			expectedAspectType: AspectType_Trine,
			expectedDegree:     120,
			expectedShouldFind: true,
		},
		testCase{
			title:              "test aspect with orb",
			inputLHS:           NewZodiacalPosFromLongitude(90),
			inputRHS:           NewZodiacalPosFromLongitude(150.5), // Slightly more than a trine
			expectedAspectType: AspectType_Sextile,                 // Assuming no aspect due to tight orb
			expectedDegree:     60.5,
			expectedShouldFind: true,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			aspect := NewAspect(tc.inputLHS, tc.inputRHS)
			if !tc.expectedShouldFind {
				assert.Nil(t, aspect)
				return
			}
			assert.NotNil(t, aspect)

			assert.Equal(t, tc.expectedAspectType, aspect.Type)
			assert.Equal(t, tc.expectedDegree, aspect.Degree)
		})
	}
}
