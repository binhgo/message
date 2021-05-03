package action

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const Limit = 100
const Url = "https://fcm.googleapis.com/fcm/send"
const api_key = "AAAA5JPsPP8:APA91bGlSNDdWQltp3s7-MolrFdCbhNa6b9ihYxQR_TlVJRd3H8IEhU-2hkMHY1dnqCuAiMvL2JZJIIHCZDkuBOppt1gSO9TSkzq5bVkBbkjwL_lBDEW9SR4q_UBDego8NCR06S0ysoF"

var fbClient *FirebaseClient

type FirebaseClient struct {
	apiKey   string
	endpoint string
	client   *http.Client
	timeout  time.Duration
	debug    bool
}

type Result struct {
	MessageID      string `json:"message_id"`
	RegistrationID string `json:"registration_id"`
	Error          string `json:"error"`
}

type FirebaseResponse struct {
	MulticastID  int64    `json:"multicast_id"`
	Success      int      `json:"success"`
	Failure      int      `json:"failure"`
	CanonicalIDs int      `json:"canonical_ids"`
	Results      []Result `json:"results"`

	// Device Group HTTP Response
	FailedRegistrationIDs []string `json:"failed_registration_ids"`

	// Topic HTTP response
	MessageID int64 `json:"message_id"`
	Error     error `json:"error"`
}

func (msg *messages) Validate() error {

	if msg == nil {
		return errors.New("ERR_INVALID_MSG")
	}

	if msg.To == "" && msg.Condition == "" && len(msg.RegistrationIDs) == 0 {
		return errors.New("ERR_INVALID_TOKEN")
	}

	if len(msg.RegistrationIDs) > Limit {
		return errors.New("ERR_TOO_MANY_REG_ID")
	}

	return nil
}

func GetDefaultFirebaseClient() (*FirebaseClient, error) {

	var err error

	if fbClient == nil {
		fbClient, err = NewClientDefault(api_key)
		return fbClient, err
	}

	return fbClient, nil
}

func NewClientDefault(apiKey string) (*FirebaseClient, error) {
	if apiKey == "" {
		return nil, errors.New("ERR_INVALID_API_KEY")
	}

	c := FirebaseClient{
		apiKey:   apiKey,
		client:   &http.Client{},
		endpoint: Url,
		timeout:  30 * time.Second,
	}

	return &c, nil
}

func NewClient(apiKey, endPoint string, timeout time.Duration) (*FirebaseClient, error) {
	if apiKey == "" {
		return nil, errors.New("ERR_INVALID_API_KEY")
	}

	c := FirebaseClient{
		apiKey:   apiKey,
		client:   &http.Client{},
		endpoint: endPoint,
		timeout:  timeout,
	}

	return &c, nil
}

// Send is func client send msg to server
func (c *FirebaseClient) Send(msg ...FirebaseMessages) (*FirebaseResponse, error) {
	var mess messages
	for _, v := range msg {
		v.apply(&mess)
	}

	if err := mess.Validate(); err != nil {
		return nil, err
	}

	data, err := json.Marshal(mess)
	if err != nil {
		return nil, err
	}

	return c.push(data)
}

func (c *FirebaseClient) push(data []byte) (*FirebaseResponse, error) {
	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("key=%s", c.apiKey))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "close")

	resp, err := c.client.Do(req)
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

	var response FirebaseResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
