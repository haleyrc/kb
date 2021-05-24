package app

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Link struct {
	Title  string
	Author string
	URL    string
	Tags   string
}

func (l *Link) Parse(args ...string) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.StringVar(&l.Title, "title", "", "The title of the link")
	fs.StringVar(&l.Author, "author", "", "The author of the link")
	fs.StringVar(&l.URL, "url", "", "The url of the link")
	fs.StringVar(&l.Tags, "tags", "", "The tags of the link")
	fs.Parse(args)
	return nil
}

func (l Link) HasTag(tag string) bool {
	tag = strings.ToLower(tag)
	allTags := strings.Split(strings.ToLower(l.Tags), ",")
	for _, t := range allTags {
		if t == tag {
			return true
		}
	}
	return false
}

type linkFilter struct {
	Tag string
}

func (lf *linkFilter) Parse(args ...string) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.StringVar(&lf.Tag, "tag", "", "The tag to search for")
	fs.Parse(args)
	return nil
}

func (lf linkFilter) Search(all []Link) ([]Link, error) {
	results := []Link{}
	for _, link := range all {
		if lf.Tag != "" && !link.HasTag(lf.Tag) {
			continue
		}
		results = append(results, link)
	}
	return results, nil
}

func handleLinks(args ...string) error {
	command := args[0]
	switch command {
	case "search":
		return searchLinks(args[1:]...)
	case "new":
		return newLink(args[1:]...)
	default:
		return fmt.Errorf("links: invalid command: %s", command)
	}
}

func newLink(args ...string) error {
	var link Link
	if err := link.Parse(args...); err != nil {
		return fmt.Errorf("new: %w", err)
	}
	// We pass a one element slice here so that we're appending a list instead
	// of a single item, which would result in incorrect YAML
	if err := appendData("links", []Link{link}); err != nil {
		return fmt.Errorf("new: %w", err)
	}
	return nil
}

func searchLinks(args ...string) error {
	var filter linkFilter
	if err := filter.Parse(args...); err != nil {
		return fmt.Errorf("search: %w", err)
	}

	var links []Link
	if err := load("links", &links); err != nil {
		return fmt.Errorf("search: %w", err)
	}

	results, err := filter.Search(links)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	for _, result := range results {
		fmt.Fprintf(os.Stderr, "%q by %s\n", result.Title, result.Author)
		fmt.Fprintf(os.Stderr, "    Tags: %s\n", result.Tags)
		fmt.Fprintf(os.Stderr, "    URL : %s\n", result.URL)
		fmt.Fprintln(os.Stderr)
	}

	return nil
}
