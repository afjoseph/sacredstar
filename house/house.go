package house

import (
	"fmt"

	"github.com/afjoseph/sacredstar/sign"
)

type House int

func (h House) String() string {
	switch h {
	case House1:
		return "1st"
	case House2:
		return "2nd"
	case House3:
		return "3rd"
	case House4:
		return "4th"
	case House5:
		return "5th"
	case House6:
		return "6th"
	case House7:
		return "7th"
	case House8:
		return "8th"
	case House9:
		return "9th"
	case House10:
		return "10th"
	case House11:
		return "11th"
	case House12:
		return "12th"
	}
	panic(fmt.Sprintf("Invalid house: %d", h))
}

func HouseFromInt(i int) (House, error) {
	switch i {
	case 1:
		return House1, nil
	case 2:
		return House2, nil
	case 3:
		return House3, nil
	case 4:
		return House4, nil
	case 5:
		return House5, nil
	case 6:
		return House6, nil
	case 7:
		return House7, nil
	case 8:
		return House8, nil
	case 9:
		return House9, nil
	case 10:
		return House10, nil
	case 11:
		return House11, nil
	case 12:
		return House12, nil
	}
	return -1, fmt.Errorf("Invalid house: %d", i)
}

var (
	House1  = House(1)
	House2  = House(2)
	House3  = House(3)
	House4  = House(4)
	House5  = House(5)
	House6  = House(6)
	House7  = House(7)
	House8  = House(8)
	House9  = House(9)
	House10 = House(10)
	House11 = House(11)
	House12 = House(12)

	HouseNone = House(-1)
)

func (h House) Int() int {
	return int(h)
}

func (h House) Opposite() House {
	i := h.Int()
	if i > 6 {
		// No way this is an error
		h, err := HouseFromInt(i - 6)
		if err != nil {
			panic(err)
		}
		return h
	}
	h, err := HouseFromInt(i + 6)
	if err != nil {
		panic(err)
	}
	return h
}

func NewHouseFromSign(targetSign sign.Sign, ascSign sign.Sign) House {
	// Ascendant's sign is the first house
	// So if the ascendant is in Taurus, Taurus is the 1st house
	// If the target sign is Gemini, then the target house is the 2nd house

	// target sign is cancer: 4
	// asc sign is Taurus: 2
	diff := targetSign.Int() - ascSign.Int()
	// XXX <18-01-2024, afjoseph> Go doesn't do well with modulo negative numbers so
	// we have to add 12 to the diff if it's negative. Python has no problem with it.
	// You can run `python3 -c "print(-1 % 12)"` and it will print 11 correctly.
	if diff < 0 {
		diff += 12
	}
	h, err := HouseFromInt((diff % 12) + 1)
	if err != nil {
		panic(err)
	}
	return h
}
