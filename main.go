package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"example.com/cache"
	"strings"
)

type clicommand struct {
	name 		string
	description string

	action 		func()
}

var (
	response = Response{"https://pokeapi.co/api/v2/location-area?offset=0&limit=20", nil, nil}
	cmds = make(map[string]clicommand)
	cached = cache.Cache{C: make(map[string]cache.CacheEntry)}
	area string
	pokemon string
)

func helpCmd() {
	fmt.Print("Welcome to the Pokedex!\n" +
		"Usage:\n\n")
	for key, value := range(cmds) {
		fmt.Print(key + ": " + value.description + "\n")
	}
}

func exitCmd() {
	os.Exit(0)
}

func mapCmd() {
	body, exists := cached.Get(response.Next)

	if (!exists) {
		resp, err := http.Get(response.Next)
		if err != nil {
			log.Fatalln(err)
		}
	
		body, err = io.ReadAll(resp.Body)
	
		if err != nil {
			log.Fatalln(err)
		}
		cached.Add(response.Next, body)
	} else {
		fmt.Println("---- Getting cache ----")
	}

	decode(body, &response)

	for _, region := range(response.Results) {
		fmt.Println(region.Name)
	}
}

func mapbCmd() {
	if response.Previous == nil {
		fmt.Println("No previous regions to show")
		return
	}

	body, exists := cached.Get(*(response.Previous))
	if (!exists) {
		resp, err := http.Get(*(response.Previous))
	
		if err != nil {
			log.Fatalln(err)
		}
	
		body, err = io.ReadAll(resp.Body)
	
		if err != nil {
			log.Fatalln(err)
		}
		cached.Add(*response.Previous, body)
	} else {
		fmt.Println("---- Getting cache ----")
	}
	decode(body, &response)

	for _, region := range(response.Results) {
		fmt.Println(region.Name)
	}
}

func exploreCmd() {
	locationURI := "https://pokeapi.co/api/v2/location-area/" + area
	body, exists := cached.Get(locationURI)
	var locationArea LocationArea

	if (!exists) {
		resp, err := http.Get(locationURI)
		if err != nil {
			log.Fatalln(err)
		}
	
		body, err = io.ReadAll(resp.Body)
	
		if err != nil {
			log.Fatalln(err)
		}
		cached.Add(locationURI, body)
	}

	err := decode(body, &locationArea)
	if err != nil {
		fmt.Println("Invalid area")
		return
	}
	for _, pokemon := range(locationArea.Pokemon_Encounters) {
		fmt.Println(pokemon.Pokemon.Name)
	}
}

func inspectCmd() {
	pokemonURI := "https://pokeapi.co/api/v2/pokemon/" + pokemon
	body, exists := cached.Get(pokemonURI)
	var poke Poke

	if (!exists) {
		resp, err := http.Get(pokemonURI)
		if err != nil {
			log.Fatalln(err)
		}
	
		body, err = io.ReadAll(resp.Body)
	
		if err != nil {
			log.Fatalln(err)
		}
		cached.Add(pokemonURI, body)
	}

	err := decode(body, &poke)
	if err != nil {
		fmt.Println("Invalid pokemon name")
		return
	}
	fmt.Println("Name: " + poke.Name)
	fmt.Printf("Height: %d\n" ,(poke.Height))
	fmt.Printf("Weight: %d\n" ,(poke.Weight))
	// fmt.Println("Weight: " + string(poke.Weight))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cmds["help"] = clicommand{"help", "Displays a help message", helpCmd}
	cmds["exit"] = clicommand{"exit", "exits the program", exitCmd}
	cmds["map"] = clicommand{"map", "displays the names of 20 location areas in the Pokemon world. Each subsequent call to map should display the next 20 locations, and so on.", mapCmd}
	cmds["mapb"] = clicommand{"mapb", "displays the previous 20 locations in the Pokemon world.", mapbCmd}
	cmds["explore"] = clicommand{"explore", "displays a list of all the Pokémon in a given area.", exploreCmd}
	cmds["inspect"] = clicommand{"inspect", "displays the details of a specific Pokémon.", inspectCmd}
	for {
		fmt.Print("Pokedex > ")
	
		scanner.Scan()
		input := strings.Fields(scanner.Text())
		
		if input[0] == "help" || input[0] == "exit" || input[0] == "map" || input[0] == "mapb" || input[0] == "explore" || input[0] == "inspect"{
			if input[0] == "explore" || input[0] == "inspect" {
				if len(input) != 2 {
					fmt.Println("Invalid arguments")
					continue
				}
				if input[0] == "explore" {
					area = input[1]
				} else {
					pokemon = input[1]
				}
			}
			cmds[input[0]].action()
		} else {
			fmt.Println("Invalid cmd")
		}
	}
}