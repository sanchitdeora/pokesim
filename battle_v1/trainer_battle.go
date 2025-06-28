//go:build exclude

package battle

import (
	"fmt"
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/errors"
	"github.com/sanchitdeora/PokeSim/pokemon"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/sanchitdeora/PokeSim/utils"
)

type TrainerBattleOpts struct {
	UserService    usermanagement.UserManager
	PokemonService pokemon.PokemonService
	BattleUser     BattleUser
	BattleOpponent BattleUser

	BattleInputChan <-chan *data.BattleAction
	BattleLogChan   chan<- string
}

type TrainerBattle struct {
	// UserTrainer
	UserTrainer             *data.User
	UserActivePokemon       *data.BattlePokemon
	UserInBattleParty       []*data.BattlePokemon
	UserUnfaintedPartyCount int

	// OpponentTrainer
	OpponentTrainer            *data.Trainer
	TrainerActivePokemon       *data.BattlePokemon
	TrainerInBattleParty       []*data.BattlePokemon
	TrainerPokemonFaced        map[*data.Pokemon][]*data.BattlePokemon
	TrainerUnfaintedPartyCount int
}

type TrainerBattleImpl struct {
	*TrainerBattleOpts
	*TrainerBattle
}

func CreateNewInBattlePokemon(pokemon data.Pokemon) *data.BattlePokemon {
	return &data.BattlePokemon{
		Pokemon:   &pokemon,
		BattleHP:  int(BattleHPCalculator(&pokemon.Stats.HP, pokemon.Level)),
		IsFainted: false,
	}
}

func NewTrainerBattle(opts *TrainerBattleOpts, trainer *data.Trainer) BattleSequence {
	// TODO: Add more validation checks
	if opts == nil || opts.BattleUser == nil || opts.BattleOpponent == nil {
		panic("missing input channel, unable to receive inputs.")
	}

	// prepare user
	var userActivePokemon *data.BattlePokemon
	userInBattlePokemonParty := make([]*data.BattlePokemon, 0, 6)

	// create trainer
	var trainerActivePokemon *data.BattlePokemon
	trainerInBattlePokemonParty := make([]*data.BattlePokemon, 0, 6)

	trainerPokemonFacedExp := make(map[*data.Pokemon][]*data.BattlePokemon, 0)
	for partyIndex, pokemon := range trainer.Party {
		if pokemon.PokemonUUID == "" {
			continue
		}
		inBattlePokemon := CreateNewInBattlePokemon(*pokemon)
		if partyIndex == 0 {
			trainerActivePokemon = inBattlePokemon
		} else if len(trainerInBattlePokemonParty) < 5 {
			trainerInBattlePokemonParty = append(trainerInBattlePokemonParty, inBattlePokemon)
		}
		trainerPokemonFacedExp[pokemon] = make([]*data.BattlePokemon, 0, 6)
	}

	// Get User
	user := opts.UserService.GetUser()

	for partyIndex, pokemon := range user.Party {
		if pokemon.PokemonUUID == "" {
			continue
		}
		inBattlePokemon := CreateNewInBattlePokemon(*pokemon)
		if partyIndex == 0 {
			userActivePokemon = inBattlePokemon
		} else if len(userInBattlePokemonParty) <= 6 {
			userInBattlePokemonParty = append(userInBattlePokemonParty, inBattlePokemon)
		}
	}

	return &TrainerBattleImpl{
		TrainerBattleOpts: opts,
		TrainerBattle: &TrainerBattle{
			UserTrainer:             user,
			UserActivePokemon:       userActivePokemon,
			UserInBattleParty:       userInBattlePokemonParty,
			UserUnfaintedPartyCount: len(userInBattlePokemonParty) + 1, // +1 for active pokemon

			OpponentTrainer:            trainer,
			TrainerActivePokemon:       trainerActivePokemon,
			TrainerInBattleParty:       trainerInBattlePokemonParty,
			TrainerPokemonFaced:        trainerPokemonFacedExp,
			TrainerUnfaintedPartyCount: len(trainerPokemonFacedExp) + 1, // +1 for active pokemon
		},
	}
}

