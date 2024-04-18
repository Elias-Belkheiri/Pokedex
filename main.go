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

func helpCmd() {
	fmt.Print("Welcome to the Pokedex!\n" +
		"Usage:\n\n" +
		"help: Displays a help message\n" +
		"exit: Exit the Pokedex\n")
}

func exitCmd() {
	os.Exit(0)
}

func main() {
	cmds := make(map[string]clicommand)
	scanner := bufio.NewScanner(os.Stdin)

	cmds["help"] = clicommand{"help", "Displays a help message", helpCmd}
	cmds["exit"] = clicommand{"exit", "exits the program", exitCmd}

	for {
		fmt.Print("Pokedex > ")
	
		scanner.Scan()
		input := scanner.Text()
	
		if input == "help" || input == "exit"{
			// cmds[input].action()
			resp, err := http.Get("https://pokeapi.co/api/v2/location/?limit=20")

			if err != nil {
				log.Fatalln(err)
			}

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				log.Fatalln(err)
			}

			respJSON := decode(body)
			// fmt.Println(body)
			fmt.Println(respJSON.Next)
		} else {
			fmt.Println("Invalid cmd")
		}
	}
}