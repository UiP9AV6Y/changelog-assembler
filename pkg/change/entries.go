package change

import (
	"sort"
	"strings"
)

type Entries []*Entry

func (e Entries) Authors() []string {
	seen := map[string]bool{}
	authors := []string{}

	for _, entry := range e {
		if entry.IsAnonymous() {
			continue
		}
		if _, ok := seen[entry.Author]; !ok {
			seen[entry.Author] = true
			authors = append(authors, entry.Author)
		}
	}

	sort.Strings(authors)
	return authors
}

func (e Entries) SortByReason() Entries {
	sort.Slice(e, func(i, j int) bool {
		return e[i].Reason < e[j].Reason
	})

	return e
}

func (e Entries) SortByMergeRequest() Entries {
	sort.Slice(e, func(i, j int) bool {
		return e[i].MergeRequest < e[j].MergeRequest
	})

	return e
}

func (e Entries) SortByAuthor() Entries {
	sort.Slice(e, func(i, j int) bool {
		return strings.Compare(e[i].Author, e[j].Author) < 0
	})

	return e
}

func (e Entries) SortByComponent() Entries {
	sort.Slice(e, func(i, j int) bool {
		return strings.Compare(e[i].Component, e[j].Component) < 0
	})

	return e
}

func (e Entries) SortByTitle() Entries {
	sort.Slice(e, func(i, j int) bool {
		return strings.Compare(e[i].Title, e[j].Title) < 0
	})

	return e
}

func (e Entries) Filter(filter func(*Entry) bool) Entries {
	entries := make(Entries, 0, len(e))

	for _, entry := range e {
		if filter(entry) {
			entries = append(entries, entry)
		}
	}

	return entries
}

func (e Entries) FilterByReason(reason Reason) Entries {
	return e.Filter(func(entry *Entry) bool {
		return entry.Reason == reason
	})
}

func (e Entries) FilterByComponent(component string) Entries {
	return e.Filter(func(entry *Entry) bool {
		return entry.Component == component
	})
}

func (e Entries) GroupByComponent() map[string]Entries {
	var entries Entries
	var ok bool

	grouped := map[string]Entries{}

	for _, entry := range e {
		if entries, ok = grouped[entry.Component]; ok {
			grouped[entry.Component] = append(entries, entry)
		} else {
			grouped[entry.Component] = []*Entry{entry}
		}
	}

	return grouped
}

func (e Entries) GroupByReason() map[Reason]Entries {
	var entries Entries
	var ok bool

	grouped := map[Reason]Entries{}

	for _, entry := range e {
		if entries, ok = grouped[entry.Reason]; ok {
			grouped[entry.Reason] = append(entries, entry)
		} else {
			grouped[entry.Reason] = []*Entry{entry}
		}
	}

	return grouped
}
