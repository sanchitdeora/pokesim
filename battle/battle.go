package battle

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/sanchitdeora/PokeSim/data"
)

type BattleIFace interface {
	InitiateBattleSequence() (*data.BattleReport, error)
	Attack(attackPokemon *data.InBattlePokemon, targetPokemon *data.InBattlePokemon, attackMove *data.Moves, isUser bool)
	CatchPokemon(targetPokemon *data.InBattlePokemon, item *data.Item)
	Run()
	IsBattleOver() bool
	BattleReport() (*data.BattleReport, error)
	GetUserActivePokemon() *data.InBattlePokemon
	GetOpponentActivePokemon() *data.InBattlePokemon
}

func getPokemonAttackOrder(userActivePokemon *data.InBattlePokemon, trainerActivePokemon *data.InBattlePokemon, userInputChan <-chan *data.BattleInput) (inputs []*data.BattleInput) {
	// func (tb *TrainerBattleImpl) GetPokemonAttackOrder() (inputs []*data.BattleInput) {
	userInput := waitForUserInput(userInputChan)
	trainerInput := waitForInput(trainerActivePokemon, userActivePokemon, false)

	userSpeed := battleStatCalculator(&userActivePokemon.Pokemon.Stats.Speed, userActivePokemon.Pokemon.Level)
	targetSpeed := battleStatCalculator(&trainerActivePokemon.Pokemon.Stats.Speed, trainerActivePokemon.Pokemon.Level)

	// if any pokemon's move is higher priority, it will go first. If equal, check speed
	if userInput.Move.Priority != trainerInput.Move.Priority {
		if userInput.Move.Priority > trainerInput.Move.Priority {
			return append(inputs, userInput, trainerInput)
		} else {
			return append(inputs, trainerInput, userInput)
		}
	}

	// if user's active pokemon speed >= opposing pokemon's, then user goes first.
	if userSpeed >= targetSpeed {
		return append(inputs, userInput, trainerInput)
	} else {
		return append(inputs, trainerInput, userInput)
	}
}

func switchPokemonWithIndex(switchingPokemonIndex int, activePokemon *data.InBattlePokemon, inBattleParty *[]*data.InBattlePokemon) *data.InBattlePokemon {
	currentPokemon := activePokemon
	nextPokemon := (*inBattleParty)[switchingPokemonIndex]
	// tb.UserActivePokemon = nil

	// only get active pokemon if unfainted pokemon available; else add to the list
	if !nextPokemon.IsFainted {
		(*inBattleParty) = append((*inBattleParty)[:switchingPokemonIndex], (*inBattleParty)[switchingPokemonIndex+1:]...)
	}
	(*inBattleParty) = append((*inBattleParty), currentPokemon)

	return nextPokemon
}

func healPokemon(targetPokemon *data.InBattlePokemon, item *data.Item) {
	targetPokemon.BattleHP += targetPokemon.BattleHP + item.Attributes

	if targetPokemon.BattleHP > int(battleHPCalculator(&targetPokemon.Pokemon.Stats.HP, targetPokemon.Pokemon.Level)) {
		targetPokemon.BattleHP = int(battleHPCalculator(&targetPokemon.Pokemon.Stats.HP, targetPokemon.Pokemon.Level))
	}
}

func calculateAttackDamage(attackPokemon *data.InBattlePokemon, targetPokemon *data.InBattlePokemon, attackMove *data.Moves, battleTypeAttackCoeff float64) int {
	var totalDamage float64

	var attackStat float64 = 0
	var defenseStat float64 = 0

	if attackMove.DamageClass == data.Physical {
		attackStat = battleStatCalculator(&attackPokemon.Pokemon.Stats.Attack, attackPokemon.Pokemon.Level) 
		defenseStat = battleStatCalculator(&targetPokemon.Pokemon.Stats.Defense, targetPokemon.Pokemon.Level)

	} else if attackMove.DamageClass == data.Special {
		attackStat = battleStatCalculator(&attackPokemon.Pokemon.Stats.SpecialAttack, attackPokemon.Pokemon.Level) 
		defenseStat = battleStatCalculator(&targetPokemon.Pokemon.Stats.SpecialDefense, targetPokemon.Pokemon.Level)
	} else {
		slog.Warn("Move Damage class not supported", "move damage class", attackMove.DamageClass)
	}

	pokemonTypes := []data.PokemonTypeName{targetPokemon.Pokemon.BasePokemon.Type1}

	if targetPokemon.Pokemon.BasePokemon.Type2 != "" {
		pokemonTypes = append(pokemonTypes, targetPokemon.Pokemon.BasePokemon.Type2)
	}
	moveEffect := data.GetMoveEffect(attackMove.Type, pokemonTypes...)

	slog.Debug(fmt.Sprintf("Calculating damage: ( ( (((2 * {%v})/5) + 2) * {%v} * ({%v} / {%v}) ) / 50 ) + 2", attackPokemon.Pokemon.Level, attackMove.Power, attackStat, defenseStat))

	totalDamage = (((((2.0 * float64(attackPokemon.Pokemon.Level)) / 5.0) + 2.0) * float64(attackMove.Power) * (attackStat / defenseStat)) / 50) + 2

	// for battles with more than 1 enemy, coefficient = 0.75
	slog.Debug(fmt.Sprintf("Battle Attack Coeff: {%f} * {%f}", totalDamage, battleTypeAttackCoeff))
	totalDamage *= battleTypeAttackCoeff

	if isCritHit() {
		slog.Info("Critical Hit!")
		slog.Debug(fmt.Sprintf("Critical hit: {%f} * 1.5", totalDamage))
		totalDamage *= 1.5
	}

	// random attack power coeffecient ranging from (0.85 - 1.00)
	slog.Debug(fmt.Sprintf("Attack Coeff 0.85 -- 1.00: {%f} * {%f}", totalDamage, attackCoefficient()))
	totalDamage *= attackCoefficient()

	if isStab(attackMove, attackPokemon) {
		slog.Debug(fmt.Sprintf("STAB: {%f} * 1.5", totalDamage))
		totalDamage *= 1.5
	}

	if moveEffect != data.NOR {
		if moveEffect == data.MNE || moveEffect == data.NVR {
			slog.Info("Not very effective!", "effect", moveEffect)
		} else if moveEffect == data.SUP || moveEffect == data.HYP {
			slog.Info("Super effective!", "effect", moveEffect)
		} else {
			slog.Info("This move has No effect to the target!", "effect", moveEffect)
		}
	}

	slog.Debug(fmt.Sprintf("Move Effect Coeff: {%f} * {%f} == {%f}", totalDamage, moveEffect, (totalDamage * float64(moveEffect))))
	return int(math.Round(totalDamage * float64(moveEffect)))
}

func battleStatCalculator(stat *data.PokemonStat, level int) float64 {
	return (((float64(2 * stat.Value) + float64(stat.IV) + float64(stat.EV / 4)) * float64(level)) / 100) + 5
}

func battleHPCalculator(HP *data.PokemonStat, level int) float64 {
	return (((float64(2 * HP.Value) + float64(HP.IV) + float64(HP.EV / 4)) * float64(level)) / 100) + 10
}
