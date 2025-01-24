package sign

import (
	"fmt"
	"math"
	"strings"

	"github.com/afjoseph/sacredstar/pointid"
)

type Sign string

var (
	Aries       = Sign("aries")
	Taurus      = Sign("taurus")
	Gemini      = Sign("gemini")
	Cancer      = Sign("cancer")
	Leo         = Sign("leo")
	Virgo       = Sign("virgo")
	Libra       = Sign("libra")
	Scorpio     = Sign("scorpio")
	Sagittarius = Sign("sagittarius")
	Capricorn   = Sign("capricorn")
	Aquarius    = Sign("aquarius")
	Pisces      = Sign("pisces")
)

// XXX <03-02-2024, afjoseph> We're only doing JSON serialization here because
// invalid signs (like `"invalid"`) should yield an error when parsed. This
// won't happen if keep the default JSON serialization
func (h Sign) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strings.ToLower(h.String()) + `"`), nil
}

func (h *Sign) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	var err error
	*h, err = NewSignFromStr(s)
	if err != nil {
		return fmt.Errorf("Invalid sign: %s", s)
	}
	return nil
}

func (h Sign) String() string {
	return string(h)
}

func (h Sign) Int() int {
	switch h {
	case Aries:
		return 1
	case Taurus:
		return 2
	case Gemini:
		return 3
	case Cancer:
		return 4
	case Leo:
		return 5
	case Virgo:
		return 6
	case Libra:
		return 7
	case Scorpio:
		return 8
	case Sagittarius:
		return 9
	case Capricorn:
		return 10
	case Aquarius:
		return 11
	case Pisces:
		return 12
	}
	return -1
}

func NewSignFromInt(i int) (Sign, error) {
	i = int(math.Mod(float64(i), 12))
	if i == 0 {
		i = 12
	}

	switch i {
	case 1:
		return Aries, nil
	case 2:
		return Taurus, nil
	case 3:
		return Gemini, nil
	case 4:
		return Cancer, nil
	case 5:
		return Leo, nil
	case 6:
		return Virgo, nil
	case 7:
		return Libra, nil
	case 8:
		return Scorpio, nil
	case 9:
		return Sagittarius, nil
	case 10:
		return Capricorn, nil
	case 11:
		return Aquarius, nil
	case 12:
		return Pisces, nil
	}
	return Sign(""), fmt.Errorf("Invalid sign: %d", i)
}

func NewSignFromStr(s string) (Sign, error) {
	switch strings.ToLower(s) {
	case "aries":
		return Aries, nil
	case "taurus":
		return Taurus, nil
	case "gemini":
		return Gemini, nil
	case "cancer":
		return Cancer, nil
	case "leo":
		return Leo, nil
	case "virgo":
		return Virgo, nil
	case "libra":
		return Libra, nil
	case "scorpio":
		return Scorpio, nil
	case "sagittarius":
		return Sagittarius, nil
	case "capricorn":
		return Capricorn, nil
	case "aquarius":
		return Aquarius, nil
	case "pisces":
		return Pisces, nil
	}
	return Sign(""), fmt.Errorf("Invalid sign: %s", s)
}

func DegreeToSign(degree float64) Sign {
	switch {
	case degree < 0:
		degree += 360
	case degree >= 360:
		degree -= 360
	}
	switch {
	case degree < 30:
		return Aries
	case degree < 60:
		return Taurus
	case degree < 90:
		return Gemini
	case degree < 120:
		return Cancer
	case degree < 150:
		return Leo
	case degree < 180:
		return Virgo
	case degree < 210:
		return Libra
	case degree < 240:
		return Scorpio
	case degree < 270:
		return Sagittarius
	case degree < 300:
		return Capricorn
	case degree < 330:
		return Aquarius
	case degree < 360:
		return Pisces
	}
	return Aries
}

func (s Sign) TraditionalRuler() pointid.PointID {
	switch s {
	case Aries:
		return pointid.Mars
	case Taurus:
		return pointid.Venus
	case Gemini:
		return pointid.Mercury
	case Cancer:
		return pointid.Moon
	case Leo:
		return pointid.Sun
	case Virgo:
		return pointid.Mercury
	case Libra:
		return pointid.Venus
	case Scorpio:
		return pointid.Mars
	case Sagittarius:
		return pointid.Jupiter
	case Capricorn:
		return pointid.Saturn
	case Aquarius:
		return pointid.Saturn
	case Pisces:
		return pointid.Jupiter
	}
	return pointid.None
}

func (s Sign) ModernRuler() pointid.PointID {
	switch s {
	case Aries:
		return pointid.Mars
	case Taurus:
		return pointid.Venus
	case Gemini:
		return pointid.Mercury
	case Cancer:
		return pointid.Moon
	case Leo:
		return pointid.Sun
	case Virgo:
		return pointid.Mercury
	case Libra:
		return pointid.Venus
	case Scorpio:
		return pointid.Pluto
	case Sagittarius:
		return pointid.Jupiter
	case Capricorn:
		return pointid.Saturn
	case Aquarius:
		return pointid.Uranus
	case Pisces:
		return pointid.Neptune
	}
	return pointid.None
}

func (s Sign) Next() Sign {
	switch s {
	case Aries:
		return Taurus
	case Taurus:
		return Gemini
	case Gemini:
		return Cancer
	case Cancer:
		return Leo
	case Leo:
		return Virgo
	case Virgo:
		return Libra
	case Libra:
		return Scorpio
	case Scorpio:
		return Sagittarius
	case Sagittarius:
		return Capricorn
	case Capricorn:
		return Aquarius
	case Aquarius:
		return Pisces
	case Pisces:
		return Aries
	}
	panic("Invalid sign")
}

func (s Sign) Previous() Sign {
	switch s {
	case Aries:
		return Pisces
	case Taurus:
		return Aries
	case Gemini:
		return Taurus
	case Cancer:
		return Gemini
	case Leo:
		return Cancer
	case Virgo:
		return Leo
	case Libra:
		return Virgo
	case Scorpio:
		return Libra
	case Sagittarius:
		return Scorpio
	case Capricorn:
		return Sagittarius
	case Aquarius:
		return Capricorn
	case Pisces:
		return Aquarius
	}
	panic("Invalid sign")
}
