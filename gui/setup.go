package gui

import (
	"fmt"
	"log/slog"

	"gioui.org/widget"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/gamestate"
	"github.com/sanchitdeora/PokeSim/trainermanagement"
	"github.com/sanchitdeora/PokeSim/utils"
)

type TrainerUI struct {
	Image     widget.Image
	Name      string
	Unlocked  bool
	Clickable *widget.Clickable
	Trainer   *data.Trainer
}

func (g *Gui) setupTrainersList() []TrainerUI {
	trainers1 := []string{
		"testfiles/trainer_files/test_trainer.json",
		"testfiles/trainer_files/test_trainer.json",
		"testfiles/trainer_files/test_trainer.json",
		"testfiles/trainer_files/test_trainer.json",
		"testfiles/trainer_files/test_trainer.json",
	}

	var trainers []TrainerUI

	for _, t := range trainers1 {
		trainerManager := trainermanagement.NewTrainerManager(trainermanagement.TrainerOpts{
			SavedTrainerPath: t,
		})
		trainers = append(trainers, g.createTrainerUI(trainerManager))
	}
	return trainers
}

func (g *Gui) createTrainerUI(trainer trainermanagement.TrainerManager) TrainerUI {
	return TrainerUI{
		Trainer:  trainer.GetTrainer(),
		Image:    g.loadPokemonImage("assets/pokemon/img/0001.png"),
		Unlocked: true,
	}
}

func (g *Gui) createLoadGamesUI() []LoadGameUI {
	loadGamesUI := make([]LoadGameUI, 0)
	loadGames := gamestate.GetGameStates("")

	if len(loadGames) > 0 {
		for _, loadGame := range loadGames {
			loadGamesUI = append(loadGamesUI, LoadGameUI{
				LoadGameButton: &widget.Clickable{},
				FileName:       loadGame,
				Image:          g.loadPokemonImage("assets/trainer/img/user_trainer_avatar.png"),
			})
		}
	}
	loadGamesUI = append(loadGamesUI, LoadGameUI{
		LoadGameButton: &widget.Clickable{},
		FileName:       "New Game",
		Image:          g.loadPokemonImage("assets/trainer/img/new_game_img.png"),
	})
	return loadGamesUI
}

func (g *Gui) fetchStarterPokemonList() []data.BasePokemon {
	var baseStarterPokemons []data.BasePokemon

	for _, i := range data.StartPokemonIds {
		basePokemon, _ := utils.ReadJsonFromFile[data.BasePokemon](fmt.Sprintf("/assets/pokemon/%04d.json", i))
		baseStarterPokemons = append(baseStarterPokemons, basePokemon)
	}

	return baseStarterPokemons
}

func (g *Gui) prepareNewGameUser(trainerName string, starterPokemon data.BasePokemon) *data.User {
	slog.Info("Preparing new game user", "trainerName", trainerName, "starterPokemon", starterPokemon)
	return &data.User{
		BaseTrainer: data.BaseTrainer{
			Name:  trainerName,
			Party: []*data.Pokemon{g.opts.PokemonService.GenerateStarterPokemon(starterPokemon)},
			Bag:   data.ItemMap{},
		},
		Stats: &data.TrainerStats{
			Badges: make([]data.BadgeType, 0),
		},
		Money: 0,
	}
}