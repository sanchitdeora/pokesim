package main

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestLoadPokemonJson(t *testing.T) {
// 	pokemon, err := migration.LoadPokemonJson("C:\\Projects\\Go-projects\\src\\PokéSim\\pokemon.json")
// 	assert.NoError(t, err)

// 	assert.NotNil(t, pokemon)
// }

// func TestLoadPokemonEvolutionJson(t *testing.T) {
// 	pokemon, err := migration.LoadEvolutionChainJson("C:\\Projects\\Go-projects\\src\\PokéSim\\pokemonEvolutionChain.json")
// 	assert.NoError(t, err)

// 	assert.NotNil(t, pokemon)
// }

func TestMigration(t *testing.T) {
	file, err := os.OpenFile("../test_manual.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        panic(err)
    }
    opts := migrationOpts{csv.NewWriter(file), file}
	defer opts.CsvFile.Close()
	defer opts.CsvWriter.Flush()

	opts.MigratePokemonToAsset("https://pokeapi.co/api/v2/pokemon/63/")
	
	assert.True(t, true)
}
