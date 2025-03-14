package pokemon

import (
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
)

//go:generate mockgen -build_flags=--mod=mod -destination=mocks/mock_levelup.go -package=mock_pokemon github.com/sanchitdeora/PokeSim/pokemon LevelUp
type LevelUp interface {
	SendEvent(action data.LevelUpEvent)
	ReceiveResponse() data.LevelUpEvent
}

type LevelUpImpl struct {
	LvlUpEventChan    chan data.LevelUpEvent
	LvlUpResponseChan chan data.LevelUpEvent
}

func NewLevelUpEvents(lvlUpEventChan chan data.LevelUpEvent, lvlUpResponseChan chan data.LevelUpEvent) LevelUp {
	return &LevelUpImpl{
		LvlUpEventChan:    lvlUpEventChan,
		LvlUpResponseChan: lvlUpResponseChan,
	}
}

func (l *LevelUpImpl) SendEvent(event data.LevelUpEvent) {
	slog.Info("Sending event", "event", event.EventType, "body", event.Body)
	l.LvlUpEventChan <- event
}

func (l *LevelUpImpl) ReceiveResponse() data.LevelUpEvent {
	response := <-l.LvlUpResponseChan
	slog.Info("Received response", "response", response.EventType)
	return response
}
