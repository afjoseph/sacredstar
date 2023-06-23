package timeandzone

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

var timeFormat = "2006-01-02T15:04"

// TimeAndZone is basically time.Time that JSON-serializes to a custom time format
// and a timezone. We **can** use a known time format that accounts for
// timezones like RFC3339, but dayjs (the JS library we use) doesn't handle
// timezones well and TBH I'm bored of timezone bullshit
type TimeAndZone struct {
	time.Time
}

func New(t time.Time) TimeAndZone {
	return TimeAndZone{t}
}

func (t TimeAndZone) MarshalJSON() ([]byte, error) {
	d := t.Time.Format(timeFormat)
	tz := t.Time.Location().String()
	return []byte(fmt.Sprintf(`{"t":"%s","z":"%s"}`, d, tz)), nil
}

func (t *TimeAndZone) UnmarshalJSON(b []byte) error {
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}
	tz, ok := obj["z"].(string)
	if !ok {
		return fmt.Errorf("No timezone in TimeAndZone")
	}
	tzLoc, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("Invalid timezone: %s", tz)
	}
	d, ok := obj["t"].(string)
	if !ok {
		return fmt.Errorf("No time in TimeAndZone")
	}
	// d, err = strconv.Unquote(d)
	// if err != nil {
	// 	return fmt.Errorf("Invalid time: %s", d)
	// }
	dAsTime, err := time.ParseInLocation(timeFormat, d, tzLoc)
	if err != nil {
		return fmt.Errorf("Invalid time: %s", d)
	}
	*t = TimeAndZone{dAsTime}
	return nil
}

func (t TimeAndZone) String() string {
	return t.Time.String()
}

func (t TimeAndZone) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *TimeAndZone) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), t)
}
