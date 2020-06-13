package change

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntrySlug(t *testing.T) {
	testCases := map[string]*Entry{
		"hello-world": {
			Title: "Hello World",
		},
		"123-hello-world": {
			Title:        "Hello World",
			MergeRequest: 123,
		},
		"hello--world": {
			Title: "Hello, World",
		},
		"10-hello--world": {
			Title:        "[Hello, World]",
			MergeRequest: 10,
		},
	}

	for want, unit := range testCases {
		assert.Equal(t, want, unit.Slug(), "slug rendered correctly")
	}
}
