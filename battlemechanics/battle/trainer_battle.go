package battle

import (
	"fmt"
	"log/slog"

	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/sanchitdeora/PokeSim/utils"
)

type TrainerBattle struct {
	opts     TrainerBattleOpts
	user     battletrainer.BattleTrainer
	opponent battletrainer.BattleTrainer
}

type TrainerBattleOpts struct {
	UserManager usermanagement.UserManager
}

func TrainerBattleManager(opts TrainerBattleOpts, user battletrainer.BattleTrainer, opponent battletrainer.BattleTrainer) PokemonBattle {
	if user == nil {
		slog.Error("user is nil")
		return nil
	}
	if opponent == nil {
		slog.Error("opponent is nil")
		return nil
	}

	return &TrainerBattle{
		opts:     opts,
		user:     user,
		opponent: opponent,
	}
}

func (t *TrainerBattle) Introduction() error {
	// Initate battle
	slog.Info("Battle Initiating...")
	t.handleBattleIntroduction()

	// Battle Processer
	t.BattleConductor()

	// Conclude battle
	t.user.HandleBattleEnd(t.user.CalculateResult(t.opponent))
	t.opponent.HandleBattleEnd(t.opponent.CalculateResult(t.user))

	return nil
}

func (t *TrainerBattle) BattleConductor() {
	// loop runs through the battle
	for {
		t.log(fmt.Sprintf("%s Health: %v", getActivePokemonName(t.user), t.user.GetActivePokemon().BattleHP))
		t.log(fmt.Sprintf("Opponent %s Health: %v", getActivePokemonName(t.opponent), t.opponent.GetActivePokemon().BattleHP))

		if t.Conclusion() {
			break
		}

		t.handleTurns()
	}
}

func (t *TrainerBattle) Conclusion() bool {
	// check if battle concluded
	return t.user.IsDefeated() || t.opponent.IsDefeated()
}

func (t *TrainerBattle) handleBattleIntroduction() {
	t.log(fmt.Sprintf("%s chooses %s!", t.getOpponenetName(), getActivePokemonName(t.opponent)))
	t.log(fmt.Sprintf("%s, I choose you!", getActivePokemonName(t.user)))

	t.updatePokemonFaced()
	t.setOpponentTargets()
}

func (t *TrainerBattle) handleTurns() {
	// get pokemon attack order
	userAction := t.user.GetAction()
	opponentAction := t.opponent.GetAction()
	battleActions := GetTurnOrder(userAction, opponentAction)

	for _, action := range battleActions {
		trainer := t.user
		target := t.opponent
		if action.ID == opponentAction.ID {
			trainer = t.opponent
			target = t.user
		}
		t.handleActions(action, trainer, target)
	}
}

func (t *TrainerBattle) handleActions(action data.BattleAction, trainer battletrainer.BattleTrainer, target battletrainer.BattleTrainer) {
	switch action.Type {
	case data.Run:
		t.log("Trainer cannot run")
		return
	case data.Switch:
		t.handleSwitch(action, trainer)
	case data.Bag:
		t.log(fmt.Sprintf("%s is using %s on %s", t.getUserName(),
			utils.ToCapitalizeFirstLetterOfEachWord(action.Item.Description),
			utils.ToCapitalizeFirstLetterOfEachWord(action.Selected.Pokemon.Name)),
		)
		if action.Item != nil && action.Item.Category == data.MedicalItems {
			err := trainer.HandleAction(action)
			if err != nil {
				//TODO: Add appropriate battle logs. Let user select again similar to Run
				slog.Error(err.Error())
			}
		}
		slog.Error("Item cannot be nil or category not medical item", "Item", action.Item)

	case data.Attack:
		slog.Info("Attack action")
		t.handleAttack(action, trainer, target)
	}
}

