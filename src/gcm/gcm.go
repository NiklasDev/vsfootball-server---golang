package gcm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const (
	// GcmSendEndpoint is the endpoint for sending messages to the GCM server.
	GcmSendEndpoint = "https://android.googleapis.com/gcm/send"
	// Initial delay before first retry, without jitter.
	backoffInitialDelay = 1000
	// Maximum delay before a retry.
	maxBackoffDelay = 1024000
)

type Sender struct {
	ApiKey string
	//    ApiKey string
	Http *http.Client
	//     Http *http.Sender
}

func New(key string) *Sender {
	return &Sender{
		ApiKey: key,
		Http:   new(http.Client),
	}
}

func (s *Sender) SendNoRetry(msg *Message) (*Response, error) {
	if err := checkSender(s); err != nil {
		return nil, err
	} else if err := checkMessage(msg); err != nil {
		return nil, err
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", GcmSendEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// fmt.Println("Api_key is ", s.ApiKey)
	req.Header.Add("Authorization", fmt.Sprintf("key=%s", s.ApiKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.Http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d error: %s", resp.StatusCode, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := new(Response)
	err = json.Unmarshal(body, response)
	return response, err
}

func (s *Sender) Send(msg *Message, retries int) (*Response, error) {
	if err := checkSender(s); err != nil {
		return nil, err
	} else if err := checkMessage(msg); err != nil {
		return nil, err
	} else if retries < 0 {
		return nil, errors.New("'retries' must not be negative.")
	}

	// Send the message for the first time
	resp, err := s.SendNoRetry(msg)
	if err != nil {
		return nil, err
	} else if resp.Failure == 0 || retries == 0 {
		return resp, nil
	}

	// One or more messages failed to send
	var regIds = msg.RegistrationIDs
	var allResults = make(map[string]Result, len(regIds))
	var backoff = backoffInitialDelay
	for i := 0; updateStatus(msg, resp, allResults) > 0 || i < retries; i++ {
		sleepTime := backoff/2 + rand.Intn(backoff)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		backoff = min(2*backoff, maxBackoffDelay)
		if resp, err = s.SendNoRetry(msg); err != nil {
			return nil, err
		}
	}

	// Bring the message back to its original state
	msg.RegistrationIDs = regIds

	// Create a Response containing the overall results
	var success, failure, canonicalIds int
	var finalResults = make([]Result, len(regIds))
	for i := 0; i < len(regIds); i++ {
		result, _ := allResults[regIds[i]]
		finalResults[i] = result
		if result.MessageID != "" {
			if result.RegistrationID != "" {
				canonicalIds++
			}
			success++
		} else {
			failure++
		}
	}

	return &Response{
		// return the most recent multicast id
		MulticastID:  resp.MulticastID,
		Success:      success,
		Failure:      failure,
		CanonicalIDs: canonicalIds,
		Results:      finalResults,
	}, nil
} //end func

func updateStatus(msg *Message, resp *Response, allResults map[string]Result) int {
	var unsentRegIds = make([]string, 0, resp.Failure)
	for i := 0; i < len(resp.Results); i++ {
		regId := msg.RegistrationIDs[i]
		allResults[regId] = resp.Results[i]
		if resp.Results[i].Error == "Unavailable" {
			unsentRegIds = append(unsentRegIds, regId)
		}
	}
	msg.RegistrationIDs = unsentRegIds
	return len(unsentRegIds)
}

// min returns the smaller of two integers. For exciting religious wars
// about why this wasn't included in the "math" package, see this thread:
// https://groups.google.com/d/topic/golang-nuts/dbyqx_LGUxM/discussion
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func checkSender(sender *Sender) error {
	if len(sender.ApiKey) == 0 {
		return errors.New("The sender's API key must not be empty.")
	}
	if sender.Http == nil {
		sender.Http = new(http.Client)
	}
	return nil
}

// checkMessage returns an error if the message is not well-formed.
func checkMessage(msg *Message) error {
	if msg == nil {
		return errors.New("The message must not be nil.")
	} else if msg.RegistrationIDs == nil {
		return errors.New("The message's RegistrationIDs field must not be nil.")
	} else if len(msg.RegistrationIDs) == 0 {
		return errors.New("The message must specify at least one registration ID.")
	} else if len(msg.RegistrationIDs) > 1000 {
		return errors.New("The message may specify at most 1000 registration IDs.")
	} else if msg.TimeToLive < 0 || 2419200 < msg.TimeToLive {
		return errors.New("The message's TimeToLive field must be an integer " +
			"between 0 and 2419200 (4 weeks).")
	}
	return nil
}

type Response struct {
	MulticastID  int64    `json:"multicast_id"`
	Success      int      `json:"success"`
	Failure      int      `json:"failure"`
	CanonicalIDs int      `json:"canonical_ids"`
	Results      []Result `json:"results"`
}

// Result represents the status of a processed message.
type Result struct {
	MessageID      string `json:"message_id"`
	RegistrationID string `json:"registration_id"`
	Error          string `json:"error"`
}
