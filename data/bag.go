package data

import "log/slog"

type Item struct {
	Count       int          `json:"count"`
	Category    ItemCategory `json:"category"`
	CostPrice   int          `json:"cost_price"`
	SellPrice   int          `json:"sell_price"`
	Attribute   float64      `json:"attribute"`
	Description string       `json:"description"`
}

type ItemName string

const (
	Potion      ItemName = "potion"
	SuperPotion ItemName = "super-potion"
	HyperPotion ItemName = "hyper-potion"
	PokeBall    ItemName = "poke-ball"
	GreatBall   ItemName = "super-ball"
	UltraBall   ItemName = "ultra-ball"
	MasterBall  ItemName = "master-ball"
)

type ItemCategory string

const (
	MedicalItems ItemCategory = "MedicalItems"
	PokeBalls    ItemCategory = "PokeBalls"
)

type BadgeType struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

var AllItems = []ItemName{
	Potion,
	SuperPotion,
	HyperPotion,
	PokeBall,
	GreatBall,
	UltraBall,
	MasterBall,
}

var ItemStore = map[ItemName]Item{
	Potion: {
		Category:    MedicalItems,
		CostPrice:   200,
		SellPrice:   50,
		Attribute:   20.0,
		Description: "Heals 20 HP",
	},
	SuperPotion: {
		Category:    MedicalItems,
		CostPrice:   400,
		SellPrice:   200,
		Attribute:   20.0,
		Description: "Heals 20 HP",
	},
	HyperPotion: {
		Category:    MedicalItems,
		CostPrice:   1000,
		SellPrice:   500,
		Attribute:   20.0,
		Description: "Heals 20 HP",
	},
	PokeBall: {
		Category:    PokeBalls,
		CostPrice:   10,
		SellPrice:   5,
		Attribute:   1.0,
		Description: "Catch Pokemon",
	},
	GreatBall: {
		Category:    PokeBalls,
		CostPrice:   10,
		SellPrice:   5,
		Attribute:   1.5,
		Description: "Great Chance of Catching a Pokemon",
	},
	UltraBall: {
		Category:    PokeBalls,
		CostPrice:   10,
		SellPrice:   5,
		Attribute:   60,
		Description: "Higher chance of catching a Pokemon",
	},
	MasterBall: {
		Category:    PokeBalls,
		CostPrice:   10,
		SellPrice:   5,
		Attribute:   100,
		Description: "Guarantees catching a Pokemon",
	},
}

func GetItemNameFromItem(item Item) ItemName {
	for name, i := range ItemStore {
		if i.Description == item.Description && i.Attribute == item.Attribute {
			return name
		}
	}
	return ""
}

func (i ItemMap) ToItemMapSave() ItemMapSave {
	itemMapSave := make(map[ItemName]int)
	for name, item := range i {
		itemMapSave[name] = item.Count
		slog.Info("item", "name", name, "item", item)
	}
	return itemMapSave
}

func (i ItemMapSave) ToItemMap() ItemMap {
	itemMap := make(map[ItemName]Item)
	for name, count := range i {
		item := ItemStore[name]
		item.Count = count
		itemMap[name] = item
		slog.Info("item", "name", name, "count", count, "item", item)
	}
	return itemMap
}