package pokemon

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"strings"

	"github.com/google/uuid"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/gamestate"
	"github.com/sanchitdeora/PokeSim/logger"
	"github.com/sanchitdeora/PokeSim/utils"
)

//go:generate mockgen -build_flags=--mod=mod -destination=mocks/mock_pokemon.go -package=mock_pokemon github.com/sanchitdeora/PokeSim/pokemon PokemonService
type PokemonService interface {
	// LevelUp(pokemon *data.Pokemon)
	Evolve(pokemon *data.Pokemon)
	LearnNewMoves(pokemon *data.Pokemon)
	ExperienceGain(pokemon *data.Pokemon, faintedPokemon data.Pokemon)
	EvGain(pokemon *data.Pokemon, evYieldToAdd data.PokemonStats)
	GetExperienceRequiredForNextLevel(pokemon *data.Pokemon) int

	GenerateStarterPokemon(basePokemon data.BasePokemon) *data.Pokemon
	SearchWildPokemon(env data.Environment) *data.Pokemon
}

type PokemonOpts struct {
	Logger           logger.Logger
	GameStateManager gamestate.GameStateManager
	LvlUpActions     LevelUp
}

type PokemonImpl struct {
	opts PokemonOpts
}

func NewPokemonService(opts PokemonOpts) PokemonService {
	if opts.Logger == nil {
		opts.Logger = logger.NewDefaultLogger()
	}
	return &PokemonImpl{opts: opts}
}

func (p *PokemonImpl) LevelUp(pokemon *data.Pokemon) {
	pokemon.Level++
	p.opts.Logger.Log(fmt.Sprintf("%s reached level %v", utils.ToCapitalizeFirstLetterOfEachWord(pokemon.Name), pokemon.Level))

	// calculate stats upgrade
	p.statUpgrades(pokemon)

	// calculate pokemon experience left
	pokemon.ExperienceLeft = p.GetExperienceRequiredForNextLevel(pokemon)

	// should pokemon learns moves; if yes which moveset.
	p.LearnNewMoves(pokemon)
}

func (p *PokemonImpl) Evolve(pokemon *data.Pokemon) {
	if !canPokemonEvolve(pokemon) {
		slog.Info("pokemon cannot evolve", "pokemon", pokemon.Name, "pokemonUUID", pokemon.PokemonUUID)
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

		responseEvolveBody := p.shouldPokemonEvolve(pokemon, evolvedBasePokemon)

		if !responseEvolveBody.AcceptEvolution {
			pokemon.EvolutionRejected = true
			return
		}

		pokemon.BasePokemon = evolvedBasePokemon
		slog.Info("user info", "user", p.opts.GameStateManager.Get().User.Party[0])
	}

	p.statUpgrades(pokemon)
	p.opts.GameStateManager.Save()
}

