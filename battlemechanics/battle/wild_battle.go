package battle

import (
	"fmt"
	"log/slog"

	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/sanchitdeora/PokeSim/utils"
)

type WildBattle struct {
	opts     WildBattleOpts
	user     battletrainer.BattleTrainer
	opponent battletrainer.BattleTrainer

	wildCaught    bool
	successfulRun bool
	runAttemtps   int
}

type WildBattleOpts struct {
	UserManager usermanagement.UserManager
}

func WildBattleManager(opts WildBattleOpts, user battletrainer.BattleTrainer, opponent battletrainer.BattleTrainer) PokemonBattle {
	if user == nil {
		slog.Error("user is nil")
		return nil
	}
	if opponent == nil {
		slog.Error("opponent is nil")
		return nil
	}
	if _, ok := opponent.(*battletrainer.BattleWild); !ok {
		slog.Error("opponent is not a wild trainer")
		return nil
	}

	return &WildBattle{
		opts:     opts,
		user:     user,
		opponent: opponent,

		wildCaught:    false,
		successfulRun: false,
		runAttemtps:   1,
	}
}

func (w *WildBattle) Introduction() error {
	// Initate battle
	slog.Info("Battle Initiating...")
	w.handleBattleIntroduction()

	// Battle Processer
	w.BattleConductor()

	// Conclude battle
	w.user.HandleBattleEnd(w.user.CalculateResult(w.opponent))
	// w.opponent.HandleBattleEnd(w.opponent.CalculateResult(w.user))

	return nil
}

func (w *WildBattle) BattleConductor() {
	// loop runs through the battle
	for {
		w.log(fmt.Sprintf("%s Health: %v", getActivePokemonName(w.user), w.user.GetActivePokemon().BattleHP))
		w.log(fmt.Sprintf("Opponent %s Health: %v", getActivePokemonName(w.opponent), w.user.GetActivePokemon().BattleHP))

		if w.Conclusion() {
			break
		}

		w.handleTurns()
	}
}

func (w *WildBattle) Conclusion() bool {
	// check if battle concluded
	return w.user.IsDefeated() || w.opponent.IsDefeated() || w.successfulRun || w.wildCaught
}

func (w *WildBattle) handleBattleIntroduction() {
	w.log(fmt.Sprintf("Wild %s appears!", getActivePokemonName(w.opponent)))
	w.log(fmt.Sprintf("%s, I choose you!", getActivePokemonName(w.user)))

	w.updatePokemonFaced()
	w.setOpponentTargets()
}

func (w *WildBattle) handleTurns() {
	// get pokemon attack order
	userAction := w.user.GetAction()
	opponentAction := w.opponent.GetAction()
	battleActions := GetTurnOrder(userAction, opponentAction)

	for _, action := range battleActions {
		if w.Conclusion() {
			break
		}
		trainer := w.user
		target := w.opponent
		if action.ID == opponentAction.ID {
			trainer = w.opponent
			target = w.user
		}
		w.handleActions(action, trainer, target)
	}
}

func (w *WildBattle) handleActions(action data.BattleAction, trainer battletrainer.BattleTrainer, target battletrainer.BattleTrainer) {
	switch action.Type {
	case data.Run:
		w.handleRun(action)
	case data.Switch:
		w.handleSwitch(action, trainer)
	case data.Bag:
		w.handleUseBag(action, trainer)
	case data.Attack:
		slog.Info("Attack action")
		w.handleAttack(action, trainer, target)
	}
}

func (w *WildBattle) handleSwitch(switchAction data.BattleAction, trainer battletrainer.BattleTrainer) {
	w.log(fmt.Sprintf("%s is switching %s for %s", w.getUserName(),
		utils.ToCapitalizeFirstLetterOfEachWord(switchAction.Selected.Pokemon.Name),
		utils.ToCapitalizeFirstLetterOfEachWord(switchAction.Target.Pokemon.Name)),
	)
	err := trainer.HandleAction(switchAction)
	if err != nil {
		//TODO: Add appropriate battle logs. Let user select again similar to Run
		slog.Error(err.Error())
	}
	w.updatePokemonFaced()
	w.setOpponentTargets()
}

