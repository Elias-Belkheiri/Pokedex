package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type clicommand struct {
	name string
	description string
	action func()
}

var response Response = Response{"https://pokeapi.co/api/v2/location?limit=20", nil, nil}

func helpCmd() {
	fmt.Print("Welcome to the Pokedex!\n" +
		"Usage:\n\n" +
		"help: Displays a help message\n" +
		"exit: Exit the Pokedex\n")
}

func exitCmd() {
	os.Exit(0)
}

func mapCmd() {
	resp, err := http.Get(response.Next)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
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
	resp, err := http.Get(*(response.Previous))

	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	decode(body, &response)

	for _, region := range(response.Results) {
		fmt.Println(region.Name)
	}
}

func main() {
	cmds := make(map[string]clicommand)
	scanner := bufio.NewScanner(os.Stdin)

	cmds["help"] = clicommand{"help", "Displays a help message", helpCmd}
	cmds["exit"] = clicommand{"exit", "exits the program", exitCmd}
	cmds["map"] = clicommand{"map", "displays the names of 20 location areas in the Pokemon world. Each subsequent call to map should display the next 20 locations, and so on.", mapCmd}
	cmds["mapb"] = clicommand{"mapb", "displays the previous 20 locations in the Pokemon world.", mapbCmd}

	for {
		fmt.Print("Pokedex > ")
	
		scanner.Scan()
		input := scanner.Text()
		
		if input == "help" || input == "exit" || input == "map" || input == "mapb" {
			cmds[input].action()
		} else {
			fmt.Println("Invalid cmd")
		}
	}
}