package change

import (
	"os/user"
	"strconv"
	"strings"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/util"
)

const slugFiller = "-"

type Entry struct {
	Title        string            `yaml:"title"`
	Author       string            `yaml:"author"`
	Reason       Reason            `yaml:"reason"`
	Component    string            `yaml:"component,omitempty"`
	MergeRequest int               `yaml:"merge_request,omitempty"`
	Annotations  map[string]string `yaml:"annotations,omitempty"`
}

func (e *Entry) String() string {
	if e.MergeRequest <= 0 {
		return e.Title
	}

	return e.Title + " #" + strconv.Itoa(e.MergeRequest)
}

func (e *Entry) MergeRequestString() string {
	if e.MergeRequest <= 0 {
		return ""
	}

	return strconv.Itoa(e.MergeRequest)
}

func (e *Entry) IsZero() bool {
	return e.MergeRequest <= 0
}

func (e *Entry) IsAnonymous() bool {
	return len(e.Author) == 0
}

func (e *Entry) Slug() string {
	if e.MergeRequest > 0 {
		return e.MergeRequestString() + slugFiller + Slug(e.Title)
	}

	return Slug(e.Title)
}

func (e *Entry) GetAnnotation(key, fallback string) string {
	if value, ok := e.Annotations[key]; ok {
		return value
	}

	return fallback
}

func Slug(title string) string {
	fill := false
	slugger := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			fill = true
			return r + 32
		case r >= 'a' && r <= 'z':
			fill = true
			return r
		case r >= '0' && r <= '9':
			fill = true
			return r
		case fill:
			fill = false
			return rune(slugFiller[0])
		}

		return -1
	}

	return strings.TrimRight(strings.Map(slugger, title), slugFiller)
}

func DefaultAuthor() string {
	if name := util.GitUsername(); len(name) > 0 {
		return name
	}

	if user, err := user.Current(); err != nil {
		return user.Username
	}

	return ""
}

func NewEntry() *Entry {
	author := DefaultAuthor()
	entry := &Entry{
		Author: author,
	}

	return entry
}
