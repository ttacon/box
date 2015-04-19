package box

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/go-querystring/query"
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

type EventQueryOptions struct {
	StreamPosition string `url:"stream_position"`
	StreamType     string `url:"stream_type"`
	Limit          int    `url:"limit"`
}

// Events retrieves events for the currently authenticated user.
//
// See: https://developers.box.com/docs/#events-get-events-for-a-user
func (e *EventService) Events(options EventQueryOptions) (*http.Response, *EventCollection, error) {
	queryString, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	req, err := e.NewRequest(
		"GET",
		"/events?"+queryString.Encode(),
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
		if eve.Source == nil {
			continue
		}
		if err = finishParsing(eve); err != nil {
			return resp, nil, err
		}
	}
	return resp, &data, err
}

type LongPollInfo struct {
	ChunkSize int                `json:"chunk_size"`
	Entries   []LongPollConnInfo `json:"entries"`
}

type LongPollConnInfo struct {
	Type         string `json:"type"`
	URL          string `json:"url"`
	TTL          string `json:"ttl"`
	MaxRetries   string `json:"max_retries"`
	RetryTimeout int    `json:"retry_timeout"`
}

func (e *EventService) LongPollURL() (*http.Response, *LongPollInfo, error) {
	req, err := e.NewRequest(
		"OPTIONS",
		"/events",
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data LongPollInfo
	resp, err := e.Do(req, &data)
	return resp, &data, err

}

func (e *EventService) ListenForEvent(i LongPollConnInfo, lastSync string) (*http.Response, []*Event, error) {
	// get events for stream_position=now for sync token
	var streamPos = lastSync
	if len(streamPos) == 0 {
		resp, events, err := e.Events(EventQueryOptions{
			StreamPosition: "now",
		})
		if err != nil {
			return resp, nil, err
		}

		streamPos = strconv.Itoa(events.NextStreamPosition)
	}

	// TODO(ttacon): timeout info
	req, err := http.NewRequest(
		"GET",
		i.URL,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var d = make(map[string]interface{})
	resp, err := e.Do(req, &d)
	// TODO(ttacon): deal with timeouts
	if err != nil {
		return resp, nil, err
	}

	// get next event with last synced note
	resp, events, err := e.Events(EventQueryOptions{
		StreamPosition: streamPos,
		Limit:          1,
	})
	if err != nil {
		return resp, nil, err
	}
	return resp, events.Entries, err
}

func (e *EventService) Channel(size int) chan *Event {
	eventStream := make(chan *Event, size)
	go e.streamEvents(eventStream)
	return eventStream
}

func (e *EventService) streamEvents(tunnel chan *Event) {
	for {
		_, longPollInfo, err := e.LongPollURL()
		if err != nil {
			close(tunnel)
			return
		}

		if len(longPollInfo.Entries) != 1 {
			close(tunnel)
			return
		}

		_, events, err := e.ListenForEvent(longPollInfo.Entries[0], "")
		if err != nil {
			close(tunnel)
			return
		}

		for _, eve := range events {
			tunnel <- eve
		}
	}
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
