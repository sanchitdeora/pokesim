package errors

import "errors"

var (
	// file errors
	ErrFileDoesNotExist     error = errors.New("file does not exist")
	ErrCouldNotReadFromFile error = errors.New("could not read from file")

	// battle error
	ErrBattleOngoing error = errors.New("battle is ongoing")

	// shop error
	ErrNotEnoughMoney     error = errors.New("not enough money")
	ErrItemCountNotEnough error = errors.New("item count not enough")

	// pokemon error
	ErrPokemonNotFound       error = errors.New("pokemon not found")
	ErrBoxPokemonCannotBeNil error = errors.New("box pokemon cannot be nil")
	ErrPartyFull             error = errors.New("party is full")
)
