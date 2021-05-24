package app

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Run(args ...string) error {
	command := args[0]
	switch command {
	case "libraries":
		return handleLibraries(args[1:]...)
	case "links":
		return handleLinks(args[1:]...)
	case "reading":
		return handleReading(args[1:]...)
	case "shopping":
		return handleShopping(args[1:]...)
	case "tools":
		return handleTools(args[1:]...)
	default:
		return fmt.Errorf("run: invalid command: %s", command)
	}
}

func appendData(key string, data interface{}) error {
	path := fmt.Sprintf("%s.yaml", key)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return fmt.Errorf("append: %w", err)
	}
	defer f.Close()

	if err := yaml.NewEncoder(f).Encode(data); err != nil {
		return fmt.Errorf("append: %w", err)
	}

	return nil
}

func load(key string, dest interface{}) error {
	path := fmt.Sprintf("%s.yaml", key)
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(dest); err != nil {
		return fmt.Errorf("load: %w", err)
	}

	return nil
}
