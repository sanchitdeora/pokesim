package usermanagement

import (
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/errors"
	"github.com/sanchitdeora/PokeSim/gamestate"
)

//go:generate mockgen -build_flags=--mod=mod -destination=mocks/mock_user_manager.go -package=mock_user_manager github.com/sanchitdeora/PokeSim/usermanagement UserManager
type UserManager interface {
	SaveUser() error
	GetUser() *data.User
	StatUpdate(res data.Result)
	ChangePokemonOrder(pokemon *data.Pokemon, order data.PokemonChangeOrder)
	AddNewPokemonToTeam(pokemon *data.Pokemon)
	UseItem(item *data.Item, count int)
	PurchaseItem(item data.Item) error
	SellItem(item data.Item) error
}

type UserOpts struct {
	GameState gamestate.GameStateManager
}

type UserImpl struct {
	opts UserOpts
	user *data.User
}

func NewUserManager(opts UserOpts) UserManager {
	if opts.GameState == nil {
		slog.Error("invalid opts")
		return nil
	}

	gameState := opts.GameState.Get()
	return &UserImpl{
		opts: opts,
		user: gameState.User,
	}
}

func (u *UserImpl) GetUser() *data.User {
	return u.user
}

func (u *UserImpl) StatUpdate(result data.Result) {
	u.user.Stats.Battles++

	if result.Status == data.Won {
		u.user.Stats.Wins++
		if result.BadgeEarned.Name != "" {
			u.user.Stats.Badges = append(u.user.Stats.Badges, result.BadgeEarned)
		}

		if len(result.BonusItems) > 0 {
			for itemName, item := range result.BonusItems {
				u.addItemToBag(u.user.Bag, itemName, item)
			}
		}

		u.user.Money += result.Money
	} else {
		u.user.Stats.Losses++

		u.user.Money -= result.Money
		if u.user.Money < 0 {
			u.user.Money = 0
		}
	}

	slog.Debug("saving user", "user", u.user, "money won/lost", result.Money)

	// Save User
	u.SaveUser()
}

func (u *UserImpl) SaveUser() error {
	return u.opts.GameState.Save()
}

func (u *UserImpl) PurchaseItem(item data.Item) error {
	cost := item.CostPrice * item.Count
	if cost > u.user.Money {
		return errors.ErrNotEnoughMoney
	}
	u.user.Money -= cost
	u.addItemToBag(u.user.Bag, data.GetNameFromItem(item), item)

	// Save User
	u.SaveUser()

	return nil
}

func (u *UserImpl) SellItem(item data.Item) error {
	cost := item.SellPrice * item.Count
	if item.Count > u.user.Bag[data.GetNameFromItem(item)].Count {
		return errors.ErrItemCountNotEnough
	}

	u.user.Money += cost
	u.UseItem(&item, item.Count)

	return nil
}

func (u *UserImpl) UseItem(item *data.Item, count int) {
	itemName := data.GetNameFromItem(*item)
	itemInBag, exists := u.user.Bag[itemName]
	if !exists {
		slog.Debug("item not found in bag")
		return
	}
	if itemInBag.Count < count {
		slog.Debug("item count not enough")
		return
	}
	itemInBag.Count -= count

	u.user.Bag[itemName] = itemInBag

	// Save User
	u.SaveUser()
}

func (u *UserImpl) ChangePokemonOrder(pokemon *data.Pokemon, order data.PokemonChangeOrder) {
	currentIdx := GetPokemonIndexInParty(u.user.Party, pokemon)
	switchIdx := currentIdx
	if currentIdx < 0 {
		slog.Debug("do nothing!")
		return
	}

	if order == data.ChangeOrderMoveUp && currentIdx > 0 {
		switchIdx = currentIdx - 1
	}
	if order == data.ChangeOrderMoveDown && currentIdx < len(u.user.Party)-1 {
		switchIdx = currentIdx + 1
	}

	u.user.Party[currentIdx], u.user.Party[switchIdx] = u.user.Party[switchIdx], u.user.Party[currentIdx]

	// Save User
	u.SaveUser()
}

func (u *UserImpl) AddNewPokemonToTeam(pokemon *data.Pokemon) {
	if len(u.user.Party) < 6 {
		u.user.Party = append(u.user.Party, pokemon)
	} else {
		// add to box
		slog.Debug("team full")
		return
	}

	// Save User
	u.SaveUser()
}

func (u *UserImpl) addItemToBag(bag data.ItemMap, itemName data.ItemName, item data.Item) {
	if itemFound, exists := bag[itemName]; !exists {
		bag[itemName] = item
	} else {
		itemFound.Count += item.Count
		bag[itemName] = itemFound
	}
}
