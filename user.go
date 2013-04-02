package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	usersShowUrl    = "http://api.fanfou.com/users/show.json"
	userTimelineUrl = "http://api.fanfou.com/statuses/user_timeline.json"
)

type User struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Protected bool   `json:"protected"`
	Following bool   `json:"following"`
}

type Status struct {
	CreatedAt string `json:"created_at"`
	Id        string `json:"id"`
	RawId     int64  `json:"rawid"`
	Text      string `json:"text"`
	Source    string `json:"source"`
}

func GetCurrentUserId(client *FanfouClient) (string, error) {
	resp, err := client.Get(usersShowUrl, map[string]string{
		"mode": "lite"}, client.Token)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var currentUser User
	json.Unmarshal(bytes, &currentUser)

	return currentUser.Id, nil
}

func GetAndStoreAllStatuses(client *FanfouClient, db *sql.DB) error {
	var maxId string
	for {
		resp, err := client.Get(userTimelineUrl,
			map[string]string{"mode": "lite", "count": "60", "max_id": maxId},
			client.Token)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// Unmarshal all json data
		statuses := []Status{}
		json.Unmarshal(bytes, &statuses)
		// If we already read all statuses, then quit
		if len(statuses) == 0 {
			break
		}

		// Insert all statuses into table
		err = InsertStatuses(db, &statuses)
		if err != nil {
			return err
		}

		// Get last status's id
		maxId = statuses[len(statuses)-1].Id
		log.Printf("id: %s OK\n", maxId)
	}

	return nil
}
