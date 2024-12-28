package data

type BattlePokemon struct {
	Pokemon
	BattleHP     int
	PokemonFaced []string
	// deprecated
	CanEvolve bool
	IsFainted bool
}

type BattleResultStatus string

const (
	Won  BattleResultStatus = "win"
	Lost BattleResultStatus = "loss"
)

type Result struct {
	Status         BattleResultStatus
	Money          int
	BonusItems     ItemMap
	BadgeEarned    BadgeType
	// deprecated
	UserWin bool
}

type BattleOpts struct {
	Type string `json:"type"`
}

type BattleAction struct {
	ID       string
	Type     BattleActionType
	Selected *BattlePokemon
	Target   *BattlePokemon
	Move     *Moves
	Item     *Item
	// deprecated
	IsUser bool
}

type BattleActionType string

const (
	Attack BattleActionType = "attack"
	Switch BattleActionType = "switch"
	Bag    BattleActionType = "bag"
	Run    BattleActionType = "run"
)

func CreateBattlePokemon(pokemon Pokemon) *BattlePokemon {
	return &BattlePokemon{
		Pokemon:      pokemon,
		BattleHP:     pokemon.Stats.HP.Value,
		PokemonFaced: make([]string, 0),
	}
}
