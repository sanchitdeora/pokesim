package battletrainer

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/logger"
	"github.com/sanchitdeora/PokeSim/pokemon"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/sanchitdeora/PokeSim/utils"
)

type BattleTrainer interface {
	GetTrainer() *data.BaseTrainer
	GetTrainerType() data.TrainerClass
	GetActivePokemon() *data.BattlePokemon
	GetParty() []*data.BattlePokemon
	GetRewards() data.Rewards
	IsDefeated() bool
	CalculateResult(opponent BattleTrainer) data.Result
	HandleTargetPokemonFainted(faintedPokemon *data.BattlePokemon)

	// Only used for testing
	SetOpponentTarget(target *data.BattlePokemon)

	HandleAction(action data.BattleAction) error
	GetAction() data.BattleAction

	HandleBattleEnd(result data.Result)
	SendBattleLog(message string)
}

type BattleTrainerOpts struct {
	PokemonService pokemon.PokemonService
	UserManager    usermanagement.UserManager
}

type BattleTrainerImpl struct {
	BattleTrainerOpts
	Trainer       *data.BaseTrainer
	ActivePokemon *data.BattlePokemon
	BattleParty   []*data.BattlePokemon
	Rewards       data.Rewards

	Logger logger.Logger

	Opponent *data.BattlePokemon
}

func (b *BattleTrainerImpl) GetTrainer() *data.BaseTrainer {
	return b.Trainer
}

func (b *BattleTrainerImpl) GetTrainerType() data.TrainerClass {
	return data.RivalPrefix
}

func (b *BattleTrainerImpl) GetActivePokemon() *data.BattlePokemon {
	return b.ActivePokemon
}

func (b *BattleTrainerImpl) GetParty() []*data.BattlePokemon {
	return b.BattleParty
}

func (b *BattleTrainerImpl) GetRewards() data.Rewards {
	return b.Rewards
}

func (b *BattleTrainerImpl) IsDefeated() bool {
	for _, p := range append([]*data.BattlePokemon{b.ActivePokemon}, b.BattleParty...) {
		if p.BattleHP > 0 {
			return false
		}
	}
	return true
}

func (b *BattleTrainerImpl) HandleTargetPokemonFainted(faintedPokemon *data.BattlePokemon) {
	for _, p := range append([]*data.BattlePokemon{b.ActivePokemon}, b.BattleParty...) {
		if utils.Contains(p.PokemonFaced, faintedPokemon.Pokemon.PokemonUUID) && p.BattleHP > 0 {
			b.PokemonService.ExperienceGain(p.Pokemon, *faintedPokemon.Pokemon)
			b.PokemonService.EvGain(p.Pokemon, faintedPokemon.BaseStats)

			// b.ActivePokemon.CanEvolve = b.PokemonService.CanEvolve(p.Pokemon)
		}
	}
}

func (b *BattleTrainerImpl) SetOpponentTarget(target *data.BattlePokemon) {
	b.Opponent = target
}

func (b *BattleTrainerImpl) SendBattleLog(message string) {
	b.Logger.Log(message)
}

func (b *BattleTrainerImpl) HandleAction(action data.BattleAction) error {
	slog.Info("Handling action", "action", action.Type)

	switch action.Type {
	case data.Switch:
		if b.HandleSwitch(action) != nil {
			return errors.New("switch failed")
		}

		b.SendBattleLog(fmt.Sprintf("%s is switching their pokemon!", b.Trainer.Name))
		b.SendBattleLog(fmt.Sprintf("%s, chooses %s!", b.Trainer.Name, action.Selected.Pokemon.Name))
		return nil

	case data.Bag:
		return b.HandleUseBag(action)

	default:
		return errors.New("invalid action type. Expected Switch or UseBag")
	}
}

func (b *BattleTrainerImpl) HandleSwitch(action data.BattleAction) error {
	if action.Selected.PokemonUUID != b.ActivePokemon.PokemonUUID {
		return errors.New("selected pokemon should be active pokemon")
	}

	b.ActivePokemon, b.BattleParty = switchPokemon(b.ActivePokemon, action.Target, b.BattleParty)
	if action.Selected.PokemonUUID == b.ActivePokemon.PokemonUUID {
		return errors.New("active pokemon switch failed. Target pokemon either fainted or not found in party")
	}

	return nil
}

func (b *BattleTrainerImpl) HandleUseBag(action data.BattleAction) error {
	slog.Info("Healing pokemon", "pokemon", action.Target.PokemonUUID)
	var target *data.BattlePokemon
	if action.Target.PokemonUUID == b.ActivePokemon.PokemonUUID {
		target = b.ActivePokemon
	} else {
		for _, p := range b.BattleParty {
			if p.PokemonUUID == action.Target.PokemonUUID {
				target = p
				break
			}
		}
	}
	if healPokemon(target, action.Item) {
		b.UserManager.UseItem(action.Item)
	}
	return nil
}

func (b *BattleTrainerImpl) CalculateResult(opponent BattleTrainer) data.Result {
	result := data.Result{}
	if b.IsDefeated() {
		result.Status = data.Lost
		result.Money = data.GetMoneyLost(b.UserManager.GetUser())

		b.SendBattleLog("You lost the battle!")
		b.SendBattleLog(fmt.Sprintf("You lost $%v!", result.Money))

	} else {
		result.Status = data.Won
		result.Money = data.GetPrizeMoney(opponent.GetTrainerType(), opponent.GetTrainer().Party)
		result.BonusItems = opponent.GetRewards().Items

		b.SendBattleLog(fmt.Sprintf("%s has won the battle!", utils.ToCapitalizeFirstLetterOfEachWord(b.GetTrainer().Name)))
		b.SendBattleLog(fmt.Sprintf("You got $%v!", result.Money))

		// if gym battle; earn badge
		if opponent.GetTrainerType() == data.GymLeaderPrefix {
			result.BadgeEarned = opponent.GetRewards().Badge
			b.SendBattleLog(fmt.Sprintf("You earned a $%v!", result.BadgeEarned.Name))
		}
	}
	return result
}

func (b *BattleTrainerImpl) HandleBattleEnd(result data.Result) {
	b.handleBattleEnd(result)
}

func (b *BattleTrainerImpl) handleBattleEnd(result data.Result) {
	// stats update
	b.UserManager.StatUpdate(result)

	// evolve all pokemons
	for _, p := range append([]*data.BattlePokemon{b.ActivePokemon}, b.BattleParty...) {
		// if p.CanEvolve {
		slog.Info("Evolving pokemon", "pokemon", p.Pokemon.PokemonUUID, "Name", p.Pokemon.Name)
		b.PokemonService.Evolve(p.Pokemon)
		// }
	}
}
