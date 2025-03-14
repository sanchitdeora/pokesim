package battle

type PokemonBattle interface {
	Introduction() error
    // ProcessAction(action *data.BattleAction, trainer *data.BattleTrainer)
    Conclusion() bool
    // GetBattleReport() *BattleState
}
