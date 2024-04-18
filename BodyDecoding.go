package main

import (
	"encoding/json"
	"log"
)

type Pokemon struct {
	Name string
}

type Response struct {
	Next string
	Previous string
	Results []Pokemon
}

func decode(body []byte) Response {
	var response Response

	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal("error parsing the JSON body")
	}

	return response
}