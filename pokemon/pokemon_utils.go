package pokemon

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/sanchitdeora/PokeSim/data"
)

const NatureCoeff = 1.0

// hp calculation is different from the other stats
func calculateHPStatUpgrade(baseHP int, iv int, ev int, level int) int {
	return int(math.Floor(((((2 * float64(baseHP)) + float64(iv) + (float64(ev) / 4)) * float64(level)) / 100) + float64(level) + 10))
}

// all other stats calculation
func calculateOtherStatUpgrade(baseStat int, iv int, ev int, level int) int {
	return int(math.Floor(((((2*float64(baseStat))+float64(iv)+(float64(ev)/4))*float64(level))/100)+5) * NatureCoeff)
}

func nextLevelErraticExp(level int) int {
	if level < 50 {
		return int(math.Round(((math.Pow(float64(level), 3)) * (100 - float64(level))) / 50))
	} else if level < 68 {
		return int(math.Round(((math.Pow(float64(level), 3)) * (100 - float64(level))) / 50))
	} else if level < 98 {
		return int(math.Round(((math.Pow(float64(level), 3)) * math.Floor((1911-(10*float64(level)))/3)) / 500))
	} else {
		return int(math.Round(((math.Pow(float64(level), 3)) * (100 - float64(level))) / 50))
	}
}

func nextLevelFastExp(level int) int {
	return ((int(math.Pow(float64(level), 3)) * 4) / 5)
}

func nextLevelMediumFastExp(level int) int {
	return int(math.Pow(float64(level), 3))
}

func nextLevelMediumSlowExp(level int) int {
	return int(math.Round((6*math.Pow(float64(level), 3))/5) - (15 * math.Pow(float64(level), 2)) + (100 * float64(level)) - 140)
}

func nextLevelSlowExp(level int) int {
	return ((int(math.Pow(float64(level), 3)) * 5) / 4)
}

func nextLevelFluctuatingExp(level int) int {
	if level < 15 {
		return int(math.Round(math.Pow(float64(level), 3) * ((math.Floor((float64(level)+1)/3) + 24) / 50)))
	} else if level < 36 {
		return int(math.Round((math.Pow(float64(level), 3) * (float64(level) + 14)) / 50))
	} else {
		return int(math.Round((math.Pow(float64(level), 3) * (math.Floor(float64(level)/2) + 32)) / 50))
	}
}

func calculateExperienceGained(level int, baseExp int, winningPokemon *data.Pokemon) int {
	var totalExp float64
	totalExp = float64(baseExp*level) / 5.0
	totalExp *= math.Pow((float64((2*level)+10)/float64(level+winningPokemon.Level+10)), 2.5) + 1.0

	if canPokemonEvolve(winningPokemon) {
		totalExp *= 1.2
	}

	return int(math.Round(totalExp))
}

// generate starter pokemon

func generatePokemonIVs() int {
	randIndex := rand.Float64()
	return int(math.Round(randIndex * (31)))
}


func generatePokemonHPStat(value, level int) data.PokemonStat {
	iv := generatePokemonIVs()
	return data.PokemonStat{
		Value: calculateHPStatUpgrade(value, iv, 0, level),
		IV:    iv,
		EV:    0,
	}
}

func generatePokemonOtherStat(value, level int) data.PokemonStat {
	iv := generatePokemonIVs()
	return data.PokemonStat{
		Value: calculateOtherStatUpgrade(value, iv, 0, level),
		IV:    iv,
		EV:    0,
	}
}

func setupMoveset(basePokemon data.BasePokemon, level int) data.Moveset {
	var moveset data.Moveset

	for i := range level {
		fmt.Printf("levelIndex:%v, move:%v\n", i, basePokemon.MovesLearned[i])
		move := basePokemon.MovesLearned[i]
		switch i % 4 {
		case 0:
			moveset.Move1 = &move
		case 1:
			moveset.Move2 = &move
		case 2:
			moveset.Move3 = &move
		case 3:
			moveset.Move4 = &move
		}
	}
	return moveset
}
