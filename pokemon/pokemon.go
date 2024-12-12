package pokemon

import (
	"fmt"
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
)

type PokemonManager interface {
	// LevelUp(pokemon *data.Pokemon)
	Evolve(pokemon *data.Pokemon)
	LearnNewMoves(move *data.Moves)
	ExperienceGain(expGain int, pokemon *data.Pokemon) bool
}

type PokemonOpts struct{}

type PokemonImpl struct {
	opts PokemonOpts
}

func NewPokemonManager(opts PokemonOpts) PokemonManager {
	return &PokemonImpl{
		opts: opts,
	}
}

func (p *PokemonImpl) LevelUp(pokemon *data.Pokemon) {
	pokemon.Level++

	// calculate stats upgrade
	statUpgrades(pokemon)

	// calculate pokemon experience left
	pokemon.ExperienceLeft = getExperienceRequiredForNextLevel(pokemon)

	// TODO: get user input if pokemon learns moves; if yes which moveset?
	// if move, exists := pokemon.MovesLearnedByLevel[pokemon.Level]; exists {
	// 	p.LearnNewMoves(&move)
	// }
}

// TODO: add should evolve method. Evolve after battle.
func (p *PokemonImpl) Evolve(pokemon *data.Pokemon) {
	if !canPokemonEvolve(pokemon) {
		return
	}

	evolvedBasePokemonPath := pokemon.EvolutionChain[pokemon.Level]
	if len(evolvedBasePokemonPath) > 1 {
		//TODO: add option to choose which pokemon to evolve to
		panic("implemenet multiple pokemon evolution")
	} else if len(evolvedBasePokemonPath) == 0 {
		slog.Error("pokemon cannot evolve", "pokemon", pokemon.Name)
	} else {
		evolvedBasePokemon, err := getBasePokemonFromPath(evolvedBasePokemonPath[0])
		if err != nil {
			slog.Error("pokemon cannot evolve", "pokemon", pokemon.Name, "error", err)
		}
		pokemon.BasePokemon = *evolvedBasePokemon
	}

	statUpgrades(pokemon)
}

func (p *PokemonImpl) LearnNewMoves(move *data.Moves) {
	// input which move should be replaced.
}

func (p *PokemonImpl) EVGain(expGain int, pokemon *data.Pokemon, evYield *data.BasePokemonStats) {
    // Calculate total EVs and determine how many to add
	stats := pokemon.Stats
	userPokemonEV := stats.HP.EV + stats.Attack.EV + stats.Defense.EV + stats.SpecialAttack.EV + stats.SpecialDefense.EV + stats.Speed.EV
	if userPokemonEV >= 510 {
		slog.Info("Pokemon is fully trained !!!", "total pokemon ev", userPokemonEV)
		return
	}

	// Calculate remaining EV space
	targetEVYield := evYield.HP + evYield.Attack + evYield.Defense + evYield.SpecialAttack + evYield.SpecialDefense + evYield.Speed
	evPointsToAdd := targetEVYield

	if userPokemonEV + targetEVYield > 510 {
		evPointsToAdd = 510 - userPokemonEV
	}

	if addPokemonEVToStats(&evPointsToAdd, &stats.HP, evYield.HP) {
		slog.Info("Adding EV to Pokemon", "pokemonStat", "HP")
		return
	}

	if addPokemonEVToStats(&evPointsToAdd, &stats.Attack, evYield.Attack) {
		slog.Info("Adding EV to Pokemon", "pokemonStat", "Attack")
		return
	}

	if addPokemonEVToStats(&evPointsToAdd, &stats.Defense, evYield.Defense) {
		slog.Info("Adding EV to Pokemon", "pokemonStat", "Defence")
		return
	}

	if addPokemonEVToStats(&evPointsToAdd, &stats.SpecialAttack, evYield.SpecialAttack) {
		slog.Info("Adding EV to Pokemon", "pokemonStat", "Special Attack")
		return
	}

	if addPokemonEVToStats(&evPointsToAdd, &stats.SpecialDefense, evYield.SpecialDefense) {
		slog.Info("Adding EV to Pokemon", "pokemonStat", "Special Defense")
		return
	}

	if addPokemonEVToStats(&evPointsToAdd, &stats.Speed, evYield.Speed) {
		slog.Info("Adding EV to Pokemon", "pokemonStat", "Speed")
		return
	}
}

