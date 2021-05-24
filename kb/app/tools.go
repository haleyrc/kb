package app

import (
	"flag"
	"fmt"
	"os"
)

type Tool struct {
	Name        string
	Description string
	Link        string
	Notes       string
}

func (t *Tool) Parse(args ...string) error {
	fs := flag.NewFlagSet("tools", flag.ContinueOnError)
	fs.StringVar(&t.Name, "name", "", "The name of the item")
	fs.StringVar(&t.Description, "desc", "", "A description of the tool")
	fs.StringVar(&t.Link, "link", "", "A link to the tool")
	fs.StringVar(&t.Notes, "notes", "", "Notes about the tool")
	fs.Parse(args)
	return nil
}

type toolFilter struct{}

func (sf *toolFilter) Parse(args ...string) error { return nil }

func (sf *toolFilter) Search(tools []Tool) ([]Tool, error) { return tools, nil }

func handleTools(args ...string) error {
	command := args[0]
	switch command {
	case "new":
		return newTool(args[1:]...)
	case "search":
		return searchTools(args[1:]...)
	default:
		return fmt.Errorf("tools: invalid command: %s", command)
	}
}

func newTool(args ...string) error {
	var tool Tool
	if err := tool.Parse(args...); err != nil {
		return fmt.Errorf("new: %w", err)
	}
	// We pass a one element slice here so that we're appending a list instead
	// of a single item, which would result in incorrect YAML
	if err := appendData("tools", []Tool{tool}); err != nil {
		return fmt.Errorf("new: %w", err)
	}
	return nil
}

func searchTools(args ...string) error {
	var filter toolFilter
	if err := filter.Parse(args...); err != nil {
		return fmt.Errorf("search: %w", err)
	}

	var tools []Tool
	if err := load("tools", &tools); err != nil {
		return fmt.Errorf("search: %w", err)
	}

	results, err := filter.Search(tools)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	for _, result := range results {
		fmt.Fprintf(os.Stderr, "%s\n", result.Name)
		fmt.Fprintf(os.Stderr, "    Description: %s\n", result.Description)
		if result.Notes != "" {
			fmt.Fprintf(os.Stderr, "    Notes      : %s\n", result.Notes)
		}
		fmt.Fprintf(os.Stderr, "    Link       : %s\n", result.Link)
		fmt.Fprintln(os.Stderr)
	}

	return nil
}
