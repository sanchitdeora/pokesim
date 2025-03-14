package data

import "github.com/sanchitdeora/PokeSim/utils"

type Environment string

const (
	Forest   Environment = "forest"
	Cave     Environment = "cave"
	Lake     Environment = "lake"
	Mountain Environment = "mountain"
	Plains   Environment = "plains"
	Beach    Environment = "beach"
	Desert   Environment = "desert"
	Swamp    Environment = "swamp"
	Volcano  Environment = "volcano"
	Sky      Environment = "sky"
	Ruins    Environment = "ruins"
	Tundra   Environment = "tundra"
)

var EnvironmentPokemonTypeMap = map[Environment][]PokemonTypeName{
	Forest:   {NormalType, GrassType, BugType, FlyingType, PoisonType, FairyType},
	Cave:     {RockType, GroundType, SteelType, DarkType, DragonType},
	Lake:     {WaterType, IceType, FairyType, DragonType},
	Mountain: {RockType, GroundType, FireType, IceType, DragonType, FlyingType},
	Plains:   {NormalType, ElectricType, FightingType, GroundType},
	Beach:    {WaterType, ElectricType, FlyingType, IceType},
	Desert:   {GroundType, RockType, FireType, SteelType},
	Swamp:    {PoisonType, WaterType, BugType, GrassType, GhostType},
	Volcano:  {FireType, RockType, GroundType, DragonType},
	Sky:      {FlyingType, ElectricType, DragonType},
	Ruins:    {PsychicType, GhostType, DarkType, RockType, SteelType},
	Tundra:   {IceType, WaterType, GroundType, FightingType},
}

type Rarity int

const (
	CommonRarity Rarity = iota
	UncommonRarity
	RareRarity
	UltraRarity
	MythicalRarity
	LegendaryRarity
)

type WildEncounter struct {
	BasePokemonID BasePokemonID   `json:"base_pokemon_id"`
	Name          string          `json:"name"`
	Type1         PokemonTypeName `json:"type1"`
	Type2         PokemonTypeName `json:"type2,omitempty"`
	EvolvedAt     int             `json:"evolved_at"`
	Environments  []Environment   `json:"environments"`
	BaseRarity    Rarity          `json:"base_rarity"`
	Weight        float64         `json:"weight,omitempty"`
}

var RaritySpeciesOverrides = map[Rarity][]BasePokemonID{
	UltraRarity: {
		132, // Ditto
		351, // Castform
		442, // Rotom
		479, // Spiritomb
	},
	RareRarity: {
		201, // Unknown
	},
}

func GetEnvironmentListFromType(type1 PokemonTypeName, type2 PokemonTypeName) []Environment {
	var envList []Environment

	for _, pType := range []PokemonTypeName{type1, type2} {
		for env, pTypes := range EnvironmentPokemonTypeMap {
			if utils.Contains(pTypes, pType) {
				if utils.Contains(envList, env) {
					continue
				}
				envList = append(envList, env)
			}
		}
	}
	return envList
}

func (bp *BasePokemon) BaseStatTotal() int {
	return bp.BaseStats.HP.Value + bp.BaseStats.Attack.Value + bp.BaseStats.Defense.Value +
		bp.BaseStats.SpecialAttack.Value + bp.BaseStats.SpecialDefense.Value + bp.BaseStats.Speed.Value
}

func (bp *BasePokemon) IsPokemonPseudoLegendary() bool {
	// #1 pokemon should be final stage of evolution chain (not including mega evolution)
	if len(bp.EvolutionChain) < 3 || utils.Contains(bp.EvolutionChain[2], bp.ID) {
		return false
	}

	// #2 pokemon has a base stat total of exactly 600
	if (bp.BaseStatTotal()) < 600 {
		return false
	}

	// #3 pokemon has a slow GrowthRate i.e. 125,000 XP at Lvl 100
	if bp.GrowthRate != Slow {
		return false
	}

	return true
}

func (bp *BasePokemon) GetEvolutionStage() EvolutionStage {
	for lvl, ids := range bp.EvolutionChain {
		if utils.Contains(ids, bp.ID) {
			if bp.FinalEvolutionLevel() == lvl {
				return FinalEvolution
			} else {
				return MidEvolution
			}
		}
	}

	return PreEvolution
}

func (bp *BasePokemon) FinalEvolutionLevel() int {
	var finalLvl int
	for lvl := range bp.EvolutionChain {
		if lvl > finalLvl {
			finalLvl = lvl
		}
	}

	return finalLvl
}

func (bp *BasePokemon) IsUniqueSpecies() bool {
	return len(bp.EvolutionChain) == 0 && bp.BaseStatTotal() > 400
}

func (w WildEncounter) PrimaryTypeMatchesEnv(env Environment) bool {
	return utils.Contains(EnvironmentPokemonTypeMap[env], w.Type1)
}

func (w WildEncounter) SecondaryTypeMatchesEnv(env Environment) bool {
	return utils.Contains(EnvironmentPokemonTypeMap[env], w.Type2)
}
