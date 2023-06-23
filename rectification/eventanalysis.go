package rectification

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/house"
	"github.com/afjoseph/sacredstar/pointid"
	"github.com/afjoseph/sacredstar/timeandzone"
	"github.com/afjoseph/sacredstar/util"
	"github.com/afjoseph/sacredstar/wrapper"
	"github.com/afjoseph/sacredstar/zodiacalpos"
)

type EventType string

const (
	EventType_MovingHome EventType = "movinghome"
	EventType_ChildBirth           = "childbirth"
	EventType_Marriage             = "marriage"
	EventType_Career               = "career"
	// EventType_Education
)

func (b EventType) String() string {
	return string(b)
}

func (b EventType) VargaChartType() chart.ChartType {
	switch b {
	case EventType_MovingHome:
		return chart.D4ChartType
	case EventType_ChildBirth:
		return chart.D7ChartType
	case EventType_Marriage:
		return chart.D9ChartType
	case EventType_Career:
		return chart.D10ChartType
	// case EventType_Education:
	//  return D24ChartType
	default:
		panic(fmt.Sprintf("Unknown birth time analysis event type: %v", b))
	}
}

type Event struct {
	Title     string                  `json:"title"     binding:"required"`
	Desc      string                  `json:"desc"      binding:"required"`
	Time      timeandzone.TimeAndZone `json:"time"      binding:"required"`
	EventType EventType               `json:"eventType" binding:"required"`
}

func (e Event) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *Event) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), e)
}

type Events []Event

func (e Events) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *Events) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), e)
}

type EventAnalysis struct {
	Event       Event        `json:"event"`
	Score       int          `json:"score"`
	Reasons     []Reason     `json:"reasons"`
	ActiveDasha chart.Dasha  `json:"activeDasha"`
	Chart       *chart.Chart `json:"chart"`
}

func (ea *EventAnalysis) Checksum() []byte {
	h32 := crc32.NewIEEE()
	for _, r := range ea.Reasons {
		h32.Write(r.Checksum())
	}
	h32.Write(util.Uint32ToByteSlice(uint32(ea.Score)))
	return h32.Sum(nil)
}

