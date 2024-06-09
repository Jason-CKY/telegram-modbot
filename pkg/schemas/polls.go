package schemas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jason-CKY/telegram-modbot/pkg/utils"
)

type Poll struct {
	PollId    string `json:"poll_id"`
	MessageId int    `json:"message_id"`
	ChatId    int64  `json:"chat_id"`
}

func (poll Poll) Create() error {
	endpoint := fmt.Sprintf("%v/items/modbot_polls", utils.DirectusHost)
	reqBody, _ := json.Marshal(poll)
	req, httpErr := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("error inserting polls to directus: %v", string(body))
	}

	return nil
}

func (poll Poll) Update() error {
	endpoint := fmt.Sprintf("%v/items/modbot_polls/%v", utils.DirectusHost, poll.PollId)
	reqBody, _ := json.Marshal(poll)
	req, httpErr := http.NewRequest(http.MethodPatch, endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("error updating chat settings to directus: %v", string(body))
	}

	return nil
}

func (poll Poll) Delete() error {
	endpoint := fmt.Sprintf("%v/items/modbot_polls/%v", utils.DirectusHost, poll.PollId)
	req, httpErr := http.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 204 {
		return fmt.Errorf("error deleting chat settings in directus: %v", string(body))
	}
	return nil
}

type PollWithChatSettings struct {
	PollId       string       `json:"poll_id"`
	MessageId    int          `json:"message_id"`
	ChatSettings ChatSettings `json:"chat_id"`
}

func GetPollWithChatSettingsByPollId(pollId string) (*PollWithChatSettings, error) {
	endpoint := fmt.Sprintf("%v/items/modbot_polls", utils.DirectusHost)
	reqBody := []byte(fmt.Sprintf(`{
		"query": {
			"filter": {
				"poll_id": {
					"_eq": "%v"
				}
			},
			"fields": [
				"poll_id",
				"message_id",
				"chat_id.chat_id",
				"chat_id.threshold",
				"chat_id.expiry_time"
			]
		}
	}`, pollId))
	req, httpErr := http.NewRequest("SEARCH", endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return nil, httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return nil, httpErr
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error getting chat settings in directus: %v", string(body))
	}
	var pollResponse map[string][]PollWithChatSettings
	jsonErr := json.Unmarshal(body, &pollResponse)
	// error handling for json unmarshaling
	if jsonErr != nil {
		return nil, jsonErr
	}

	if len(pollResponse["data"]) == 0 {
		return nil, nil
	}

	return &pollResponse["data"][0], nil
}
