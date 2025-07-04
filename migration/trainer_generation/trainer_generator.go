package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/pokemon"
	"github.com/sanchitdeora/PokeSim/utils"
)

func main() {
	// flags
	name := flag.String("n", "Trainer", "Trainer name")
	isGymLeader := flag.Bool("g", false, "Indicate Gym Leader")
	badge := flag.String("badge-name", "", "Badge name (required if -g is present)")
	region := flag.String("region", "", "Region (required if -g is present)")
	pokemonList := flag.String("p", "", "Pokémon list as ID:Level pairs, e.g. '95:14,74:12'")
	flag.Parse()

	// Get trainer name from command-line argument
	slog.Info("Args", "args", flag.Args(), "isGymLeader", *isGymLeader, "badge", *badge, "region", *region, "pokemonList", *pokemonList)

	trainer := &data.Trainer{
		BaseTrainer: data.BaseTrainer{
			Name: *name,
			// default add 2 potions in bag
			Bag: data.ItemMapSave{data.Potion: 2}.ToItemMap(),
		},
		// default trainer
		Type:      data.TrainerPrefix,
		Rewards:   &data.Rewards{},
		ImagePath: "assets/trainer/img/user_trainer_avatar.png",
	}

	if *isGymLeader {
		if *badge == "" || *region == "" {
			panic("Badge name and region are required if -g is present")
		}
		trainer.Rewards.Badge.Name = data.BadgeName(*badge)
		trainer.Rewards.Badge.Region = data.RegionName(*region)

		trainer.Type = data.GymLeaderPrefix
		trainer.ImagePath = fmt.Sprintf("assets/trainer/gym/img/gym_trainer_%s.png", strings.ToLower(trainer.Name))
	}

	// add pokemons
	var party []*data.Pokemon
	if *pokemonList != "" {
		pairs := strings.Split(*pokemonList, ",")

		if len(pairs) > 6 {
			log.Fatalf("Too many Pokémon. Max allowed is 6.")
		}

		minIV := 0
		if *badge != "" {
			minIV = data.MinIVByBadges[data.BadgeName(*badge)]
		}
		slog.Info("Min IV", "minIV", minIV)

		for _, pair := range pairs {
			parts := strings.Split(pair, ":")
			if len(parts) != 2 {
				log.Fatalf("Invalid pair format: %v (expected ID:Level)", pair)
			}

			id, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				log.Fatalf("Invalid Pokémon ID in: %v", pair)
			}
			level, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				log.Fatalf("Invalid level in: %v", pair)
			}

			basePokemon, _ := data.GetBasePokemonByID(id)
			party = append(party, pokemon.GeneratePokemon(basePokemon, level, minIV))
		}
	}

	slog.Info("Generated trainer", "trainer", trainer, "rewards", *trainer.Rewards)
	trainer.Party = party

	// write to file
	slog.Info("Writing trainer to file", "name", trainer.Name)
	err := utils.WriteJsonToFile(fmt.Sprintf("assets/trainer/gym/%s.json", strings.ToLower(trainer.Name)), trainer.ToTrainerSave())
	if err != nil {
		panic(err)
	}
}
