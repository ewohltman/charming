package client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ewohltman/charming/pkg/model"
)

type item struct {
	name        string
	ingredients string
}

func (i item) String() string {
	return fmt.Sprintf("Name: %s\nIngredients: %s",
		i.name,
		i.ingredients,
	)
}

type Fake struct {
	list  []model.Item
	items map[string]item
}

func NewFake() *Fake {
	list := []model.Item{
		"Ramen",
		"Tomato Soup",
		"Hamburgers",
		"Cheeseburgers",
		"Currywurst",
		"Okonomiyaki",
		"Pasta",
		"Fillet Mignon",
		"Caviar",
		"Just Wine",
	}

	items := make(map[string]item)

	for _, listItem := range list {
		items[string(listItem)] = item{
			name:        string(listItem),
			ingredients: "",
		}
	}

	return &Fake{
		list:  list,
		items: items,
	}
}

func (f *Fake) List() ([]model.Item, error) {
	return f.list, nil
}

func (f *Fake) Get(name string) (model.Item, error) {
	name = strings.TrimSpace(name)

	i, found := f.items[name]
	if !found {
		return "", errors.New("item not found")
	}

	return model.Item(i.String()), nil
}
