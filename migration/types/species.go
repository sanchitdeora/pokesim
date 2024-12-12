package types

type PokemonSpecies struct {
	EvolutionChain struct {
		Url string `json:"url"`
	} `json:"evolution_chain"`
	GrowthRate     BaseStruct       `json:"growth_rate"`
	PokedexNumbers []PokedexNumbers `json:"pokedex_numbers"`
}

const NationalPokedex = "national"

type PokedexNumbers struct {
	EntryNumber int        `json:"entry_number"`
	Pokedex     BaseStruct `json:"pokedex"`
}
