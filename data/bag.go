package data

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

var AllItems = []ItemName{
	Potion,
	SuperPotion,
	HyperPotion,
	PokeBall,
	GreatBall,
	UltraBall,
	MasterBall,
}

// Store
type StoreItem struct {
	Name ItemName `json:"name"`
	Item Item     `json:"item"`
}

var StoreItems = []StoreItem{
	{
		Name: Potion,
		Item: Item{
			Count:       1,
			Category:    MedicalItems,
			CostPrice:   200,
			SellPrice:   50,
			Attribute:   20.0,
			Description: "Heals 20 HP",
		},
	},
	{
		Name: SuperPotion,
		Item: Item{
			Count:       1,
			Category:    MedicalItems,
			CostPrice:   700,
			SellPrice:   175,
			Attribute:   60.0,
			Description: "Heals 60 HP",
		},
	},
	{
		Name: HyperPotion,
		Item: Item{
			Count:       1,
			Category:    MedicalItems,
			CostPrice:   1500,
			SellPrice:   375,
			Attribute:   120.0,
			Description: "Heals 120 HP",
		},
	},
	{
		Name: PokeBall,
		Item: Item{
			Count:       1,
			Category:    PokeBalls,
			CostPrice:   200,
			SellPrice:   50,
			Attribute:   1.0,
			Description: "Catching a Pokemon",
		},
	},
	{
		Name: GreatBall,
		Item: Item{
			Count:       1,
			Category:    PokeBalls,
			CostPrice:   600,
			SellPrice:   150,
			Attribute:   1.5,
			Description: "High Chance of Catching a Pokemon",
		},
	},
	{
		Name: UltraBall,
		Item: Item{
			Count:       1,
			Category:    PokeBalls,
			CostPrice:   800,
			SellPrice:   200,
			Attribute:   2.0,
			Description: "Higher Chance of Catching a Pokemon",
		},
	},
	{
		Name: MasterBall,
		Item: Item{
			Count:       1,
			Category:    PokeBalls,
			CostPrice:   10,
			SellPrice:   5,
			Attribute:   255.0,
			Description: "Guarantees Catching a Pokemon",
		},
	},
}

func GetNameFromItem(item Item) ItemName {
	for _, i := range StoreItems {
		if i.Item.Description == item.Description && i.Item.Attribute == item.Attribute {
			return i.Name
		}
	}
	return ""
}

func GetItemFromName(name ItemName) Item {
	for _, i := range StoreItems {
		if i.Name == name {
			return i.Item
		}
	}
	return Item{}
}

func (i ItemMap) ToItemMapSave() ItemMapSave {
	itemMapSave := make(map[ItemName]int)
	for name, item := range i {
		itemMapSave[name] = item.Count
	}
	return itemMapSave
}

func (i ItemMapSave) ToItemMap() ItemMap {
	itemMap := make(map[ItemName]Item)
	for name, count := range i {

		item := GetItemFromName(name)

		item.Count = count
		itemMap[name] = item
	}
	return itemMap
}
