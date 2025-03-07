package main

import (
	"fmt"
	"log"

	"github.com/galihrivanto/go-questionary/prompt"
)

func main() {
	// Text input example
	name := prompt.NewText("What is your name?")
	nameResult, err := name.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hello, %s!\n\n", nameResult)

	// Confirmation example
	confirm := prompt.NewConfirm("Do you want to continue?")
	confirm.DefaultValue = true
	confirmResult, err := confirm.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("You chose: %v\n\n", confirmResult)

	// List selection example
	options := []string{
		"Option 1",
		"Option 2",
		"Option 3",
		"Option 4",
		"Option 5",
	}
	list := prompt.NewList("Choose an option:", options)
	listResult, err := list.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("You selected: %s\n\n", listResult)

	// Text input with validation
	password := prompt.NewPassword("Enter a password (min 8 characters):")
	password.Validator = func(input string) error {
		if len(input) < 8 {
			return fmt.Errorf("password must be at least 8 characters long")
		}
		return nil
	}
	passwordResult, err := password.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Password accepted: %s\n", passwordResult)
}