func (w *WildBattle) handleRun(action data.BattleAction) {
	w.successfulRun = IsRunSuccessful(action.Selected, action.Target, w.runAttemtps)
	if w.successfulRun {
		w.log("Got away safely!")
		return
	} else {
		w.runAttemtps++
		w.log("Cannot run")
	}
}

func (w *WildBattle) handleUseBag(action data.BattleAction, trainer battletrainer.BattleTrainer) {
	w.log(fmt.Sprintf("%s is using %s on %s", w.getUserName(),
		utils.ToCapitalizeFirstLetterOfEachWord(string(data.GetItemNameFromItem(*action.Item))),
		utils.ToCapitalizeFirstLetterOfEachWord(action.Selected.Pokemon.Name)),
	)
	if action.Item != nil {
		err := trainer.HandleAction(action)
		if err != nil {
			//TODO: Add appropriate battle logs. Let user select again similar to Run
			slog.Error(err.Error())
		}

		if action.Item.Category == data.PokeBalls {
			w.wildCaught = IsPokemonInParty(trainer.GetParty(), action.Selected)
		}
	} else {
		slog.Error("Item cannot be nil", "Item", action.Item)
	}
}

func (w *WildBattle) handleAttack(attackAction data.BattleAction, trainer battletrainer.BattleTrainer, target battletrainer.BattleTrainer) {
	if attackAction.Selected.BattleHP == 0 {
		slog.Info(fmt.Sprintf("%s fainted and cannot attack", utils.ToCapitalizeFirstLetterOfEachWord(attackAction.Selected.Pokemon.Name)))
		return
	}

	w.log(fmt.Sprintf("%s used %s", getActivePokemonName(trainer), utils.ToCapitalizeFirstLetterOfEachWord(attackAction.Move.Name)))

	damagePts := w.performAttack(attackAction, trainer, target)
	w.log(fmt.Sprintf("%s did %v points of damage to %s", getActivePokemonName(trainer), damagePts, getActivePokemonName(target)))

	if target.GetActivePokemon().BattleHP == 0 {
		w.log(fmt.Sprintf("%s has fainted!", getActivePokemonName(target)))

		trainer.HandleTargetPokemonFainted(target.GetActivePokemon())

		nextUnfaintedPokemonIndex, count := GetNextUnfaintedPokemonAndCount(target.GetParty())
		if count > 0 {
			// switch active pokemon with next unfainted in party and push fainted pokemon at end
			switchAction := data.BattleAction{
				Type:     data.Switch,
				Selected: target.GetActivePokemon(),
				Target:   target.GetParty()[nextUnfaintedPokemonIndex],
			}
			w.handleSwitch(switchAction, target)
		}
	}
}

func (w *WildBattle) performAttack(attackAction data.BattleAction, trainer battletrainer.BattleTrainer, target battletrainer.BattleTrainer) int {
	//TODO: ADD VALIDTION -- confirm action selected and target are trainer and target trainer's active pokemon
	damagePoints, effect, isCrit := data.CalculateAttackDamage(trainer.GetActivePokemon(), target.GetActivePokemon(), attackAction.Move, 1.0)
	if isCrit {
		w.log("Critical hit!")
	}
	if effect != data.NOR {
		if effect == data.MNE || effect == data.NVR {
			w.log("Not very effective!")
		} else if effect == data.SUP || effect == data.HYP {
			w.log("Super effective!")
		} else {
			w.log("This move has No effect to the target!")
		}
	}

	target.GetActivePokemon().BattleHP = max(0, target.GetActivePokemon().BattleHP-damagePoints)
	return damagePoints
}

func (w *WildBattle) updatePokemonFaced() {
	w.user.GetActivePokemon().PokemonFaced = append(w.user.GetActivePokemon().PokemonFaced, w.opponent.GetActivePokemon().PokemonUUID)
}

func (w *WildBattle) setOpponentTargets() {
	// set opponent targets
	w.user.SetOpponentTarget(w.opponent.GetActivePokemon())
	w.opponent.SetOpponentTarget(w.user.GetActivePokemon())
}

func (w *WildBattle) log(msg string) {
	w.user.SendBattleLog(msg)
	// t.opponent.SendBattleLog(msg)
}

// log formatter

func (w *WildBattle) getUserName() string {
	return utils.ToCapitalizeFirstLetterOfEachWord(w.user.GetTrainer().Name)
}
