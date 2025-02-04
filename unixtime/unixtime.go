package unixtime

import (
	"database/sql/driver"
	"encoding/json"
	"log/slog"
	"time"
)

// UnixTime is a wrapper around time.Time that serializes as a Unix timestamp
type UnixTime struct {
	time.Time
}

func (ut UnixTime) LogValue() slog.Value {
	return slog.StringValue(ut.Format(time.RFC3339))
}

// MarshalJSON implements json.Marshaler
func (ut UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ut.Unix())
}

// UnmarshalJSON implements json.Unmarshaler
func (ut *UnixTime) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	ut.Time = time.Unix(timestamp, 0)
	return nil
}

// Value implements driver.Valuer for database serialization
func (ut UnixTime) Value() (driver.Value, error) {
	return ut.Unix(), nil
}

// Scan implements sql.Scanner for database deserialization
func (ut *UnixTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case int64:
		ut.Time = time.Unix(v, 0)
	case float64:
		ut.Time = time.Unix(int64(v), 0)
	case []byte:
		// Handle string/[]byte timestamp
		if timestamp, err := time.Parse(time.RFC3339, string(v)); err == nil {
			ut.Time = timestamp
		} else if unix, err := time.Parse("2006-01-02 15:04:05", string(v)); err == nil {
			ut.Time = unix
		}
	case string:
		if timestamp, err := time.Parse(time.RFC3339, v); err == nil {
			ut.Time = timestamp
		} else if unix, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			ut.Time = unix
		}
	case time.Time:
		ut.Time = v
	}
	return nil
}

// NewUnixTime creates a new UnixTime from a time.Time
func New(t time.Time) UnixTime {
	return UnixTime{Time: t}
}

// UnixTimeNow returns a new UnixTime with the current time
func Now() UnixTime {
	return UnixTime{Time: time.Now()}
}
