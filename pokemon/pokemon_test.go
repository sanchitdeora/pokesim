package pokemon_test

import (
	"math"
	"testing"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/pokemon"
	"github.com/stretchr/testify/assert"
)

func createPokemonManager() pokemon.PokemonManager {
	return pokemon.NewPokemonManager(pokemon.PokemonOpts{})
}

func TestExperienceGain(t *testing.T) {
	pm := createPokemonManager()
	pokemon := &data.Pokemon{
		BasePokemon: data.BasePokemon{
			Name: "Ratatta",
			GrowthRate: data.MediumFast,
		},
		Level: 30,
		ExperienceLeft: 100,
	}

	canEvovle := pm.ExperienceGain(100, pokemon)

	assert.False(t, canEvovle)
	assert.Equal(t, 31, pokemon.Level)
	assert.Equal(t, math.Pow(31, 3), float64(pokemon.ExperienceLeft))
}
