package main

import (
	"bytes"
	"courseWork/shared"
	"encoding/json"
	"errors"
	"net/http"
)

func createAthlete(a shared.Athlete) (string, error) {
	const url = "http://localhost:1323/athlete/create"
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

func fetchBestAthletes() ([]shared.Athlete, error) {
	const url = "http://localhost:1323/athlete/fetch/best"
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

func fetchAthletes() ([]shared.Athlete, error) {
	const url = "http://localhost:1323/athlete/fetch/all"
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

func updateAthlete(a shared.Athlete) error {
	const url = "http://localhost:1323/athlete/update"
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
