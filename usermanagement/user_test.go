package usermanagement_test

import (
	"testing"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/stretchr/testify/assert"
)

const (
	TestUserPath  = "C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\user\\test_user.json"
	TestUserPath1 = "C:\\Projects\\Go-projects\\src\\PokéSim\\saved\\user.json"
)

func createUserService(savedPath string) usermanagement.UserManager {
	return usermanagement.NewUserService(usermanagement.UserOpts{SavedUserPath: savedPath})
}

// LoadUser

func TestLoadUser(t *testing.T) {
	userService := createUserService(TestUserPath)

	user := userService.GetUser()

	assert.Equal(t, "John Cena", user.Name)
	assert.Equal(t, "bulbasaur", user.Party[0].Name)
	assert.Equal(t, "squirtle", user.Party[1].Name)
}

func TestPostBattleUpdate_BattleWon(t *testing.T) {
	userService := createUserService(TestUserPath)

	userBeforeUpdate := userService.GetUser()
	user := userService.GetUser()

	report := &data.Result{
		UserWin: true,
		Money:   100,
		BonusItems: data.ItemMap{
			data.Potion: {Category: data.MedicalItems, Count: 1},
			data.PokeBall: {Category: data.PokeBalls, Count: 1},
		},
		BadgeEarned: data.BadgeType{
			Name: "test",
		},
	}

	err := userService.PostBattleUpdate(user, report)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.Stats.Battles)
	assert.Equal(t, 1, user.Stats.Wins)
	assert.Equal(t, 0, user.Stats.Losses)
	assert.Equal(t, 1, len(user.Stats.Badges))
	assert.Equal(t, 100, user.Money)
	assert.Equal(t, 2, len(user.Bag))
	assert.Equal(t, 3, user.Bag[data.Potion].Count)
	assert.Equal(t, 1, user.Bag[data.PokeBall].Count)

	// clean up after
	userService.SaveUser(userBeforeUpdate)
}

func TestPostBattleUpdate_BattleLost_NoMoney(t *testing.T) {
	userService := createUserService(TestUserPath)

	userBeforeUpdate := userService.GetUser()
	user := userService.GetUser()

	report := &data.Result{
		UserWin: false,
		Money:   100,
	}

	err := userService.PostBattleUpdate(user, report)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.Stats.Battles)
	assert.Equal(t, 0, user.Stats.Wins)
	assert.Equal(t, 1, user.Stats.Losses)
	assert.Equal(t, 0, len(user.Stats.Badges))
	assert.Equal(t, 0, user.Money)
	assert.Equal(t, 1, len(user.Bag))
	assert.Equal(t, 2, user.Bag[data.Potion].Count)

	// clean up after
	userService.SaveUser(userBeforeUpdate)
}

func TestPostBattleUpdate_BattleLost_MoneyLost(t *testing.T) {
	userService := createUserService(TestUserPath)

	userBeforeUpdate := userService.GetUser()
	user := userService.GetUser()
	user.Money = 100
	
	report := &data.Result{
		UserWin: false,
		Money:   100,
	}

	err := userService.PostBattleUpdate(user, report)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.Stats.Battles)
	assert.Equal(t, 0, user.Stats.Wins)
	assert.Equal(t, 1, user.Stats.Losses)
	assert.Equal(t, 0, len(user.Stats.Badges))
	assert.Equal(t, 0, user.Money)
	assert.Equal(t, 1, len(user.Bag))
	assert.Equal(t, 2, user.Bag[data.Potion].Count)

	// clean up after
	userService.SaveUser(userBeforeUpdate)
}


func TestPostWildUpdate(t *testing.T) {
	userService := createUserService(TestUserPath)

	user := userService.GetUser()

	err := userService.PostWildUpdate(user, true, &data.Pokemon{})
	assert.NoError(t, err)
	assert.Equal(t, 1, user.Stats.Battles)
	assert.Equal(t, 1, user.Stats.Wins)
	assert.Equal(t, 1, user.Stats.Catches)

	// clean up after
	user.Stats = &data.TrainerStats{}
	userService.SaveUser(user)
}
