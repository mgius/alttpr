package bps

import (
	"bytes"
	"os"
	"testing"
)

func compare_bps(expected *BPSPatch, actual *BPSPatch, t *testing.T) {
	if actual.SourceSize != expected.SourceSize {
		t.Fatalf("source_size mismatch %d = %d", actual.SourceSize, expected.SourceSize)
	}

	if actual.TargetSize != expected.TargetSize {
		t.Fatalf("target_size mismatch %d = %d", actual.TargetSize, expected.TargetSize)
	}

	if actual.MetadataSize != expected.MetadataSize {
		t.Fatalf("metadata_size mismatch %d = %d", actual.MetadataSize, expected.MetadataSize)
	}

	if actual.Metadata != expected.Metadata {
		t.Fatalf("metadata mismatch %s = %s", actual.Metadata, expected.Metadata)
	}

	if actual.SourceChecksum != expected.SourceChecksum {
		t.Fatalf("SourceChecksum mismatch %d = %d", actual.SourceChecksum, expected.SourceChecksum)
	}

	if actual.TargetChecksum != expected.TargetChecksum {
		t.Fatalf("TargetChecksum mismatch %d = %d", actual.TargetChecksum, expected.TargetChecksum)
	}

	if actual.PatchChecksum != expected.PatchChecksum {
		t.Fatalf("PatchChecksum mismatch %d = %d", actual.PatchChecksum, expected.PatchChecksum)
	}

}

func TestReadBPSFile(t *testing.T) {
	// Test against a trivial text diff BPS patch
	expected_bps := BPSPatch{
		SourceSize:     45,
		TargetSize:     92,
		MetadataSize:   0,
		Metadata:       "",
		SourceChecksum: 0x133070d,
		TargetChecksum: 0x76c91265,
		PatchChecksum:  0xc18e4db1,
	}

	f, _ := os.Open("testpatch.bps")
	bps, _ := read_bps_patch_file(f)

	compare_bps(&expected_bps, &bps, t)

}

func TestReadALTTPRBPSFile(t *testing.T) {
	// Test against a real ALTTPR bps patch.  This data confirmed by two alternative python based implementations
	expected_bps := BPSPatch{
		SourceSize:     1048576,
		TargetSize:     2097152,
		MetadataSize:   66,
		Metadata:       `{"created":"2021-09-18","hash":"7f2e1606616492d7dfb589e8dfb70027"}`,
		SourceChecksum: 0x3322effc,
		TargetChecksum: 0xe7565629,
		PatchChecksum:  0xb2b9ef4b,
		Actions:        make([]byte, 126299),
	}

	f, _ := os.Open("7f2e1606616492d7dfb589e8dfb70027.bps")
	bps, err := read_bps_patch_file(f)
	if err != nil {
		t.Fatalf("%s", err)
	}

	compare_bps(&expected_bps, &bps, t)
}

func TestEncodeOneByte(t *testing.T) {
	const encode_one_byte uint64 = 0b1011     // decimal 11
	const expected_encoding byte = 0b10001011 // decimal 11 with highest bit flagged
	var writeBuffer bytes.Buffer

	err := bps_write_num(&writeBuffer, encode_one_byte)

	if err != nil {
		t.Fatalf("bps_write_num returned an error: %s", err)
	}

	if writeBuffer.Len() != 1 {
		t.Fatalf("bps_write_num wrote too many bytes")
	}

	if writeBuffer.Bytes()[0] != expected_encoding {
		t.Fatalf("bps_write_num did not encode correctly")
	}
}

// TODO: fix encoded value
func _TestEncodeTwoBytes(t *testing.T) {
	const encode_two_bytes uint64 = 0b101_0001011 // 651
	expected_encoding := []byte{0b0_0001011, 0b1_0000101}

	var writeBuffer bytes.Buffer

	err := bps_write_num(&writeBuffer, encode_two_bytes)

	if err != nil {
		t.Fatalf("bps_write_num returned an error: %s", err)
	}

	if writeBuffer.Len() != 2 {
		t.Fatalf("bps_write_num wrote too many bytes")
	}

	if writeBuffer.Bytes()[0] != expected_encoding[0] {
		t.Fatalf("bps_write_num did not encode correctly")
	}

	if writeBuffer.Bytes()[1] != expected_encoding[1] {
		t.Fatalf("bps_write_num did not encode correctly %b != %b", writeBuffer.Bytes()[1], expected_encoding[1])
	}

}

func TestDecodeOneByte(t *testing.T) {
	var encoded []byte = []byte{0b10001011} // decimal 11 with highest bit flagged
	const expected_decode uint64 = 0b1011   // decimal 11

	readBuffer := bytes.NewBuffer(encoded)

	decoded, err := bps_read_num(readBuffer)

	if err != nil {
		t.Fatalf("bps_read_num threw an error")
	}

	if decoded != expected_decode {
		t.Fatalf("bps_read_num did not decode correctly")
	}
}

// TODO: fix the encoded value
func _TestDecodeTwoBytes(t *testing.T) {
	encoded := []byte{0b0_0001011, 0b1_0000101}
	const expected_decode uint64 = 0b101_0001011 // 651

	readBuffer := bytes.NewBuffer(encoded)

	decoded, err := bps_read_num(readBuffer)

	if err != nil {
		t.Fatalf("bps_read_num threw an error")
	}

	if decoded != expected_decode {
		t.Fatalf("bps_read_num did not decode correctly")
	}
}

func TestCanDecodeEncodedNumbers(t *testing.T) {
	const encode_big_num uint64 = 0xdeadbeefdeadbeef

	var writeBuffer bytes.Buffer

	err := bps_write_num(&writeBuffer, encode_big_num)

	if err != nil {
		t.Fatalf("bps_write_num returned an error: %s", err)
	}

	read_num, err := bps_read_num(&writeBuffer)

	if err != nil {
		t.Fatalf("bps_read_num returned an error: %s", err)
	}

	if read_num != encode_big_num {
		t.Fatalf("Number did not Encode/Decode to the same number. %x != %x", read_num, encode_big_num)
	}
}
