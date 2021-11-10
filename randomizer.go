package alttpr

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mgius/bps"
)

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

func (patch *Randomizer) CreatePatchedROM(customization CustomizationConfig, sourcefile *os.File) (data []byte, err error) {
	// TODO: cache the BPS file
	// TODO: work on bytes instead of os.File
	err = patch.PopulateBasePatch()
	if err != nil {
		return
	}

	base_patched_rom, err := patch.base_patch.PatchSourceFile(sourcefile)
	if err != nil {
		err = fmt.Errorf("failed to patch source file: %w", err)
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
				err = fmt.Errorf("unable to parse %s as an int during randomization patches: %w", start_byte_s, err)
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
		return fmt.Errorf("error fetching BPS base patch: %w", err)
	}
	base_patch, err := bps.FromBytes(bps_data)
	if err != nil {
		return fmt.Errorf("error when parsing BPS base patch: %w", err)
	}

	patch.base_patch = &base_patch
	return nil
}
