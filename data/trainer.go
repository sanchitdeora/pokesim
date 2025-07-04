package data

type ItemMap map[ItemName]Item
type ItemMapSave map[ItemName]int

// save models
type BaseTrainerSave struct {
	Name  string         `json:"name"`
	Party []*PokemonSave `json:"party"`
	Bag   ItemMapSave    `json:"bag"`
}

type UserSave struct {
	BaseTrainerSave
	Stats *TrainerStats `json:"stats"`
	Money int           `json:"money"`
}

type TrainerSave struct {
	BaseTrainerSave
	// TrainerID string       `json:"trainer_id"`
	Type      TrainerClass `json:"type"`
	Rewards   *Rewards     `json:"rewards"`
	ImagePath string       `json:"image_path"`
}

type User struct {
	BaseTrainer
	Stats *TrainerStats
	Money int
}

type BaseTrainer struct {
	Name  string     `json:"name"`
	Party []*Pokemon `json:"party"`
	Bag   ItemMap    `json:"bag"`
}

type Trainer struct {
	BaseTrainer
	Type      TrainerClass
	Rewards   *Rewards
	ImagePath string
}

type Rewards struct {
	Items ItemMap   `json:"items"`
	Badge BadgeType `json:"badge_type,omitempty"`
}

type TrainerClass string

const (
	TrainerPrefix    TrainerClass = "trainer"
	GymLeaderPrefix  TrainerClass = "gym-leader"
	TournamentPrefix TrainerClass = "tournament-trainer"
	RivalPrefix      TrainerClass = "rival"
	WildPrefix       TrainerClass = "wild"
)

type TrainerStats struct {
	Badges  []BadgeType `json:"badges,omitempty"`
	Battles int         `json:"battles"`
	Wins    int         `json:"wins"`
	Losses  int         `json:"losses"`
	Catches int         `json:"catches"`
	PokeDEX int         `json:"pokedex"`
}

type BadgeName string
type RegionName string

const (
	// Gen 1
	RegionKanto  BadgeName = "kanto"
	
	BadgeBoulder BadgeName = "boulder-badge"
	BadgeCascade BadgeName = "cascade-badge"
	BadgeThunder BadgeName = "thunder-badge"
	BadgeRainbow BadgeName = "rainbow-badge"
	BadgeSoul    BadgeName = "soul-badge"
	BadgeMarsh   BadgeName = "marsh-badge"
	BadgeVolcano BadgeName = "volcano-badge"
	BadgeEarth   BadgeName = "earth-badge"
)

var MinIVByBadges map[BadgeName]int = map[BadgeName]int{
	BadgeBoulder: 15,
	BadgeCascade: 17,
	BadgeThunder: 20,
	BadgeRainbow: 22,
	BadgeSoul:    25,
	BadgeMarsh:   27,
	BadgeVolcano: 29,
	BadgeEarth:   30,
}

type BadgeType struct {
	Name   BadgeName  `json:"name"`
	Region RegionName `json:"region"`
}

var BasePayoutTable map[TrainerClass]int = map[TrainerClass]int{
	TrainerPrefix:    80,
	GymLeaderPrefix:  160,
	TournamentPrefix: 160,
	RivalPrefix:      160,
	// Add more as needed
}

var BlackOutPayoutTable map[int]int = map[int]int{0: 8, 1: 16, 2: 24, 3: 36, 4: 48, 5: 64, 6: 80, 7: 100, 8: 120}

// Utils
func GetPrizeMoney(class TrainerClass, party []*Pokemon) int {
	highestLevel := 0
	for _, pokemon := range party {
		if pokemon.Level > highestLevel {
			highestLevel = pokemon.Level
		}
	}
	return BasePayoutTable[class] * highestLevel
}

func GetMoneyLost(user *User) int {
	numBadges := 0
	if user.Stats != nil && user.Stats.Badges != nil {
		numBadges = len(user.Stats.Badges)
	}

	highestLevel := 0
	for _, pokemon := range user.Party {
		if pokemon.Level > highestLevel {
			highestLevel = pokemon.Level
		}
	}
	return BlackOutPayoutTable[numBadges] * highestLevel
}

func (t *TrainerSave) ToTrainer() *Trainer {
	var party []*Pokemon
	for _, pokemon := range t.Party {
		party = append(party, pokemon.ToPokemon())
	}

	return &Trainer{
		BaseTrainer: BaseTrainer{
			Name:  t.Name,
			Party: party,
			Bag:   t.Bag.ToItemMap(),
		},
		Type:      t.Type,
		Rewards:   t.Rewards,
		ImagePath: t.ImagePath,
	}
}

func (t *Trainer) ToTrainerSave() *TrainerSave {
	var party []*PokemonSave
	for _, pokemon := range t.Party {
		party = append(party, pokemon.ToPokemonSave())
	}

	return &TrainerSave{
		BaseTrainerSave: BaseTrainerSave{
			Name:  t.Name,
			Party: party,
			Bag:   t.Bag.ToItemMapSave(),
		},
		Type:      t.Type,
		Rewards:   t.Rewards,
		ImagePath: t.ImagePath,
	}
}

func (u *UserSave) ToUser() *User {
	var party []*Pokemon
	for _, pokemon := range u.Party {
		party = append(party, pokemon.ToPokemon())
	}

	return &User{
		BaseTrainer: BaseTrainer{
			Name:  u.Name,
			Party: party,
			Bag:   u.Bag.ToItemMap(),
		},
		Stats: u.Stats,
		Money: u.Money,
	}
}

func (u *User) ToUserSave() *UserSave {
	var party []*PokemonSave
	for _, pokemon := range u.Party {
		party = append(party, pokemon.ToPokemonSave())
	}

	return &UserSave{
		BaseTrainerSave: BaseTrainerSave{
			Name:  u.Name,
			Party: party,
			Bag:   u.Bag.ToItemMapSave(),
		},
		Stats: u.Stats,
		Money: u.Money,
	}
}

func (u *User) GetBagItemCount() int {
	count := 0
	for _, b := range u.Bag {
		count += b.Count
	}
	return count
}
