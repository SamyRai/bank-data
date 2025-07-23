package dataset

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var keyMap = map[string]string{
	"Name":    "N",
	"Country": "C",
	"Code":    "R",
}
var reverseKeyMap = map[string]string{
	"N": "Name",
	"C": "Country",
	"R": "Code",
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	records := []map[string]string{
		{"Name": "Bank1", "Country": "DE", "Code": "123"},
		{"Name": "Bank2", "Country": "FR", "Code": "456"},
	}
	buf := &bytes.Buffer{}
	if err := Encode(buf, records, keyMap); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	decoded, err := Decode(buf, reverseKeyMap, "", "")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if !reflect.DeepEqual(records, decoded) {
		t.Errorf("Round-trip mismatch:\nwant: %#v\ngot: %#v", records, decoded)
	}
}

func TestDecodeHandlesOddTokens(t *testing.T) {
	data := "N\x1FBank1\x1FC\x1FDE\x1FR\x1F123\x1FX" // Odd token count, X has no value
	buf := bytes.NewBufferString(data + "\n")
	decoded, err := Decode(buf, reverseKeyMap, "", "")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	// Expect zero records, as line is malformed and should be skipped
	if len(decoded) != 0 {
		t.Errorf("expected 0 records for malformed line, got %d", len(decoded))
	}
}

