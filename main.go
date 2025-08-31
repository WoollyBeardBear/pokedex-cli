package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}

map[string]cliCommand{
	"exit": {
		name:			"exit",
		description:	"Exit the Pokedex",
		callback:		commandExit,
	},
	"help": {
		name:			"help",
		description:	"Find help",
		callback:		commandHelp,
	},
}


func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		fmt.Printlnf("Your command was: %s \n", input[0])
	}
}


func cleanInput(text string) []string {
	lowerFields := strings.Fields(strings.ToLower(text))
	//for _, l := range lowerFields {
	//	strings.TrimSpace(l)
	//}
	return lowerFields
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("Did not exit program")
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\n
	help: Displays a help message\nexit: Exit the Pokedex"\n)
}	
