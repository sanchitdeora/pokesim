package data

type ItemMap map[ItemName]Item

// save models
type BaseTrainerSave struct {
	Name  string         `json:"name"`
	Party []*PokemonSave `json:"party"`
	Bag   ItemMap        `json:"bag"`
}

type UserSave struct {
	BaseTrainerSave
	Stats *TrainerStats `json:"stats"`
	Money int
}

type TrainerSave struct {
	BaseTrainerSave
	TrainerID string
	Type      TrainerClass
	Rewards   *Rewards
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
	Type    TrainerClass
	Rewards *Rewards
}

type Rewards struct {
	Items ItemMap   `json:"items"`
	Badge BadgeType `json:"badge_type,omitempty"`
}

type TrainerClass string

const (
	TrainerPrefix    TrainerClass = "Trainer"
	GymLeaderPrefix  TrainerClass = "Gym Leader"
	TournamentPrefix TrainerClass = "Tournament Trainer"
	RivalPrefix      TrainerClass = "Rival"
	WildPrefix      TrainerClass = "Wild"
)

type TrainerStats struct {
	Badges  []BadgeType `json:"badges,omitempty"`
	Battles int         `json:"battles"`
	Wins    int         `json:"wins"`
	Losses  int         `json:"losses"`
	Catches int         `json:"catches"`
	PokeDEX int         `json:"pokedex"`
}

var BasePayoutTable map[TrainerClass]int = map[TrainerClass]int{
	TrainerPrefix:    80,
	GymLeaderPrefix:  160,
	TournamentPrefix: 160,
	RivalPrefix:      160,
	// Add more as needed
}

var BlackOutPayoutTable map[int]int = map[int]int{0: 8, 1: 16, 2: 24, 3: 36, 4: 48, 5: 64, 6: 80, 7: 100, 8: 120}

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
			Bag:   t.Bag,
		},
		Type:    t.Type,
		Rewards: t.Rewards,
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
			Bag:   u.Bag,
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
			Bag:   u.Bag,
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
