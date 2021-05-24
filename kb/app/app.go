package app

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func Run(args ...string) error {
	command := args[0]
	switch command {
	case "links":
		return handleLinks(args[1:]...)
	case "shopping":
		return handleShopping(args[1:]...)
	default:
		return fmt.Errorf("run: invalid command: %s", command)
	}
}

func handleShopping(args ...string) error {
	command := args[0]
	switch command {
	case "search":
		return searchShoppingList(args[1:]...)
	default:
		return fmt.Errorf("shopping: invalid command: %s", command)
	}
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
	if err := appendData("links", []Link{link}); err != nil {
		return fmt.Errorf("new: %w", err)
	}
	return nil
}

func searchShoppingList(args ...string) error {
	var filter shoppingFilter
	if err := filter.Parse(args...); err != nil {
		return fmt.Errorf("search: %w", err)
	}

	var items []ShoppingItem
	if err := load("shopping", &items); err != nil {
		return fmt.Errorf("search: %w", err)
	}

	results, err := filter.Search(items)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	for _, result := range results {
		fmt.Fprintf(os.Stderr, "%s\n", result.Name)
		fmt.Fprintf(os.Stderr, "%s\n", result.Category)
		fmt.Fprintf(os.Stderr, "$%.2f\n", result.Price)
		fmt.Fprintf(os.Stderr, "%s\n", result.Link)
	}

	return nil
}

type shoppingFilter struct{}

func (sf *shoppingFilter) Parse(args ...string) error {
	return nil
}

func (sf *shoppingFilter) Search(items []ShoppingItem) ([]ShoppingItem, error) {
	return items, nil
}

type ShoppingItem struct {
	Name     string
	Category string
	Link     string
	Price    float64
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
		fmt.Fprintf(os.Stderr, "%s\n", result.Title)
		fmt.Fprintf(os.Stderr, "%s\n", result.Author)
		fmt.Fprintf(os.Stderr, "%s\n", result.Tags)
		fmt.Fprintf(os.Stderr, "%s\n", result.URL)
		fmt.Fprintln(os.Stderr)
	}

	return nil
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
