package alttprclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/mgius/bps"
)

const (
	baseurl             = "https://alttpr.com"
	randomizer_endpoint = "api/randomizer"
	settings_endpoint   = "randomizer/settings"
	bps_endpoint        = "bps"
)

type Client struct {
	// Right now, this is just an empty struct, but it'll be nice to be able to
	// add things to this in the future
}

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

func decodeRandomizerResponse(body io.Reader) (decoded Randomizer, err error) {
	err = json.NewDecoder(body).Decode(&decoded)
	if err != nil {
		err = fmt.Errorf("Error during JSON decode: %w", err)
		return
	}

	return

}

type CustomizationConfig struct {
	// Eventually this is where thigns like heart speed and sprite config will live
}

type CrystalsConfig struct {
	Ganon string `json:"ganon"`
	Tower string `json:"tower"`
}

type EnemizerConfig struct {
	BossShuffle  string `json:"boss_shuffle"`
	Damage       string `json:"enemy_damage"`
	Health       string `json:"enemy_health"`
	EnemyShuffle string `json:"enemy_shuffle"`
}

type ItemConfig struct {
	Functionality string `json:"functionality"`
	Pool          string `json:"pool"`
}

// What you submit to the server
type RandomizerConfig struct {
	Accessibility string         `json:"accessibility"`
	Crystals      CrystalsConfig `json:"crystals"`
	DungeonItems  string         `json:"dungeon_items"`
	Entrances     string         `json:"entrances"`
	Enemizer      EnemizerConfig `json:"enemizer"`
	Glitches      string         `json:"glitches"`
	Goal          string         `json:"goal"`
	Hints         string         `json:"hints"`
	Item          ItemConfig     `json:"item"`
	ItemPlacement string         `json:"item_placement"`
	Lang          string         `json:"lang"`
	Mode          string         `json:"mode"` // world state
	Spoilers      string         `json:"spoilers"`
	Tournament    bool           `json:"tournament"`
	Weapons       string         `json:"weapons"`
}

type SpoilerResponse struct {
}

type Patch map[string][]byte

type Randomizer struct {
	Logic          string          `json:"logic"`
	Hash           string          `json:"hash"`             // Randomizer hash, this is what's displayed in game
	CurrentRomHash string          `json:"current_rom_hash"` // Hash of the BPS patch to apply to the base rom
	SizeInMB       int             `json:"size"`
	GeneratedTime  string          `json:"generated"`
	Spoiler        SpoilerResponse `json:"spoiler"`
	Patches        []Patch         `json:"patch"` // Series of byte patches to apply to various byte nums
	base_patch     *bps.BPSPatch
	client         *Client
}

// Fetch a new Randomizer patch for the given config
func (client *Client) GetRandomizer(config RandomizerConfig) (response Randomizer, err error) {
	request_url := fmt.Sprintf("%s/%s", baseurl, randomizer_endpoint)

	body, err := json.Marshal(config)
	if err != nil {
		err = fmt.Errorf("Error serializing config to json: %w", err)
		return
	}

	resp, err := http.Post(request_url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		err = fmt.Errorf("Error during HTTP POST: %w", err)
		return
	}
	defer resp.Body.Close()

	response, err = decodeRandomizerResponse(resp.Body)
	if err != nil {
		err = fmt.Errorf("Error decoding randomizer response: %w", err)
		return
	}

	response.client = client

	return
}

func (client *Client) GetBasePatch(hash string) (data []byte, err error) {
	request_url := fmt.Sprintf("%s/%s/%s.bps", baseurl, bps_endpoint, hash)

	resp, err := http.Get(request_url)
	if err != nil {
		err = fmt.Errorf("Error on base patch get request: %w", err)
		return
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)

	if err != nil {
		err = fmt.Errorf("Error when reading base patched response body: %w", err)
		return
	}

	return

}

func (patch *Randomizer) CreatePatchedROM(customization CustomizationConfig, sourcefile *os.File) (data []byte, err error) {
	// TODO: cache the BPS file
	err = patch.PopulateBasePatch()
	if err != nil {
		return
	}

	base_patched_rom, err := patch.base_patch.PatchSourceFile(sourcefile)
	if err != nil {
		err = fmt.Errorf("Failed to patch source file: %w", err)
		return
	}

	data = make([]byte, patch.SizeInMB*1024*1024)
	copy(data, base_patched_rom)

	for _, p := range patch.Patches {
		for start_byte_s, p_bytes := range p {
			var start_byte int

			patch_len := len(p_bytes)

			start_byte, err = strconv.Atoi(start_byte_s)
			if err != nil {
				err = fmt.Errorf("Unable to parse %s as an int during randomization patches: %w", start_byte_s, err)
				return
			}
			copy(data[start_byte:start_byte+patch_len], p_bytes)
		}
	}

	return

}

func (patch *Randomizer) PopulateBasePatch() error {
	if patch.base_patch != nil {
		return nil
	}

	bps_data, err := patch.client.GetBasePatch(patch.CurrentRomHash)
	if err != nil {
		return fmt.Errorf("Error fetching BPS base patch: %w", err)
	}
	base_patch, err := bps.FromBytes(bps_data)
	if err != nil {
		return fmt.Errorf("Error when parsing BPS base patch: %w", err)
	}

	patch.base_patch = &base_patch
	return nil
}

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
