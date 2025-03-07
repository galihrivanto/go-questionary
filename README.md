<div align="center">
  <h1>go-questionary</h1>
  
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/galihrivanto/go-questionary)
  
A Go library for building command line user prompts

Inspired by [questionary](https://github.com/tmbo/questionary) Python package.
Built using [bubbletea](https://github.com/charmbracelet/bubbletea)
</div>


## Features

- [x] Struct Tags
- [x] Text Prompt
- [x] Password Prompt
- [x] List / Select Prompt
- [x] Boolean Prompt

## Installation

```bash
go get github.com/galihrivanto/go-questionary
```

## Usage

### Basic Prompts

```go
package main

import (
    "fmt"
    "log"
    "github.com/galihrivanto/go-questionary/prompt"
)

func main() {
    // Text prompt
    name, err := prompt.NewTextPrompt("What's your name?").Run()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Hello, %s!\n", name)

    // Confirmation prompt
    continue, err := prompt.NewConfirmPrompt("Do you want to continue?").Run()
    if err != nil {
        log.Fatal(err)
    }

    // List prompt
    options := []string{"Option 1", "Option 2", "Option 3"}
    choice, err := prompt.NewListPrompt("Choose an option:", options).Run()
    if err != nil {
        log.Fatal(err)
    }

    // Password prompt
    password, err := prompt.NewPasswordPrompt("Enter your password:").Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

### Using Struct Tags

You can use struct tags to create prompts automatically from a struct:

```go
package main

import (
    "fmt"
    "log"
    "github.com/galihrivanto/go-questionary/prompt"
)

type UserConfig struct {
    Username string `prompt:"text"`
    Password string `prompt:"password"`
    Role     string `prompt:"list[admin,user,guest]"`
    Active   bool   `prompt:"confirm"`
}

func main() {
    config := &UserConfig{}
    err := prompt.PromptFromStruct(config)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("User Configuration:\n")
    fmt.Printf("Username: %s\n", config.Username)
    fmt.Printf("Role: %s\n", config.Role)
    fmt.Printf("Active: %v\n", config.Active)
}
```

### Supported Struct Tags

- `prompt:"text"` - Text input prompt
- `prompt:"password"` - Password input prompt (masked input)
- `prompt:"confirm"` - Yes/No confirmation prompt
- `prompt:"list[option1,option2,...]"` - List selection prompt with options

### Customizing Prompts

You can customize prompts by setting additional options:

```go
package main

import (
    "github.com/galihrivanto/go-questionary/prompt"
)

func main() {
    // Text prompt with validation
    namePrompt := prompt.NewTextPrompt("What's your name?")
    namePrompt.Validator = func(input string) error {
        if len(input) < 2 {
            return fmt.Errorf("name must be at least 2 characters")
        }
        return nil
    }

    // Prompt with default value
    emailPrompt := prompt.NewTextPrompt("Enter your email:")
    emailPrompt.Default = "user@example.com"

    // Custom styling
    style := prompt.DefaultStyle()
    style.QuestionStyle = style.QuestionStyle.Foreground(lipgloss.Color("86"))
    
    listPrompt := prompt.NewListPrompt("Select option:", []string{"A", "B", "C"})
    listPrompt.Style = style
}
```






