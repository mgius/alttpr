package alttpr

import (
	"fmt"
	"os"
)

func main() {
	client := Client{}
	randomizer, err := client.GetRandomizer(DefaultRandomizerConfig())
	if err != nil {
		panic(err.Error())
	}

	base_rom, _ := os.Open("Zelda.sfc")

	patched_bytes, err := randomizer.CreatePatchedROM(CustomizationConfig{}, base_rom)
	if err != nil {
		panic(err.Error())
	}

	outfile, err := os.Create("ZeldaPatched.sfc")
	if err != nil {
		panic(err.Error())
	}
	defer outfile.Close()

	outfile.Write(patched_bytes)

	fmt.Println("Hello World")
}
