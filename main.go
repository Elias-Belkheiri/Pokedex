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

var response = Response{"https://pokeapi.co/api/v2/location-area?offset=0&limit=20", nil, nil}
var	cmds = make(map[string]clicommand)
var	cached = cache.Cache{C: make(map[string]cache.CacheEntry)}
var area string

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

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cmds["help"] = clicommand{"help", "Displays a help message", helpCmd}
	cmds["exit"] = clicommand{"exit", "exits the program", exitCmd}
	cmds["map"] = clicommand{"map", "displays the names of 20 location areas in the Pokemon world. Each subsequent call to map should display the next 20 locations, and so on.", mapCmd}
	cmds["mapb"] = clicommand{"mapb", "displays the previous 20 locations in the Pokemon world.", mapbCmd}
	cmds["explore"] = clicommand{"explore", "displays a list of all the Pokémon in a given area.", exploreCmd}
	for {
		fmt.Print("Pokedex > ")
	
		scanner.Scan()
		input := strings.Fields(scanner.Text())
		
		if input[0] == "help" || input[0] == "exit" || input[0] == "map" || input[0] == "mapb" || input[0] == "explore" {
			if input[0] == "explore" {
				if len(input) != 2 {
					fmt.Println("Invalid cmd")
					continue
				}
				area = input[1]
			}
			cmds[input[0]].action()
		} else {
			fmt.Println("Invalid cmd")
		}
	}
}