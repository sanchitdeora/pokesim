package data

type InBattlePokemon struct {
	Pokemon   *Pokemon
	BattleHP  int
	IsFainted bool
	CanEvolve bool
}

type BattleReport struct {
	UserWin     bool
	Money       int
	BonusItems  ItemMap
	BadgeEarned *BadgeType
}

type BattleOpts struct {
	Type string `json:"type"`
}

type BattleInput struct {
	Type           BattleInputType
	CurrentPokemon *InBattlePokemon
	Target         *InBattlePokemon
	Move           *Moves
	Item           *Item
	IsUser         bool
}

type BattleInputType string

const (
	Attack BattleInputType = "attack"
	Switch BattleInputType = "switch"
	Bag    BattleInputType = "bag"
	Run    BattleInputType = "run"
)
