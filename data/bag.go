package data

type Item struct {
	Count       int          `json:"count"`
	Category    ItemCategory `json:"category"`
	Cost        int          `json:"cost"`
	Attributes  int          `json:"attributes"`
	Description string       `json:"description"`
}

type ItemName string

const (
	Potion      ItemName = "potion"
	SuperPotion ItemName = "super-potion"
	HyperPotion ItemName = "hyper-potion"
	PokeBall    ItemName = "poke-ball"
	SuperBall   ItemName = "super-ball"
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

func AddItemToBag(bag ItemMap, itemNameToAdd ItemName, itemToAdd Item) {
	if item, exists := bag[itemNameToAdd]; !exists {
		bag[itemNameToAdd] = itemToAdd
	} else {
		item.Count += itemToAdd.Count
		bag[itemNameToAdd] = item
	}
}
