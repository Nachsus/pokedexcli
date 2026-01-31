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

		first := strings.Fields(text)[0]
		cmd, ok := supportedCommands[first]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := cmd.callback()
		if err != nil {
			fmt.Printf("Error in function: %s", err)
			continue
		}
	}
}
