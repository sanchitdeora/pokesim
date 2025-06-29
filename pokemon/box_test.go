package pokemon

import (
	"testing"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/gamestate"
	"github.com/sanchitdeora/PokeSim/logger"
	"github.com/sanchitdeora/PokeSim/utils"
)

func TestBoxImpl_AddToBox(t *testing.T) {
	boxMgr := createBoxManager()

	pokemon := &data.Pokemon{
		PokemonUUID: "test-pokemon-uuid",
	}

	err := boxMgr.AddToBox(pokemon)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// verify that the pokemon was added to the box
	boxContents := boxMgr.GetBoxPokemon()
	if len(boxContents) != 1 {
		t.Errorf("expected 1 pokemon in box, got %d", len(boxContents))
	}

	if boxContents[0].PokemonUUID != pokemon.PokemonUUID {
		t.Errorf("expected pokemon %q in box, got %q", pokemon.PokemonUUID, boxContents[0].PokemonUUID)
	}

	err = boxMgr.Release(pokemon)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	boxContents = boxMgr.GetBoxPokemon()
	if len(boxContents) != 0 {
		t.Errorf("expected 1 pokemon in box, got %d", len(boxContents))
	}
}

func TestBoxImpl_Swap(t *testing.T) {
	boxMgr := createBoxManager()

	pokemon1 := &data.Pokemon{
		PokemonUUID: "test-pokemon-uuid-1",
	}
	err := boxMgr.AddToBox(pokemon1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	pokemon2 := getTestUser().User.Party[0]

	err = boxMgr.Swap(pokemon1, pokemon2)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// verify that the pokemons were swapped
	boxContents := boxMgr.GetBoxPokemon()

	if len(boxContents) != 1 {
		t.Errorf("expected 1 pokemon in box, got %d", len(boxContents))
	}

	if boxContents[0].PokemonUUID != pokemon2.PokemonUUID {
		t.Errorf("expected pokemon %q in box, got %q", pokemon2.PokemonUUID, boxContents[0].PokemonUUID)
	}

	err = boxMgr.Swap(pokemon2, pokemon1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if boxContents[0].PokemonUUID != pokemon1.PokemonUUID {
		t.Errorf("expected pokemon %q in box, got %q", pokemon1.PokemonUUID, boxContents[0].PokemonUUID)
	}

	err = boxMgr.Release(pokemon1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	boxContents = boxMgr.GetBoxPokemon()
	if len(boxContents) != 0 {
		t.Errorf("expected 1 pokemon in box, got %d", len(boxContents))
	}
}

func TestBoxImpl_Swap_OnePokemonMove(t *testing.T) {
	boxMgr := createBoxManager()

	pokemon1 := &data.Pokemon{
		PokemonUUID: "test-pokemon-uuid-1",
	}
	err := boxMgr.AddToBox(pokemon1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	err = boxMgr.Swap(pokemon1, nil)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// verify that the pokemons were swapped
	boxContents := boxMgr.GetBoxPokemon()
	if len(boxContents) != 0 {
		t.Errorf("expected 1 pokemon in box, got %d", len(boxContents))
	}

	err = boxMgr.Swap(nil, pokemon1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	boxContents = boxMgr.GetBoxPokemon()
	if len(boxContents) != 1 {
		t.Errorf("expected 1 pokemon in box, got %d", len(boxContents))
	}
	if boxContents[0].PokemonUUID != pokemon1.PokemonUUID {
		t.Errorf("expected pokemon %q in box, got %q", pokemon1.PokemonUUID, boxContents[0].PokemonUUID)
	}

	err = boxMgr.Release(pokemon1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	boxContents = boxMgr.GetBoxPokemon()
	if len(boxContents) != 0 {
		t.Errorf("expected 1 pokemon in box, got %d", len(boxContents))
	}
}

func createBoxManager() BoxManager {	
	return NewBoxManager(BoxOpts{
		Logger:    logger.NewDefaultLogger(),
		GameStateManager: gamestate.NewGameStateManager(nil, "testfiles/gametest", "unittestgame"),
	})
}

func getTestUser() *gamestate.GameState {
	user, _ := utils.ReadJsonFromFile[gamestate.GameStateSave]("/testfiles/gametest/unittestgame.json")
	return user.ToGameState()
}