package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("")
		fmt.Print("Pokedex > ")

		ok := scanner.Scan()
		if !ok {
			fmt.Println("Error in scan")
			continue
		}

		text := scanner.Text()
		text = strings.ToLower(text)
		if text == "" {
			fmt.Println("Please enter a command")
			continue
		}

		fields := strings.Fields(text)
		commandName := fields[0]
		args := []string{}
		if len(fields) > 1 {
			args = fields[1:]
		}

		cmd, ok := supportedCommands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := cmd.callback(args)
		if err != nil {
			fmt.Printf("Error in function: %s", err)
			continue
		}
	}
}
