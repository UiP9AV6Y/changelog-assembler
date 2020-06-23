package change

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAnnotations(t *testing.T) {
	testCases := []string{
		"",
		"=",
		"test",
		"=test",
		"unit=",
		"unit=test",
		"unit=test=test",
	}
	wantKeys := []string{
		"",
		"",
		"test",
		"",
		"unit",
		"unit",
		"unit",
	}
	wantValues := []string{
		"",
		"",
		"",
		"test",
		"",
		"test",
		"test=test",
	}

	for i, test := range testCases {
		have := []string{test}
		got := ParseAnnotations(have)
		wantKey := wantKeys[i]
		wantValue := wantValues[i]
		value, ok := got[wantKey]

		assert.Equal(t, len(got), 1)
		assert.Truef(t, ok, "found key %q", wantKey)
		assert.Equal(t, value, wantValue, "got value")
	}
}
