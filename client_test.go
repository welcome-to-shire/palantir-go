package palantir

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func createTestServer(statusCode int, body string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, body)
	}))

	return ts
}

func createTestMessage(title, content string) *Message {
	msg := new(Message)
	msg.Title = title
	msg.Content = content
	msg.CreatedAt = ISO8601Time(time.Now())

	return msg
}

func marshalTestMessage(title, content string, t *testing.T) string {
	msg := createTestMessage(title, content)

	m, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	return string(m)
}

func marshalTestTicket(value string, t *testing.T) string {
	ticket := Ticket{value}

	r, err := json.Marshal(ticket)
	if err != nil {
		t.Fatal(err)
	}

	return string(r)
}

func TestGetMessage(t *testing.T) {
	tests := []struct {
		code       int
		body       string
		hasMessage bool
	}{
		{200, marshalTestMessage("", "", t), true},
		{400, "", false},
		{500, "", false},
	}

	for _, test := range tests {
		server := createTestServer(test.code, test.body)
		client := MakeClient(server.URL[7:])
		_, err := client.GetMessage("test")
		if test.hasMessage && err != nil {
			t.Errorf("expected get message : %v", err)
		}
		if !test.hasMessage && err == nil {
			t.Errorf("should not return message")
		}
		server.Close()
	}
}

func TestCreateMessage(t *testing.T) {
	tm := createTestMessage("", "")
	tests := []struct {
		code      int
		body      string
		msg       *Message
		isSuccess bool
	}{
		{201, marshalTestTicket("", t), tm, true},
		{400, "", tm, false},
		{500, "", tm, false},
	}

	for _, test := range tests {
		server := createTestServer(test.code, test.body)
		client := MakeClient(server.URL[7:])
		_, err := client.CreateMessage("test", *test.msg)
		if test.isSuccess && err != nil {
			t.Errorf("expected create message: %v", err)
		}
		if !test.isSuccess && err == nil {
			t.Errorf("should not create message")
		}
		server.Close()
	}
}
