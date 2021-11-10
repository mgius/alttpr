package alttpr

// Return a randomizer config that reflects the defaults provided by ALTTPR
func DefaultRandomizerConfig() RandomizerConfig {
	return RandomizerConfig{
		Accessibility: "items",
		Crystals: CrystalsConfig{
			Tower: "7",
			Ganon: "7",
		},
		DungeonItems: "standard",
		Entrances:    "none",
		Enemizer: EnemizerConfig{
			BossShuffle:  "none",
			Damage:       "default",
			Health:       "default",
			EnemyShuffle: "none",
		},
		Glitches: "none",
		Goal:     "ganon",
		Hints:    "on",
		Item: ItemConfig{
			Functionality: "normal",
			Pool:          "normal",
		},
		ItemPlacement: "advanced",
		Lang:          "en",
		Mode:          "open",
		Spoilers:      "on",
		Tournament:    false,
		Weapons:       "randomized",
	}
}
