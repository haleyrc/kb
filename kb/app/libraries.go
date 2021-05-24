package app

import (
	"flag"
	"fmt"
	"os"
)

type Library struct {
	Name        string
	Language    string
	URL         string
	Description string
}

func (l *Library) Parse(args ...string) error {
	fs := flag.NewFlagSet("library", flag.ContinueOnError)
	fs.StringVar(&l.Name, "name", "", "The name of the library")
	fs.StringVar(&l.Language, "lang", "", "The language the library is written in")
	fs.StringVar(&l.URL, "url", "", "The URL to the library")
	fs.StringVar(&l.Description, "desc", "", "A description of the library's function")
	fs.Parse(args)
	return nil
}

type libraryFilter struct{}

func (lf *libraryFilter) Parse(args ...string) error                    { return nil }
func (lf *libraryFilter) Search(libraries []Library) ([]Library, error) { return libraries, nil }

func handleLibraries(args ...string) error {
	command := args[0]
	switch command {
	case "search":
		return searchLibraries(args[1:]...)
	case "new":
		return newLibrary(args[1:]...)
	default:
		return fmt.Errorf("libraries: invalid command: %s", command)
	}
}

func newLibrary(args ...string) error {
	var library Library
	if err := library.Parse(args...); err != nil {
		return fmt.Errorf("new: %w", err)
	}
	// We pass a one element slice here so that we're appending a list instead
	// of a single item, which would result in incorrect YAML
	if err := appendData("libraries", []Library{library}); err != nil {
		return fmt.Errorf("new: %w", err)
	}
	return nil
}

func searchLibraries(args ...string) error {
	var filter libraryFilter
	if err := filter.Parse(args...); err != nil {
		return fmt.Errorf("search: %w", err)
	}

	var libraries []Library
	if err := load("libraries", &libraries); err != nil {
		return fmt.Errorf("search: %w", err)
	}

	results, err := filter.Search(libraries)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	for _, result := range results {
		fmt.Fprintf(os.Stderr, "%s - %s\n", result.Name, result.Description)
		fmt.Fprintf(os.Stderr, "    Language: %s\n", result.Language)
		fmt.Fprintf(os.Stderr, "    URL     : %s\n", result.URL)
	}

	return nil
}