func (p *PokemonImpl) LearnNewMoves(pokemon *data.Pokemon) {
	// input which move should be replaced.
	newMove, exists := pokemon.MovesLearned[pokemon.Level]
	if !exists {
		return
	}

	if pokemon.Moveset.Move1 == nil {
		pokemon.Moveset.Move1 = &newMove
		p.opts.Logger.Log(fmt.Sprintf("%s learned move %s", utils.ToCapitalizeFirstLetterOfEachWord(pokemon.Name), utils.ToCapitalizeFirstLetterOfEachWord(newMove.Name)))
		return
	}
	if pokemon.Moveset.Move2 == nil {
		pokemon.Moveset.Move2 = &newMove
		p.opts.Logger.Log(fmt.Sprintf("%s learned move %s", utils.ToCapitalizeFirstLetterOfEachWord(pokemon.Name), utils.ToCapitalizeFirstLetterOfEachWord(newMove.Name)))
		return
	}
	if pokemon.Moveset.Move3 == nil {
		pokemon.Moveset.Move3 = &newMove
		p.opts.Logger.Log(fmt.Sprintf("%s learned move %s", utils.ToCapitalizeFirstLetterOfEachWord(pokemon.Name), utils.ToCapitalizeFirstLetterOfEachWord(newMove.Name)))
		return
	}
	if pokemon.Moveset.Move4 == nil {
		pokemon.Moveset.Move4 = &newMove
		p.opts.Logger.Log(fmt.Sprintf("%s learned move %s", utils.ToCapitalizeFirstLetterOfEachWord(pokemon.Name), utils.ToCapitalizeFirstLetterOfEachWord(newMove.Name)))
		return
	}

	// replace move request
	p.opts.LvlUpActions.SendEvent(data.LevelUpEvent{
		EventType: data.LevelUpEventLearnMove,
		Body: data.EventLearnMoveBody{
			Pokemon: *pokemon,
			NewMove: newMove,
		},
	})

	// receive response from user.
	var learnMoveBody data.ResponseLearnMoveBody
	for {
		action := p.opts.LvlUpActions.ReceiveResponse()
		if action.EventType == data.LevelUpEventLearnMove {
			learnMoveBody = action.Body.(data.ResponseLearnMoveBody)
			break
		}
	}
	if learnMoveBody.NewMove == learnMoveBody.ReplacedMove {
		p.opts.Logger.Log(fmt.Sprintf("%s did not learn move %s", utils.ToCapitalizeFirstLetterOfEachWord(pokemon.Name), utils.ToCapitalizeFirstLetterOfEachWord(learnMoveBody.NewMove.Name)))
	} else {
		pokemon.Moveset = learnMoveBody.UpdatedMoveset
		p.opts.Logger.Log(fmt.Sprintf("%s learned move %s", utils.ToCapitalizeFirstLetterOfEachWord(pokemon.Name), utils.ToCapitalizeFirstLetterOfEachWord(learnMoveBody.NewMove.Name)))
		p.opts.Logger.Log(fmt.Sprintf("%s forgot move %s", utils.ToCapitalizeFirstLetterOfEachWord(pokemon.Name), utils.ToCapitalizeFirstLetterOfEachWord(learnMoveBody.ReplacedMove.Name)))
	}
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

	p.opts.Logger.Log(fmt.Sprintf("%s gained %v experience points", utils.ToCapitalizeFirstLetterOfEachWord(pokemon.Name), expGain))

	for {
		if pokemon.ExperienceLeft > expGain {
			pokemon.ExperienceLeft -= expGain
			break
		}
		expGain -= pokemon.ExperienceLeft
		p.LevelUp(pokemon)
	}
}

func (p *PokemonImpl) statUpgrades(pokemon *data.Pokemon) {
	hp := calculateHPStatUpgrade(pokemon.BaseStats.HP.Value, pokemon.Stats.HP.IV, pokemon.Stats.HP.EV, pokemon.Level)

	attack := calculateOtherStatUpgrade(pokemon.BaseStats.Attack.Value, pokemon.Stats.Attack.IV, pokemon.Stats.Attack.EV, pokemon.Level)
	defense := calculateOtherStatUpgrade(pokemon.BaseStats.Defense.Value, pokemon.Stats.Defense.IV, pokemon.Stats.Defense.EV, pokemon.Level)

	spAttack := calculateOtherStatUpgrade(pokemon.BaseStats.SpecialAttack.Value, pokemon.Stats.SpecialAttack.IV, pokemon.Stats.SpecialAttack.EV, pokemon.Level)
	spDefense := calculateOtherStatUpgrade(pokemon.BaseStats.SpecialDefense.Value, pokemon.Stats.SpecialDefense.IV, pokemon.Stats.SpecialDefense.EV, pokemon.Level)

	speed := calculateOtherStatUpgrade(pokemon.BaseStats.Speed.Value, pokemon.Stats.Speed.IV, pokemon.Stats.Speed.EV, pokemon.Level)

	slog.Info("Stat upgrade for pokemon:")
	p.opts.Logger.Log(fmt.Sprintf("\nHP: +%v", hp-pokemon.Stats.HP.Value))
	p.opts.Logger.Log(fmt.Sprintf("Attack: +%v", attack-pokemon.Stats.Attack.Value))
	p.opts.Logger.Log(fmt.Sprintf("Defense: +%v", defense-pokemon.Stats.Defense.Value))
	p.opts.Logger.Log(fmt.Sprintf("Special Attack: +%v", spAttack-pokemon.Stats.SpecialAttack.Value))
	p.opts.Logger.Log(fmt.Sprintf("Special Defence: +%v", spDefense-pokemon.Stats.SpecialDefense.Value))
	p.opts.Logger.Log(fmt.Sprintf("Speed: +%v", speed-pokemon.Stats.Speed.Value))

	pokemon.Stats.HP.Value = hp
	pokemon.Stats.Attack.Value = attack
	pokemon.Stats.Defense.Value = defense
	pokemon.Stats.SpecialAttack.Value = spAttack
	pokemon.Stats.SpecialDefense.Value = spDefense
	pokemon.Stats.Speed.Value = speed
}

