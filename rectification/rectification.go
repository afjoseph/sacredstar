package rectification

// The idea of birth time analysis is the following:
// - Each Varga chart (i.e., a 'D' chart) pertains to a specific area of life
// - When that area of life occurs, the dasha and sub-dasha lords of that chart
//   will have a relevant part in that varga chart
// - Birth time analysis works by
//   - Finding the dasha and sub-dasha lords at the time of the event
//   - Calculating the varga charts (for that area of life) around the time of birth
//     - So, for example, if the event is marriage, we calculate multiple D9 charts
//       around the time of birth (maybe a few charts before and after the birth)
//   - Then we see how the dasha and sub-dasha lords of the event chart are placed
//     in the varga charts: the chart that makes the "most" sense is the one we'll pick
//   - Do the same for other areas of life and compare those varga charts
//   - If this method is correct, there should be an overlap between the varga charts
//     and that overlap will be the time of birth
//
// Varga charts:
// - D4: Moving home
//   - Ascendant lasting 30 min
// - D7: Children
//   - Ascendant lasting 17 min
// - D9: Marriage
//   - Ascendant lasting 13 min
// - D10: Career
//   - Ascendant lasting 12 min
// - D24: Education
//   - Ascendant lasting 5 min
//
// For now, we'll only focus on D4, D7, D9 and D10
//
// The "predictive factors" we're looking for, in relation to the dasha and sub-dasha lords, are:
// - House
// - House Lord
// - Karaka
// - Ascendant in the varga chart
// - Natal chart house lord in the varga chart
// - Bhavat Bhavam
//
// The input to the analysis depends on the varga chart, but it is always a singular event time (time.Time).
// The processing depends on the varga chart, but they all follow the same rules. Here's some very dump pseudocode:
//
//				func analyzeBirthTime(
//					birthTime time.Time,
//					eventTime time.Time,
//					vargaChartType VargaChartType,
//				) (ret []BirthTimeRectificationReason, score int) {
//					vargaChart := make a varga chart from the time of birth using vargaChartType
//					dashaTree := from birthTime
//					dasha, subdasha := dashaTree.GetDashaForTime(eventTime)
//
//					reasons := []BirthTimeRectificationReason{}
//					impHouse := vargaChartType.ImportantHouse()
//					impHouseLord := vargaChart.GetHouseLordFor(impHouse)
//					importHouseLordZodPos := vargaChart.GetZodiacalPosFor(impHouseLord)
//					dashaLordHouse := vargaChart.GetHouseFor(dasha.Lord)
//					dashaLordZodPos := vargaChart.GetZodiacalPosFor(dasha.Lord)
//					subdashaLordHouse := vargaChart.GetHouseFor(subdasha.Lord)
//					subdashaLordZodPos := vargaChart.GetZodiacalPosFor(subdasha.Lord)
//					karaka := vargaChartType.Karaka()
//					karakaHouse := vargaChart.GetHouseFor(karaka)
//					// Calculate the birth's Rashi chart
//					rashiChart := make a rashi chart from the time of birth
//					rashiChartImpHouseLord := rashiChart.GetHouseLordFor(impHouse)
//
//					// Rule 1: dasha or subdasha lord are in the important house of the varga chart type
//					if dashaLordHouse == impHouse || subdashaLordHouse == impHouse {
//						reasons += BirthTimeRectificationReason_House
//					}
//
//					// Rule 2: dasha or subdasha lords are aspecting the lord
//					// of an important house
//					asp := dashaLordZodPos.GetAspect(importHouseLordZodPos)
//					if asp {
//						reasons += BirthTimeRectificationReason_HouseLord
//					}
//					asp := subDashaLordZodPos.GetAspect(importHouseLordZodPos)
//					if asp {
//						reasons += BirthTimeRectificationReason_HouseLord
//					}
//
//					// Rule 3: The karaka is aspecting a dasha or subdasha lord
//					asp := karaka.GetAspect(dasha.Lord)
//					if asp {
// 						reasons += BirthTimeRectificationReason_KarakaAspect
//					}
//					asps := karaka.GetAspect(subdasha.Lord)
//					if asp {
//						reasons += BirthTimeRectificationReason_KarakaAspect
//					}
//
//					// Minor reason to consider:
//					// - The karaka is in the important house for the varga chart
//					//   - This is a minor reason because the karaka is always in the important house
//					//     and the dasha doesn't change that: it just means that event is more
//					//     likely to be in that person's life
//
//					// Rule 4: The ascendant is in the important house
//					if dashaLordHouse == House1 || subdashaLordHouse == House1 {
//						reasons += BirthTimeRectificationReason_Ascendant
//					}
//
//					// Rule 5: The Rashi chart house lord becomes a dasha or sub-dasha lord
//					if rashiChartImpHouseLord == dasha.Lord || rashiChartImpHouseLord == subdasha.Lord {
//						reasons += BirthTimeRectificationReason_RashiChartHouseLord
//					}
//
//					// Rule 6: Bhavat Bhavam: the dasha or sub-dasha lord is placed in the xth house
//					// from the xth important house of the varga chart.
//					// For example, if the dasha or sub-dasha lord is placed in the 5th house from the 5th
//					// house (9th house) of a D7 chart (childbirth varga
//					// chart), there's a good chance that the event (i.e., childbirth)
//					// would happen
//					secondaryImpHouse := impHouse.Int() + (impHouse.Int()-1)
//					if dashaLordHouse == secondaryImpHouse || subdashaLordHouse == secondaryImpHouse {
//						reasons += BirthTimeRectificationReason_BhavatBhavam
//					}
//
//					return len(reasons), reasons
//				}
//
//  			Now we can run the above analysis every 5 minutes in a range
//
//  			bt := possible birth time
//  			range := maybe 1 hour before and after the birth time
//  			inputs := []EventInput{
//  				EventInput{
//  					EventType: MovingHouse,
//  					Time: ...,
//  				},
//  				EventInput{
//  					EventType: Marriage,
//  					Time: ...,
//  				},
//  				EventInput{
//  					EventType: ChildBirth,
//  					Time: ...,
//  				},
//  				EventInput{
//  					EventType: Education,
//  					Time: ...,
//  				},
//  				EventInput{
//  					EventType: Career,
//  					Time: ...,
//  				},
//  			}
//
//  			We want a data structure here that can hold the birth time and the number of reasons
//  			for each varga chart
//
//  			type birthTimeAnalysisChart struct {
//  				BirthTime time.Time
//  				ReasonsPerChartType map[VargaChartType][]BirthTimeRectificationReason
//  				Score int
//  			}
//
//  			analysisResults := []birthTimeAnalysisChart{}
//  			for range between bt - range and bt + range {
//                	step := can be 5 minutes
//  				potentialBirthTime := bt + step
//  				reasonsPerChartType := map[VargaChartType][]BirthTimeRectificationReason{}
//  				sumOfScores := 0
//  				for input in range inputs {
//  					reasons, score := analyzeBirthTime(
//  						potentialBirthTime,
//  						input.Time,
//  						input.EventType.VargaChartType(),
//  					)
//  					if reasons < 0 {
//  						continue
//  					}
//  					sumOfScores += score
//  					reasonsPerChartType[input.EventType.VargaChartType()] = reasons
//  				}
//  				if sumOfScores == 0 {
//  					continue
//  				}
//  				analysisResults = append(analysisResults, birthTimeAnalysisChart{
//  					BirthTime: potentialBirthTime,
//  					ReasonsPerChartType: map[VargaChartType][]BirthTimeRectificationReason{},
//  					Score: sumOfScores,
//  				})
//  			}
//
//  			Just organizing the results by score should be good enough for now
//
//  			Now, analysisResults contains a list of potential birth times and the reasons for each
//  			We need a way to display them
//  			See here https://pkg.go.dev/sort

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"slices"
	"strconv"
	"time"

	"github.com/afjoseph/sacredstar/chart"
	"github.com/afjoseph/sacredstar/wrapper"
)

