package usermanagement

import (
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/errors"
	"github.com/sanchitdeora/PokeSim/utils"
)

type UserManager interface {
	SaveUser(user *data.User) error
	GetUser() *data.User
	PostBattleUpdate(user *data.User, report *data.Result) error
	PostWildUpdate(user *data.User, win bool, caught *data.Pokemon) error
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

func (u *UserImpl) PostBattleUpdate(user *data.User, result *data.Result) error {
	user.Stats.Battles++

	if result.UserWin {
		user.Stats.Wins++
		// if result.BadgeEarned != nil {
		// 	user.Stats.Badges = append(user.Stats.Badges, *result.BadgeEarned)
		// }

		if len(result.BonusItems) > 0 {
			for itemName, item := range result.BonusItems {
				data.AddItemToBag(user.Bag, itemName, item)
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

	if err := u.SaveUser(user); err != nil {
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

	if err := u.SaveUser(user); err != nil {
		slog.Error("error updating user after wild pokemon battle", "error", err)
		return err
	}

	return nil
}

func (u *UserImpl) SaveUser(user *data.User) error {
	err := utils.WriteJsonToFile(u.opts.SavedUserPath, user.ToUserSave())
	if err != nil {
		slog.Error("could not write to saved file", "error", err)
		return errors.ErrCouldNotReadFromFile
	}

	return nil
}

func (u *UserImpl) GetUser() *data.User {
	return u.user
}
