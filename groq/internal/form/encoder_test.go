package form

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/ZaguanLabs/groq-go/groq/option"
)

type TestForm struct {
	String string                   `json:"string"`
	Int    int                      `json:"int"`
	File   io.Reader                `json:"file"`
	Opt    *option.Optional[string] `json:"opt,omitempty"`
	OptSet *option.Optional[string] `json:"opt_set,omitempty"`
}

func TestEncoder_Encode(t *testing.T) {
	f, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("test content")
	f.Seek(0, 0)

	form := TestForm{
		String: "hello",
		Int:    42,
		File:   f,
		Opt:    option.Ptr(option.None[string]()),
		OptSet: option.Ptr(option.Some("set")),
	}

	enc := NewEncoder()
	ct, body, err := enc.Encode(form)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}

	if !strings.Contains(ct, "multipart/form-data") {
		t.Errorf("Content-Type mismatch: %s", ct)
	}

	b, _ := io.ReadAll(body)
	s := string(b)

	if !strings.Contains(s, `name="string"`) || !strings.Contains(s, "hello") {
		t.Error("Missing string field")
	}
	if !strings.Contains(s, `name="int"`) || !strings.Contains(s, "42") {
		t.Error("Missing int field")
	}
	if !strings.Contains(s, `name="file"`) || !strings.Contains(s, `filename="`) {
		// filename might be test... or file.bin depending on type check
		// os.CreateTemp returns *os.File, so it should match *os.File case and use filename
		t.Error("Missing file field")
	}
	if strings.Contains(s, `name="opt"`) {
		t.Error("Opt field should be omitted")
	}
	if !strings.Contains(s, `name="opt_set"`) || !strings.Contains(s, "set") {
		t.Error("Missing opt_set field")
	}
}

func TestEncoder_EncodeStringReader(t *testing.T) {
	form := TestForm{
		String: "test",
		File:   strings.NewReader("file content"),
	}

	enc := NewEncoder()
	ct, body, err := enc.Encode(form)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}

	if !strings.Contains(ct, "multipart/form-data") {
		t.Errorf("Content-Type mismatch: %s", ct)
	}

	b, _ := io.ReadAll(body)
	s := string(b)

	if !strings.Contains(s, "file content") {
		t.Error("Missing file content")
	}
	if !strings.Contains(s, `filename="file.bin"`) {
		t.Error("Expected default filename for io.Reader")
	}
}

func TestEncoder_EncodeEmptyStruct(t *testing.T) {
	type EmptyForm struct{}

	enc := NewEncoder()
	ct, body, err := enc.Encode(EmptyForm{})
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}

	if !strings.Contains(ct, "multipart/form-data") {
		t.Errorf("Content-Type mismatch: %s", ct)
	}

	b, _ := io.ReadAll(body)
	if len(b) == 0 {
		t.Error("Expected non-empty body even for empty struct")
	}
}

func TestEncoder_EncodeNilOptional(t *testing.T) {
	form := TestForm{
		String: "test",
		Opt:    nil, // nil pointer to Optional
	}

	enc := NewEncoder()
	_, body, err := enc.Encode(form)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}

	b, _ := io.ReadAll(body)
	s := string(b)

	if strings.Contains(s, `name="opt"`) {
		t.Error("Nil optional should be omitted")
	}
}

func TestEncoder_EncodeMultipleFiles(t *testing.T) {
	type MultiFileForm struct {
		File1 io.Reader `json:"file1"`
		File2 io.Reader `json:"file2"`
	}

	form := MultiFileForm{
		File1: strings.NewReader("content1"),
		File2: strings.NewReader("content2"),
	}

	enc := NewEncoder()
	_, body, err := enc.Encode(form)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}

	b, _ := io.ReadAll(body)
	s := string(b)

	if !strings.Contains(s, "content1") {
		t.Error("Missing file1 content")
	}
	if !strings.Contains(s, "content2") {
		t.Error("Missing file2 content")
	}
}

func TestEncoder_EncodeZeroValues(t *testing.T) {
	type ZeroForm struct {
		String string `json:"string"`
		Int    int    `json:"int"`
		Bool   bool   `json:"bool"`
	}

	form := ZeroForm{} // All zero values

	enc := NewEncoder()
	_, body, err := enc.Encode(form)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}

	b, _ := io.ReadAll(body)
	s := string(b)

	// Zero values should still be included
	if !strings.Contains(s, `name="string"`) {
		t.Error("Missing string field")
	}
	if !strings.Contains(s, `name="int"`) {
		t.Error("Missing int field")
	}
}