// IntervalAnalysis is the analysis of different events in a time interval
type IntervalAnalysis struct {
	TotalScore    int                  `json:"totalScore"`
	Interval      chart.Interval       `json:"interval"`
	EventAnalyses []*EventAnalysis     `json:"eventAnalyses"`
	Ascendants    []*AscendantInterval `json:"ascendants"`
}

func NewIntervalAnalysis(
	totalScore int,
	interval chart.Interval,
	eventAnalyses []*EventAnalysis,
	ascendants []*AscendantInterval,
) *IntervalAnalysis {
	return &IntervalAnalysis{
		TotalScore:    totalScore,
		Interval:      interval,
		EventAnalyses: eventAnalyses,
		Ascendants:    ascendants,
	}
}

func (bta *IntervalAnalysis) Checksum() []byte {
	h32 := crc32.NewIEEE()
	h32.Write([]byte(strconv.Itoa(bta.TotalScore)))
	// XXX <15-02-2024, afjoseph> Interval sometimes gives a different checksum
	// for the same value
	// h32.Write(bta.Interval.Checksum())
	for _, ea := range bta.EventAnalyses {
		h32.Write(ea.Checksum())
	}
	return h32.Sum(nil)
}

type Rectification struct {
	Analyses []*IntervalAnalysis `json:"analyses"`
}

func (r *Rectification) Checksum() (string, error) {
	h32 := crc32.NewIEEE()
	for _, a := range r.Analyses {
		h32.Write(a.Checksum())
	}
	return hex.EncodeToString(h32.Sum(nil)), nil
}

