package bps

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ALTTPR_BASEURL             = "https://alttpr.com"
	ALTTPR_RANDOMIZER_ENDPOINT = "/api/randomizer"
	ALTTPR_SETTINGS_ENDPOINT   = "/randomizer/settings"
)

type SettingsResponse struct {
	// Presets map[string]string `json:"presets"`
	GlitchesRequired map[string]string `json:"glitches_required"`
}

// What you submit to the server
type RandomizerConfig struct {
	Accessibility string `json:"accessibility"`
	DungeonItems  string `json:"dungeon_items"`
	Entrances     string `json:"entrances"`
	Glitches      string `json:"glitches"`
	Goal          string `json:"goal"`
	Hints         string `json:"hints"`
	ItemPlacement string `json:"item_placement"`
	Lang          string `json:"lang"`
	Mode          string `json:"mode"`
	Spoilers      string `json:"spoilers"`
	Tournament    bool   `json:"tournament"`
	Weapons       string `json:"weapons"`
}

type SpoilerResponse struct {
}

type PatchData []int
type Patch map[string]PatchData

type RandomizerResponse struct {
	Logic          string          `json:"logic"`
	Hash           string          `json:"hash"`
	CurrentRomHash string          `json:"current_rom_hash"`
	SizeInMB       int             `json:"size"`
	GeneratedTime  string          `json:"generated"`
	Spoiler        SpoilerResponse `json:"spoiler"`
	Patches        []Patch         `json:"patch"`
}

func alttpr_randomizer_request() (RandomizerResponse, error) {
	request_url := fmt.Sprintf("%s%s", ALTTPR_BASEURL, ALTTPR_RANDOMIZER_ENDPOINT)
	resp, err := http.Post(request_url, "application/json", nil)
	if err != nil {
		return RandomizerResponse{}, err
	}
	defer resp.Body.Close()

	var decoded RandomizerResponse

	err = json.NewDecoder(resp.Body).Decode(&decoded)

	fmt.Println(decoded)

	return decoded, nil
}

func alttpr_settings_request() (SettingsResponse, error) {

	// Better URL construction method?
	request_url := fmt.Sprintf("%s%s", ALTTPR_BASEURL, ALTTPR_SETTINGS_ENDPOINT)
	resp, err := http.Get(request_url)
	defer resp.Body.Close()

	if err != nil {
		// I'm sure there's a more idiomatic way to do this
		fmt.Println("Error")
		fmt.Println(err)
		return SettingsResponse{}, err
	}

	var decoded SettingsResponse
	// Not OK because err is already defined above
	err = json.NewDecoder(resp.Body).Decode(&decoded)

	fmt.Println(resp)
	fmt.Println(err)
	fmt.Println(decoded)

	return decoded, nil

}

func main() {
	fmt.Println("Hello World")
}
