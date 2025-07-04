package main

import (
	"os"
	"testing"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
)

func TestMainTrainerGenerator(t *testing.T) {
	// Set up the command-line arguments
	os.Args = []string{
		"trainer_generator",
		"-n", "Test Trainer",
		"-g",
		"-badge-name", "thunder-badge",
		"-region", "kanto",
		"-p", "100:16,25:17,26:18",
	}

	// Run the main function
	main()

	// Check that the output file was generated
	trainer, _ := utils.ReadJsonFromFile[data.TrainerSave]("/assets/trainer/gym/test trainer.json")

	// Check that the trainer data is correct
	if trainer.Name != "Test Trainer" {
		t.Errorf("Trainer name incorrect: want %s, got %s", "Test Trainer", trainer.Name)
	}
	if len(trainer.Party) != 2 {
		t.Errorf("Trainer party size incorrect: want %d, got %d", 1, len(trainer.Party))
	}
	if trainer.Party[0].BasePokemonID != 74 {
		t.Errorf("Trainer party Pokémon ID incorrect: want %d, got %d", 1, trainer.Party[0].BasePokemonID)
	}
}
