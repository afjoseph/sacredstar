package zodiacalpos

import (
	"fmt"
	"math"

	"github.com/afjoseph/sacredstar/sign"
)

type ZodiacalPos struct {
	Sign    sign.Sign `json:"sign"`
	Degrees int       `json:"degrees"`
	Minutes int       `json:"minutes"`
}

func NewZodiacalPos(sign sign.Sign, deg, min int) *ZodiacalPos {
	return &ZodiacalPos{
		Sign:    sign,
		Degrees: deg,
		Minutes: min,
	}
}

func NewZodiacalPosFromLongitude(longitude float64) *ZodiacalPos {
	var s sign.Sign
	var signDeg float64
	switch {
	case longitude < 0:
		longitude += 360
	case longitude >= 360:
		longitude -= 360
	}
	switch {
	case longitude < 30:
		s = sign.Aries
		signDeg = longitude
	case longitude < 60:
		s = sign.Taurus
		signDeg = longitude - 30
	case longitude < 90:
		s = sign.Gemini
		signDeg = longitude - 60
	case longitude < 120:
		s = sign.Cancer
		signDeg = longitude - 90
	case longitude < 150:
		s = sign.Leo
		signDeg = longitude - 120
	case longitude < 180:
		s = sign.Virgo
		signDeg = longitude - 150
	case longitude < 210:
		s = sign.Libra
		signDeg = longitude - 180
	case longitude < 240:
		s = sign.Scorpio
		signDeg = longitude - 210
	case longitude < 270:
		s = sign.Sagittarius
		signDeg = longitude - 240
	case longitude < 300:
		s = sign.Capricorn
		signDeg = longitude - 270
	case longitude < 330:
		s = sign.Aquarius
		signDeg = longitude - 300
	case longitude < 360:
		s = sign.Pisces
		signDeg = longitude - 330
	}

	// The whole number part of the degree will be the degrees
	signDegAsInt := int(signDeg)
	// Subtract the whole number from the original degrees to get the fraction,
	// then multiply by 60 to convert to minutes
	// 60 -> 1
	// 30 -> 0.5
	// 15 -> 0.25
	// y -> x
	// y = x * 60
	// Round the minutes to the nearest whole number if needed
	fractionalPart := signDeg - math.Floor(signDeg)
	roundedMinutes := int(fractionalPart * 60)
	// If the minutes round up to 60, add one to degrees and set minutes to 0
	if roundedMinutes == 60 {
		signDegAsInt++
		roundedMinutes = 0
	}
	return &ZodiacalPos{
		Sign:    s,
		Degrees: signDegAsInt,
		Minutes: roundedMinutes,
	}
}

func (zp *ZodiacalPos) String() string {
	return fmt.Sprintf("%s %d %d", zp.Sign.String(), zp.Degrees, zp.Minutes)
}

func (zp *ZodiacalPos) LessThan(rhs *ZodiacalPos) bool {
	if zp.Sign.Int() < rhs.Sign.Int() {
		return true
	}
	if zp.Sign.Int() == rhs.Sign.Int() && zp.Degrees < rhs.Degrees {
		return true
	}
	if zp.Sign.Int() == rhs.Sign.Int() && zp.Degrees == rhs.Degrees &&
		zp.Minutes < rhs.Minutes {
		return true
	}
	return false
}

func (zp *ZodiacalPos) LessThanOrEqualTo(rhs *ZodiacalPos) bool {
	if zp.Sign.Int() < rhs.Sign.Int() {
		return true
	}
	if zp.Sign.Int() == rhs.Sign.Int() && zp.Degrees < rhs.Degrees {
		return true
	}
	if zp.Sign.Int() == rhs.Sign.Int() && zp.Degrees == rhs.Degrees &&
		zp.Minutes <= rhs.Minutes {
		return true
	}
	return false
}

func (zp *ZodiacalPos) GreaterThan(rhs *ZodiacalPos) bool {
	if zp.Sign.Int() > rhs.Sign.Int() {
		return true
	}
	if zp.Sign.Int() == rhs.Sign.Int() && zp.Degrees > rhs.Degrees {
		return true
	}
	if zp.Sign.Int() == rhs.Sign.Int() && zp.Degrees == rhs.Degrees &&
		zp.Minutes > rhs.Minutes {
		return true
	}
	return false
}

func (zp *ZodiacalPos) GreaterThanOrEqualTo(rhs *ZodiacalPos) bool {
	if zp.Sign.Int() > rhs.Sign.Int() {
		return true
	}
	if zp.Sign.Int() == rhs.Sign.Int() && zp.Degrees >= rhs.Degrees {
		return true
	}
	if zp.Sign.Int() == rhs.Sign.Int() && zp.Degrees == rhs.Degrees &&
		zp.Minutes >= rhs.Minutes {
		return true
	}
	return false
}

func (zp *ZodiacalPos) EqualTo(rhs *ZodiacalPos) bool {
	return zp.Sign.Int() == rhs.Sign.Int() && zp.Degrees == rhs.Degrees &&
		zp.Minutes == rhs.Minutes
}

func (zp *ZodiacalPos) Diff(rhs *ZodiacalPos) *ZodiacalPos {
	var diff *ZodiacalPos
	if zp.LessThan(rhs) {
		diff = rhs.Diff(zp)
	} else {
		diff = NewZodiacalPos(zp.Sign, zp.Degrees-rhs.Degrees, zp.Minutes-rhs.Minutes)
	}
	return diff
}

func (zp *ZodiacalPos) DiffInAbsDegrees(rhs *ZodiacalPos) float64 {
	llhs := zp.AbsDegrees()
	rrhs := rhs.AbsDegrees()
	ret := math.Abs(llhs - rrhs)
	// If the difference is greater than 180, then we need to subtract the
	// difference from 360 to get the actual difference.
	// In this way, the difference between 0 and 359 is 1, not 359.
	if ret > 180 {
		ret = 360 - ret
	}
	return ret
}

func (zp *ZodiacalPos) Add(rhs *ZodiacalPos) *ZodiacalPos {
	var sum *ZodiacalPos
	if zp.GreaterThan(rhs) {
		sum = rhs.Add(zp)
	} else {
		sum = NewZodiacalPos(zp.Sign, zp.Degrees+rhs.Degrees, zp.Minutes+rhs.Minutes)
	}
	return sum
}

func (zp *ZodiacalPos) AbsDegrees() float64 {
	s := float64((zp.Sign.Int() - 1) * 30)
	d := float64(zp.Degrees)
	m := float64(zp.Minutes) / float64(60.0)
	return float64(s + d + m)
}

func (zp *ZodiacalPos) SignDegrees() float64 {
	d := float64(zp.Degrees)
	m := float64(zp.Minutes) / float64(60.0)
	return float64(d + m)
}

func (zp *ZodiacalPos) Opposite() *ZodiacalPos {
	deg := zp.AbsDegrees() + 180
	if deg >= 360 {
		deg -= 360
	}
	return NewZodiacalPosFromLongitude(deg)
}
