package rectification

import (
	"fmt"
	"hash/crc32"

	"github.com/afjoseph/sacredstar/util"
)

type ReasonType uint32

const (
	ReasonType_House ReasonType = iota
	ReasonType_DashaLordAsHouseLord
	ReasonType_HouseLordAspect_Hard
	ReasonType_HouseLordAspect_Soft
	ReasonType_KarakaAspect_Hard
	ReasonType_KarakaAspect_Soft
	ReasonType_Ascendant
	ReasonType_RashiChartHouseLord
	ReasonType_BhavatBhavam
)

type Reason struct {
	Type        ReasonType `json:"type"`
	Description string     `json:"desc"`
}

func (r Reason) Checksum() []byte {
	h32 := crc32.NewIEEE()
	h32.Write(util.Uint32ToByteSlice(uint32(r.Type)))
	// XXX <08-02-2024, afjoseph> Would be wise to include this for better
	// accuracy but thsi also means any change in aspect degree (which is
	// included in the description) would yield a "different" reason, even
	// though it's the same. I'll test removing this for now and we can restore
	// it with a less-detailed description later
	// h32.Write([]byte(r.Description))
	return h32.Sum(nil)
}

func NewReason(
	type_ ReasonType,
	description string,
) Reason {
	return Reason{
		Type:        type_,
		Description: description,
	}
}

func (b ReasonType) String() string {
	switch b {
	case ReasonType_House:
		return "House"
	case ReasonType_DashaLordAsHouseLord:
		return "Dasha Lord As House Lord"
	case ReasonType_HouseLordAspect_Hard:
		return "House Lord Aspect - Hard"
	case ReasonType_HouseLordAspect_Soft:
		return "House Lord Aspect - Soft"
	case ReasonType_KarakaAspect_Hard:
		return "Karaka Aspect - Hard"
	case ReasonType_KarakaAspect_Soft:
		return "Karaka Aspect - Soft"
	case ReasonType_Ascendant:
		return "Ascendant"
	case ReasonType_RashiChartHouseLord:
		return "Rashi Chart House Lord"
	case ReasonType_BhavatBhavam:
		return "Bhavat Bhavam"
	default:
		return "Unknown"
	}
}

func (b Reason) String() string {
	return fmt.Sprintf("%s: %s", b.Type, b.Description)
}
