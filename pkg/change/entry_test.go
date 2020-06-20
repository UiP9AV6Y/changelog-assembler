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
		"hello-duplicates": {
			Title: "Hello, Duplicates",
		},
		"10-hello-world": {
			Title:        "[Hello, World]",
			MergeRequest: 10,
		},
	}

	for want, unit := range testCases {
		assert.Equal(t, want, unit.Slug(), "slug rendered correctly")
	}
}

func TestSlug(t *testing.T) {
	testCases := map[string]string{
		"hello-world":    "Hello World",
		"spaces-trimmed": "  Spaces=trimmed  ",
		"a-b":            "/a, B\\",
		"a":              "__a//////////////",
		"":               "",
	}

	for want, unit := range testCases {
		assert.Equal(t, want, Slug(unit), "slug rendered correctly")
	}
}
