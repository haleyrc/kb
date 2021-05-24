package app

import (
	"flag"
	"fmt"
	"os"
)

type ReadingListItem struct {
	Name   string
	Author string
	Link   string
}

func (i *ReadingListItem) Parse(args ...string) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.StringVar(&i.Name, "name", "", "The name of the reading list item")
	fs.StringVar(&i.Author, "author", "", "The author's name")
	fs.StringVar(&i.Link, "link", "", "The link to the item")
	fs.Parse(args)
	return nil
}

type readingListFilter struct{}

func (rlf *readingListFilter) Parse(args ...string) error { return nil }
func (rlf *readingListFilter) Search(items []ReadingListItem) ([]ReadingListItem, error) {
	return items, nil
}

func handleReading(args ...string) error {
	command := args[0]
	switch command {
	case "search":
		return searchReadingList(args[1:]...)
	case "new":
		return newReadingListItem(args[1:]...)
	default:
		return fmt.Errorf("reading: invalid command: %s", command)
	}
}

func newReadingListItem(args ...string) error {
	var item ReadingListItem
	if err := item.Parse(args...); err != nil {
		return fmt.Errorf("new: %w", err)
	}
	// We pass a one element slice here so that we're appending a list instead
	// of a single item, which would result in incorrect YAML
	if err := appendData("reading", []ReadingListItem{item}); err != nil {
		return fmt.Errorf("new: %w", err)
	}
	return nil
}

func searchReadingList(args ...string) error {
	var filter readingListFilter
	if err := filter.Parse(args...); err != nil {
		return fmt.Errorf("search: %w", err)
	}

	var items []ReadingListItem
	if err := load("reading", &items); err != nil {
		return fmt.Errorf("search: %w", err)
	}

	results, err := filter.Search(items)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	for _, result := range results {
		fmt.Fprintf(os.Stderr, "%q", result.Name)
		if result.Author != "" {
			fmt.Fprintf(os.Stderr, " by %s", result.Author)
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "    Link: %s\n", result.Link)
		fmt.Fprintln(os.Stderr)
	}

	return nil
}