func TestEncodeSkipsUnknownFields(t *testing.T) {
	record := map[string]string{"Name": "Bank1", "Country": "DE", "Unknown": "foo"}
	buf := &bytes.Buffer{}
	if err := Encode(buf, []map[string]string{record}, keyMap); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	out := buf.String()
	if out == "" || strings.Contains(out, "Unknown"+Sep+"foo") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestDecodeEmptyInput(t *testing.T) {
	buf := bytes.NewBufferString("")
	decoded, err := Decode(buf, reverseKeyMap, "", "")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if len(decoded) != 0 {
		t.Errorf("expected 0 records, got %d", len(decoded))
	}
}

func TestEncodeDecodeUnicodeAndSpaces(t *testing.T) {
	records := []map[string]string{
		{"Name": "Bänk Ω", "Country": "DE FR", "Code": ""},
	}
	buf := &bytes.Buffer{}
	if err := Encode(buf, records, keyMap); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	decoded, err := Decode(buf, reverseKeyMap, "", "")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded[0]["Name"] != "Bänk Ω" || decoded[0]["Country"] != "DE FR" {
		t.Errorf("unexpected decoded values: %#v", decoded[0])
	}
}

func TestDecodeUnknownShortKey(t *testing.T) {
	data := "N\x1FBank1\x1FZ\x1FUnknownKey\x1FC\x1FDE\x1FR\x1F123\n"
	buf := bytes.NewBufferString(data)
	decoded, err := Decode(buf, reverseKeyMap, "", "")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if len(decoded) != 1 {
		t.Fatalf("expected 1 record, got %d", len(decoded))
	}
	if _, ok := decoded[0]["UnknownKey"]; ok {
		t.Errorf("unexpected key decoded: %v", decoded[0])
	}
}

func TestEncodeUnknownLongKey(t *testing.T) {
	record := map[string]string{"Name": "Bank1", "Country": "DE", "Unknown": "foo"}
	buf := &bytes.Buffer{}
	if err := Encode(buf, []map[string]string{record}, keyMap); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "Unknown foo") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestDecodeLargeInput(t *testing.T) {
	var b bytes.Buffer
	for i := 0; i < 1000; i++ {
		b.WriteString(fmt.Sprintf("N\x1FBank%d\x1FC\x1FDE\x1FR\x1F%d\n", i, i))
	}
	decoded, err := Decode(&b, reverseKeyMap, "", "")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if len(decoded) != 1000 {
		t.Errorf("expected 1000 records, got %d", len(decoded))
	}
	if decoded[999]["Name"] != "Bank999" || decoded[999]["Code"] != "999" {
		t.Errorf("unexpected last record: %#v", decoded[999])
	}
}

func TestDecodeIncompletePairs(t *testing.T) {
	data := "N\x1FBank1\x1FC\x1FDE\x1FR\n" // R has no value
	buf := bytes.NewBufferString(data)
	decoded, err := Decode(buf, reverseKeyMap, "", "")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	// Expect zero records, as line is malformed and should be skipped
	if len(decoded) != 0 {
		t.Errorf("expected 0 records for malformed line, got %d", len(decoded))
	}
}

func TestDecodeMalformedLines(t *testing.T) {
	data := "N\x1FBank1\x1FC\x1FDE\x1FR\nN\x1FBank2\x1FC\x1FDE\x1FR\x1F123\x1F" // first line missing value, second line odd tokens
	buf := bytes.NewBufferString(data + "\n")
	decoded, err := Decode(buf, reverseKeyMap, "", "")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	// Expect zero records for malformed input
	if len(decoded) != 0 {
		t.Errorf("expected 0 records for malformed lines, got %d", len(decoded))
	}
}

func TestDecodeUnknownShortKeyLogging(t *testing.T) {
	data := "N\x1FBank1\x1FZ\x1FUnknownKey\x1FC\x1FDE\x1FR\x1F123\n"
	buf := bytes.NewBufferString(data)
	// Should log unknown key, but not include in record
	decoded, err := Decode(buf, reverseKeyMap, "", "")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if len(decoded) != 1 {
		t.Fatalf("expected 1 record, got %d", len(decoded))
	}
	if _, ok := decoded[0]["UnknownKey"]; ok {
		t.Errorf("unexpected key decoded: %v", decoded[0])
	}
}

func TestDecodeUnicodeAndControlBytes(t *testing.T) {
	data := "N\x1FBänk Ω\x1FC\x1FDE\x1FR\x1F\x00\n"
	buf := bytes.NewBufferString(data)
	decoded, err := Decode(buf, reverseKeyMap, "", "")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded[0]["Name"] != "Bänk Ω" || decoded[0]["Code"] != "\x00" {
		t.Errorf("unexpected decoded values: %#v", decoded[0])
	}
}

func TestDecodeChecksumValidation(t *testing.T) {
	data := "N\x1FBank1\x1FC\x1FDE\x1FR\x1F123\n"
	buf := bytes.NewBufferString(data)
	// Compute correct checksum from the exact bytes passed to Decode
	checksumBytes := buf.Bytes()
	h := sha256.Sum256(checksumBytes)
	checksum := hex.EncodeToString(h[:])
	// Use a fresh buffer for Decode
	buf2 := bytes.NewBuffer(checksumBytes)
	decoded, err := Decode(buf2, reverseKeyMap, "", checksum)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if len(decoded) != 1 || decoded[0]["Name"] != "Bank1" {
		t.Errorf("unexpected decoded: %#v", decoded)
	}
	// Wrong checksum
	buf3 := bytes.NewBuffer(checksumBytes)
	_, err2 := Decode(buf3, reverseKeyMap, "", "deadbeef")
	if err2 == nil || !strings.Contains(err2.Error(), "checksum mismatch") {
		t.Errorf("expected checksum error, got: %v", err2)
	}
}

func TestDecodeStreamMode(t *testing.T) {
	var b bytes.Buffer
	for i := 0; i < 10; i++ {
		b.WriteString(fmt.Sprintf("N\x1FBank%d\x1FC\x1FDE\x1FR\x1F%d\n", i, i))
	}
	count := 0
	var lastName string
	err := DecodeStream(&b, reverseKeyMap, "", func(rec map[string]string, lineNum int) {
		count++
		lastName = rec["Name"]
		if lineNum != count {
			t.Errorf("lineNum mismatch: got %d, want %d", lineNum, count)
		}
	})
	if err != nil {
		t.Fatalf("DecodeStream failed: %v", err)
	}
	if count != 10 || lastName != "Bank10" && lastName != "Bank9" {
		t.Errorf("unexpected stream results: count=%d lastName=%s", count, lastName)
	}
}

func TestDecodeLargeDatasetMemory(t *testing.T) {
	var b bytes.Buffer
	for i := 0; i < 10000; i++ {
		b.WriteString(fmt.Sprintf("N\x1FBank%d\x1FC\x1FDE\x1FR\x1F%d\n", i, i))
	}
	decoded, err := Decode(&b, reverseKeyMap, "", "")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if len(decoded) != 10000 {
		t.Errorf("expected 10000 records, got %d", len(decoded))
	}
	if decoded[9999]["Name"] != "Bank9999" || decoded[9999]["Code"] != "9999" {
		t.Errorf("unexpected last record: %#v", decoded[9999])
	}
}
