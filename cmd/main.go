package main

import (
	"github.com/sanchitdeora/PokeSim/gui"

	"github.com/sanchitdeora/PokeSim/logger"
)

const (
	SAVED_USER_PATH = "\\saved\\user.json"
)

func main() {
	logger.InitLogger()
	// logger.InitTestLogger()
	
	// initialize GUI
	gui.InitializeGUI()
}

