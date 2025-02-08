package gui

import (
	"gioui.org/widget"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/trainermanagement"
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
		Trainer:   trainer.GetTrainer(),
		Image:     g.loadPokemonImage("assets/pokemon/img/0001.png"),
		Unlocked:  true,
	}
}
