package change

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ReasonInvalid Reason = iota
	ReasonAdded
	ReasonFixed
	ReasonChanged
	ReasonDeprecated
	ReasonRemoved
	ReasonSecurity
	ReasonPerformance
	ReasonOther
)

var (
	ReasonNames = []string{
		"",
		"added",
		"fixed",
		"changed",
		"deprecated",
		"removed",
		"security",
		"performance",
		"other",
	}
	ReasonDescriptions = []string{
		"",
		"New Features",
		"Fixed Issues",
		"Changed Features",
		"Deprecated Features",
		"Removed Features",
		"Security Enhancements",
		"Performance Enhancements",
		"Other Changes",
	}
	ReasonInstances = []Reason{
		ReasonInvalid,
		ReasonAdded,
		ReasonFixed,
		ReasonChanged,
		ReasonDeprecated,
		ReasonRemoved,
		ReasonSecurity,
		ReasonPerformance,
		ReasonOther,
	}
)

type Reason int

func (r Reason) String() string {
	return ReasonNames[r]
}

func (r Reason) Description() string {
	return ReasonDescriptions[r]
}

func (r Reason) CapsString() string {
	return strings.ToUpper(ReasonNames[r])
}

func (r Reason) TitleString() string {
	return strings.Title(ReasonNames[r])
}

func (r Reason) Alias() string {
	if r.IsZero() {
		return ""
	}

	return ReasonNames[r][0:1]
}

func (r Reason) IsZero() bool {
	return r <= 0
}

func (r Reason) MarshalYAML() (interface{}, error) {
	return ReasonNames[r], nil
}

func (r *Reason) UnmarshalYAML(value *yaml.Node) error {
	if v, err := ParseReason(value.Value); err != nil {
		return err
	} else {
		*r = v
	}

	return nil
}

func (r *Reason) UnmarshalJSON(value []byte) error {
	if v, err := ParseReason(string(value)); err != nil {
		return err
	} else {
		*r = v
	}

	return nil
}

func (r *Reason) UnmarshalText(value []byte) error {
	return r.UnmarshalJSON(value)
}

func ParseReason(value string) (Reason, error) {
	value = strings.ToLower(value) // case insensitive comparison

	for i, name := range ReasonNames {
		if name == value {
			return ReasonInstances[i], nil
		}
	}

	return ReasonInvalid, fmt.Errorf("invalid entry reason %q", value)
}
