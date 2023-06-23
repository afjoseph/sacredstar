package sign

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignJSONSerialization(t *testing.T) {
	s := Aries
	b, err := json.Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t, `"aries"`, string(b))
	var s2 Sign
	err = json.Unmarshal(b, &s2)
	assert.NoError(t, err)
	assert.Equal(t, s, s2)

	err = json.Unmarshal([]byte(`"invalid"`), &s2)
	assert.Error(t, err)
}
