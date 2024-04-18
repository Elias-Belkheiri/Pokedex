package main

import (
	"encoding/json"
	"log"
)

type Regions struct {
	Name string
}

type Response struct {
	Next string
	Previous *string
	Results []Regions
}

func decode(body []byte, response *Response) {
	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal("error parsing the JSON body")
	}
}