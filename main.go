package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/woollybeardbear/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
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

type LocationArea struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Config struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
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
	"explore": {
		name:        "explore",
		description: "find pokemon in each area",
		callback:    commandExplore,
	},
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var mapConfig Config
	url := "https://pokeapi.co/api/v2/location-area/"
	mapConfig.Next = &url
	const baseTime = 5 * time.Millisecond
	mapConfig.Cache = pokecache.NewCache(baseTime)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		cmd := input[0]

		if val, ok := commands[cmd]; !ok {
			fmt.Print("Unknown command\n")
		} else {
			//fmt.Printf("Your command was: %s %s\n", cmd, input[1:])
			val.callback(&mapConfig, input[1:])
		}

	}
}

func cleanInput(text string) []string {
	lowerFields := strings.Fields(strings.ToLower(text))
	return lowerFields
}

func commandMap(c *Config, args []string) error {
	if c.Next == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	body, ok := c.Cache.Get(*c.Next)
	if !ok {
		res, err := http.Get(*c.Next)
		if err != nil {
			fmt.Println("error getting response")
			return err
		}
		defer res.Body.Close()
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
		c.Cache.Add(*c.Next, body)
		c.Next = locations.Next
		c.Previous = locations.Previous

		for _, location := range locations.Results {
			fmt.Printf("%s\n", location.Name)
		}
		return nil
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

func commandMapBack(c *Config, args []string) error {
	if c.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	body, ok := c.Cache.Get(*c.Previous)
	if !ok {
		res, err := http.Get(*c.Previous)
		if err != nil {
			fmt.Println("error getting response")
			return err
		}
		defer res.Body.Close()
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
		c.Cache.Add(*c.Previous, body)
		c.Next = locations.Next
		c.Previous = locations.Previous
		for _, location := range locations.Results {
			fmt.Printf("%s\n", location.Name)
		}
		return nil
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

func commandExit(c *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("Did not exit program")
}

func commandHelp(c *Config, args []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex\nmap: Displays a list of locations in the pokedex\nmapb: Displays the previous list of locations")
	return nil
}

func commandExplore(c *Config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Please include a location area after the explore command")
		return nil
	}

	fmt.Printf("Exploring %s \n", args[0])
	url := "https://pokeapi.co/api/v2/location-area/" + args[0]
	body, ok := c.Cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("Error getting response")
			return err
		}

		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("error reading body")
			return err
		}
		c.Cache.Add(url, body)
	}
	var locationArea LocationArea
	if err := json.Unmarshal(body, &locationArea); err != nil {
		fmt.Printf("error unmarshalling data: %v\n", err)
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf("- %s \n", encounter.Pokemon.Name)
	}

	return nil
}
