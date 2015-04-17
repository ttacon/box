package box

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type EventService struct {
	*Client
}

type Event struct {
	Type          string          `json:"type"`
	EventID       string          `json:"event_id"`
	CreatedBy     *User           `json:"created_by"`
	EventType     string          `json:"event_type"`
	SessionID     string          `json:"session_id"`
	Source        json.RawMessage `json:"source"`
	File          *File           `json:"-"`
	Folder        *Folder         `json:"-"`
	Comment       *Comment        `json:"-"`
	Collaboration *Collaboration  `json:"-"`
	// Test and add the rest
}

type EventCollection struct {
	Entries            []*Event `json:"entries"`
	ChunkSize          int      `json:"chunk_size"`
	NextStreamPosition int      `json:"next_stream_position"`
}

func (e *EventService) Events() (*http.Response, *EventCollection, error) {
	req, err := e.NewRequest(
		"GET",
		"/events",
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	var data EventCollection
	resp, err := e.Do(req, &data)
	if err != nil {
		return resp, nil, err
	}
	for _, eve := range data.Entries {
		if err = finishParsing(eve); err != nil {
			return resp, nil, err
		}
	}
	return resp, &data, err
}

////////// types //////////

////////// functions //////////

////////// enum funcs, whooo! //////////

func finishParsing(ev *Event) error {
	switch eventSourceType(ev) {
	case "file":
		ev.File = &File{}
		err := json.Unmarshal(ev.Source, ev.File)
		if err != nil {
			return err
		}
		ev.Source = nil
	case "folder":
		ev.Folder = &Folder{}
		err := json.Unmarshal(ev.Source, ev.Folder)
		if err != nil {
			return err
		}
		ev.Source = nil
	case "comment":
		ev.Comment = &Comment{}
		err := json.Unmarshal(ev.Source, ev.Comment)
		if err != nil {
			return err
		}
		ev.Source = nil
	case "collaboration":
		ev.Collaboration = &Collaboration{}
		err := json.Unmarshal(ev.Source, ev.Collaboration)
		if err != nil {
			return err
		}
		ev.Source = nil
	default:
		fmt.Println(eventSourceType(ev))
		return errors.New("not implemented yet (read I'm lazy :P )")
	}
	return nil
}

func eventSourceType(ev *Event) string {
	// ugly hack, but performant enough until we generate and then
	// manipulate unmarshaling code
	type EventSourceType struct {
		Type string `json:"type"`
	}

	var est EventSourceType
	// explicit swallow the error
	_ = json.Unmarshal(ev.Source, &est)
	return est.Type
}
