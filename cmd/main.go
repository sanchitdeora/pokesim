package main

import (
	// "fmt"

	"github.com/sanchitdeora/PokeSim/battle"
	// "github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/gui"
	"github.com/sanchitdeora/PokeSim/pokemon"

	// "github.com/sanchitdeora/PokeSim/trainermanagement"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	// "github.com/sanchitdeora/PokeSim/utils"
)

type Services struct {
	UserService usermanagement.User
	PokemonService pokemon.Service
	BattleService battle.BattleIFace
}

const (
	SAVED_USER_PATH = "C:\\Projects\\Go-projects\\src\\PokéSim\\saved\\user.json"
)

func main() {

	userService, pokemonService := initializeService()

	// initialize GUI
	gui.InitializeGUI(gui.GuiOpts{UserService: userService, PokemonService: pokemonService})

}

func initializeService() (usermanagement.User, pokemon.Service) {
	userService := usermanagement.NewUserService(usermanagement.UserOpts{
		SavedUserPath: SAVED_USER_PATH,
	})

	pokemonService := pokemon.NewPokemonService(pokemon.PokemonOpts{})
	
	return userService, pokemonService
}