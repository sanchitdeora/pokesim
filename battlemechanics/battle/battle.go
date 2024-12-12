package battle

import (
	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/data"
)

type PokemonBattle interface {
	Introduction() error
    // ProcessAction(action *data.BattleInput, trainer *data.BattleTrainer)
    Conclusion() (result data.Result, concluded bool)
    // GetBattleReport() *BattleState
}

type WildBattle struct{
	user battletrainer.BattleTrainer
	wild data.BattlePokemon
}
