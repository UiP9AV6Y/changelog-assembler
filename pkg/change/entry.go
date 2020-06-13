package change

import (
	"os/user"
	"strconv"
	"strings"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/util"
)

const slugFiller = "-"

type Entry struct {
	Title        string `yaml:"title"`
	Author       string `yaml:"author"`
	Reason       Reason `yaml:"reason"`
	Component    string `yaml:"component,omitempty"`
	MergeRequest int    `yaml:"merge_request,omitempty"`
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
	var slug string
	slugger := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return r + 32
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		}
		return rune(slugFiller[0])
	}
	slug = strings.Trim(strings.Map(slugger, e.Title), slugFiller)

	if e.MergeRequest > 0 {
		return e.MergeRequestString() + slugFiller + slug
	}

	return slug
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
	entry := &Entry{}

	return entry
}
