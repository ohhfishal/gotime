package event

import (
	"encoding/json"
	"time"
)

type Tag string
type Status string

func (s Status) Valid() error {
	return nil
}

type Event struct {
	Title       string
	Description string
	Start       time.Time
	End         time.Time
	Tags        []Tag
	Status      Status
}

func Decode(line string) (Event, error) {
	var event Event
	if err := json.Unmarshal([]byte(line), &event); err != nil {
		return Event{}, err
	}
	return event, nil
}

func (event Event) Encode() (string, error) {
	bytes, err := json.Marshal(event)
	if err != nil {
		return ``, err
	}
	return string(bytes), nil
}
