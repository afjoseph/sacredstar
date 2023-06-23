package chart

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hako/durafmt"
	"github.com/stretchr/testify/assert"
)

// TestGetAntardashas is a sanity check to make sure the sum of Antardasha
// periods is equal to the Mahadasha period
func TestGetAntardashas(t *testing.T) {
	for _, mdl := range MahadashaList {
		totalDuration := mdl.MahadashaDuration()

		sum := time.Duration(0)
		for _, adl := range mdl.GetAntardashas() {
			sum += AntardashaDuration(mdl, adl)
		}

		diff := totalDuration - sum
		// 2 day epsilon max
		epsillon := time.Hour * 24 * 2
		// fmt.Printf("diff: %v\n", durafmt.Parse(diff))
		assert.True(
			t,
			diff < epsillon,
			"expected %s to be equal to %s",
			durafmt.Parse(totalDuration),
			durafmt.Parse(sum),
		)
	}
}

func TestJSONSerialization(t *testing.T) {
	for _, mdl := range MahadashaList {
		b, err := json.Marshal(mdl)
		assert.NoError(t, err)
		assert.Equal(t, `"`+mdl.ToPointID().String()+`"`, string(b))
		var mdl2 DashaLord
		err = json.Unmarshal(b, &mdl2)
		assert.NoError(t, err)
		assert.Equal(t, mdl, mdl2)
	}
}
