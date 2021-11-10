package alttpr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseurl             = "https://alttpr.com"
	randomizer_endpoint = "api/randomizer"
	settings_endpoint   = "randomizer/settings"
	bps_endpoint        = "bps"
)

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

type Client struct {
	// Right now, this is just an empty struct, but it'll be nice to be able to
	// add things to this in the future
}

func decodeRandomizerResponse(body io.Reader) (decoded Randomizer, err error) {
	err = json.NewDecoder(body).Decode(&decoded)
	if err != nil {
		err = fmt.Errorf("error during JSON decode: %w", err)
		return
	}

	return

}

// Fetch a new Randomizer patch for the given config
func (client *Client) GetRandomizer(config RandomizerConfig) (response Randomizer, err error) {
	request_url := fmt.Sprintf("%s/%s", baseurl, randomizer_endpoint)

	body, err := json.Marshal(config)
	if err != nil {
		err = fmt.Errorf("error serializing config to json: %w", err)
		return
	}

	resp, err := http.Post(request_url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		err = fmt.Errorf("error during HTTP POST: %w", err)
		return
	}
	defer resp.Body.Close()

	response, err = decodeRandomizerResponse(resp.Body)
	if err != nil {
		err = fmt.Errorf("error decoding randomizer response: %w", err)
		return
	}

	response.client = client

	return
}

func (client *Client) GetBasePatch(hash string) (data []byte, err error) {
	request_url := fmt.Sprintf("%s/%s/%s.bps", baseurl, bps_endpoint, hash)

	resp, err := http.Get(request_url)
	if err != nil {
		err = fmt.Errorf("error on base patch get request: %w", err)
		return
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)

	if err != nil {
		err = fmt.Errorf("error when reading base patched response body: %w", err)
		return
	}

	return

}
