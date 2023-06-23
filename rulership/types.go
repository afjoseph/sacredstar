package rulership

import "github.com/go-playground/errors/v5"

type Strength int

const (
	StrengthStrong Strength = iota
	StrengthNeutral
	StrengthWeak
)

func (rs Strength) String() string {
	switch rs {
	case StrengthStrong:
		return "strong"
	case StrengthNeutral:
		return "neutral"
	case StrengthWeak:
		return "weak"
	default:
		panic(errors.Newf("Unknown ruler strength: %d", rs))
	}
}

type StrengthReason interface {
	String() string
	Weight() int
}

type StrengthReasonDomicile struct{}

func (sr StrengthReasonDomicile) String() string {
	return "domicile"
}

func (sr StrengthReasonDomicile) Weight() int {
	return 3
}

type StrengthReasonDetriment struct{}

func (sr StrengthReasonDetriment) String() string {
	return "detriment"
}

func (sr StrengthReasonDetriment) Weight() int {
	return -3
}

type StrengthReasonExalted struct{}

func (sr StrengthReasonExalted) String() string {
	return "exalted"
}

func (sr StrengthReasonExalted) Weight() int {
	return 2
}

type StrengthReasonFall struct{}

func (sr StrengthReasonFall) String() string {
	return "fall"
}

func (sr StrengthReasonFall) Weight() int {
	return -2
}

type StrengthReasonAspectToMalefic struct{}

func (sr StrengthReasonAspectToMalefic) String() string {
	return "aspect to malefic"
}

func (sr StrengthReasonAspectToMalefic) Weight() int {
	return -1
}

type StrengthReasonCombust struct{}

func (sr StrengthReasonCombust) String() string {
	return "combust"
}

func (sr StrengthReasonCombust) Weight() int {
	return -2
}
