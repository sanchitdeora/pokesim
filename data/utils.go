package data

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand"
)

func CalculateAttackDamage(attackPokemon *BattlePokemon, targetPokemon *BattlePokemon, attackMove *Moves, battleTypeAttackCoeff float64) (damage int, moveEffect TypeEffective, isCritHit bool) {
	var totalDamage float64

	var attackStat float64 = 0
	var defenseStat float64 = 0

	if attackMove.DamageClass == Physical {
		attackStat = float64(attackPokemon.Pokemon.Stats.Attack.Value)
		defenseStat = float64(attackPokemon.Pokemon.Stats.Defense.Value)
	} else if attackMove.DamageClass == Special {
		attackStat = float64(attackPokemon.Pokemon.Stats.SpecialAttack.Value)
		defenseStat = float64(attackPokemon.Pokemon.Stats.SpecialDefense.Value)
	} else {
		slog.Warn("Move Damage class not supported", "move damage class", attackMove.DamageClass)
	}

	pokemonTypes := []PokemonTypeName{targetPokemon.Pokemon.BasePokemon.Type1}

	if targetPokemon.Pokemon.BasePokemon.Type2 != "" {
		pokemonTypes = append(pokemonTypes, targetPokemon.Pokemon.BasePokemon.Type2)
	}
	moveEffect = GetMoveEffect(attackMove.Type, pokemonTypes...)

	slog.Debug(fmt.Sprintf("Calculating damage: ( ( (((2 * {%v})/5) + 2) * {%v} * ({%v} / {%v}) ) / 50 ) + 2", attackPokemon.Pokemon.Level, attackMove.Power, attackStat, defenseStat))

	totalDamage = (((((2.0 * float64(attackPokemon.Pokemon.Level)) / 5.0) + 2.0) * float64(attackMove.Power) * (attackStat / defenseStat)) / 50) + 2

	// for battles with more than 1 enemy, coefficient = 0.75
	slog.Debug(fmt.Sprintf("Battle Attack Coeff: {%f} * {%f}", totalDamage, battleTypeAttackCoeff))
	totalDamage *= battleTypeAttackCoeff

	isCritHit = calculateCritHit()

	if isCritHit {
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

	slog.Debug(fmt.Sprintf("Move Effect Coeff: {%f} * {%f} == {%f}", totalDamage, moveEffect, (totalDamage * float64(moveEffect))))
	damage = int(math.Round(totalDamage * float64(moveEffect)))

	return
}


func calculateCritHit() bool {
	return randomGenerator(0, 1) < (1.0 / 24.0)
}

func attackCoefficient() float64 {
	return randomGenerator(0.85, 1)
}

func isStab(attackMove *Moves, attackPokemon *BattlePokemon) bool {
	return attackMove.Type == attackPokemon.Pokemon.BasePokemon.Type1 || attackMove.Type == attackPokemon.Pokemon.BasePokemon.Type2
}

func randomGenerator(min float64, max float64) float64 {
	randIndex := rand.Float64()
	return (min + randIndex*(max-min))
}
