package usermanagement

import "github.com/sanchitdeora/PokeSim/data"

func GetPokemonIndexInParty(party []*data.Pokemon, pokemonName *data.Pokemon) int {
	for i, p := range party {
		if p.ID == pokemonName.ID {
			return i
		}
	}
	return -1
}
