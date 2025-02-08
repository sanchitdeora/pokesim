package main

import (
	battle "github.com/sanchitdeora/PokeSim/battle_v1"
	"github.com/sanchitdeora/PokeSim/gamestate"
	"github.com/sanchitdeora/PokeSim/gui"

	// "github.com/sanchitdeora/PokeSim/logger"
	"github.com/sanchitdeora/PokeSim/pokemon"

	"github.com/sanchitdeora/PokeSim/usermanagement"
)

type Services struct {
	UserService    usermanagement.UserManager
	PokemonService pokemon.PokemonService
	BattleService  battle.BattleSequence
}

const (
	SAVED_USER_PATH = "\\saved\\user.json"
)

func main() {
	// logger.InitLogger()

	gameManager, userManager, pokemonManager := initializeService()

	// initialize GUI
	gui.InitializeGUI(gui.GuiOpts{GameManager: gameManager, UserManager: userManager, PokemonService: pokemonManager})

}

func initializeService() (gamestate.GameStateManager, usermanagement.UserManager, pokemon.PokemonService) {
	gameManager := gamestate.GameStateManager(gamestate.NewGameStateManager(nil, "saved", "user"))

	userService := usermanagement.NewUserManager(usermanagement.UserOpts{
		GameState: gameManager,
	})

	pokemonService := pokemon.NewPokemonService(pokemon.PokemonOpts{})

	return gameManager, userService, pokemonService
}
