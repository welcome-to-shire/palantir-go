// API client.
package palantir

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	baseUrl string
}

func MakeClient(host string) *Client {
	client := new(Client)

	client.baseUrl = fmt.Sprintf("http://%s/api/%s", host, ApiVersion)

	return client
}

func (c Client) fullUrl(part string) string {
	return fmt.Sprintf("%s/%s", c.baseUrl, part)
}

func (c Client) GetMessage(subject string) (*Message, error) {
	fullUrl := c.fullUrl(subject)
	resp, err := http.Get(fullUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 == 2 {
		message, err := makeMessageFromServerResponse(body)
		if err != nil {
			return nil, err
		}

		return message, nil
	}

	errorMessage, err := makeErrorFromServerResponse(body)
	if err != nil {
		return nil, err
	}
	return nil, errorMessage
}

func (c Client) CreateMessage(subject string, message Message) (*Ticket, error) {
	payload := message.MustMarshal()
	req, err := http.NewRequest("POST", c.fullUrl(subject), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 == 2 {
		ticket, err := makeTicketFromServerResponse(body)
		if err != nil {
			return nil, err
		}

		return ticket, nil
	}

	errorMessage, err := makeErrorFromServerResponse(body)
	if err != nil {
		return nil, err
	}
	return nil, errorMessage
}
