package gamestate

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/errors"
	"github.com/sanchitdeora/PokeSim/utils"
)

//go:generate mockgen -build_flags=--mod=mod -destination=mocks/mock_game_state_manager.go -package=mock_game_state_manager github.com/sanchitdeora/PokeSim/gamestate GameStateManager
type GameStateManager interface {
	Save() error
	// Load() (*GameState, error)
	Get() *GameState
	AddTrainerProgress(trainerID string)
}

type GameStateImpl struct {
	*GameState
	Filepath string
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

func GetGameStates(prefix string) []string {
	if prefix == "" {
		prefix = "/saved"
	}

	return utils.GetListOfFilesInDirectory(prefix)
}

func NewGameStateManager(user *data.User, prefixPath string, filename string) GameStateManager {
	slog.Info("Creating game state manager...", "user", user, "filename", filename, "prefixPath", prefixPath)

	if filename == "" {
		filename = fmt.Sprintf("pokesim_%v", time.Now().Unix())
	}

	var filepath string
	if prefixPath != "" {
		filepath = fmt.Sprintf("%s/%s.json", prefixPath, filename)
	} else {
		filepath = fmt.Sprintf("saved/%s.json", filename)
	}

	// Load Game
	if utils.CheckPathExists(filepath) {
		slog.Info("Loading game state...", "filepath", filepath)
		game, err := LoadGame(filepath)
		if err != nil {
			slog.Error("error while loading game state", "error", err)
			panic(err)
		}

		return &GameStateImpl{
			GameState: game,
			Filepath:  filepath,
		}
	}

	// New Game
	gameState := &GameStateImpl{
		GameState: &GameState{
			User:            user,
			TrainerProgress: make([]string, 0),
			// Pokedex:         pokedex,
		},
		Filepath: filepath,
	}

	err := gameState.Save()
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
	slog.Debug("Saving game state...", "filepath", g.Filepath, "game", g.GameState)
	return utils.WriteJsonToFile(g.Filepath, g.ToGameStateSave())
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

func (g *GameStateImpl) Get() *GameState {
	return g.GameState
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
