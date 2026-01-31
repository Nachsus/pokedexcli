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
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			fmt.Println("Error in scan")
			break
		}
		text := scanner.Text()
		text = strings.ToLower(text)
		first := strings.Fields(text)[0]
		fmt.Println("Your command was: " + first)
	}
}
