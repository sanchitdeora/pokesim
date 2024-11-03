package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sanchitdeora/PokeSim/battle"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
)

type BattleArena struct {
	*GuiOpts
	Battle             battle.BattleIFace
	BattleChan         chan<- *data.BattleInput
	BattleInputButtons *fyne.Container
}

func LoadBattleScreen(opts *GuiOpts, battleChan chan<- *data.BattleInput, battle battle.BattleIFace) fyne.CanvasObject {
	battleArena := &BattleArena{
		GuiOpts:            opts,
		Battle:             battle,
		BattleInputButtons: container.NewGridWithColumns(1), // Initialize with expected layout
		BattleChan:         battleChan,
	}

	// Populate initial buttons in the full space of BattleActionButtons
	battleArena.UpdateBattleInputButtons(battleArena.getDefaultBattleInputButtons())

	// Arrange the sections with opponent, user, and action areas
	battleArenaContainer := container.NewGridWithRows(3,
		getOpponentPokemonInfo(battle.GetOpponentActivePokemon()),
		getUserPokemonInfo(battle.GetUserActivePokemon()),
		battleArena.BattleInputButtons,
	)

	label := widget.NewLabel("Battle Arena")
	battleScreen := container.NewBorder(label, nil, nil, nil,
		addBorder(battleArenaContainer),
	)

	go battleArena.InitiateBattle()
	go battleArena.GuiOpts.LogListener()

	return battleScreen
}

func (b *BattleArena) InitiateBattle() {
	b.Battle.InitiateBattleSequence()
	close(b.BattleChan)
}

func getOpponentPokemonInfo(inBattlePokemon *data.InBattlePokemon) fyne.CanvasObject {

	// Opponent Pokémon Section
	opponentName := widget.NewLabel(utils.ToCapitalizeFirstLetterOfEachWord(inBattlePokemon.Pokemon.Name))
	opponentLevel := widget.NewLabel(fmt.Sprintf("Lv. %v", inBattlePokemon.Pokemon.Level))
	opponentHPBar := widget.NewProgressBar()
	opponentHPBar.Max = float64(inBattlePokemon.Pokemon.Stats.HP.Value)
	opponentHPBar.SetValue(float64(inBattlePokemon.BattleHP)) // Example HP, 80/100

	opponentInfoBox := container.NewGridWithRows(3,
		container.NewGridWithColumns(2,
			opponentName, opponentLevel,
		),
		opponentHPBar,
	)
	opponentImage := canvas.NewImageFromFile("C:\\Projects\\Go-projects\\src\\PokéSim\\assets\\pokemon\\img\\1.png") // Update path as needed
	opponentImage.FillMode = 2

	opponentSection := container.NewPadded(container.NewGridWithColumns(2,
		opponentInfoBox,
		opponentImage,
	))

	return container.NewVBox(opponentSection, widget.NewSeparator())
}

func getUserPokemonInfo(inBattlePokemon *data.InBattlePokemon) fyne.CanvasObject {
	// User Pokémon Section
	userName := widget.NewLabel(utils.ToCapitalizeFirstLetterOfEachWord(inBattlePokemon.Pokemon.Name))
	userLevel := widget.NewLabel(fmt.Sprintf("Lv. %v", inBattlePokemon.Pokemon.Level))

	userHPBar := widget.NewProgressBar()
	userHPBar.Max = float64(inBattlePokemon.Pokemon.Stats.HP.Value)
	userHPBar.SetValue(float64(inBattlePokemon.BattleHP))

	userExpBar := widget.NewProgressBar()
	userExpBar.Max = float64(100)
	userExpBar.SetValue(80) // Example HP, 80/100

	userInfoBox := container.NewVBox(
		container.NewGridWithColumns(2,
			userName,
			userLevel,
		),
		container.NewGridWithRows(2,
			userHPBar,
			widget.NewLabel(fmt.Sprintf("%v/ %v", inBattlePokemon.BattleHP, inBattlePokemon.Pokemon.Stats.HP.Value)),
		),
		container.NewGridWithColumns(2, widget.NewLabel("EXP"), userExpBar),
	)

	userImage := canvas.NewImageFromFile("C:\\Projects\\Go-projects\\src\\PokéSim\\assets\\pokemon\\img\\back\\1.png") // Update path as needed
	userImage.FillMode = 2

	userSection := container.NewGridWithColumns(2,
		userImage,
		userInfoBox,
	)

	return container.NewVBox(userSection, widget.NewSeparator())
}

func (b *BattleArena) getDefaultBattleInputButtons() *fyne.Container {
	return container.NewGridWithColumns(1,
		container.NewGridWithColumns(2, // This layout will be forced onto BattleActionButtons
			widget.NewButton("Attack", func() { b.UpdateBattleInputButtons(b.handleAttackSelection()) }),
			widget.NewButton("Switch", func() {}),
			widget.NewButton("Bag", func() {}),
			widget.NewButton("Run", func() {}),
		),
		container.New(layout.NewCenterLayout()),
	)
}

func (b *BattleArena) handleAttackSelection() fyne.CanvasObject {
	backButton := widget.NewButton("Back", func() { b.UpdateBattleInputButtons(b.getDefaultBattleInputButtons()) })
	backButton.Importance = widget.HighImportance

	return container.NewGridWithColumns(1,
		container.NewGridWithColumns(2, // Same column count to maintain consistency
			widget.NewButton(utils.ToCapitalizeFirstLetterOfEachWord(b.Battle.GetUserActivePokemon().Pokemon.Moveset.Move1.Name), func() { b.HandleAttack(b.Battle.GetUserActivePokemon().Pokemon.Moveset.Move1) }),
			widget.NewButton(utils.ToCapitalizeFirstLetterOfEachWord(b.Battle.GetUserActivePokemon().Pokemon.Moveset.Move2.Name), func() { b.HandleAttack(b.Battle.GetUserActivePokemon().Pokemon.Moveset.Move2) }),
			widget.NewButton(utils.ToCapitalizeFirstLetterOfEachWord(b.Battle.GetUserActivePokemon().Pokemon.Moveset.Move3.Name), func() { b.HandleAttack(b.Battle.GetUserActivePokemon().Pokemon.Moveset.Move3) }),
			widget.NewButton(utils.ToCapitalizeFirstLetterOfEachWord(b.Battle.GetUserActivePokemon().Pokemon.Moveset.Move4.Name), func() { b.HandleAttack(b.Battle.GetUserActivePokemon().Pokemon.Moveset.Move4) }),
		),
		container.New(layout.NewCenterLayout(), backButton),
	)
}

func (b *BattleArena) HandleAttack(move *data.Moves) {
	// slog.Info("move chose", "move", move.Name)

	input := &data.BattleInput{
		Type:           data.Attack,
		CurrentPokemon: b.Battle.GetUserActivePokemon(),
		Target:         b.Battle.GetOpponentActivePokemon(),
		Move:           move,
		Item:           nil,
		IsUser:         true,
	}

	b.BattleChan <- input
	// b.UpdateBattleInputButtons(b.getDefaultBattleInputButtons())
}

// UpdateBattleInputButtons updates the BattleActionButtons container with new content
func (b *BattleArena) UpdateBattleInputButtons(newContent fyne.CanvasObject) {
	b.BattleInputButtons.Objects = []fyne.CanvasObject{newContent}
	b.BattleInputButtons.Refresh()
}
