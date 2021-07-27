package domain

import (
	"errors"
	"fmt"
)

func MakeSet() *customSet {
	return &customSet{
		container: make(map[string]float64),
	}
}

type customSet struct {
	container map[string] float64
	keys []string
}

func (c *customSet) GetAll() map[string]float64 {
	return c.container
}

func (c *customSet) GetKeys() []string {
	return c.keys
}

func (c *customSet) Exists(key string) bool {
	_, exists := c.container[key]
	return exists
}

func (c *customSet) Get(key string) (float64, error) {
	if !c.Exists(key) {
		return  0, errors.New("key not exist")
	}
	return c.container[key], nil
}

func (c *customSet) Add(key string, confidence float64) {
	if c.Exists(key) {
		return
	}
	c.keys = append(c.keys, key)

	c.container[key] = confidence
}

func (c *customSet) Remove(key string) error {
	_, exists := c.container[key]
	if !exists {
		return fmt.Errorf("Remove Error: Item doesn't exist in set")
	}
	delete(c.container, key)
	return nil
}

func (c *customSet) Size() int {
	return len(c.container)
}

