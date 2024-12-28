package usermanagement

import (
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/errors"
	"github.com/sanchitdeora/PokeSim/utils"
)

//go:generate mockgen -build_flags=--mod=mod -destination=mocks/mock_user_manager.go -package=mock_user_manager github.com/sanchitdeora/PokeSim/usermanagement UserManager
type UserManager interface {
	SaveUser() error
	GetUser() *data.User
	PostBattleUpdate(user *data.User, report *data.Result) error
	PostWildUpdate(user *data.User, win bool, caught *data.Pokemon) error
	StatUpdate(res data.Result)
}

type UserOpts struct {
	SavedUserPath string
}

type UserImpl struct {
	opts UserOpts
	user *data.User
}

func NewUserService(opts UserOpts) UserManager {
	user, err := LoadUser(opts.SavedUserPath)
	if err != nil {
		slog.Error("could not load user", "error", err)
		return nil
	}
	return &UserImpl{
		opts: opts,
		user: user,
	}
}

func LoadUser(filepath string) (*data.User, error) {
	var savedUser data.UserSave
	var err error

	if utils.CheckPathExists(filepath) {
		savedUser, err = utils.ReadJsonFromFile[data.UserSave](filepath)
		if err != nil {
			slog.Error("could not read from saved file", "error", err)
			return nil, errors.ErrCouldNotReadFromFile
		}
	} else {
		slog.Error("file does not exist...")
		return nil, errors.ErrFileDoesNotExist
	}
	return savedUser.ToUser(), nil
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

	if err := u.SaveUser(); err != nil {
		slog.Error("error updating user after trainer battle", "error", err)
	}
}

func (u *UserImpl) SaveUser() error {
	err := utils.WriteJsonToFile(u.opts.SavedUserPath, u.user.ToUserSave())
	if err != nil {
		slog.Error("could not write to saved file", "error", err)
		return errors.ErrCouldNotReadFromFile
	}

	return nil
}

func (u *UserImpl) GetUser() *data.User {
	return u.user
}

func (u *UserImpl) addItemToBag(bag data.ItemMap, itemName data.ItemName, item data.Item) {
	if itemFound, exists := bag[itemName]; !exists {
		bag[itemName] = item
	} else {
		itemFound.Count += item.Count
		bag[itemName] = itemFound
	}
}

// v1 code
func (u *UserImpl) PostBattleUpdate(user *data.User, result *data.Result) error {
	user.Stats.Battles++

	if result.UserWin {
		user.Stats.Wins++
		if result.BadgeEarned.Name != "" {
			user.Stats.Badges = append(user.Stats.Badges, result.BadgeEarned)
		}

		if len(result.BonusItems) > 0 {
			for itemName, item := range result.BonusItems {
				u.addItemToBag(u.user.Bag, itemName, item)
			}
		}

		user.Money += result.Money

	} else {
		user.Stats.Losses++

		if (user.Money - result.Money) < 0 {
			user.Money = 0
		} else {
			user.Money -= result.Money
		}
	}

	if err := u.SaveUser(); err != nil {
		slog.Error("error updating user after trainer battle", "error", err)
		return err
	}

	return nil
}

func (u *UserImpl) PostWildUpdate(user *data.User, win bool, caught *data.Pokemon) error {
	user.Stats.Battles++
	if win {
		user.Stats.Wins++
		if caught != nil {
			user.Stats.Catches++
			//TODO: update pokedex with caught.
		}
	} else {
		user.Stats.Losses++
	}

	if err := u.SaveUser(); err != nil {
		slog.Error("error updating user after wild pokemon battle", "error", err)
		return err
	}

	return nil
}