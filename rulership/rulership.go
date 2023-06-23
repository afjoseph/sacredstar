package rulership

import (
	"github.com/afjoseph/sacredstar/astropoint"
	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/house"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/zodiacalpos"
	"github.com/go-playground/errors/v5"
)

// MustFindRuler finds the ruler for a house h in a chart chrt.
// It assumes traditional rulership (i.e., Anything after Saturn is out)
func MustFindRuler(
	chrt *chart.Chart,
	h house.House,
) *astropoint.AstroPoint {
	// Get the sign of the ruler
	sign := chrt.GetSignOfHouse(h)
	rulerPid := sign.TraditionalRuler()
	if rulerPid == pointid.None {
		panic(errors.Newf("No ruler for house %s: %s", h, chrt.String()))
	}
	ruler := chrt.MustGetPoint(rulerPid)
	return ruler
}

func GetStrength(
	chrt *chart.Chart,
	p *astropoint.AstroPoint,
) (Strength, int, []StrengthReason) {
	reasons := []StrengthReason{}

	score := 0
	if p.IsDomicile() {
		reasons = append(reasons, StrengthReasonDomicile{})
		score += reasons[len(reasons)-1].Weight()
	}
	if p.IsDetriment() {
		reasons = append(reasons, StrengthReasonDetriment{})
		score += reasons[len(reasons)-1].Weight()
	}
	if p.IsExalted() {
		reasons = append(reasons, StrengthReasonExalted{})
		score += reasons[len(reasons)-1].Weight()
	}
	if p.IsFall() {
		reasons = append(reasons, StrengthReasonFall{})
		score += reasons[len(reasons)-1].Weight()
	}
	// Check hard aspects to malefics
	for _, maleficPointID := range []pointid.PointID{pointid.Mars, pointid.Saturn} {
		maleficPoint := chrt.MustGetPoint(maleficPointID)
		asp := p.ZodiacalPos.GetAspect(maleficPoint.ZodiacalPos)
		if asp != nil && asp.IsHard() {
			reasons = append(reasons, StrengthReasonAspectToMalefic{})
			score += reasons[len(reasons)-1].Weight()
		}
	}

	// Check if combust
	sunPoint := chrt.MustGetPoint(pointid.Sun)
	asp := p.ZodiacalPos.GetAspect(sunPoint.ZodiacalPos)
	if asp != nil && asp.Type == zodiacalpos.AspectType_Conjunction {
		reasons = append(reasons, StrengthReasonCombust{})
		score += reasons[len(reasons)-1].Weight()
	}

	// Determine score
	const threshold = 0
	switch {
	case score > threshold:
		return StrengthStrong, score, reasons
	case score < -threshold:
		return StrengthWeak, score, reasons
	default:
		return StrengthNeutral, score, reasons
	}
}