// Initiate starts the battle sequence.
//
// The sequence is as follows:
// 1. Display the names of the starting pokemon.
// 2. Get the attack order from the user and the opponent.
// 3. Each pokemon takes turns attacking each other.
// 4. If the user or opponent runs out of unfainted pokemon, the battle is over.
// 5. Get the battle report and return it.
// 6. Update the user's information after the battle.
// 7. Evolve pokemon if necessary.
// 8. Close all channels.
func (tb *TrainerBattleImpl) Initiate() (*data.Result, error) {
	slog.Info("Battle Starts...")

	tb.BattleLog(fmt.Sprintf("%s chooses %s!", tb.getTrainerName(false), utils.ToCapitalizeFirstLetterOfEachWord(tb.TrainerActivePokemon.Pokemon.Name)))
	tb.BattleLog(fmt.Sprintf("%s, I choose you!\n", tb.getActivePokemonName(true)))
	tb.UpdateTrainerPokemonFaced(tb.UserActivePokemon)

	for {
		tb.BattleLog(fmt.Sprintf("%s Health: %v", tb.getActivePokemonName(true), tb.UserActivePokemon.BattleHP))
		tb.BattleLog(fmt.Sprintf("%s Health: %v\n", tb.getActivePokemonName(false), tb.TrainerActivePokemon.BattleHP))

		if tb.IsOver() {
			tb.BattleLog("Battle is completed!")
			break
		}

		battleInputs := getPokemonAttackOrder(tb.UserActivePokemon, tb.TrainerActivePokemon, tb.BattleInputChan)
		for _, input := range battleInputs {
			tb.Turn(input)
		}
	}

	report, err := tb.Report()
	if err != nil {
		slog.Error("error getting battle report", "error", err)
	}

	// Update User and save after battle completed
	tb.UserService.StatUpdate(*report)

	// Pokemon Evolutions
	tb.evolvePokemon()

	// close all channels
	close(tb.BattleLogChan)
	return report, err
}

// evolvePokemon iterates through the user's in battle party and checks if they can evolve.
// If so, it calls the PokemonService's Evolve method to evolve the Pokemon.
func (tb *TrainerBattleImpl) evolvePokemon() {
	for _, inBattlePokemon := range tb.UserInBattleParty {
		if inBattlePokemon.CanEvolve {
			tb.PokemonService.Evolve(inBattlePokemon.Pokemon)
		}
	}
}

func (tb *TrainerBattleImpl) Turn(userInput *data.BattleAction) {
	switch userInput.Type {
	case data.Switch:
		tb.BattleLog(fmt.Sprintf("%s is switching %s for %s", tb.getTrainerName(userInput.IsUser), userInput.Selected.Pokemon.Name, userInput.Target.Pokemon.Name))
		tb.SwitchPokemon(userInput.Target.Pokemon, userInput.IsUser)

	case data.Bag:
		if userInput.Item != nil && userInput.Item.Category == data.MedicalItems {
			tb.UseItem(userInput.Target, userInput.Item)
		}

	case data.Attack:
		tb.HandleAttack(userInput.Selected, userInput.Target, userInput.Move, userInput.IsUser)

	// TODO: turn should not be counted
	case data.Run:
		tb.Run()
	}
}

func (tb *TrainerBattleImpl) HandleAttack(attackPokemon *data.BattlePokemon, targetPokemon *data.BattlePokemon, attackMove *data.Moves, isUser bool) {
	tb.Attack(attackPokemon, targetPokemon, attackMove, isUser)

	if tb.UserActivePokemon.IsFainted {
		slog.Info(fmt.Sprintf("Before %s has fainted!", tb.getActivePokemonName(true)), "User Unfaint Count", tb.UserUnfaintedPartyCount)
		tb.BattleLog(fmt.Sprintf("%s has fainted!", tb.getActivePokemonName(true)))
		tb.UserUnfaintedPartyCount = tb.UserUnfaintedPartyCount - 1
		slog.Info(fmt.Sprintf("After %s has fainted!", tb.getActivePokemonName(true)), "User Unfaint Count", tb.UserUnfaintedPartyCount)

		if tb.UserUnfaintedPartyCount > 0 {
			// switch active pokemon with first in party and push fainted pokemon at end
			nextUnfaintedPokemonIndex := 0
			for i, pokemon := range tb.TrainerInBattleParty {
				if !pokemon.IsFainted {
					nextUnfaintedPokemonIndex = i
					break
				}
			}
			switchPokemonWithIndex(nextUnfaintedPokemonIndex, tb.UserActivePokemon, &tb.UserInBattleParty)
			// tb.SwitchPokemonWithIndex(nextUnfaintedPokemonIndex, true)
		}
	}

	if tb.TrainerActivePokemon.IsFainted {
		tb.HandleTargetPokemonFaint()
	}
}

