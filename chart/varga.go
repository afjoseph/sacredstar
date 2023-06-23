package chart

import (
	"math"

	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/sign"
	"github.com/afjoseph/sacredstar/zodiacalpos"
)

func transformZodiacalPosToVarga(
	pid pointid.PointID,
	zp *zodiacalpos.ZodiacalPos,
	chartType ChartType,
) (*zodiacalpos.ZodiacalPos, error) {
	if chartType == D1ChartType {
		// No division in D1 chart
		return zp, nil
	}

	switch chartType {
	case D4ChartType:
		return transformZodiacalPosToD4(pid, zp)
	case D7ChartType:
		return transformZodiacalPosToD7(pid, zp)
	case D9ChartType:
		return transformZodiacalPosToD9(pid, zp)
	case D10ChartType:
		return transformZodiacalPosToD10(pid, zp)
	// case D24ChartType:
	// 	return transformZodiacalPosToD24(pid, zp)
	default:
		panic("unreachable")
	}
}

func transformZodiacalPosToD4(
	pid pointid.PointID,
	zp *zodiacalpos.ZodiacalPos,
) (*zodiacalpos.ZodiacalPos, error) {
	// - So assume Mars is in Sagittarius 3d7m
	//   - that's 3.116 degrees
	//   - and it's in the first divisions
	//     - if you divide the sign into 4 divisions
	//   - 30 / 4 = 7.5 for each divisions
	//     - Let's call 7.5 the "cusp"
	//   - To figure out which division it's in, divide the degrees by 7.5
	//   - div := int(3.116 / 7.5) = 0
	//     - so the first division
	// - Now, figure out how much percentage of the division is
	//   Mars occupying
	//   - For the first division, its 3.116 / 7.5 = 0.415
	//     - for any other division, the equation will be
	//       perc := (deg%7.5) / 7.5
	//       Or perc := (deg%cusp) / cusp
	// - Now we have the "div" and "cusp" and `perc`
	//   - All we need is to multiply perc with 30
	//     to get the degree in the new sign
	//   - So 0.415 * 30 = 12.45
	//     - which is 12 degrees and 27 minutes
	//     - which is the new zodiacal position
	// - Also, an easy calculation for the sign is to add 3*div to the
	//   original sign

	// newSign := sign.NewSignFromInt(zp.Sign.Int() + 3)
	totalDivisions := 4
	cusp := 30.0 / float64(totalDivisions)
	signDeg := zp.SignDegrees()
	currDivision := int(signDeg / cusp)
	perc := math.Mod(signDeg, cusp) / cusp
	newDegrees := perc * 30.0
	newMinutes := math.Mod(newDegrees, 1) * 60.0

	newSignInt := zp.Sign.Int() + ((totalDivisions - 1) * currDivision)
	// if newSignInt > 12 {
	// 	newSignInt -= 12
	// }
	// newSignInt = int(math.Mod(float64(newSignInt), 12))
	newSign, err := sign.NewSignFromInt(newSignInt)
	if err != nil {
		return nil, err
	}
	return zodiacalpos.NewZodiacalPos(
		newSign,
		int(newDegrees),
		int(newMinutes),
	), nil
}

func transformZodiacalPosToD7(
	pid pointid.PointID,
	zp *zodiacalpos.ZodiacalPos,
) (*zodiacalpos.ZodiacalPos, error) {
	totalDivisions := 7
	cusp := 30.0 / float64(totalDivisions)
	signDeg := zp.SignDegrees()
	currDivision := int(signDeg / cusp)
	perc := math.Mod(signDeg, cusp) / cusp
	newDegrees := perc * 30.0
	newMinutes := math.Mod(newDegrees, 1) * 60.0

	newSignInt := zp.Sign.Int()
	if zp.Sign.Int()%2 == 0 {
		newSignInt += 6
	}
	// if newSignInt > 12 {
	// 	newSignInt -= 12
	// }
	newSignInt += currDivision
	// newSignInt = int(math.Mod(float64(newSignInt), 12))
	newSign, err := sign.NewSignFromInt(newSignInt)
	if err != nil {
		return nil, err
	}
	return zodiacalpos.NewZodiacalPos(
		newSign,
		int(newDegrees),
		int(newMinutes),
	), nil
}

