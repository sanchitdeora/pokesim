package pokemon

import (
	"fmt"
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/logger"
	"github.com/sanchitdeora/PokeSim/utils"
)

//go:generate mockgen -build_flags=--mod=mod -destination=mocks/mock_pokemon_manager.go -package=mock_pokemon_manager github.com/sanchitdeora/PokeSim/pokemon PokemonManager
type PokemonManager interface {
	// LevelUp(pokemon *data.Pokemon)
	Evolve(pokemon *data.Pokemon)
	LearnNewMoves(move *data.Moves)
	ExperienceGain(pokemon *data.Pokemon, faintedPokemon data.Pokemon)
	EvGain(pokemon *data.Pokemon, evYieldToAdd data.PokemonStats)
}

type PokemonOpts struct {
	Logger logger.Logger
}

type PokemonImpl struct {
	opts PokemonOpts
}

func NewPokemonManager(opts PokemonOpts) PokemonManager {
	if opts.Logger == nil {
		opts.Logger = logger.NewDefaultLogger()
	}
	return &PokemonImpl{opts: opts}
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

func (p *PokemonImpl) Evolve(pokemon *data.Pokemon) {
	if !canPokemonEvolve(pokemon) {
		slog.Debug("pokemon cannot evolve", "pokemon", pokemon.Name, "pokemonUUID", pokemon.PokemonUUID)
		return
	}

	evolvedBasePokemonPath := pokemon.EvolutionChain[pokemon.Level]
	if len(evolvedBasePokemonPath) > 1 {
		//TODO: add option to choose which pokemon to evolve to
		panic("implemenet multiple pokemon evolution")

	} else if len(evolvedBasePokemonPath) == 0 {
		slog.Error("pokemon cannot evolve", "pokemon", pokemon.Name)
	} else {
		evolvedBasePokemon, err := getBasePokemonByID(evolvedBasePokemonPath[0])
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

func (p *PokemonImpl) EvGain(pokemon *data.Pokemon, evYieldToAdd data.PokemonStats) {
	// Calculate total EVs and determine how many to add
	stats := pokemon.Stats
	userPokemonEv := stats.HP.EV + stats.Attack.EV + stats.Defense.EV + stats.SpecialAttack.EV + stats.SpecialDefense.EV + stats.Speed.EV
	if userPokemonEv >= 510 {
		slog.Info("Pokemon is fully trained !!!", "total pokemon ev", userPokemonEv)
		return
	}

	// Calculate remaining EV space
	targetYield := evYieldToAdd.HP.EV + evYieldToAdd.Attack.EV + evYieldToAdd.Defense.EV + evYieldToAdd.SpecialAttack.EV + evYieldToAdd.SpecialDefense.EV + evYieldToAdd.Speed.EV
	evPointsToAdd := targetYield

	if userPokemonEv+targetYield > 510 {
		evPointsToAdd = 510 - userPokemonEv
	}

	if addPokemonEvToStats(&evPointsToAdd, &stats.HP, evYieldToAdd.HP.EV) {
		slog.Debug("Adding EV to Pokemon", "pokemonStat", "HP")
		return
	}

	if addPokemonEvToStats(&evPointsToAdd, &stats.Attack, evYieldToAdd.Attack.EV) {
		slog.Debug("Adding EV to Pokemon", "pokemonStat", "Attack")
		return
	}

	if addPokemonEvToStats(&evPointsToAdd, &stats.Defense, evYieldToAdd.Defense.EV) {
		slog.Debug("Adding EV to Pokemon", "pokemonStat", "Defence")
		return
	}

	if addPokemonEvToStats(&evPointsToAdd, &stats.SpecialAttack, evYieldToAdd.SpecialAttack.EV) {
		slog.Debug("Adding EV to Pokemon", "pokemonStat", "Special Attack")
		return
	}

	if addPokemonEvToStats(&evPointsToAdd, &stats.SpecialDefense, evYieldToAdd.SpecialDefense.EV) {
		slog.Debug("Adding EV to Pokemon", "pokemonStat", "Special Defense")
		return
	}

	if addPokemonEvToStats(&evPointsToAdd, &stats.Speed, evYieldToAdd.Speed.EV) {
		slog.Debug("Adding EV to Pokemon", "pokemonStat", "Speed")
		return
	}
}

func (p *PokemonImpl) ExperienceGain(pokemon *data.Pokemon, faintedPokemon data.Pokemon) {
	expGain := calculateExperienceGained(faintedPokemon.Level, faintedPokemon.BaseExperience, pokemon)

	slog.Info("Pokemon Gains Experience", "expGain", expGain, "Pokemon", pokemon.Name)
	// tb.BattleLog(fmt.Sprintf("%s gained %v experience points", utils.ToCapitalizeFirstLetterOfEachWord(inBattlePokemon.Pokemon.Name), expGain))

	for {
		if pokemon.ExperienceLeft > expGain {
			pokemon.ExperienceLeft -= expGain
			break
		}
		expGain -= pokemon.ExperienceLeft
		p.LevelUp(pokemon)
	}
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

// TODO: It is possible that pokemon has evolved by more than a level and surpassed the expected evolution level.
// Also check lower levels to confirm if it can evolve.
func canPokemonEvolve(pokemon *data.Pokemon) bool {
	_, exists := pokemon.EvolutionChain[pokemon.Level]
	return exists
}

func addPokemonEvToStats(pointsToAdd *int, pokemonStat *data.PokemonStat, evYield int) bool {
	if evYield == 0 {
		return false
	}

	if *pointsToAdd >= evYield {
		slog.Info("Successfully added EV points", "points added", evYield)

		*pointsToAdd -= evYield
		pokemonStat.EV = evYield

		return true

	} else {
		slog.Info("Successfully added EV points", "points added", *pointsToAdd)

		pokemonStat.EV = *pointsToAdd

		return true
	}
}

func getBasePokemonByID(ID data.BasePokemonId) (*data.BasePokemon, error) {
	path := fmt.Sprintf("/assets/pokemon/%04d.json", ID)
	pokemon, err := utils.ReadJsonFromFile[data.BasePokemon](path)
	if err != nil {
		return nil, err
	}
	return &pokemon, nil
}