// NewRectification returns a new Rectification
// The function loops over each step in the range and analyzes the birth time.
// It calculates a checksum per event analysis and checks if the checksum for a
// time interval is the same. If it is, it extends the interval. If it isn't,
// it creates a new IntervalAnalysis with the data we've collected from
// the previous iterations and starts a new interval.
func NewRectification(
	swe *wrapper.SwissEph,
	birthTime time.Time,
	lon, lat float64,
	_range time.Duration,
	step time.Duration,
	events []Event,
) (*Rectification, error) {
	// Every 'step' minutes, starting from '_range' minutes before and after
	// the birth time
	lastChecksum := uint32(0)
	lastEventAnalyses := []*EventAnalysis{}
	lastTotalScore := 0
	startOfInterval := birthTime.Add(-_range)
	endOfInterval := startOfInterval
	mbtas := []*IntervalAnalysis{}
	for t := birthTime.Add(-_range); t.Before(birthTime.Add(_range)); t = t.Add(step) {
		// Analyze events and collect score
		totalScore := 0
		analyses := []*EventAnalysis{}
		for _, e := range events {
			ret, err := e.analyze(swe, t, lon, lat)
			if err != nil {
				return nil, fmt.Errorf(
					"while analyzing birth time %v: %w",
					t,
					err,
				)
			}
			// If this event chart had no reasons, skip it
			if ret.Score == 0 {
				continue
			}
			totalScore += ret.Score
			analyses = append(analyses, ret)
		}

		// FOR DEBUGGING
		// fmt.Printf("==================== Time: %s | Score: %d\n", t, totalScore)
		// for _, ea := range analyses {
		// 	fmt.Printf("Event: %s | At: %s\n", ea.Event.Desc, ea.Event.Time)
		// 	for _, r := range ea.Reasons {
		// 		fmt.Printf("  - %s\n", r.String())
		// 	}
		// }
		// fmt.Println("====================")

		// Case 1: No relevant data?
		// There's nothing to collect: reset the interval
		// and other data.
		if totalScore == 0 {
			// If there's nothing in this analysis, reset the interval and
			// other data
			startOfInterval = t
			endOfInterval = t
			lastTotalScore = 0
			lastEventAnalyses = []*EventAnalysis{}
			lastChecksum = 0
			continue
		}

		// Calculate checksum (since at this stage we definitely have data)
		h32 := crc32.NewIEEE()
		for _, ea := range analyses {
			h32.Write(ea.Checksum())
		}
		currChecksum := h32.Sum32()

		// Case 2: 1st step in the interval?
		// Record the data and the checksum: it'll be used to compare with the
		// next step
		if startOfInterval == t {
			// Just collect the data
			lastTotalScore = totalScore
			lastEventAnalyses = analyses
			lastChecksum = currChecksum
			continue
		}

		// Case 3: Same checksum?
		// Extend the interval. We don't need to collect the data again because
		// it's the same as the last step
		if currChecksum == lastChecksum {
			// Extend the interval
			endOfInterval = t
			continue
		}

		// Case 4: Different checksum?
		// This means we're done with the old data and can process it since
		// this is a new interval with a new checksum.
		//
		// - Process the old data
		interval := chart.NewInterval(startOfInterval, endOfInterval)
		ascendants, err := analyzeAscendantsForTimeInterval(
			swe,
			events,
			interval,
			lon,
			lat,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"while analyzing ascendants for interval %v: %w",
				interval,
				err,
			)
		}
		mbtas = append(mbtas, NewIntervalAnalysis(
			lastTotalScore,
			interval,
			lastEventAnalyses,
			ascendants,
		))
		// - Start a new interval
		startOfInterval = t
		endOfInterval = t
		lastTotalScore = totalScore
		lastEventAnalyses = analyses
		lastChecksum = currChecksum
	}

	slices.SortFunc[[]*IntervalAnalysis](
		mbtas,
		func(lhs *IntervalAnalysis, rhs *IntervalAnalysis) int {
			// Always favor the analysis that covers all events
			lhsTotalEvents := len(lhs.EventAnalyses)
			rhsTotalEvents := len(rhs.EventAnalyses)
			if lhsTotalEvents < rhsTotalEvents {
				return -1
			}
			if lhsTotalEvents > rhsTotalEvents {
				return 1
			}

			// If both has the same number of events, sort by highest average
			lhsAvg := float64(lhs.TotalScore) / float64(len(lhs.EventAnalyses))
			rhsAvg := float64(rhs.TotalScore) / float64(len(rhs.EventAnalyses))
			if lhsAvg < rhsAvg {
				return -1
			}
			if lhsAvg > rhsAvg {
				return 1
			}

			// If the average is the same, sort by total score
			if lhs.TotalScore < rhs.TotalScore {
				return -1
			}
			if lhs.TotalScore > rhs.TotalScore {
				return 1
			}
			return 0
		},
	)
	slices.Reverse[[]*IntervalAnalysis](mbtas)
	return &Rectification{mbtas}, nil
}
