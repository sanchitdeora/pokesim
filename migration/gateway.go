package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sanchitdeora/PokeSim/migration/types"
)

func GetPokemon(url string) (*types.Pokemon, error) {
	return getFromAPI[types.Pokemon](url)
}

func GetPokemonSpecies(url string) (*types.PokemonSpecies, error) {
	return getFromAPI[types.PokemonSpecies](url)
}

func GetPokemonMoves(url string) (*types.PokemonMoves, error) {
	return getFromAPI[types.PokemonMoves](url)
}

func GetEvolutionChain(url string) (*types.Evolution, error) {
	return getFromAPI[types.Evolution](url)
}

func DownloadPNG(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func getFromAPI[T any](url string) (*T, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	var body T
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}
	return &body, nil
}
