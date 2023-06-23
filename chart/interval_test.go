package chart

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntervalJSONSerialization(t *testing.T) {
	i := NewInterval(
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
	)
	b, err := json.Marshal(i)
	assert.NoError(t, err)
	assert.Equal(
		t,
		`{"start":{"t":"2020-01-01T00:00","z":"UTC"},"end":{"t":"2020-01-02T00:00","z":"UTC"}}`,
		string(b),
	)
	var i2 Interval
	err = json.Unmarshal(b, &i2)
	assert.NoError(t, err)
	assert.Equal(t, i, i2)

	err = json.Unmarshal([]byte(`"invalid"`), &i2)
	assert.Error(t, err)
}
