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
	result := t.BattleConductor()

	// conclude battle
	t.handleBattleConclusion(result)
	return nil
}

func (t *TrainerBattle) BattleConductor() (result data.Result) {
	// loop runs through the battle
	for {
		t.log(fmt.Sprintf("%s Health: %v", t.getUserActiveName(), t.user.GetActivePokemon().BattleHP))
		t.log(fmt.Sprintf("%s Health: %v", t.getOpponentActiveName(), t.opponent.GetActivePokemon().BattleHP))

		if res, concluded := t.Conclusion(); concluded {
			return res
		}

		t.handleTurns()
	}
}

func (t *TrainerBattle) Conclusion() (result data.Result, concluded bool) {
	// check if battle concluded
	if concluded = t.user.IsDefeated() || t.opponent.IsDefeated(); !concluded {
		return data.Result{}, false
	}

	// check if user lost
	if t.user.IsDefeated() {
		result.UserWin = false
		result.Money = data.GetMoneyLost(t.opts.UserManager.GetUser())

		t.log(fmt.Sprintf("You lost the battle to %s!", t.getOpponenetName()))
		t.log(fmt.Sprintf("You lost $%v!", result.Money))
	} else {
		result.UserWin = true
		result.Money = data.GetPrizeMoney(t.opponent.GetTrainerType(), t.opponent.GetTrainer().Party)
		result.BonusItems = t.opponent.GetRewards().Items

		t.log(fmt.Sprintf("%s has won the battle!", t.getUserName()))
		t.log(fmt.Sprintf("You got $%v!", result.Money))

		// if gym battle; earn badge
		if t.opponent.GetTrainerType() == data.GymLeaderPrefix {
			result.BadgeEarned = t.opponent.GetRewards().Badge
			t.log(fmt.Sprintf("You earned a $%v!", result.BadgeEarned.Name))
		}
	}
	return result, true
}

func (t *TrainerBattle) handleBattleIntroduction() {
	t.log(fmt.Sprintf("%s chooses %s!", t.getOpponenetName(), t.getOpponentActiveName()))
	t.log(fmt.Sprintf("%s, I choose you!\n", t.getUserActiveName()))

	t.updatePokemonFaced()

	t.user.SetOpponentTarget(t.opponent.GetActivePokemon())
	t.opponent.SetOpponentTarget(t.user.GetActivePokemon())
}

func (t *TrainerBattle) handleTurns() {
	// get pokemon attack order
	userInput := t.user.GetInput()
	opponentInput := t.opponent.GetInput()
	battleInputs := GetTurnOrder(userInput, opponentInput)

	for _, input := range battleInputs {
		trainer := t.user.(*battletrainer.BattleTester)
		if input == opponentInput {
			trainer = t.opponent.(*battletrainer.BattleTester)
		}
		t.handleActions(input, trainer)
	}
}

func (t *TrainerBattle) handleActions(input data.BattleInput, trainer battletrainer.BattleTrainer) {
	switch input.Type {
	case data.Run:
		t.log("Trainer cannot run")
		return

	case data.Switch:
		t.log(fmt.Sprintf("%s is switching %s for %s", t.getUserName(),
			utils.ToCapitalizeFirstLetterOfEachWord(input.Selected.Pokemon.Name),
			utils.ToCapitalizeFirstLetterOfEachWord(input.Target.Pokemon.Name)),
		)
		trainer.HandleSwitch(input)
		t.updatePokemonFaced()

	case data.Bag:
		t.log(fmt.Sprintf("%s is using %s on %s", t.getUserName(),
			utils.ToCapitalizeFirstLetterOfEachWord(input.Item.Description),
			utils.ToCapitalizeFirstLetterOfEachWord(input.Selected.Pokemon.Name)),
		)
		if input.Item != nil && input.Item.Category == data.MedicalItems {
			trainer.HandleUseBag(input)
		}
		slog.Error("Item cannot be nil or category not medical item", "Item", input.Item)

	case data.Attack:
		slog.Info("Attack action")
		t.handleAttack(input.Selected.Pokemon)
	}
}

func (t *TrainerBattle) handleAttack(targetPokemon data.Pokemon) {
	slog.Info("Attack action")
}

func (t *TrainerBattle) handleBattleConclusion(result data.Result) {
	// update user stats and rewards

	// pokemon evolution

	// close all channels
}

func (t *TrainerBattle) updatePokemonFaced() {
	t.user.GetActivePokemon().PokemonFaced = append(t.user.GetActivePokemon().PokemonFaced, *t.opponent.GetActivePokemon())
	t.opponent.GetActivePokemon().PokemonFaced = append(t.user.GetActivePokemon().PokemonFaced, *t.user.GetActivePokemon())
}

func (t *TrainerBattle) log(msg string) {
	if err := t.user.SendBattleLog(msg); err != nil {
		slog.Error("error sending battle log", "error", err)
		panic(err)
	}
	if err := t.opponent.SendBattleLog(msg); err != nil {
		slog.Error("error sending battle log", "error", err)
		panic(err)
	}
	// log to console
	slog.Info(msg)
}

func (t *TrainerBattle) getUserName() string {
	return utils.ToCapitalizeFirstLetterOfEachWord(t.user.GetTrainer().Name)
}

func (t *TrainerBattle) getOpponenetName() string {
	return fmt.Sprintf("%s %s", t.opponent.GetTrainerType(), utils.ToCapitalizeFirstLetterOfEachWord(t.opponent.GetTrainer().Name))
}

func (t *TrainerBattle) getUserActiveName() string {
	return utils.ToCapitalizeFirstLetterOfEachWord(t.user.GetActivePokemon().Name)
}

func (t *TrainerBattle) getOpponentActiveName() string {
	return utils.ToCapitalizeFirstLetterOfEachWord(t.opponent.GetActivePokemon().Name)
}
