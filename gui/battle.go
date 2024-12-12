package gui

import (
	"fmt"
	"log/slog"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	battle "github.com/sanchitdeora/PokeSim/battle_v1"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
)

type BattleArena struct {
	*GuiOpts
	Battle             battle.BattleSequence
	BattleChan         chan<- *data.BattleInput
	BattleInputButtons *fyne.Container
}

func NewBattleArena(opts *GuiOpts, battleChan chan<- *data.BattleInput, battle battle.BattleSequence) *BattleArena {
	battleArena := &BattleArena{
		GuiOpts:            opts,
		Battle:             battle,
		BattleInputButtons: container.NewGridWithColumns(1), // Initialize with expected layout
		BattleChan:         battleChan,
	}

	go battleArena.InitiateBattle()
	go battleArena.GuiOpts.LogListener()

	return battleArena
}

func (b *BattleArena) LoadBattleScreen() fyne.CanvasObject {
	// Populate initial buttons in the full space of BattleActionButtons
	b.UpdateBattleInputButtons(b.getDefaultBattleInputButtons())

	// Arrange the sections with opponent, user, and action areas
	battleArenaContainer := container.NewGridWithRows(3,
		getOpponentPokemonInfo(b.Battle.GetUserActive(false)),
		getUserPokemonInfo(b.Battle.GetUserActive(true)),
		b.BattleInputButtons,
	)

	label := widget.NewLabel("Battle Arena")
	return container.NewBorder(label, nil, nil, nil,
		addBorder(battleArenaContainer),
	)
}

func (b *BattleArena) BattleComplete(report *data.Result) fyne.CanvasObject {
	var title string
	var subtitle string
	content := widget.NewLabel("")
	if report.UserWin {
		title = "You won the battle!"
		subtitle = fmt.Sprintf("You earned $%v", report.Money)
		if report.BadgeEarned != nil {
			content.Text = fmt.Sprintf("You also earned a %s badge", report.BadgeEarned.Name)
		}
	} else {
		title = "You lost the battle!"
		subtitle = fmt.Sprintf("You lost $%v", report.Money)
	}

	return container.NewBorder(widget.NewLabel("Battle Complete"), nil, nil, nil,
		container.NewCenter(container.NewVBox(
			widget.NewCard(title, subtitle, content),
		)),
	)
}

func (b *BattleArena) InitiateBattle() {
	report, err := b.Battle.Initiate()
	if err != nil {
		slog.Error("error during battle", err)
		panic("battle threw an error")
	}

	b.UpdateActionContent(b.BattleComplete(report))

	close(b.BattleChan)
}

func getOpponentPokemonInfo(inBattlePokemon *data.BattlePokemon) fyne.CanvasObject {
	slog.Info("DEBUG====Opponent", "pokemonName", inBattlePokemon.Pokemon.Name, "Current HP", inBattlePokemon.BattleHP, "Fainted?", inBattlePokemon.IsFainted)

	// Opponent Pokémon Section
	opponentName := widget.NewLabel(utils.ToCapitalizeFirstLetterOfEachWord(inBattlePokemon.Pokemon.Name))
	opponentLevel := widget.NewLabel(fmt.Sprintf("Lv. %v", inBattlePokemon.Pokemon.Level))

	opponentMaxHP := battle.BattleHPCalculator(&inBattlePokemon.Pokemon.Stats.HP, inBattlePokemon.Pokemon.Level)
	opponentHPBar := widget.NewProgressBar()
	opponentHPBar.Max = opponentMaxHP
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

func getUserPokemonInfo(inBattlePokemon *data.BattlePokemon) fyne.CanvasObject {
	slog.Info("DEBUG====User", "pokemonName", inBattlePokemon.Pokemon.Name, "Current HP", inBattlePokemon.BattleHP, "Fainted?", inBattlePokemon.IsFainted)

	// User Pokémon Section
	userName := widget.NewLabel(utils.ToCapitalizeFirstLetterOfEachWord(inBattlePokemon.Pokemon.Name))
	userLevel := widget.NewLabel(fmt.Sprintf("Lv. %v", inBattlePokemon.Pokemon.Level))

	userMaxHP := battle.BattleHPCalculator(&inBattlePokemon.Pokemon.Stats.HP, inBattlePokemon.Pokemon.Level)
	userHPBar := widget.NewProgressBar()
	userHPBar.Max = userMaxHP
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
			widget.NewLabel(fmt.Sprintf("%v/ %v", inBattlePokemon.BattleHP, math.Round(userMaxHP))),
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
			widget.NewButton("Switch", func() { b.HandleRun() }),
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
			widget.NewButton(utils.ToCapitalizeFirstLetterOfEachWord(b.Battle.GetUserActive(true).Pokemon.Moveset.Move1.Name), func() { b.HandleAttack(b.Battle.GetUserActive(true).Pokemon.Moveset.Move1) }),
			widget.NewButton(utils.ToCapitalizeFirstLetterOfEachWord(b.Battle.GetUserActive(true).Pokemon.Moveset.Move2.Name), func() { b.HandleAttack(b.Battle.GetUserActive(true).Pokemon.Moveset.Move2) }),
			widget.NewButton(utils.ToCapitalizeFirstLetterOfEachWord(b.Battle.GetUserActive(true).Pokemon.Moveset.Move3.Name), func() { b.HandleAttack(b.Battle.GetUserActive(true).Pokemon.Moveset.Move3) }),
			widget.NewButton(utils.ToCapitalizeFirstLetterOfEachWord(b.Battle.GetUserActive(true).Pokemon.Moveset.Move4.Name), func() { b.HandleAttack(b.Battle.GetUserActive(true).Pokemon.Moveset.Move4) }),
		),
		container.New(layout.NewCenterLayout(), backButton),
	)
}

func (b *BattleArena) HandleAttack(move *data.Moves) {
	// slog.Info("move chose", "move", move.Name)

	input := &data.BattleInput{
		Type:           data.Attack,
		Selected: b.Battle.GetUserActive(true),
		Target:         b.Battle.GetUserActive(false),
		Move:           move,
		Item:           nil,
		IsUser:         true,
	}

	b.BattleChan <- input
	b.UpdateActionContent(b.LoadBattleScreen())
	// b.UpdateBattleInputButtons(b.getDefaultBattleInputButtons())
}

func (b *BattleArena) HandleRun() {
	// slog.Info("move chose", "move", move.Name)

	input := &data.BattleInput{
		Type:   data.Run,
		IsUser: true,
	}

	b.BattleChan <- input
	b.UpdateActionContent(b.LoadBattleScreen())
	// b.UpdateBattleInputButtons(b.getDefaultBattleInputButtons())
}

// UpdateBattleInputButtons updates the BattleActionButtons container with new content
func (b *BattleArena) UpdateBattleInputButtons(newContent fyne.CanvasObject) {
	b.BattleInputButtons.Objects = []fyne.CanvasObject{newContent}
	b.BattleInputButtons.Refresh()
}
