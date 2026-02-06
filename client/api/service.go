package api

import (
	"bytes"
	"courseWork/shared"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Service struct {
	BaseURL string
}

func NewService(baseURL string) *Service {
	return &Service{BaseURL: baseURL}
}

func (s *Service) CreateAthlete(a shared.Athlete) (string, error) {
	url := fmt.Sprintf("%s/athlete/create", s.BaseURL)
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

func (s *Service) FetchBestAthletes() ([]shared.Athlete, error) {
	url := fmt.Sprintf("%s/athlete/fetch/best", s.BaseURL)
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

func (s *Service) FetchAthletes() ([]shared.Athlete, error) {
	url := fmt.Sprintf("%s/athlete/fetch/all", s.BaseURL)
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

func (s *Service) UpdateAthlete(a shared.Athlete) error {
	url := fmt.Sprintf("%s/athlete/update", s.BaseURL)
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

func (s *Service) FetchAthletesSortedByRun100m() ([]shared.Athlete, error) {
	url := fmt.Sprintf("%s/athlete/fetch/sorted", s.BaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var athletes []shared.Athlete
	err = json.NewDecoder(resp.Body).Decode(&athletes)
	if err != nil {
		return nil, err
	}

	return athletes, nil
}
