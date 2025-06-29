package pokemon

import (
	"slices"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/errors"
	"github.com/sanchitdeora/PokeSim/gamestate"
	"github.com/sanchitdeora/PokeSim/logger"
)

type BoxManager interface {
	AddToBox(pokemon *data.Pokemon) error
	Release(pokemon *data.Pokemon) error
	Swap(boxPokemon *data.Pokemon, partyPokemon *data.Pokemon) error
	GetBoxPokemon() data.Box
}

type BoxOpts struct {
	Logger    logger.Logger
	GameStateManager gamestate.GameStateManager
}

type BoxImpl struct {
	opts BoxOpts
}

func NewBoxManager(opts BoxOpts) BoxManager {
	if opts.Logger == nil {
		opts.Logger = logger.NewDefaultLogger()
	}
	return &BoxImpl{opts: opts}
}

func (b *BoxImpl) AddToBox(pokemon *data.Pokemon) error {
	box := b.opts.GameStateManager.GetBox()
	(*box) = append((*box), pokemon.ToPokemonSave())

	err := b.opts.GameStateManager.Save()
	if err != nil {
		return err
	}
	return nil
}

func (b *BoxImpl) Release(pokemon *data.Pokemon) error {
	box := b.opts.GameStateManager.GetBox()

	for i := range *box {
		if (*box)[i].PokemonUUID == pokemon.PokemonUUID {
			(*box) = slices.Delete((*box), i, i+1)
			err := b.opts.GameStateManager.Save()
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.ErrPokemonNotFound
}

func (b *BoxImpl) Swap(boxPokemon *data.Pokemon, partyPokemon *data.Pokemon) error {
	box := b.opts.GameStateManager.GetBox()
	gamestate := b.opts.GameStateManager.Get()

	// party pokemon index
	if partyPokemon != nil {
		for i := range gamestate.User.Party {
			if gamestate.User.Party[i].PokemonUUID == partyPokemon.PokemonUUID {
				gamestate.User.Party = slices.Delete(gamestate.User.Party, i, i+1)
				break
			}
		}
		b.AddToBox(partyPokemon)
	}

	// box pokemon index
	if boxPokemon != nil {
		for i := range *box {
			if (*box)[i].PokemonUUID == boxPokemon.PokemonUUID {
				if len(gamestate.User.Party) >= 6 {
					return errors.ErrPartyFull
				}
				(*box) = slices.Delete((*box), i, i+1)
				break
			}
		}
		gamestate.User.Party = append(gamestate.User.Party, boxPokemon)
	}

	b.opts.GameStateManager.Save()
	return nil
}

func (b *BoxImpl) GetBoxPokemon() data.Box {
	return *b.opts.GameStateManager.GetBox()
}
