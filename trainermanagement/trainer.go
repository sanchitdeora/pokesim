package trainermanagement

import (
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/errors"
	"github.com/sanchitdeora/PokeSim/utils"
)

type TrainerManager interface {
	GetTrainer() *data.Trainer
}

type TrainerOpts struct {
	SavedTrainerPath string
}

type trainerImpl struct {
	opts TrainerOpts
	Trainer *data.Trainer
}

func NewTrainerManager(opts TrainerOpts) TrainerManager {
	trainer, err := loadTrainer(opts.SavedTrainerPath)
	if err != nil {
		slog.Error("failed to load trainer", "error", err)
		return nil
	}

	return &trainerImpl{
		opts: opts,
		Trainer: trainer,
	}
}

func (t *trainerImpl) GetTrainer() *data.Trainer {
	return t.Trainer
}

func loadTrainer(path string) (*data.Trainer, error) {
	var savedTrainer data.TrainerSave
	var err error

	if utils.CheckPathExists(path) {
		savedTrainer, err = utils.ReadJsonFromFile[data.TrainerSave](path)
		if err != nil {
			slog.Error("could not read from saved file", "error", err)
			return nil, errors.ErrCouldNotReadFromFile
		}
	} else {
		slog.Error("file does not exist...")
		return nil, errors.ErrFileDoesNotExist
	}
	return savedTrainer.ToTrainer(), nil
}