func (tb *TrainerBattleImpl) HandleTargetPokemonFaint() {
	slog.Info(fmt.Sprintf("Before %s has fainted!", tb.getActivePokemonName(false)), "Opponent Unfaint Count", tb.TrainerUnfaintedPartyCount)

	tb.BattleLog(fmt.Sprintf("%s has fainted!", tb.getActivePokemonName(false)))
	tb.TrainerUnfaintedPartyCount = tb.TrainerUnfaintedPartyCount - 1
	slog.Info(fmt.Sprintf("After %s has fainted!", tb.getActivePokemonName(false)), "Opponent Unfaint Count", tb.TrainerUnfaintedPartyCount)

	tb.UpdateInvolvedUserPokemon(tb.TrainerActivePokemon.Pokemon, tb.TrainerPokemonFaced[tb.TrainerActivePokemon.Pokemon])

	if tb.TrainerUnfaintedPartyCount > 0 {
		// switch active pokemon with first in party and push fainted pokemon at end
		nextUnfaintedPokemonIndex := 0
		for i, pokemon := range tb.TrainerInBattleParty {
			if !pokemon.IsFainted {
				nextUnfaintedPokemonIndex = i
				break
			}
		}
		switchPokemonWithIndex(nextUnfaintedPokemonIndex, tb.TrainerActivePokemon, &tb.TrainerInBattleParty)
		// tb.SwitchPokemonWithIndex(nextUnfaintedPokemonIndex, false)
	}
}

func (tb *TrainerBattleImpl) Attack(attackPokemon *data.BattlePokemon, targetPokemon *data.BattlePokemon, attackMove *data.Moves, isUser bool) {

	tb.BattleLog(fmt.Sprintf("%s used %s", tb.getActivePokemonName(isUser), utils.ToCapitalizeFirstLetterOfEachWord(attackMove.Name)))

	damagePoints := calculateAttackDamage(attackPokemon, targetPokemon, attackMove, 1.0)
	// slog.Info("====DEBUG====attack info", "is it user?", isUser, "attackPokemon", attackPokemon, "move", attackMove, "damagePoints", damagePoints)

	if damagePoints >= targetPokemon.BattleHP {
		targetPokemon.BattleHP = 0
		targetPokemon.IsFainted = true
	} else {
		targetPokemon.BattleHP -= damagePoints
	}

	tb.BattleLog(fmt.Sprintf("%s did %v points of damage to %s", tb.getActivePokemonName(isUser), damagePoints, tb.getActivePokemonName(!isUser)))
}

func (tb *TrainerBattleImpl) SwitchPokemon(switchingPokemon *data.Pokemon, isUser bool) {
	var activePokemon *data.BattlePokemon
	var inBattleParty []*data.BattlePokemon

	if isUser {
		activePokemon = tb.UserActivePokemon
		inBattleParty = tb.UserInBattleParty
	} else {
		activePokemon = tb.TrainerActivePokemon
		inBattleParty = tb.TrainerInBattleParty
	}

	for i, pokemon := range inBattleParty {
		if pokemon.Pokemon.PokemonUUID == switchingPokemon.PokemonUUID && !pokemon.IsFainted {
			currentPokemon := activePokemon
			nextPokemon := inBattleParty[i]
			// tb.UserActivePokemon = nil

			// only get active pokemon if unfainted pokemon available; else add to the list
			if !nextPokemon.IsFainted {
				(inBattleParty) = append(inBattleParty[:i], inBattleParty[i+1:]...)
			}
			inBattleParty = append(inBattleParty, currentPokemon)

			break
		}
	}

	tb.BattleLog(fmt.Sprintf("%s is switching their pokemon!", tb.getTrainerName(isUser)))

	if isUser {
		tb.BattleLog(fmt.Sprintf("%s, I choose you!", tb.getActivePokemonName(true)))
	} else {
		tb.BattleLog(fmt.Sprintf("%s, chooses %s!", tb.getTrainerName(false), utils.ToCapitalizeFirstLetterOfEachWord(tb.TrainerActivePokemon.Pokemon.Name)))
	}

	// update map of pokemon facing
	tb.UpdateTrainerPokemonFaced(tb.UserActivePokemon)
}

