package types

type PokemonMoves struct {
	Accuracy    int        `json:"accuracy"`
	DamageClass BaseStruct `json:"damage_class"`
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Priority    int        `json:"priority"`
	Power       int        `json:"power"`
	Type        BaseStruct `json:"type"`
}
