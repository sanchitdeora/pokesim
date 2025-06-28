package gamestate

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/errors"
	"github.com/sanchitdeora/PokeSim/utils"
)

const (
	SavedGamePrefix = "/saved/game"
	SavedBoxPrefix  = "/saved/box"
)

//go:generate mockgen -build_flags=--mod=mod -destination=mocks/mock_game_state_manager.go -package=mock_game_state_manager github.com/sanchitdeora/PokeSim/gamestate GameStateManager
type GameStateManager interface {
	Save() error
	Get() *GameState
	GetBox() *data.Box
	AddTrainerProgress(trainerID string)
}

type GameStateImpl struct {
	// gamestate
	*GameState
	GamePath string

	// boxstate
	Box     *data.Box
	BoxPath string
}

type GameState struct {
	User            *data.User
	TrainerProgress []string
	// Pokedex         *data.Pokedex
}

type GameStateSave struct {
	User            *data.UserSave `json:"user"`
	TrainerProgress []string       `json:"trainer_progress"`
	// Pokedex         *data.Pokedex
}

func GetGameStatesPath(prefix string) []string {
	if prefix == "" {
		prefix = SavedGamePrefix
	}

	return utils.GetListOfFilesInDirectory(prefix)
}

func NewGameStateManager(user *data.User, prefixPath string, filename string) GameStateManager {
	slog.Info("Creating game state manager...", "user", user, "filename", filename, "prefixPath", prefixPath)

	if filename == "" {
		filename = fmt.Sprintf("pokesim_%v", time.Now().Unix())
	}

	var gamePath string
	if prefixPath != "" {
		gamePath = fmt.Sprintf("%s/%s.json", prefixPath, filename)
	} else {
		gamePath = fmt.Sprintf("%s/%s.json", SavedGamePrefix, filename)
	}

	var boxPath string
	if prefixPath != "" {
		boxPath = fmt.Sprintf("%s/box/%s.json", prefixPath, filename)
	} else {
		boxPath = fmt.Sprintf("%s/%s.json", SavedBoxPrefix, filename)
	}

	game, err := LoadGame(gamePath)

	// Load Game
	if game != nil {
		slog.Info("Loading game state...", "gamePath", gamePath)

		box, err := LoadBox(boxPath)
		if box == nil {
			box = make(data.Box, 0)
		} else if err != nil {
			slog.Error("error while loading box state", "error", err)
		}

		return &GameStateImpl{
			GameState: game,
			GamePath:  gamePath,

			Box:     &box,
			BoxPath: boxPath,
		}
	}

	// New Game
	box := make(data.Box, 0)
	gameState := &GameStateImpl{
		GameState: &GameState{
			User:            user,
			TrainerProgress: make([]string, 0),
			// Pokedex:         pokedex,
		},
		GamePath: gamePath,

		Box:     &box,
		BoxPath: boxPath,
	}

	err = gameState.Save()
	if err != nil {
		slog.Error("error while saving game state", "error", err)
	}

	return gameState
}

func (g *GameStateImpl) AddTrainerProgress(trainerID string) {
	g.TrainerProgress = append(g.TrainerProgress, trainerID)
	g.Save()
}

func (g *GameStateImpl) Save() error {
	slog.Info("Saving game state...", "filepath", g.GamePath, "game", g.GameState)
	err := utils.WriteJsonToFile(g.GamePath, g.ToGameStateSave())
	if err != nil {
		slog.Error("error while saving game state", "error", err)
		return err
	}

	err = utils.WriteJsonToFile(g.BoxPath, g.Box)
	if err != nil {
		slog.Error("error while saving box state", "error", err)
		return err
	}

	return nil
}

func LoadGame(filepath string) (*GameState, error) {
	var gameState *GameStateSave
	var err error

	if utils.CheckPathExists(filepath) {
		gameState, err = utils.ReadJsonFromFile[*GameStateSave](filepath)
		if err != nil {
			slog.Error("could not read from saved file", "error", err)
			return nil, errors.ErrCouldNotReadFromFile
		}
	} else {
		slog.Error("file does not exist...")
		return nil, errors.ErrFileDoesNotExist
	}

	return gameState.ToGameState(), nil
}

func LoadBox(filepath string) (data.Box, error) {
	var boxState data.Box
	var err error

	if utils.CheckPathExists(filepath) {
		boxState, err = utils.ReadJsonFromFile[data.Box](filepath)
		if err != nil {
			slog.Error("could not read from saved file", "error", err)
			return nil, errors.ErrCouldNotReadFromFile
		}
	} else {
		slog.Error("file does not exist...")
		return nil, errors.ErrFileDoesNotExist
	}

	return boxState, nil
}

func (g *GameStateImpl) Get() *GameState {
	return g.GameState
}

func (g *GameStateImpl) GetBox() *data.Box {
	return g.Box
}

func (g *GameState) ToGameStateSave() *GameStateSave {

	return &GameStateSave{
		User:            g.User.ToUserSave(),
		TrainerProgress: g.TrainerProgress,
		// Pokedex:         g.Pokedex,
	}
}

func (g *GameStateSave) ToGameState() *GameState {
	return &GameState{
		User:            g.User.ToUser(),
		TrainerProgress: g.TrainerProgress,
		// Pokedex:         g.Pokedex,
	}
}
