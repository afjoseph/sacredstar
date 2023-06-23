package timeandzone

import (
	"encoding/json"
	"fmt"
	"strings"
)

var Zones zones

type zones struct {
	momentTimezoneDB map[string]interface{}
}

func init() {
	if err := json.Unmarshal(
		[]byte(MomentDBZonesJSONLiteral),
		&Zones.momentTimezoneDB,
	); err != nil {
		panic(fmt.Errorf("while unmarshalling moment-timezone DB: %v", err))
	}
}

func (z zones) GetLonLatFromTimezone(
	timezone string,
) (float64, float64, error) {
	if strings.ToLower(timezone) == "utc" {
		return 0, 0, nil
	}
	if len(z.momentTimezoneDB) == 0 {
		panic("moment-timezone DB not initialized")
	}

	zs, ok := z.momentTimezoneDB["zones"]
	if !ok {
		panic("timezones not found in moment-timezone DB")
	}
	zsMap, ok := zs.(map[string]interface{})
	if !ok {
		panic("timezones not a map in moment-timezone DB")
	}

	targetTZ, ok := zsMap[timezone]
	if !ok {
		return 0, 0, fmt.Errorf(
			"timezone %s not found in moment-timezone DB",
			timezone,
		)
	}
	targetTZAsMap, ok := targetTZ.(map[string]interface{})
	if !ok {
		return 0, 0, fmt.Errorf(
			"timezone %s not a map in moment-timezone DB",
			timezone,
		)
	}
	lon, ok := targetTZAsMap["long"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf(
			"lon not found in timezone %s in moment-timezone DB",
			timezone,
		)
	}
	lat, ok := targetTZAsMap["lat"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf(
			"lat not found in timezone %s in moment-timezone DB",
			timezone,
		)
	}
	return lon, lat, nil
}
