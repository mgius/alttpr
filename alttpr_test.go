package alttpr

import (
	"os"
	"testing"

	"github.com/mgius/bps"
)

func TestReadRandomizerJSON(t *testing.T) {
	resp, _ := os.Open("test/randomizer_response.json")

	decoded, err := decodeRandomizerResponse(resp)

	if err != nil {
		t.Errorf(err.Error())
	}

	if decoded.CurrentRomHash != "7f2e1606616492d7dfb589e8dfb70027" {
		t.Errorf("Rom hash does not match")
	}

}

func TestGetBasePatch(t *testing.T) {
	client := Client{}
	patch_bytes, err := client.GetBasePatch("7f2e1606616492d7dfb589e8dfb70027")
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = bps.FromBytes(patch_bytes)
	if err != nil {
		t.Errorf(err.Error())
	}

}

func TestEndToEnd(t *testing.T) {
	base_rom, err := os.Open("test/Zelda.sfc")
	if err != nil {
		t.Skipf("Could not read test/Zelda.sfc.  Skipping this test")
	}

	client := Client{}
	randomizer, err := client.GetRandomizer(DefaultRandomizerConfig())
	if err != nil {
		t.Errorf(err.Error())
	}

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
}