func (t *TrainerBattle) handleSwitch(switchAction data.BattleAction, trainer battletrainer.BattleTrainer) {
	t.log(fmt.Sprintf("%s is switching %s for %s", t.getUserName(),
		utils.ToCapitalizeFirstLetterOfEachWord(switchAction.Selected.Pokemon.Name),
		utils.ToCapitalizeFirstLetterOfEachWord(switchAction.Target.Pokemon.Name)),
	)
	err := trainer.HandleAction(switchAction)
	if err != nil {
		//TODO: Add appropriate battle logs. Let user select again similar to Run
		slog.Error(err.Error())
	}
	t.updatePokemonFaced()
	t.setOpponentTargets()
}

func (t *TrainerBattle) handleAttack(attackAction data.BattleAction, trainer battletrainer.BattleTrainer, target battletrainer.BattleTrainer) {
	if attackAction.Selected.BattleHP == 0 {
		slog.Info(fmt.Sprintf("%s fainted and cannot attack", utils.ToCapitalizeFirstLetterOfEachWord(attackAction.Selected.Pokemon.Name)))
		return
	}

	t.log(fmt.Sprintf("%s used %s", getActivePokemonName(trainer), utils.ToCapitalizeFirstLetterOfEachWord(attackAction.Move.Name)))

	damagePts := t.performAttack(attackAction, trainer, target)
	t.log(fmt.Sprintf("%s did %v points of damage to %s", getActivePokemonName(trainer), damagePts, getActivePokemonName(target)))

	if target.GetActivePokemon().BattleHP == 0 {
		t.log(fmt.Sprintf("%s has fainted!", getActivePokemonName(target)))

		trainer.HandleTargetPokemonFainted(target.GetActivePokemon())

		nextUnfaintedPokemonIndex, count := GetNextUnfaintedPokemonAndCount(target.GetParty())
		if count > 0 {
			// switch active pokemon with next unfainted in party and push fainted pokemon at end
			switchAction := data.BattleAction{
				Type:     data.Switch,
				Selected: target.GetActivePokemon(),
				Target:   target.GetParty()[nextUnfaintedPokemonIndex],
			}
			t.handleSwitch(switchAction, target)
		}
	}
}

func (t *TrainerBattle) performAttack(attackAction data.BattleAction, trainer battletrainer.BattleTrainer, target battletrainer.BattleTrainer) int {
	//TODO: ADD VALIDTION -- confirm action selected and target are trainer and target trainer's active pokemon
	damagePoints := data.CalculateAttackDamage(trainer.GetActivePokemon(), target.GetActivePokemon(), attackAction.Move, 1.0)

	target.GetActivePokemon().BattleHP = utils.Max(0, target.GetActivePokemon().BattleHP-damagePoints)
	return damagePoints
}

func (t *TrainerBattle) updatePokemonFaced() {
	t.user.GetActivePokemon().PokemonFaced = append(t.user.GetActivePokemon().PokemonFaced, t.opponent.GetActivePokemon().Pokemon.PokemonUUID)
	t.opponent.GetActivePokemon().PokemonFaced = append(t.opponent.GetActivePokemon().PokemonFaced, t.user.GetActivePokemon().Pokemon.PokemonUUID)
}

func (t *TrainerBattle) setOpponentTargets() {
	// set opponent targets
	t.user.SetOpponentTarget(t.opponent.GetActivePokemon())
	t.opponent.SetOpponentTarget(t.user.GetActivePokemon())
}

func (t *TrainerBattle) log(msg string) {
	t.user.SendBattleLog(msg)

	// t.opponent.SendBattleLog(msg)
}

// log formatters
func (t *TrainerBattle) getUserName() string {
	return utils.ToCapitalizeFirstLetterOfEachWord(t.user.GetTrainer().Name)
}

func (t *TrainerBattle) getOpponenetName() string {
	return fmt.Sprintf("%s %s", t.opponent.GetTrainerType(), utils.ToCapitalizeFirstLetterOfEachWord(t.opponent.GetTrainer().Name))
}

func getActivePokemonName(t battletrainer.BattleTrainer) string {
	return utils.ToCapitalizeFirstLetterOfEachWord(t.GetActivePokemon().Name)
}