func (tb *TrainerBattleImpl) UseItem(targetPokemon *data.BattlePokemon, item *data.Item) {
	switch item.Category {
	case data.MedicalItems:
		healPokemon(targetPokemon, item)
	case data.PokeBalls:
		tb.CatchPokemon(targetPokemon, item)
	}
}

// TODO: turn should not be counted
func (tb *TrainerBattleImpl) CatchPokemon(targetPokemon *data.BattlePokemon, item *data.Item) {
	tb.BattleLog("Cannot catch a pokemon in trainer battle")
}

func (tb *TrainerBattleImpl) Run() {
	tb.BattleLog("No! There's no running from a trainer battle!")
}

func (tb *TrainerBattleImpl) IsOver() bool {
	if tb.UserUnfaintedPartyCount == 0 || tb.TrainerUnfaintedPartyCount == 0 {
		return true
	}
	return false
}

func (tb *TrainerBattleImpl) Report() (*data.Result, error) {
	if !tb.IsOver() {
		return nil, errors.ErrBattleOngoing
	}

	var report data.Result
	if tb.UserUnfaintedPartyCount == 0 {
		report.Money = data.GetMoneyLost(tb.UserTrainer)

		tb.BattleLog(fmt.Sprintf("You lost %v ₽!", report.Money))
		tb.BattleLog(fmt.Sprintf("You lost the battle to %s!", tb.OpponentTrainer.Name))
		report.UserWin = false

	} else {
		report.UserWin = true
		report.Money = data.GetPrizeMoney(tb.OpponentTrainer.Type, tb.OpponentTrainer.Party)
		report.BonusItems = tb.OpponentTrainer.Rewards.Items

		tb.BattleLog(fmt.Sprintf("%s has won the battle!", tb.UserTrainer.Name))
		tb.BattleLog(fmt.Sprintf("You got %v ₽!", report.Money))

		// if gym battle; earn badge
		if tb.OpponentTrainer.Type == data.GymLeaderPrefix {
			report.BadgeEarned = tb.OpponentTrainer.Rewards.Badge
			tb.BattleLog(fmt.Sprintf("You earned a %v ₽!", tb.OpponentTrainer.Rewards.Badge.Name))
		}
	}

	return &report, nil
}

func (tb *TrainerBattleImpl) UpdateInvolvedUserPokemon(faintedPokemon *data.Pokemon, pokemonFaced []*data.BattlePokemon) {
	for _, inBattlePokemon := range pokemonFaced {
		if !inBattlePokemon.IsFainted {
			expGain := calculateExperienceGained(faintedPokemon.Level, faintedPokemon.BasePokemon.BaseExperience, inBattlePokemon.Pokemon.Level)
			tb.BattleLog(fmt.Sprintf("%s gained %v experience points", utils.ToCapitalizeFirstLetterOfEachWord(inBattlePokemon.Pokemon.Name), expGain))

			tb.PokemonService.ExperienceGain(inBattlePokemon.Pokemon, *inBattlePokemon.Pokemon)
		}
	}
}

func (tb *TrainerBattleImpl) UpdateTrainerPokemonFaced(pokemon *data.BattlePokemon) {
	tb.TrainerPokemonFaced[tb.TrainerActivePokemon.Pokemon] =
		append(tb.TrainerPokemonFaced[tb.TrainerActivePokemon.Pokemon], pokemon)
}

func (tb *TrainerBattleImpl) GetUserActive(isUser bool) *data.BattlePokemon {
	if isUser {
		return tb.UserActivePokemon
	} else {
		return tb.TrainerActivePokemon
	}
}

func (tb *TrainerBattleImpl) getTrainerName(isUser bool) string {
	if isUser {
		return tb.UserTrainer.Name
	} else {
		return fmt.Sprintf("%s %s", tb.OpponentTrainer.Type, tb.OpponentTrainer.Name)
	}
}

func (tb *TrainerBattleImpl) getActivePokemonName(isUser bool) string {
	if isUser {
		return utils.ToCapitalizeFirstLetterOfEachWord(tb.UserActivePokemon.Pokemon.Name)
	} else {
		return fmt.Sprintf("The opposing %s", utils.ToCapitalizeFirstLetterOfEachWord(tb.TrainerActivePokemon.Pokemon.Name))
	}
}

func (tb *TrainerBattleImpl) BattleLog(log string) {
	// slog.Info("sending battle log to container...", "log", log, "to channel", tb.BattleLogChan)
	tb.BattleLogChan <- log
}
