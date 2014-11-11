package palantir

import (
	"encoding/json"
	"fmt"
	"time"
)

type ISO8601Time time.Time

func (t ISO8601Time) MarshalJSON() ([]byte, error) {
	timeString := fmt.Sprintf(
		"\"%s\"",
		time.Time(t).Format(time.RFC3339Nano),
	)
	return []byte(timeString), nil
}

func (t *ISO8601Time) UnmarshalJSON(raw []byte) error {
	timeFmt := fmt.Sprintf("\"%s\"", time.RFC3339Nano)

	parsed, err := time.Parse(timeFmt, string(raw))
	if err != nil {
		return err
	}
	*t = ISO8601Time(parsed)
	return nil
}

type Jsonable interface {
	Marshal() ([]byte, error)
	MustMarshal() []byte
}

type Message struct {
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	CreatedAt ISO8601Time `json:"created_at,omitempty"`
}

func makeMessageFromServerResponse(raw []byte) (*Message, error) {
	m := new(Message)
	err := json.Unmarshal(raw, m)

	return m, err
}

func (m Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m Message) MustMarshal() []byte {
	b, err := m.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}

type Ticket struct {
	Id string `json:"id"`
}

func makeTicketFromServerResponse(raw []byte) (*Ticket, error) {
	t := new(Ticket)
	err := json.Unmarshal(raw, t)

	return t, err
}

func (t Ticket) Marshal() ([]byte, error) {
	return json.Marshal(t)
}

func (t Ticket) MustMarshal() []byte {
	b, err := t.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}

type Error struct {
	Reason string `json"error"`
}

func makeErrorFromServerResponse(raw []byte) (*Error, error) {
	e := new(Error)
	err := json.Unmarshal(raw, e)

	return e, err
}

func (e Error) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e Error) MustMarshal() []byte {
	b, err := e.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}

func (e Error) Error() string {
	return e.Reason
}
