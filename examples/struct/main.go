package main

import (
	"fmt"
	"log"

	"github.com/galihrivanto/go-questionary/prompt"
)

type Option struct {
	Name     string `prompt:"text"`
	Password string `prompt:"password"`
	Continue bool   `prompt:"confirm"`
	Choice   string `prompt:"list[option 1,option2,option 3]"`
}

func main() {
	opt := &Option{}
	err := prompt.PromptFromStruct(opt)
	if err != nil {
		log.Fatal(err)
	}

	// Now opt will be filled with user responses
	fmt.Printf("Name: %s\n", opt.Name)
	fmt.Printf("Password: %s\n", opt.Password)
	fmt.Printf("Continue: %v\n", opt.Continue)
	fmt.Printf("Choice: %s\n", opt.Choice)
}
