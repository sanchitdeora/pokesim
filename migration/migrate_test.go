package main

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	file, err := os.OpenFile("../test_manual.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	opts := migrationOpts{
		CsvWriter: csv.NewWriter(file),
		CsvFile:   file,

		WildEncounters: make(map[data.BasePokemonID]data.WildEncounter, 0),
	}
	defer opts.CsvFile.Close()
	defer opts.CsvWriter.Flush()

	opts.MigratePokemonToAsset("https://pokeapi.co/api/v2/pokemon/133/")

	assert.True(t, true)
}

func TestWildEncounters(t *testing.T) {
	file, err := os.OpenFile("../test_manual.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	opts := migrationOpts{
		CsvWriter: csv.NewWriter(file),
		CsvFile:   file,

		WildEncounters: make(map[data.BasePokemonID]data.WildEncounter, 0),
	}

	defer opts.CsvFile.Close()
	defer opts.CsvWriter.Flush()

	opts.MigratePokemonToAsset("https://pokeapi.co/api/v2/pokemon/1/")

	assert.True(t, true)
	assert.Equal(t, 1, len(opts.WildEncounters))
}
