package app

import (
	"flag"
	"fmt"
	"os"
)

type ShoppingListItem struct {
	Name     string
	Category string
	Link     string
	Price    float64
}

func (i *ShoppingListItem) Parse(args ...string) error {
	fs := flag.NewFlagSet("shopping", flag.ContinueOnError)
	fs.StringVar(&i.Name, "name", "", "The name of the item")
	fs.StringVar(&i.Category, "cat", "", "The category for the item")
	fs.StringVar(&i.Link, "link", "", "A link to the item")
	fs.Float64Var(&i.Price, "price", 0, "The price of the item")
	fs.Parse(args)
	return nil
}

type shoppingFilter struct{}

func (sf *shoppingFilter) Parse(args ...string) error {
	return nil
}

func (sf *shoppingFilter) Search(items []ShoppingListItem) ([]ShoppingListItem, error) {
	return items, nil
}

func handleShopping(args ...string) error {
	command := args[0]
	switch command {
	case "new":
		return newShoppingListItem(args[1:]...)
	case "search":
		return searchShoppingList(args[1:]...)
	default:
		return fmt.Errorf("shopping: invalid command: %s", command)
	}
}

func newShoppingListItem(args ...string) error {
	var item ShoppingListItem
	if err := item.Parse(args...); err != nil {
		return fmt.Errorf("new: %w", err)
	}
	// We pass a one element slice here so that we're appending a list instead
	// of a single item, which would result in incorrect YAML
	if err := appendData("shopping", []ShoppingListItem{item}); err != nil {
		return fmt.Errorf("new: %w", err)
	}
	return nil
}

func searchShoppingList(args ...string) error {
	var filter shoppingFilter
	if err := filter.Parse(args...); err != nil {
		return fmt.Errorf("search: %w", err)
	}

	var items []ShoppingListItem
	if err := load("shopping", &items); err != nil {
		return fmt.Errorf("search: %w", err)
	}

	results, err := filter.Search(items)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	for _, result := range results {
		fmt.Fprintf(os.Stderr, "%s\n", result.Name)
		fmt.Fprintf(os.Stderr, "    Category: %s\n", result.Category)
		fmt.Fprintf(os.Stderr, "    Price   : $%.2f\n", result.Price)
		fmt.Fprintf(os.Stderr, "    Link    : %s\n", result.Link)
	}

	return nil
}
