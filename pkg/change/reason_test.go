package change

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsZero(t *testing.T) {
	testCases := map[Reason]bool{
		ReasonInvalid:     true,
		ReasonAdded:       false,
		ReasonFixed:       false,
		ReasonChanged:     false,
		ReasonDeprecated:  false,
		ReasonRemoved:     false,
		ReasonSecurity:    false,
		ReasonPerformance: false,
		ReasonOther:       false,
	}

	for unit, want := range testCases {
		assert.Equal(t, want, unit.IsZero(), "zeroes detected")
	}
}
