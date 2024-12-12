package data

type BattlePokemon struct {
	Pokemon
	BattleHP int
	// deprecated
	IsFainted    bool
	CanEvolve    bool
	PokemonFaced []BattlePokemon
}

type Result struct {
	UserWin     bool
	Money       int
	BonusItems  ItemMap
	BadgeEarned BadgeType
}

type BattleOpts struct {
	Type string `json:"type"`
}

type BattleInput struct {
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

type BattleAction struct {
	Type   BattleActionType
	Action interface{}
}

type ActionAttack struct {
	Move   Moves
	Target BattlePokemon
}

type ActionItem struct {
	Item   Item
	Target BattlePokemon
}

type ActionSwitch struct {
	Target BattlePokemon
}
