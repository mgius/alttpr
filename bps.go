package bps

import (
	"bytes"
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

func read_bps_patch_file(patchfile *os.File) (BPSPatch, error) {
	bps_header := make([]byte, len(BPS_HEADER))
	read_bytes, _ := patchfile.Read(bps_header)
	if read_bytes != len(BPS_HEADER) || !bytes.Equal(bps_header, BPS_HEADER) {
		return BPSPatch{}, errors.New("Magic Header Incorrect")
	}

	source_size, _ := bps_read_num(patchfile)
	target_size, _ := bps_read_num(patchfile)
	metadata_size, _ := bps_read_num(patchfile)

	return BPSPatch{
		SourceSize:   source_size,
		TargetSize:   target_size,
		MetadataSize: metadata_size,
	}, nil

}

func convert_byte(b byte) (uint64, error) {
	return uint64(b), nil
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

func bps_read_num(reader io.Reader) (uint64, error) {
	var (
		data  uint64 = 0
		shift uint64 = 1
	)

	for true {
		// Read a byte
		var x = make([]byte, 1)
		bytes_read, err := reader.Read(x)
		if err != nil || bytes_read != 1 {
			return 0, err
		}

		// Mask off the eigth bit.  Multiply the remaining 7 bits by the shift,
		// which will increase with each byte we read
		data += uint64((x[0] & 0x7f)) * shift

		// If the 8th bit is set, we've reached end of number
		if (x[0] & 0x80) == 0x80 {
			break
		}
		// Increase the shift so that further reads are larger
		shift <<= 7

		// I think this has to do with the way the encoding subtracts one from
		// the data as it goes, so you add "one" as we go?  But what about the first one?
		data += shift
	}

	return data, nil
}
