package data

import (
	"fmt"

	"github.com/sanchitdeora/PokeSim/utils"
)

type BasePokemonID int

type PokemonSave struct {
	BasePokemonID     BasePokemonID `json:"base_pokemon_id"`
	PokemonUUID       string        `json:"pokemon_uuid"`
	Stats             PokemonStats  `json:"stats"`
	Level             int           `json:"level"`
	ExperienceLeft    int           `json:"experience_left"`
	Moveset           Moveset       `json:"moveset"`
	EvolutionRejected bool          `json:"evolution_rejected"`
}

type Pokemon struct {
	BasePokemon
	PokemonUUID       string       `json:"pokemon_uuid"`
	Level             int          `json:"level"`
	ExperienceLeft    int          `json:"experience_left"`
	Stats             PokemonStats `json:"stats"`
	Moveset           Moveset      `json:"moveset"`
	EvolutionRejected bool         `json:"evolution_rejected"`
}

type BasePokemon struct {
	ID             BasePokemonID           `json:"id"`
	Name           string                  `json:"name"`
	BaseExperience int                     `json:"base_experience"`
	GrowthRate     GrowthRateTypes         `json:"growth_rate"`
	CaptureRate    int                     `json:"capture_rate"`
	MovesLearned   map[int]Moves           `json:"moves_learned_by_level"`
	EvolutionChain map[int][]BasePokemonID `json:"evolution_chain"`
	SpritesURL     Sprites                 `json:"sprites"`
	BaseStats      PokemonStats            `json:"base_stats"`
	Type1          PokemonTypeName         `json:"type1"`
	Type2          PokemonTypeName         `json:"type2,omitempty"`
	IsLegendary    bool                    `json:"is_legendary"`
	IsMythical     bool                    `json:"is_mythical"`
}

type Sprites struct {
	BackPath  string `json:"back_default"`
	FrontPath string `json:"front_default"`
}

type BasePokemonStats struct {
	HP             int `json:"hp"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	Speed          int `json:"speed"`
}

type PokemonStats struct {
	HP             PokemonStat `json:"hp"`
	Attack         PokemonStat `json:"attack"`
	Defense        PokemonStat `json:"defense"`
	SpecialAttack  PokemonStat `json:"special_attack"`
	SpecialDefense PokemonStat `json:"special_defense"`
	Speed          PokemonStat `json:"speed"`
}

type PokemonStat struct {
	Value int `json:"value"`
	IV    int `json:"iv"`
	EV    int `json:"ev"`
}

type Moves struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Accuracy    int             `json:"accuracy"`
	Priority    int             `json:"priority"`
	Power       int             `json:"power"`
	DamageClass MoveDamageClass `json:"damage_class"`
	Type        PokemonTypeName `json:"type"`
}

type Moveset struct {
	Move1 *Moves `json:"move1"`
	Move2 *Moves `json:"move2"`
	Move3 *Moves `json:"move3"`
	Move4 *Moves `json:"move4"`
}

type GrowthRateTypes string

const (
	Erratic     GrowthRateTypes = "erratic"
	Fast        GrowthRateTypes = "fast"
	MediumFast  GrowthRateTypes = "medium-fast"
	MediumSlow  GrowthRateTypes = "medium-slow"
	Slow        GrowthRateTypes = "slow"
	Fluctuating GrowthRateTypes = "fluctuating"
)

type MoveDamageClass string

const (
	Physical MoveDamageClass = "physical"
	Special  MoveDamageClass = "special"
	Status   MoveDamageClass = "status"
)

type EvolutionStage int

const (
	PreEvolution EvolutionStage = iota
	MidEvolution
	FinalEvolution
)

type LevelUpEvent struct {
	EventType LevelUpEventType `json:"event_type"`
	Body      interface{}      `json:"body"`
}

type LevelUpEventType string

const (
	LevelUpEventLearnMove LevelUpEventType = "learn_move"
	LevelUpEventEvolve    LevelUpEventType = "evolve"
)

type EventEvolveBody struct {
	PokemonUUID        string      `json:"pokemon_uuid"`
	CurrentBasePokemon BasePokemon `json:"current_base_pokemon"`
	EvolvedBasePokemon BasePokemon `json:"evolved_base_pokemon"`
}

type ResponseEvolveBody struct {
	AcceptEvolution bool `json:"accept_evolution"`
}

type EventLearnMoveBody struct {
	Pokemon Pokemon `json:"pokemon"`
	NewMove Moves   `json:"new_move"`
}

type ResponseLearnMoveBody struct {
	UpdatedMoveset Moveset `json:"moveset"`
	ReplacedMove   Moves   `json:"replaced_move"`
	NewMove        Moves   `json:"new_move"`
}

var StartPokemonIds = []BasePokemonID{
	1, 4, 7,
}

func (s *PokemonSave) ToPokemon() *Pokemon {
	path := fmt.Sprintf("/assets/pokemon/%04d.json", s.BasePokemonID)
	basePokemon, _ := utils.ReadJsonFromFile[BasePokemon](path)
	return &Pokemon{
		PokemonUUID:       s.PokemonUUID,
		BasePokemon:       basePokemon,
		Stats:             s.Stats,
		Level:             s.Level,
		ExperienceLeft:    s.ExperienceLeft,
		Moveset:           s.Moveset,
		EvolutionRejected: s.EvolutionRejected,
	}
}

type PokemonChangeOrder bool

const (
	ChangeOrderMoveDown PokemonChangeOrder = false
	ChangeOrderMoveUp   PokemonChangeOrder = true
)

func (p *Pokemon) ToPokemonSave() *PokemonSave {
	return &PokemonSave{
		PokemonUUID:       p.PokemonUUID,
		BasePokemonID:     p.BasePokemon.ID,
		Stats:             p.Stats,
		Level:             p.Level,
		ExperienceLeft:    p.ExperienceLeft,
		Moveset:           p.Moveset,
		EvolutionRejected: p.EvolutionRejected,
	}
}
