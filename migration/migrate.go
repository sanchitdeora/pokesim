package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/migration/types"
	"github.com/sanchitdeora/PokeSim/utils"
)

const (
	PokeApiBaseUrl  = "https://pokeapi.co/api/v2"
	PokemonResource = "pokemon"

	manualAdjustmentFile = "./manual_adjustment.csv"
	wildEncountersPath   = "./assets/wild_encounters.json"
)

type migrationOpts struct {
	CsvWriter *csv.Writer
	CsvFile   *os.File

	WildEncounters map[data.BasePokemonID]data.WildEncounter
}

func main() {
	if len(os.Args) < 3 {
		slog.Error("Usage: migration.go [pokemonID begin] [pokemonID end]")
		return
	}

	beginID, _ := strconv.Atoi(os.Args[1])
	endID, _ := strconv.Atoi(os.Args[2])
	slog.Info("Migrating pokemon", "beginID", beginID, "endID", endID)

	file, err := os.OpenFile(manualAdjustmentFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	wildEncounters, err := utils.ReadJsonFromFile[map[data.BasePokemonID]data.WildEncounter](wildEncountersPath)
	if err != nil {
		slog.Info("Error reading wild encounters json", "error", err)
		wildEncounters = map[data.BasePokemonID]data.WildEncounter{}
	}

	opts := migrationOpts{
		CsvWriter: csv.NewWriter(file),
		CsvFile:   file,

		WildEncounters: wildEncounters,
	}

	defer opts.CsvFile.Close()
	defer opts.CsvWriter.Flush()

	for i := beginID; i <= endID; i++ {
		// get base pokemon
		fmt.Println("==============================================================================")
		slog.Info("Migrating pokemon", "id", i)
		opts.MigratePokemonToAsset(fmt.Sprintf("%s/%s/%d/", PokeApiBaseUrl, PokemonResource, i))
		fmt.Println("==============================================================================")
	}

	// save wild encounters
	err = utils.WriteJsonToFile(wildEncountersPath, opts.WildEncounters)
	if err != nil {
		panic(err)
	}
}

func (opts *migrationOpts) MigratePokemonToAsset(pokemonUrl string) {

	// get pokemon
	loadedPokemon, err := GetPokemon(pokemonUrl)
	if err != nil {
		slog.Error("error loading pokemon json", "error", err)
	}

	// get species
	loadedSpecies, err := GetPokemonSpecies(loadedPokemon.Species.Url)
	if err != nil {
		slog.Error("error loading pokemon species json", "error", err)
	}

	// download sprites and save. Get save url and store
	backSpritePath := DownloadAndSaveSprite(loadedPokemon.Sprites.BackUrl, loadedSpecies.PokedexNumbers[0].EntryNumber, true)
	frontSpritePath := DownloadAndSaveSprite(loadedPokemon.Sprites.FrontUrl, loadedSpecies.PokedexNumbers[0].EntryNumber, false)

	// fetch base pokemon types
	var pType1, pType2 data.PokemonTypeName
	pType1 = data.PokemonTypeName(loadedPokemon.Types[0].Type.Name)
	if len(loadedPokemon.Types) > 1 {
		pType2 = data.PokemonTypeName(loadedPokemon.Types[1].Type.Name)
	}

	// build base pokemon struct
	basePokemon := data.BasePokemon{
		ID:             data.BasePokemonID(loadedSpecies.PokedexNumbers[0].EntryNumber),
		Name:           loadedPokemon.Name,
		BaseExperience: loadedPokemon.BaseExperience,
		GrowthRate:     data.GrowthRateTypes(loadedSpecies.GrowthRate.Name),
		MovesLearned:   MovesLearnedMapper(loadedPokemon.Moves),
		EvolutionChain: opts.EvolutionChainMapper(loadedSpecies.EvolutionChain.Url),
		SpritesURL: data.Sprites{
			BackPath:  backSpritePath,
			FrontPath: frontSpritePath,
		},
		BaseStats: data.PokemonStats{
			HP: data.PokemonStat{
				Value: loadedPokemon.Stats[0].BaseStat,
				EV:    loadedPokemon.Stats[0].Effort,
			},
			Attack: data.PokemonStat{
				Value: loadedPokemon.Stats[1].BaseStat,
				EV:    loadedPokemon.Stats[1].Effort,
			},
			Defense: data.PokemonStat{
				Value: loadedPokemon.Stats[2].BaseStat,
				EV:    loadedPokemon.Stats[2].Effort,
			},
			SpecialAttack: data.PokemonStat{
				Value: loadedPokemon.Stats[3].BaseStat,
				EV:    loadedPokemon.Stats[3].Effort,
			},
			SpecialDefense: data.PokemonStat{
				Value: loadedPokemon.Stats[4].BaseStat,
				EV:    loadedPokemon.Stats[4].Effort,
			},
			Speed: data.PokemonStat{
				Value: loadedPokemon.Stats[5].BaseStat,
				EV:    loadedPokemon.Stats[5].Effort,
			},
		},
		Type1:       pType1,
		Type2:       pType2,
		IsLegendary: loadedSpecies.IsLegendary,
		IsMythical:  loadedSpecies.IsMythical,
	}

	// save pokemon to file
	savePokemonFilePath := "./assets/pokemon/" + fmt.Sprintf("%04d", loadedSpecies.PokedexNumbers[0].EntryNumber) + ".json"
	err = utils.WriteJsonToFile(savePokemonFilePath, basePokemon)
	if err != nil {
		slog.Error("could not write to saved file", "error", err)
	}

	// add wild encounter entries
	opts.CreateWildEncounters(basePokemon)
}

func DownloadAndSaveSprite(downloadURL string, id int, isBack bool) string {
	var saveFilePath string
	if isBack {
		saveFilePath = "./assets/pokemon/img/back/" + fmt.Sprintf("%04d", id) + ".png"
	} else {
		saveFilePath = "./assets/pokemon/img/" + fmt.Sprintf("%04d", id) + ".png"
	}

	if utils.CheckPathExists(saveFilePath) {
		slog.Debug("File already exists", "filepath", saveFilePath)
		return saveFilePath
	}

	// download sprite
	err := DownloadPNG(downloadURL, saveFilePath)
	if err != nil {
		slog.Error("error downloading sprite", "url", downloadURL, "error", err)
		return ""
	}
	return saveFilePath
}

func MovesLearnedMapper(moves []types.Moves) map[int]data.Moves {
	movesLearned := make(map[int]data.Moves)

	for _, move := range moves {
		lastVG := move.VersionGroup[len(move.VersionGroup)-1]
		if lastVG.LearnMethod.Name == "level-up" {
			// move, err := LoadMovesJson("../pokemonMoves.json")
			move, err := GetPokemonMoves(move.Move.Url)
			if err != nil {
				slog.Error("error loading pokemon move json", "error", err)
			}

			if data.MoveDamageClass(move.DamageClass.Name) != data.Status && move.Power > 0 {
				movesLearned[lastVG.LevelLearnedAt] = data.Moves{
					ID:          move.ID,
					Name:        move.Name,
					Accuracy:    move.Accuracy,
					Priority:    move.Priority,
					Power:       move.Power,
					DamageClass: data.MoveDamageClass(move.DamageClass.Name),
					Type:        data.PokemonTypeName(move.Type.Name),
				}
			}
		}
	}
	return movesLearned
}

func (opts *migrationOpts) EvolutionChainMapper(EvolutionChainUrl string) map[int][]data.BasePokemonID {
	evolutionMap := make(map[int][]data.BasePokemonID)

	loadedEvolution, err := GetEvolutionChain(EvolutionChainUrl)
	if err != nil {
		slog.Error("error loading pokemon evolution chain json", "error", err)
	}

	stage := data.PreEvolution

	evolutionChain := loadedEvolution.Chain
	for {
		evolvesTo := evolutionChain.EvolvesTo
		if len(evolvesTo) < 1 {
			break
		}

		if len(evolvesTo) > 1 {
			speciesListStr := ""
			for _, chain := range evolvesTo {
				speciesListStr += chain.Species.Name + ", "
			}
			opts.reportManualAdjustmentReq(
				evolutionChain.Species.Name,
				"evolvesTo is more than one pokemon species: "+speciesListStr,
				"plan how to split",
			)
			panic("evolves to more than one pokemon")
		}

		evolvesToChain := evolvesTo[0]

		if evolvesToChain.EvolutionDetails[0].Trigger.Name != "level-up" || evolvesToChain.EvolutionDetails[0].MinLevel == 0 {
			opts.reportManualAdjustmentReq(
				evolutionChain.Species.Name,
				"missing min level with trigger: "+evolvesToChain.EvolutionDetails[0].Trigger.Name,
				"Decide a level",
			)
		}

		minLevel := GetAdjustedEvolutionLevel(evolvesToChain.EvolutionDetails[0], stage)
		loadedSpecies, err := GetPokemonSpecies(evolvesToChain.Species.Url)
		if err != nil {
			slog.Error("error loading pokemon species json", "error", err)
		}

		evolutionMap[minLevel] = append(evolutionMap[minLevel], data.BasePokemonID(loadedSpecies.PokedexNumbers[0].EntryNumber))

		evolutionChain = evolvesToChain
		stage++
	}

	return evolutionMap
}

func (opts *migrationOpts) CreateWildEncounters(basePokemon data.BasePokemon) {
	evolvesAt := 1
	for evolutionLvl, basePokemonId := range basePokemon.EvolutionChain {
		if utils.Contains(basePokemonId, basePokemon.ID) {
			evolvesAt = evolutionLvl
			break
		}
	}

	baseRarity := opts.GetBaseRarity(basePokemon)

	wildEncounter := data.WildEncounter{
		BasePokemonID: basePokemon.ID,
		Name:          basePokemon.Name,
		Type1:         basePokemon.Type1,
		Type2:         basePokemon.Type2,
		EvolvedAt:     evolvesAt,
		Environments:  data.GetEnvironmentListFromType(basePokemon.Type1, basePokemon.Type2),
		BaseRarity:    baseRarity,
	}

	opts.WildEncounters[basePokemon.ID] = wildEncounter
}

func (opts *migrationOpts) GetBaseRarity(basePokemon data.BasePokemon) data.Rarity {

	// override rarity
	for rarity, ids := range data.RaritySpeciesOverrides {
		if utils.Contains(ids, basePokemon.ID) {
			return rarity
		}
	}

	switch {
	// legendary pokemons
	case basePokemon.IsLegendary:
		return data.LegendaryRarity

	// mythical and pseudo legendary pokemons
	case basePokemon.IsMythical || basePokemon.IsPokemonPseudoLegendary():
		return data.MythicalRarity

	// BST 550 - 600 or Late Evolutions or Unique Species
	case basePokemon.BaseStatTotal() > 550 || basePokemon.FinalEvolutionLevel() >= 40 || basePokemon.IsUniqueSpecies():
		return data.UltraRarity

	// BST 500 - 550 or Final Forms
	case basePokemon.BaseStatTotal() > 450 || basePokemon.GetEvolutionStage() == data.FinalEvolution:
		return data.RareRarity

	// BST 500 - 550 or Mid Forms
	case basePokemon.BaseStatTotal() > 350 || basePokemon.GetEvolutionStage() == data.MidEvolution || (basePokemon.GetEvolutionStage() == data.FinalEvolution && basePokemon.BaseStatTotal() <= 450):
		return data.UncommonRarity

	default:
		return data.CommonRarity
	}
}

func (opts *migrationOpts) reportManualAdjustmentReq(pokemonName string, reason string, solution string) {
	// for testing
	if opts.CsvWriter == nil || opts.CsvFile == nil {
		return
	}
	// Check if the CSV file exists
	if _, err := os.Stat(opts.CsvFile.Name()); err != nil {
		slog.Error("error opening CSV file", "error", err)
		return
	}

	opts.CsvWriter.Write([]string{pokemonName, reason, solution})
}

func GetAdjustedEvolutionLevel(details types.EvolutionDetails, stage data.EvolutionStage) int {

	if details.Trigger.Name == "level-up" {
		// Handle level-based evolutions
		if details.MinLevel > 0 {
			return details.MinLevel
		}

		// Handle friendship-based evolutions (e.g., Riolu → Lucario)
		if details.MinHappiness > 0 {
			if stage == data.PreEvolution {
				return 20 // Riolu → Lucario
			}
			return 30 // Chansey → Blissey (or similar)
		}

		// Handle time of day evolutions (e.g., Eevee → Espeon)
		if details.TimeOfDay != "" {
			if stage == data.PreEvolution {
				return 20 // Happiny → Chansey
			}
			return 30 // Other evolutions (e.g., Espeon, Umbreon)
		}

		// Handle held item evolutions (e.g., Sneasel with Razor Claw)
		if details.HeldItem.Name != "" {
			if stage == data.PreEvolution {
				return 20 // Happiny → Chansey
			}
			return 35 // Other evolutions (e.g., Sneasel → Weavile)
		}

		// Handle location-based evolutions (e.g., Magneton → Magnezone)
		if details.Location != "" {
			return 35
		}
	}

	// Handle trade-based evolutions (e.g., Scyther → Scizor)
	if details.Trigger.Name == "trade" {
		return 40
	}

	// Handle use-item evolutions (e.g., Poliwhirl → Politoed with King's Rock)
	if details.Trigger.Name == "use-item" {
		return 25
	}

	return details.MinLevel
}