func transformZodiacalPosToD9(
	pid pointid.PointID,
	zp *zodiacalpos.ZodiacalPos,
) (*zodiacalpos.ZodiacalPos, error) {
	totalDivisions := 9
	cusp := 30.0 / float64(totalDivisions)
	signDeg := zp.SignDegrees()
	currDivision := int(signDeg / cusp)
	perc := math.Mod(signDeg, cusp) / cusp
	newDegrees := perc * 30.0
	newMinutes := math.Mod(newDegrees, 1) * 60.0

	// Sign calculation
	// - if sign is divisible by 4 (or == 1), add nothing
	// - if sign is divisible by 3, add 3 to the sign
	// - if sign is divisible by 2, add 8 to the sign
	newSignInt := zp.Sign.Int()
	if newSignInt == 1 || newSignInt == 4 || newSignInt == 7 ||
		newSignInt == 10 {
		// Do nothing
	} else if newSignInt == 2 || newSignInt == 5 || newSignInt == 8 || newSignInt == 11 {
		newSignInt += 8
	} else if newSignInt == 3 || newSignInt == 6 || newSignInt == 9 || newSignInt == 12 {
		newSignInt += 4
	}
	newSignInt += currDivision
	// if newSignInt > 12 {
	// 	newSignInt -= 12
	// }
	// if newSignInt == 12 {
	// 	fmt.Printf("FUCK")
	// }
	// newSignInt = int(math.Mod(float64(newSignInt), 12))
	newSign, err := sign.NewSignFromInt(newSignInt)
	if err != nil {
		return nil, err
	}
	return zodiacalpos.NewZodiacalPos(
		newSign,
		int(newDegrees),
		int(newMinutes),
	), nil
}

func transformZodiacalPosToD10(
	pid pointid.PointID,
	zp *zodiacalpos.ZodiacalPos,
) (*zodiacalpos.ZodiacalPos, error) {
	totalDivisions := 10
	cusp := 30.0 / float64(totalDivisions)
	signDeg := zp.SignDegrees()
	currDivision := int(signDeg / cusp)
	perc := math.Mod(signDeg, cusp) / cusp
	newDegrees := perc * 30.0
	newMinutes := math.Mod(newDegrees, 1) * 60.0

	// Sign calculation
	// - odd signs increase by 1
	// - even signs increase by 8+1
	newSignInt := zp.Sign.Int()
	if newSignInt%2 == 0 {
		newSignInt += 8
	}
	newSignInt += currDivision
	// if newSignInt > 12 {
	// 	newSignInt -= 12
	// }
	// newSignInt = int(math.Mod(float64(newSignInt), 12))
	newSign, err := sign.NewSignFromInt(newSignInt)
	if err != nil {
		return nil, err
	}
	return zodiacalpos.NewZodiacalPos(
		newSign,
		int(newDegrees),
		int(newMinutes),
	), nil
}

// func transformZodiacalPosToD24(
// 	pid pointid.PointID,
// 	zp *zodiacalpos.ZodiacalPos,
// ) (*zodiacalpos.ZodiacalPos, error) {
// 	totalDivisions := 24
// 	cusp := 30.0 / float64(totalDivisions)
// 	signDeg := zp.SignDegrees()
// 	currDivision := int(signDeg / cusp)
// 	perc := math.Mod(signDeg, cusp) / cusp
// 	newDegrees := perc * 30.0
// 	newMinutes := math.Mod(newDegrees, 1) * 60.0

// 	// Sign calculation
// 	// - odd signs increase by 1
// 	// - even signs increase by 8+1
// 	newSignInt := zp.Sign.Int()
// 	if newSignInt%2 == 0 {
// 		newSignInt += 8
// 	}
// 	newSignInt += currDivision
// 	if newSignInt > 12 {
// 		newSignInt -= 12
// 	}
// 	newSignInt = int(math.Mod(float64(newSignInt), 12))
// 	newSign := sign.NewSignFromInt(newSignInt)
// 	return zodiacalpos.NewZodiacalPos(newSign, int(newDegrees), int(newMinutes)), nil
// }
