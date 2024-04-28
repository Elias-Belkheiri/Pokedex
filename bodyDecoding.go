package main

import (
	"encoding/json"
)

type Regions struct {
	Name string
}

type Response struct {
	Next string
	Previous *string
	Results []Regions
}

type Poke struct {
	Name	string
	Url		string
	Height	int
	Weight	int
}

type Pokemon struct {
	Pokemon Poke
}
type LocationArea struct {
	Name				string
	Pokemon_Encounters	[]Pokemon
}

func decode[k Response | LocationArea | Poke](body []byte, response *k) error{
	err := json.Unmarshal(body, response)
	return err
}