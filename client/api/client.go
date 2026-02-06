package api

import (
	"bytes"
	"courseWork/shared"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) CreateAthlete(a shared.Athlete) (string, error) {
	url := fmt.Sprintf("%s/athlete/create", c.BaseURL)
	jsonData, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return resp.Status, nil
}

func (c *Client) FetchBestAthletes() ([]shared.Athlete, error) {
	url := fmt.Sprintf("%s/athlete/fetch/best", c.BaseURL)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, err
	}

	var athletes []shared.Athlete
	err = json.NewDecoder(response.Body).Decode(&athletes)
	if err != nil {
		return nil, err
	}

	return athletes, nil
}

func (c *Client) FetchAthletes() ([]shared.Athlete, error) {
	url := fmt.Sprintf("%s/athlete/fetch/all", c.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	var athletes []shared.Athlete
	err = json.NewDecoder(response.Body).Decode(&athletes)
	if err != nil {
		return nil, err
	}

	return athletes, nil
}

func (c *Client) UpdateAthlete(a shared.Athlete) error {
	url := fmt.Sprintf("%s/athlete/update", c.BaseURL)
	marshalAthlete, err := json.Marshal(a)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(marshalAthlete)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
