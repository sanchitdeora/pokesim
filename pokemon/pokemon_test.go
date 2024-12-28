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
			BaseStats: data.PokemonStats{
				HP: data.PokemonStat{
					Value: 30,
				},
				Attack: data.PokemonStat{
					Value: 56,
				},
				Defense: data.PokemonStat{
					Value: 35,
				},
				SpecialAttack: data.PokemonStat{
					Value: 25,
				},
				SpecialDefense: data.PokemonStat{
					Value: 35,
				},
				Speed: data.PokemonStat{
					Value: 72,
				},
			},
		},
		Stats: data.PokemonStats{
			HP: data.PokemonStat{
				Value: 58,
			},
			Attack: data.PokemonStat{
				Value: 38,
			},
			Defense: data.PokemonStat{
				Value: 26,
			},
			SpecialAttack: data.PokemonStat{
				Value: 20,
			},
			SpecialDefense: data.PokemonStat{
				Value: 26,
			},
			Speed: data.PokemonStat{
				Value: 48,
			},
		},
		Level: 30,
		ExperienceLeft: 612,
	}

	faintedPokemon := data.Pokemon{
		BasePokemon: data.BasePokemon{
			Name: "Ratatta",
			GrowthRate: data.MediumFast,
			BaseExperience: 51,
		},
		Level: 30,
	}

	pm.ExperienceGain(pokemon, faintedPokemon)

	assert.Equal(t, 31, pokemon.Level)
	assert.Equal(t, math.Pow(31, 3), float64(pokemon.ExperienceLeft))
}
