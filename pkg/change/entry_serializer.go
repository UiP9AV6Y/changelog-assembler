package change

import (
	"gopkg.in/yaml.v3"
)

type EntrySerializer struct{}

func (s *EntrySerializer) Serialize(entry interface{}) ([]byte, error) {
	return yaml.Marshal(entry)
}

func (s *EntrySerializer) Deserialize(data []byte) (interface{}, error) {
	entry := NewEntry()

	if err := yaml.Unmarshal(data, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

func NewEntrySerializer() *EntrySerializer {
	entrySerializer := &EntrySerializer{}

	return entrySerializer
}
