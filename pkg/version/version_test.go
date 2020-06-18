package version

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInfoString(t *testing.T) {
	version = "1.2.3"
	commit = "mock"
	date = "2009-02-13T23:31:30Z"

	unit := NewInfo()
	reference := time.Unix(1234567890, 0)

	assert.Equal(t, "1.2.3-mock", unit.String())
	assert.Equal(t, reference.String(), unit.BuildDate.String())
}

func TestInfoInvalidDate(t *testing.T) {
	version = "1.2.3"
	commit = "mock"
	date = "bac gaff"

	unit := NewInfo()
	reference := time.Unix(0, 0)

	assert.Equal(t, reference.String(), unit.BuildDate.String())
}
