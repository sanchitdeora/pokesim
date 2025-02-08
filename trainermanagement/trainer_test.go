package trainermanagement_test

import (
	"testing"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/trainermanagement"
	"github.com/stretchr/testify/assert"
)

func createTrainerManager() trainermanagement.TrainerManager {
	return trainermanagement.NewTrainerManager(trainermanagement.TrainerOpts{
		SavedTrainerPath: "testfiles\\trainer_files\\test_trainer.json",
	})
}

func TestGetTrainer(t *testing.T) {
	tm := createTrainerManager()
	trainer := tm.GetTrainer()

	assert.Equal(t, "Steve Austin", trainer.Name)
	assert.Equal(t, data.TrainerPrefix, trainer.Type)
	assert.Equal(t, "test-badge", trainer.Rewards.Badge.Name)
	assert.Equal(t, 1, len(trainer.Bag))
}