func (e Event) analyze(
	swe *wrapper.SwissEph,
	birthTime time.Time,
	lon, lat float64,
) (*EventAnalysis, error) {
	birthTimeInJulian := swe.GoTimeToJulianDay(birthTime.UTC())
	rashiChart, err := chart.NewChartFromJulianDay(
		swe,
		birthTimeInJulian,
		lon, lat,
		chart.D1ChartType,
		pointid.VedicPlanets,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"while calculating rashi chart for birth time %v: %w",
			birthTime,
			err,
		)
	}
	moonAstroPoint := rashiChart.GetPoint(pointid.Moon)
	if moonAstroPoint == nil {
		return nil, fmt.Errorf(
			"moon not found in rashi chart for birth time %v",
			birthTime,
		)
	}
	dashaTree, err := chart.NewDashaTree(
		swe,
		birthTime,
		moonAstroPoint,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"while calculating dasha tree for birth time %v: %w",
			birthTime,
			err,
		)
	}

	// Get the varga chart and dashas for the event time
	vargaChart, err := chart.NewChartFromJulianDay(
		swe,
		birthTimeInJulian, lon, lat,
		e.EventType.VargaChartType(),
		pointid.VedicPlanets,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error calculating varga chart %v for birth time %v: %w",
			e.EventType.VargaChartType(),
			birthTime,
			err,
		)
	}
	dasha, didFind := dashaTree.GetDashaForTime(e.Time.Time)
	if !didFind {
		return nil, fmt.Errorf(
			"could not find dasha for event time %v",
			e.Time,
		)
	}

	reasons := []Reason{}
	impHouses := e.EventType.VargaChartType().ImportantHouses()
	dashaLordHouse := vargaChart.GetHouseFor(dasha.Mahadasha.ToPointID())
	dashaLordZodPos := vargaChart.GetZodiacalPosFor(dasha.Mahadasha.ToPointID())
	subdashaLordHouse := vargaChart.GetHouseFor(dasha.Antardasha.ToPointID())
	subdashaLordZodPos := vargaChart.GetZodiacalPosFor(
		dasha.Antardasha.ToPointID(),
	)
	karakas := e.EventType.VargaChartType().Karakas()
	var karakaHouses []house.House
	karakaLord := map[pointid.PointID]*zodiacalpos.ZodiacalPos{}
	for _, karaka := range karakas {
		karakaLord[karaka] = vargaChart.GetZodiacalPosFor(karaka)
		karakaHouses = append(karakaHouses, vargaChart.GetHouseFor(karaka))
	}

	// Rule 1: dasha or subdasha lord are in the important house of the varga
	// chart type
	for _, impHouse := range impHouses {
		if dashaLordHouse == impHouse {
			reasons = append(reasons, NewReason(
				ReasonType_House,
				fmt.Sprintf("Dasha %s lord is in the important house %s",
					dasha.Mahadasha,
					impHouse,
				),
			))
		} else if subdashaLordHouse == impHouse {
			reasons = append(reasons, NewReason(
				ReasonType_House,
				fmt.Sprintf("Subdasha %s lord is in the important house %s",
					dasha.Antardasha,
					impHouse,
				),
			))
		}
	}

	// Rule 2: dasha or subdasha lords are aspecting the lord of an important
	// house OR the lord of an important house is the dasha or subdasha lord
	for _, impHouse := range impHouses {
		impHouseLord, err := vargaChart.GetHouseLordFor(
			impHouse,
			chart.HouseLordPlacement_Traditional,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"while getting house lord for %s house: %w",
				impHouse,
				err,
			)
		}
		if impHouseLord == dasha.Mahadasha.ToPointID() {
			reasons = append(reasons, NewReason(
				ReasonType_DashaLordAsHouseLord,
				fmt.Sprintf(
					"Dasha lord (%s) is the lord of an important house (the %s house)",
					dasha.Mahadasha,
					impHouse,
				),
			))
			// No need to check for aspects if the dasha or subdasha lord is
			// the lord of the important house: they are the same planet
			continue
		} else if impHouseLord == dasha.Antardasha.ToPointID() {
			reasons = append(reasons, NewReason(
				ReasonType_DashaLordAsHouseLord,
				fmt.Sprintf(
					"Subdasha lord (%s) is the lord of an important house (the %s house)",
					dasha.Antardasha,
					impHouse,
				),
			))
			// No need to check for aspects if the dasha or subdasha lord is
			// the lord of the important house: they are the same planet
			continue
		}

		impHouseLordZodPos := vargaChart.GetZodiacalPosFor(impHouseLord)
		if asp := dashaLordZodPos.GetAspect(impHouseLordZodPos); asp != nil {
			if asp.IsHard() {
				reasons = append(reasons, NewReason(
					ReasonType_HouseLordAspect_Hard,
					fmt.Sprintf(
						"Dasha %s lord is aspecting lord of an important house (%s %s house lord) with a hard aspect (%s)",
						dasha.Mahadasha,
						impHouseLord,
						impHouse,
						asp,
					),
				))
			}
			// } else if asp.IsSoft() {
			// 	reasons = append(reasons, NewReason(
			// 		ReasonType_HouseLordAspect_Soft,
			// 		fmt.Sprintf(
			// 			"Dasha %s lord is aspecting lord of an important house (%s %s house lord) with a soft aspect (%s)",
			// 			dasha.Mahadasha,
			// 			impHouseLord,
			// 			impHouse, asp,
			// 		),
			// 	))
			// }
		}
		if asp := subdashaLordZodPos.GetAspect(impHouseLordZodPos); asp != nil {
			if asp.IsHard() {
				reasons = append(reasons, NewReason(
					ReasonType_HouseLordAspect_Hard,
					fmt.Sprintf(
						"Subdasha %s lord is aspecting lord of an important house (%s %s house lord) with a hard aspect (%s)",
						dasha.Antardasha,
						impHouseLord,
						impHouse,
						asp,
					),
				))
			}
			// } else if asp.IsSoft() {
			// 	reasons = append(reasons, NewReason(
			// 		ReasonType_HouseLordAspect_Soft,
			// 		fmt.Sprintf(
			// 			"Subdasha %s lord is aspecting lord of an important house (%s %s house lord) with a soft aspect (%s)",
			// 			dasha.Antardasha,
			// 			impHouseLord, impHouse, asp,
			// 		),
			// 	))
			// }
		}
	}

	// Rule 3: The karaka is aspecting a dasha or subdasha lord
	for kl, karakaZodPos := range karakaLord {
		if asp := dashaLordZodPos.GetAspect(karakaZodPos); asp != nil {
			if asp.IsHard() {
				reasons = append(reasons, NewReason(
					ReasonType_KarakaAspect_Hard,
					fmt.Sprintf(
						"Dasha %s lord is aspecting the karaka lord %s with a hard aspect (%s)",
						dasha.Mahadasha,
						kl,
						asp,
					),
				))
			}
			// } else if asp.IsSoft() {
			// 	reasons = append(reasons, NewReason(
			// 		ReasonType_KarakaAspect_Soft,
			// 		fmt.Sprintf(
			// 			"Dasha %s lord is aspecting the karaka lord %s with a soft aspect (%s)",
			// 			dasha.Mahadasha,
			// 			kl, asp,
			// 		),
			// 	))
			// }
		}
		if asp := subdashaLordZodPos.GetAspect(karakaZodPos); asp != nil {
			if asp.IsHard() {
				reasons = append(reasons, NewReason(
					ReasonType_KarakaAspect_Hard,
					fmt.Sprintf(
						"Subdasha %s lord is aspecting the karaka lord %s with a hard aspect (%s)",
						dasha.Antardasha,
						kl,
						asp,
					),
				))
			}
			// } else if asp.IsSoft() {
			// 	reasons = append(reasons, NewReason(
			// 		ReasonType_KarakaAspect_Soft,
			// 		fmt.Sprintf(
			// 			"Subdasha %s lord is aspecting the karaka lord %s with a soft aspect (%s)",
			// 			dasha.Antardasha,
			// 			kl, asp,
			// 		),
			// 	))
			// }
		}
	}

	// Rule 4: The ascendant is in the important house
	if dashaLordHouse == house.House1 {
		reasons = append(reasons, NewReason(
			ReasonType_Ascendant,
			fmt.Sprintf("Dasha %s lord is in the ascendant", dasha.Mahadasha),
		))
	} else if subdashaLordHouse == house.House1 {
		reasons = append(reasons, NewReason(
			ReasonType_Ascendant,
			fmt.Sprintf("Subdasha %s lord is in the ascendant", dasha.Antardasha),
		))
	}

	// Rule 5: The Rashi chart house lord becomes a dasha or sub-dasha lord
	for _, impHouse := range impHouses {
		rashiChartImpHouseLord, err := rashiChart.GetHouseLordFor(
			impHouse,
			chart.HouseLordPlacement_Traditional,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"while getting house lord for %s house: %w",
				impHouse,
				err,
			)
		}

		if rashiChartImpHouseLord == dasha.Mahadasha.ToPointID() {
			reasons = append(
				reasons,
				NewReason(
					ReasonType_RashiChartHouseLord,
					fmt.Sprintf(
						"Dasha lord %s is the Rashi chart lord of the important house %s",
						dasha.Mahadasha,
						impHouse,
					),
				),
			)
		} else if rashiChartImpHouseLord == dasha.Antardasha.ToPointID() {
			reasons = append(
				reasons,
				NewReason(ReasonType_RashiChartHouseLord,
					fmt.Sprintf("Subdasha lord %s is the Rashi chart lord of the important house %s",
						dasha.Antardasha,
						impHouse,
					),
				),
			)
		}
	}

	// Rule 6: Bhavat Bhavam: the dasha or sub-dasha lord is placed in the xth house
	// from the xth important house of the varga chart.
	// For example, if the dasha or sub-dasha lord is placed in the 5th house from the 5th
	// house (9th house) of a D7 chart (childbirth varga
	// chart), there's a good chance that the event (i.e., childbirth)
	// would happen
	for _, impHouse := range impHouses {
		secondaryImpHouseAsInt := impHouse.Int() + (impHouse.Int() - 1)
		// So the 12th house from the 12th house is the 11th house
		// 12 + (12 - 1) = 23
		// 23 % 12 = 11
		if secondaryImpHouseAsInt > 12 {
			secondaryImpHouseAsInt = secondaryImpHouseAsInt % 12
		}
		secondaryImpHouse, err := house.HouseFromInt(secondaryImpHouseAsInt)
		if err != nil {
			panic("This can't be null")
		}

		if dashaLordHouse == secondaryImpHouse {
			reasons = append(
				reasons,
				NewReason(
					ReasonType_BhavatBhavam,
					fmt.Sprintf(
						"Dasha lord (%s) is %d houses from the important %dth house (i.e., the %s house)",
						dasha.Mahadasha,
						impHouse.Int(),
						impHouse.Int(),
						secondaryImpHouse,
					),
				),
			)
		} else if subdashaLordHouse == secondaryImpHouse {
			reasons = append(
				reasons,
				NewReason(
					ReasonType_BhavatBhavam,
					fmt.Sprintf(
						"Subdasha lord (%s) is %d houses from the important %dth house (i.e., the %s house)",
						dasha.Antardasha,
						impHouse.Int(),
						impHouse.Int(),
						secondaryImpHouse,
					),
				),
			)
		}
	}

	// TODO <26-01-2024, afjoseph> We can do better, but now the score is just
	// the number of reasons
	score := 0
	for _, r := range reasons {
		if r.Type == ReasonType_House ||
			r.Type == ReasonType_DashaLordAsHouseLord {
			score += 2
		} else {
			score++
		}
	}

	return &EventAnalysis{
		Event:       e,
		ActiveDasha: dasha,
		Reasons:     reasons,
		Score:       score,
		Chart:       vargaChart,
	}, nil
}
