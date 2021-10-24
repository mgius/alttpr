package bps

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

var (
	BPS_HEADER = []byte("BPS1")
)

type BPSPatch struct {
	SourceSize     uint64
	TargetSize     uint64
	MetadataSize   uint64
	Metadata       string
	Actions        []byte
	SourceChecksum uint32
	TargetChecksum uint32
	PatchChecksum  uint32
}

func (*BPSPatch) PatchSourceFile(sourcefile *os.File) (targetfiledata []byte, err error) {
	return

}

func FromFile(patchfile *os.File) (BPSPatch, error) {
	filestat, _ := patchfile.Stat()
	filesize := filestat.Size()

	full_file := make([]byte, filesize)
	patchfile.Read(full_file)

	if !bytes.Equal(full_file[:len(BPS_HEADER)], BPS_HEADER) {
		return BPSPatch{}, errors.New("Magic Header Incorrect")
	}

	remaining := full_file[len(BPS_HEADER):]

	// TODO: error handling
	source_size, remaining, _, _ := bps_read_num(remaining)

	target_size, remaining, _, _ := bps_read_num(remaining)
	metadata_size, remaining, _, _ := bps_read_num(remaining)
	metadata, remaining := string(remaining[:metadata_size]), remaining[metadata_size:]

	action_len := len(remaining) - 12
	actions, remaining := remaining[:action_len], remaining[action_len:]

	source_checksum := binary.LittleEndian.Uint32(remaining[:4])
	target_checksum := binary.LittleEndian.Uint32(remaining[4:8])
	patch_checksum := binary.LittleEndian.Uint32(remaining[8:12])

	// TODO: validate patch_checksum

	return BPSPatch{
		SourceSize:     source_size,
		TargetSize:     target_size,
		MetadataSize:   metadata_size,
		Metadata:       metadata,
		Actions:        actions,
		SourceChecksum: source_checksum,
		TargetChecksum: target_checksum,
		PatchChecksum:  patch_checksum,
	}, nil

}

func bps_write_num(bytewriter io.ByteWriter, num uint64) error {
	for true {
		// slice off the lowest 7 bits of num
		x := byte(num & 0x7f)
		// shift the lowest 7 bits out of the num
		num >>= 7

		// If we've encoded all bits of the number into either x or the byte
		// stream, write out x with the end of number bit set
		if num == 0 {
			err := bytewriter.WriteByte(0x80 | x)
			if err != nil {
				return err
			}
			break
		}

		// Otherwise, write out the byte and loop around
		bytewriter.WriteByte(x)

		// weird optimization for "one"?
		// I don't understand the purpose of this optimization, and the
		// reference decode implementation doesn't seem to handle this
		// optimization at all, and every other bps impl I've seen doesn't do
		// this either
		num--
	}

	return nil
}

func bps_read_num(stream []byte) (data uint64, remainder []byte, bytes_read int, err error) {
	var (
		// data  uint64 = 0
		shift uint64 = 1
	)

	for bytes_read < len(stream) {
		// Grab the next byte and indicate we read one.
		var x = stream[bytes_read]
		bytes_read++

		// Mask off the eigth bit.  Multiply the remaining 7 bits by the shift,
		// and add into our data parameter.
		data += uint64((x & 0x7f)) * shift

		// If the 8th bit is set, we've reached end of number
		if (x & 0x80) == 0x80 {
			remainder = stream[bytes_read:]
			return
		}
		// Increase the shift so that further reads represent higher bits in the read number
		shift <<= 7

		// I think this has to do with the way the encoding subtracts one from
		// the data as it goes, so you add "one" as we go?  But what about the first one?
		data += shift
	}

	err = errors.New("bps_read_num: Ran out of bytes before termination bit was set")

	return
}
