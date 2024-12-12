package main

import (

	battle "github.com/sanchitdeora/PokeSim/battle_v1"
	"github.com/sanchitdeora/PokeSim/gui"
	"github.com/sanchitdeora/PokeSim/logger"
	"github.com/sanchitdeora/PokeSim/pokemon"

	"github.com/sanchitdeora/PokeSim/usermanagement"
)

type Services struct {
	UserService    usermanagement.UserManager
	PokemonService pokemon.PokemonManager
	BattleService  battle.BattleSequence
}

const (
	SAVED_USER_PATH = "C:\\Projects\\Go-projects\\src\\PokéSim\\saved\\user.json"
)

func main() {
	logger.InitLogger()

	userService, pokemonService := initializeService()

	// initialize GUI
	gui.InitializeGUI(gui.GuiOpts{UserService: userService, PokemonService: pokemonService})

}

func initializeService() (usermanagement.UserManager, pokemon.PokemonManager) {
	userService := usermanagement.NewUserService(usermanagement.UserOpts{
		SavedUserPath: SAVED_USER_PATH,
	})

	pokemonService := pokemon.NewPokemonManager(pokemon.PokemonOpts{})

	return userService, pokemonService
}
