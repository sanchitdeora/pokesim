package types

type Evolution struct {
	ID    int   `json:"id"`
	Chain Chain `json:"chain"`
}

type Chain struct {
	EvolutionDetails []EvolutionDetails `json:"evolution_details"`
	EvolvesTo        []Chain            `json:"evolves_to"`
	Species          BaseStruct         `json:"species"`
}

type EvolutionDetails struct {
	MinLevel     int        `json:"min_level"`
	Trigger      BaseStruct `json:"trigger"`
	MinHappiness int        `json:"min_happiness"`
	HeldItem     BaseStruct `json:"held_item"`
	TimeOfDay    string     `json:"time_of_day"`
	Location     string     `json:"location"`
}
