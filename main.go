package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Locations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	Next     *string
	Previous *string
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Find help",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "page through locations",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapBack",
		description: "page back through locations",
		callback:    commandMapBack,
	},
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var mapConfig Config
	url := "https://pokeapi.co/api/v2/location-area"
	mapConfig.Next = &url
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		cmd := input[0]

		if val, ok := commands[cmd]; !ok {
			fmt.Print("Unknown command\n")
		} else {
			fmt.Printf("Your command was: %s \n", cmd)
			val.callback(&mapConfig)
		}

	}
}

func cleanInput(text string) []string {
	lowerFields := strings.Fields(strings.ToLower(text))
	return lowerFields
}

func commandMap(c *Config) error {
	if c.Next == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	res, err := http.Get(*c.Next)
	if err != nil {
		fmt.Println("error getting response")
		return err
	}
	defer res.Body.Close()
	fmt.Println("right after response")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading body")
		return err
	}

	var locations Locations
	if err := json.Unmarshal(body, &locations); err != nil {
		fmt.Printf("error unmarshalling data: %v\n", err)
		return err
	}
	c.Next = locations.Next
	c.Previous = locations.Previous
	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}
	return nil
}

func commandMapBack(c *Config) error {
	if c.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	res, err := http.Get(*c.Previous)
	if err != nil {
		fmt.Println("error getting response")
		return err
	}
	defer res.Body.Close()
	fmt.Println("right after response")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading body")
		return err
	}

	var locations Locations
	if err := json.Unmarshal(body, &locations); err != nil {
		fmt.Printf("error unmarshalling data: %v\n", err)
		return err
	}
	c.Next = locations.Next
	c.Previous = locations.Previous
	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}
	return nil
}

func commandExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("Did not exit program")
}

func commandHelp(c *Config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}
