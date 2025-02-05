package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/glopez94/pokedex/internal/pokecache"
)

// Definir la estructura cliCommand
type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

// Definir la estructura config para manejar la paginación
type config struct {
	NextURL     string
	PreviousURL string
	Cache       *pokecache.Cache
	Pokedex     map[string]Pokemon
}

// Definir la estructura para un Pokémon capturado
type Pokemon struct {
	Name           string
	BaseExperience int
	Height         int
	Weight         int
	Stats          []struct {
		BaseStat int
		Name     string
	}
	Types []string
}

// Definir la estructura para la respuesta de la API
type locationAreaResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

// Definir la estructura para la respuesta de explorar área
type exploreAreaResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// Definir la estructura para la respuesta de la API de Pokémon
type pokemonResponse struct {
	BaseExperience int    `json:"base_experience"`
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

// Declarar el mapa de comandos sin inicializarlo
var commands map[string]cliCommand

func main() {
	// Inicializar el mapa de comandos dentro de la función main
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore a specific location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a specific Pokémon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a captured Pokémon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all captured Pokémon",
			callback:    commandPokedex,
		},
	}

	cfg := &config{
		NextURL: "https://pokeapi.co/api/v2/location-area/",
		Cache:   pokecache.NewCache(5 * time.Minute),
		Pokedex: make(map[string]Pokemon),
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		// Esperar la entrada del usuario
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		words := cleanInput(input)
		// Verificar si hay al menos una palabra en la entrada
		if len(words) > 0 {
			commandName := words[0]
			// Buscar el comando en el registro de comandos
			if command, found := commands[commandName]; found {
				// Ejecutar la función de devolución de llamada del comando y manejar cualquier error
				if err := command.callback(cfg, words[1:]); err != nil {
					fmt.Printf("Error: %v\n", err)
				}
			} else {
				// Si el comando no se encuentra en el registro
				fmt.Println("Unknown command")
			}
		}
	}
}

// Función para limpiar la entrada del usuario
func cleanInput(text string) []string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return strings.Fields(text)
}

// Devolución de llamada para el comando "exit"
func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// Devolución de llamada para el comando "help"
func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

// Devolución de llamada para el comando "map"
func commandMap(cfg *config, args []string) error {
	if cfg.NextURL == "" {
		fmt.Println("No more locations to display.")
		return nil
	}
	return fetchAndDisplayLocations(cfg.NextURL, cfg)
}

// Devolución de llamada para el comando "mapb"
func commandMapBack(cfg *config, args []string) error {
	if cfg.PreviousURL == "" {
		fmt.Println("You're on the first page.")
		return nil
	}
	return fetchAndDisplayLocations(cfg.PreviousURL, cfg)
}

// Devolución de llamada para el comando "explore"
func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a location area to explore")
	}
	areaName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", areaName)
	return fetchAndDisplayPokemon(url, cfg)
}

// Devolución de llamada para el comando "catch"
func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a Pokémon to catch")
	}
	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch Pokémon: %v", err)
	}
	defer resp.Body.Close()

	var data pokemonResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Determinar la probabilidad de atrapar al Pokémon con un factor de ajuste
	adjustmentFactor := 90.0 // Puedes ajustar este valor para aumentar la probabilidad de captura
	catchProbability := adjustmentFactor * (100.0 / float64(data.BaseExperience))
	if rand.Float64()*100 < catchProbability {
		fmt.Printf("%s was caught!\n", pokemonName)
		cfg.Pokedex[pokemonName] = Pokemon{
			Name:           data.Name,
			BaseExperience: data.BaseExperience,
			Height:         data.Height,
			Weight:         data.Weight,
			Stats: make([]struct {
				BaseStat int
				Name     string
			}, len(data.Stats)),
			Types: make([]string, len(data.Types)),
		}
		for i, stat := range data.Stats {
			cfg.Pokedex[pokemonName].Stats[i] = struct {
				BaseStat int
				Name     string
			}{
				BaseStat: stat.BaseStat,
				Name:     stat.Stat.Name,
			}
		}
		for i, t := range data.Types {
			cfg.Pokedex[pokemonName].Types[i] = t.Type.Name
		}
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

// Devolución de llamada para el comando "inspect"
func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a Pokémon to inspect")
	}
	pokemonName := args[0]
	pokemon, found := cfg.Pokedex[pokemonName]
	if !found {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t)
	}

	return nil
}

// Devolución de llamada para el comando "pokedex"
func commandPokedex(cfg *config, args []string) error {
	if len(cfg.Pokedex) == 0 {
		fmt.Println("Your Pokedex is empty.")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range cfg.Pokedex {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}

// Función para obtener y mostrar las áreas de ubicación
func fetchAndDisplayLocations(url string, cfg *config) error {
	// Verificar si la URL ya está en la caché
	if data, found := cfg.Cache.Get(url); found {
		return displayLocations(data, cfg)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch locations: %v", err)
	}
	defer resp.Body.Close()

	var data locationAreaResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Convertir la respuesta a JSON y agregarla a la caché
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	cfg.Cache.Add(url, jsonData)

	return displayLocations(jsonData, cfg)
}

// Función para mostrar las áreas de ubicación
func displayLocations(data []byte, cfg *config) error {
	var locationData locationAreaResponse
	if err := json.Unmarshal(data, &locationData); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	for _, result := range locationData.Results {
		fmt.Println(result.Name)
	}

	cfg.NextURL = setURL(locationData.Next)
	cfg.PreviousURL = setURL(locationData.Previous)

	return nil
}

// Función para obtener y mostrar los Pokémon en un área de ubicación
func fetchAndDisplayPokemon(url string, cfg *config) error {
	// Verificar si la URL ya está en la caché
	if data, found := cfg.Cache.Get(url); found {
		return displayPokemon(data)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch location area: %v", err)
	}
	defer resp.Body.Close()

	var data exploreAreaResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Convertir la respuesta a JSON y agregarla a la caché
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	cfg.Cache.Add(url, jsonData)

	return displayPokemon(jsonData)
}

// Función para mostrar los Pokémon en un área de ubicación
func displayPokemon(data []byte) error {
	var exploreData exploreAreaResponse
	if err := json.Unmarshal(data, &exploreData); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range exploreData.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func setURL(url *string) string {
	if url != nil {
		return *url
	}
	return ""
}
