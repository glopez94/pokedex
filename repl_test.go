package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/glopez94/pokedex/internal/pokecache"
)

// TestCleanInput prueba la función cleanInput
func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "   Go is Awesome   ",
			expected: []string{"go", "is", "awesome"},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "SingleWord",
			expected: []string{"singleword"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("cleanInput(%q) == %v, expected %v", c.input, actual, c.expected)
		}
	}
}

// TestSetURL prueba la función setURL
func TestSetURL(t *testing.T) {
	cases := []struct {
		input    *string
		expected string
	}{
		{
			input:    nil,
			expected: "",
		},
		{
			input:    stringPtr("https://example.com"),
			expected: "https://example.com",
		},
	}

	for _, c := range cases {
		actual := setURL(c.input)
		if actual != c.expected {
			t.Errorf("setURL(%v) == %v, expected %v", c.input, actual, c.expected)
		}
	}
}

// stringPtr es una función auxiliar para obtener un puntero a una cadena
func stringPtr(s string) *string {
	return &s
}

// TestCommandHelp prueba la función commandHelp
func TestCommandHelp(t *testing.T) {
	cfg := &config{}
	err := commandHelp(cfg, nil)
	if err != nil {
		t.Errorf("commandHelp() returned an error: %v", err)
	}
}

// TestCommandExit prueba la función commandExit
// Nota: Esta prueba está comentada porque commandExit llama a os.Exit, lo que terminaría el proceso de prueba
// func TestCommandExit(t *testing.T) {
// 	cfg := &config{}
// 	err := commandExit(cfg, nil)
// 	if err != nil {
// 		t.Errorf("commandExit() returned an error: %v", err)
// 	}
// }

// TestCommandMap prueba la función commandMap
func TestCommandMap(t *testing.T) {
	cfg := &config{
		NextURL: "https://pokeapi.co/api/v2/location-area/",
		Cache:   pokecache.NewCache(5 * time.Minute),
	}
	err := commandMap(cfg, nil)
	if err != nil {
		t.Errorf("commandMap() returned an error: %v", err)
	}
}

// TestCommandMapBack prueba la función commandMapBack
func TestCommandMapBack(t *testing.T) {
	cfg := &config{
		PreviousURL: "https://pokeapi.co/api/v2/location-area/",
		Cache:       pokecache.NewCache(5 * time.Minute),
	}
	err := commandMapBack(cfg, nil)
	if err != nil {
		t.Errorf("commandMapBack() returned an error: %v", err)
	}
}

// TestCommandExplore prueba la función commandExplore
func TestCommandExplore(t *testing.T) {
	cfg := &config{
		Cache: pokecache.NewCache(5 * time.Minute),
	}
	err := commandExplore(cfg, []string{"pastoria-city-area"})
	if err != nil {
		t.Errorf("commandExplore() returned an error: %v", err)
	}
}

// TestCommandCatch prueba la función commandCatch
func TestCommandCatch(t *testing.T) {
	cfg := &config{
		Cache:   pokecache.NewCache(5 * time.Minute),
		Pokedex: make(map[string]Pokemon),
	}
	err := commandCatch(cfg, []string{"pikachu"})
	if err != nil {
		t.Errorf("commandCatch() returned an error: %v", err)
	}
	if _, found := cfg.Pokedex["pikachu"]; !found {
		t.Errorf("commandCatch() did not catch the Pokémon")
	}
}

// TestCommandInspect prueba la función commandInspect
func TestCommandInspect(t *testing.T) {
	cfg := &config{
		Pokedex: map[string]Pokemon{
			"pikachu": {
				Name:   "pikachu",
				Height: 4,
				Weight: 60,
				Stats: []struct {
					BaseStat int
					Name     string
				}{
					{BaseStat: 35, Name: "hp"},
					{BaseStat: 55, Name: "attack"},
					{BaseStat: 40, Name: "defense"},
					{BaseStat: 50, Name: "special-attack"},
					{BaseStat: 50, Name: "special-defense"},
					{BaseStat: 90, Name: "speed"},
				},
				Types: []string{"electric"},
			},
		},
	}
	err := commandInspect(cfg, []string{"pikachu"})
	if err != nil {
		t.Errorf("commandInspect() returned an error: %v", err)
	}
}

// TestCommandPokedex prueba la función commandPokedex
func TestCommandPokedex(t *testing.T) {
	cfg := &config{
		Pokedex: map[string]Pokemon{
			"pikachu":   {},
			"bulbasaur": {},
		},
	}
	err := commandPokedex(cfg, nil)
	if err != nil {
		t.Errorf("commandPokedex() returned an error: %v", err)
	}
}
