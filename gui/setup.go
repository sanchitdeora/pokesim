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

type LoadGameUI struct {
	LoadGameBtn *widget.Clickable
	FileName    string
	Image       widget.Image
}

type TrainerUI struct {
	Image    widget.Image
	Name     string
	Unlocked bool
	Trainer  *data.Trainer
}

type EnvironmentUI struct {
	Image          widget.Image
	Name           data.Environment
	EnvironmentBtn *widget.Clickable
}

type ItemShopUI struct {
	ImageURL    string
	StoreItem   data.StoreItem
	Unlocked    bool
	IncreaseBtn *widget.Clickable
	DecreaseBtn *widget.Clickable
}

func createLoadGamesUI() []LoadGameUI {
	loadGamesUI := make([]LoadGameUI, 0)
	loadGames := gamestate.GetGameStatesPath("")

	slog.Info("loadGames", "loadGames", loadGames)

	for _, loadGame := range loadGames {
		loadGamesUI = append(loadGamesUI, LoadGameUI{
			LoadGameBtn: new(widget.Clickable),
			FileName:    loadGame,
			Image:       loadImage("assets/trainer/img/user_trainer_avatar.png"),
		})
	}

	loadGamesUI = append(loadGamesUI, LoadGameUI{
		LoadGameBtn: new(widget.Clickable),
		FileName:    "New Game",
		Image:       loadImage("assets/trainer/img/new_game_img.png"),
	})

	return loadGamesUI
}

func setupTrainersList() []TrainerUI {
	trainers1 := []string{
		"assets/trainer/brock.json",
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
		trainers = append(trainers, createTrainerUI(trainerManager))
	}
	return trainers
}

func createTrainerUI(trainer trainermanagement.TrainerManager) TrainerUI {
	var path string
	slog.Info("Trainer", "name", trainer.GetTrainer().Name, "trainer image path", trainer.GetTrainer().ImagePath)
	if trainer.GetTrainer().ImagePath == "" {
		path = "assets/pokemon/img/0001.png"
	} else {
		path = trainer.GetTrainer().ImagePath
	}

	return TrainerUI{
		Trainer:  trainer.GetTrainer(),
		Image:    loadImage(path),
		Unlocked: true,
	}
}

func fetchStarterPokemonList() []data.BasePokemon {
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

func getItemShopList() []ItemShopUI {
	items := make([]ItemShopUI, 0)

	for _, i := range data.StoreItems {
		if i.Name == data.MasterBall {
			continue
		}
		items = append(items, ItemShopUI{
			StoreItem:   i,
			ImageURL:    "/assets/items/pokeball.png",
			Unlocked:    true,
			IncreaseBtn: new(widget.Clickable),
			DecreaseBtn: new(widget.Clickable),
		})
	}
	return items
}

func getEnvironments() []EnvironmentUI {
	return []EnvironmentUI{
		{loadImage("/assets/environments/forest.png"), data.Forest, new(widget.Clickable)},
		{loadImage("/assets/environments/cave.png"), data.Cave, new(widget.Clickable)},
		{loadImage("/assets/environments/lake.png"), data.Lake, new(widget.Clickable)},
		{loadImage("/assets/environments/mountain.png"), data.Mountain, new(widget.Clickable)},
		{loadImage("/assets/environments/plains.png"), data.Plains, new(widget.Clickable)},
		{loadImage("/assets/environments/beach.png"), data.Beach, new(widget.Clickable)},
		{loadImage("/assets/environments/desert.png"), data.Desert, new(widget.Clickable)},
		{loadImage("/assets/environments/swamp.png"), data.Swamp, new(widget.Clickable)},
		{loadImage("/assets/environments/volcano.png"), data.Volcano, new(widget.Clickable)},
		{loadImage("/assets/environments/sky.jpg"), data.Sky, new(widget.Clickable)},
		{loadImage("/assets/environments/ruins.png"), data.Ruins, new(widget.Clickable)},
		{loadImage("/assets/environments/tundra.png"), data.Tundra, new(widget.Clickable)},
	}
}