// calculate next level exp required
func (p *PokemonImpl) GetExperienceRequiredForNextLevel(pokemon *data.Pokemon) int {
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

func (p *PokemonImpl) GenerateStarterPokemon(basePokemon data.BasePokemon) *data.Pokemon {
	return generatePokemon(basePokemon, 5)
}

func (p *PokemonImpl) SearchWildPokemon(env data.Environment) *data.Pokemon {
	minLvl, maxLvl := p.getWildPokemonLevelRange()

	wildEncounters, _ := utils.ReadJsonFromFile[map[data.BasePokemonID]data.WildEncounter]("./testfiles/wild_encounters.json")

	filterByLvlAndEnvs := make([]data.WildEncounter, 0)

	// Step 1: filter pokemons by level and environment
	for _, encounter := range wildEncounters {
		if encounter.EvolvedAt > maxLvl {
			continue
		}

		if utils.Contains(encounter.Environments, env) {
			filterByLvlAndEnvs = append(filterByLvlAndEnvs, encounter)
		}
	}

	// Step 2: Group by rarity and avoid duplicates
	wildEncounterMap := make(map[data.Rarity]map[int]data.WildEncounter)

	for _, encounter := range filterByLvlAndEnvs {
		if _, exists := wildEncounterMap[encounter.BaseRarity]; !exists {
			wildEncounterMap[encounter.BaseRarity] = make(map[int]data.WildEncounter)
		}
		wildEncounterMap[encounter.BaseRarity][int(encounter.BasePokemonID)] = encounter
	}

	// Step 3: Calculate weights
	weightedEncounters := make([]data.WildEncounter, 0)
	totalWeight := 0.0

	for rarity, encounters := range wildEncounterMap {
		for _, encounter := range encounters {
			baseWeight := getBaseRarityWeight(rarity)

			// Apply dynamic rarity adjustment
			dynamicWeight := adjustRarityByProgress(encounter, minLvl, maxLvl)

			// Apply environment weight adjustment
			typeMatchWeight := getEnvironmentMatchWeight(env, encounter)

			// Final weight calculation
			finalWeight := baseWeight * dynamicWeight * typeMatchWeight

			encounter.Weight = finalWeight
			totalWeight += finalWeight
			weightedEncounters = append(weightedEncounters, encounter)
		}
	}

	// Step 4: Weighted random selection
	if totalWeight == 0 {
		return nil
	}

	randWeight := rand.Float64() * totalWeight
	currentWeight := 0.0

	for _, encounter := range weightedEncounters {
		currentWeight += encounter.Weight
		if randWeight <= currentWeight {

			if minLvl < encounter.EvolvedAt {
				minLvl = encounter.EvolvedAt
			}
			lvl := int(math.Round(utils.RandomGenerator(float64(minLvl), float64(maxLvl))))
			basePokemon, err := getBasePokemonByID(encounter.BasePokemonID)
			if err != nil {
				slog.Error("failed to get base pokemon", "error", err)
				return nil
			}

			slog.Info("min level: %d, max level: %d", "minLvl", minLvl, "maxLvl", maxLvl, "level", lvl)

			return generatePokemon(basePokemon, lvl)
		}
	}
	return nil
}

func canPokemonEvolve(pokemon *data.Pokemon) bool {
	if pokemon.EvolutionRejected {
		return false
	}

	startLevel := 1
	for level, evolutions := range pokemon.EvolutionChain {
		if evolutions[0] == pokemon.BasePokemon.ID && pokemon.Level > level {
			startLevel = level + 1
			break
		}
	}
	for i := startLevel; i <= pokemon.Level; i++ {
		if _, exists := pokemon.EvolutionChain[i]; exists && pokemon.EvolutionChain[i][0] != pokemon.BasePokemon.ID {
			return true
		}
	}

	return false
}

func (p *PokemonImpl) shouldPokemonEvolve(pokemon *data.Pokemon, evolvedBasePokemon data.BasePokemon) data.ResponseEvolveBody {
	slog.Info("pokemon can evolve", "pokemon", pokemon.Name, "evolving to", evolvedBasePokemon.Name)

	// ask user if evolve?
	p.opts.LvlUpActions.SendEvent(data.LevelUpEvent{
		EventType: data.LevelUpEventEvolve,
		Body: data.EventEvolveBody{
			PokemonUUID:        pokemon.PokemonUUID,
			CurrentBasePokemon: pokemon.BasePokemon,
			EvolvedBasePokemon: evolvedBasePokemon,
		}})

	// receive response from user. TODO: add UI to accept and send back response
	var evolveBody data.ResponseEvolveBody
	for {
		action := p.opts.LvlUpActions.ReceiveResponse()
		if action.EventType == data.LevelUpEventEvolve {
			evolveBody = action.Body.(data.ResponseEvolveBody)
			break
		}
	}
	return evolveBody
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

func (p *PokemonImpl) getWildPokemonLevelRange() (int, int) {
	gamestate := p.opts.GameStateManager.Get()

	coreLevels := getCorePartyLvl(gamestate.User.Party)
	averagePartyLevel := getAvgLevel(coreLevels)

	// avgBoxLevel := getAvgLevel(box)
	avgBoxLevel := 0

	// Base level range
	minLevel := int(float64(averagePartyLevel) * 0.8)
	maxLevel := int(float64(averagePartyLevel) * 1.2)

	// Consider box Pokémon for slight influence
	if avgBoxLevel > 0 {
		minLevel = (minLevel + avgBoxLevel) / 2
		maxLevel = (maxLevel + avgBoxLevel) / 2
	}

	// Adjust based on progression
	trainersDefeated, gymLeadersDefeated := 0, 0
	for _, trainer := range gamestate.TrainerProgress {
		if strings.Split(trainer, "-")[0] == "tr" {
			trainersDefeated++
		} else if strings.Split(trainer, "-")[0] == "gym" {
			gymLeadersDefeated++
		}
	}
	progressionBoost := (trainersDefeated / 10) + (gymLeadersDefeated * 2)

	// Widen the range as player progresses
	minLevel -= progressionBoost / 2
	maxLevel += progressionBoost

	// Handle early game
	if averagePartyLevel <= 10 {
		minLevel = max(1, averagePartyLevel-2)
		maxLevel = averagePartyLevel + 2
	}

	// Clamp levels
	if minLevel < 1 {
		minLevel = 1
	}
	if maxLevel > 100 {
		maxLevel = 100
	}

	return minLevel, maxLevel
}

func adjustRarityByProgress(encounter data.WildEncounter, minLvl, maxLvl int) float64 {
	// Early evolutions become more common later
	if encounter.EvolvedAt > 0 && encounter.EvolvedAt <= minLvl+5 {
		return 1.5 // Becomes more common if the trainer is beyond evolution level
	}

	// Rare Pokémon may become slightly less rare later in the game
	if encounter.BaseRarity == data.RareRarity && maxLvl > 30 {
		return 1.0 + math.Min(0.2, float64(maxLvl-30)/100)
	}

	// Ultra-rare Pokémon may become slightly more common in late-game
	if encounter.BaseRarity == data.UltraRarity && maxLvl > 50 {
		return 0.8
	}

	return 1.0
}

func getBasePokemonByID(ID data.BasePokemonID) (data.BasePokemon, error) {
	path := fmt.Sprintf("/assets/pokemon/%04d.json", ID)
	pokemon, err := utils.ReadJsonFromFile[data.BasePokemon](path)
	if err != nil {
		return data.BasePokemon{}, err
	}
	return pokemon, nil
}

func generatePokemon(basePokemon data.BasePokemon, level int) *data.Pokemon {
	if level < 1 {
		slog.Info("level is less than 1", "level", level)
		level = 1
	}

	slog.Info("Generating pokemon...", "basepokemon", basePokemon.Name, "level", level)

	pokemon := &data.Pokemon{
		BasePokemon:    basePokemon,
		PokemonUUID:    uuid.NewString(),
		Level:          level,
		ExperienceLeft: 0,
		Stats: data.PokemonStats{
			HP:             generatePokemonHPStat(basePokemon.BaseStats.HP.Value, level),
			Attack:         generatePokemonOtherStat(basePokemon.BaseStats.Attack.Value, level),
			Defense:        generatePokemonOtherStat(basePokemon.BaseStats.Defense.Value, level),
			SpecialAttack:  generatePokemonOtherStat(basePokemon.BaseStats.SpecialAttack.Value, level),
			SpecialDefense: generatePokemonOtherStat(basePokemon.BaseStats.SpecialDefense.Value, level),
			Speed:          generatePokemonOtherStat(basePokemon.BaseStats.Speed.Value, level),
		},
		Moveset: setupMoveset(basePokemon, level),
	}

	slog.Info("Generated pokemon", "pokemon", pokemon, "level", level)

	return pokemon
}
