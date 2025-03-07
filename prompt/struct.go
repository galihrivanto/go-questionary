package prompt

import (
	"fmt"
	"reflect"
	"strings"
)

// PromptTag represents the parsed prompt struct tag
type PromptTag struct {
	Type    string   // text, confirm, list
	Options []string // options for list type
}

// ParsePromptTag parses the "prompt" struct tag
func ParsePromptTag(tag string) (*PromptTag, error) {
	if tag == "" {
		return nil, nil
	}

	content := strings.TrimSpace(tag)

	// Handle list type with options
	if strings.HasPrefix(content, "list[") && strings.HasSuffix(content, "]") {
		options := strings.Split(content[5:len(content)-1], ",")
		// Trim spaces from options
		for i := range options {
			options[i] = strings.TrimSpace(options[i])
		}
		return &PromptTag{Type: "list", Options: options}, nil
	}

	// Handle simple types
	return &PromptTag{Type: content}, nil
}

// PromptFromStruct generates prompts based on struct tags
func PromptFromStruct(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("input must be a pointer to struct")
	}

	val = val.Elem()
	typ := val.Type()

	for i := range val.NumField() {
		field := val.Field(i)
		structField := typ.Field(i)

		tag, err := ParsePromptTag(structField.Tag.Get("prompt"))
		if err != nil {
			return err
		}
		if tag == nil {
			continue
		}

		// Create appropriate prompt based on tag type
		var result interface{}
		switch tag.Type {
		case "text":
			textPrompt := NewText(structField.Name)
			result, err = textPrompt.Run()

		case "password":
			passwordPrompt := NewPassword(structField.Name)
			result, err = passwordPrompt.Run()

		case "confirm":
			confirmPrompt := NewConfirm(structField.Name)
			result, err = confirmPrompt.Run()

		case "list":
			listPrompt := NewList(structField.Name, tag.Options)
			result, err = listPrompt.Run()

		default:
			return fmt.Errorf("unknown prompt type: %s", tag.Type)
		}

		if err != nil {
			return err
		}

		// Set the result back to the struct field
		field.Set(reflect.ValueOf(result))
	}

	return nil
}
