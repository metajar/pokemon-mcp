package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// PokemonResponse represents the Pokemon data we want from the API
type PokemonResponse struct {
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Types  []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"Pokemon Server 🎮",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Add Pokemon info tool
	pokemonTool := mcp.NewTool("get_pokemon",
		mcp.WithDescription("Get information about a Pokemon"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the Pokemon (lowercase)"),
		),
	)

	// Add Pokemon compare tool
	compareTool := mcp.NewTool("compare_pokemon",
		mcp.WithDescription("Compare two Pokemon"),
		mcp.WithString("pokemon1",
			mcp.Required(),
			mcp.Description("First Pokemon to compare (lowercase)"),
		),
		mcp.WithString("pokemon2",
			mcp.Required(),
			mcp.Description("Second Pokemon to compare (lowercase)"),
		),
	)

	// Register tools
	s.AddTool(pokemonTool, getPokemonHandler)
	s.AddTool(compareTool, comparePokemonHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func getPokemonHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	name, ok := arguments["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name must be a string"), nil
	}

	pokemon, err := fetchPokemon(name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error fetching Pokemon: %v", err)), nil
	}

	// Format the response
	types := make([]string, len(pokemon.Types))
	for i, t := range pokemon.Types {
		types[i] = t.Type.Name
	}

	response := fmt.Sprintf("🔍 Pokemon Information for %s:\n\n"+
		"Height: %d decimeters\n"+
		"Weight: %d hectograms\n"+
		"Types: %s\n\n"+
		"Base Stats:\n",
		strings.Title(pokemon.Name),
		pokemon.Height,
		pokemon.Weight,
		strings.Join(types, ", "))

	for _, stat := range pokemon.Stats {
		response += fmt.Sprintf("%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	return mcp.NewToolResultText(response), nil
}

func comparePokemonHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	pokemon1, ok := arguments["pokemon1"].(string)
	if !ok {
		return mcp.NewToolResultError("pokemon1 must be a string"), nil
	}

	pokemon2, ok := arguments["pokemon2"].(string)
	if !ok {
		return mcp.NewToolResultError("pokemon2 must be a string"), nil
	}

	p1, err := fetchPokemon(pokemon1)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error fetching %s: %v", pokemon1, err)), nil
	}

	p2, err := fetchPokemon(pokemon2)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error fetching %s: %v", pokemon2, err)), nil
	}

	// Create comparison text
	comparison := fmt.Sprintf("⚔️ Pokemon Comparison: %s vs %s\n\n",
		strings.Title(p1.Name),
		strings.Title(p2.Name))

	// Compare stats
	comparison += "Base Stats Comparison:\n"
	for i := 0; i < len(p1.Stats); i++ {
		stat1 := p1.Stats[i]
		stat2 := p2.Stats[i]
		comparison += fmt.Sprintf("%s: %d vs %d\n",
			stat1.Stat.Name,
			stat1.BaseStat,
			stat2.BaseStat)
	}

	return mcp.NewToolResultText(comparison), nil
}

func fetchPokemon(name string) (*PokemonResponse, error) {
	resp, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", strings.ToLower(name)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var pokemon PokemonResponse
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return nil, err
	}

	return &pokemon, nil
}