func (p *PokemonImpl) ExperienceGain(expGain int, pokemon *data.Pokemon) bool {
	slog.Info("Pokemon Gains Experience", "expGain", expGain, "Pokemon", pokemon.Name)
	var canEvolve bool

	for {
		if pokemon.ExperienceLeft > expGain {
			pokemon.ExperienceLeft -= expGain
			break
		}
		expGain -= pokemon.ExperienceLeft
		p.LevelUp(pokemon)

		if canPokemonEvolve(pokemon) {
			canEvolve = true
		}
	}
	return canEvolve
}

func statUpgrades(pokemon *data.Pokemon) {
	hp := calculateHPStatUpgrade(pokemon.BaseStats.HP.Value, &pokemon.Stats.HP, pokemon.Level)

	attack := calculateOtherStatUpgrade(pokemon.BaseStats.Attack.Value, &pokemon.Stats.Attack, pokemon.Level)
	defense := calculateOtherStatUpgrade(pokemon.BaseStats.Defense.Value, &pokemon.Stats.Defense, pokemon.Level)

	spAttack := calculateOtherStatUpgrade(pokemon.BaseStats.SpecialAttack.Value, &pokemon.Stats.SpecialAttack, pokemon.Level)
	spDefense := calculateOtherStatUpgrade(pokemon.BaseStats.SpecialDefense.Value, &pokemon.Stats.SpecialDefense, pokemon.Level)

	speed := calculateOtherStatUpgrade(pokemon.BaseStats.Speed.Value, &pokemon.Stats.Speed, pokemon.Level)

	slog.Info("Stat upgrade for pokemon:")
	slog.Info(fmt.Sprintf("HP: +%v", hp-pokemon.Stats.HP.Value))
	slog.Info(fmt.Sprintf("Attack: +%v", attack-pokemon.Stats.Attack.Value))
	slog.Info(fmt.Sprintf("Defense: +%v", defense-pokemon.Stats.Defense.Value))
	slog.Info(fmt.Sprintf("Special Attack: +%v", spAttack-pokemon.Stats.SpecialAttack.Value))
	slog.Info(fmt.Sprintf("Special Defence: +%v", spDefense-pokemon.Stats.SpecialDefense.Value))
	slog.Info(fmt.Sprintf("Speed: +%v", speed-pokemon.Stats.Speed.Value))

	pokemon.Stats.HP.Value = hp
	pokemon.Stats.Attack.Value = attack
	pokemon.Stats.Defense.Value = defense
	pokemon.Stats.SpecialAttack.Value = spAttack
	pokemon.Stats.SpecialDefense.Value = spDefense
	pokemon.Stats.Speed.Value = speed
}

// calculate next level exp required
func getExperienceRequiredForNextLevel(pokemon *data.Pokemon) int {
	switch pokemon.GrowthRate {
	case data.Erratic:
		return nextLevelErraticExp(pokemon.Level)
	case data.Fast:
		return nextLevelFastExp(pokemon.Level)
	case data.MediumFast:
		return nextLevelMediumFastExp(pokemon.Level)
	case data.MediumSlow:
		return nextLevelMediumSlowExp(pokemon.Level)
	case data.Slow:
		return nextLevelSlowExp(pokemon.Level)
	case data.Fluctuating:
		return nextLevelFluctuatingExp(pokemon.Level)
	default:
		slog.Error("invalid growth rate type found, defaulting to MediumFast", "growth rate type", pokemon.GrowthRate)
		return nextLevelMediumFastExp(pokemon.Level)
	}
}

func canPokemonEvolve(pokemon *data.Pokemon) bool {
	_, exists := pokemon.EvolutionChain[pokemon.Level]
	return exists
}

func addPokemonEVToStats(pointsToAdd *int, pokemonStat *data.PokemonStat, evYield int) bool {
	if evYield == 0 {
		return false
	}

	if *pointsToAdd >= evYield {
		slog.Info("Successfully added EV points", "points added", evYield)

		*pointsToAdd -= evYield
		pokemonStat.EV = evYield

		return false

	} else {
		slog.Info("Successfully added EV points", "points added", *pointsToAdd)

		pokemonStat.EV = *pointsToAdd

		return true
	}
}

func getBasePokemonFromPath(path data.BasePokemonId) (*data.BasePokemon, error) {
	pokemon, err := utils.ReadJsonFromFile[data.BasePokemon](string(path))
	if err != nil {
		return nil, err
	}
	return &pokemon, nil
}
