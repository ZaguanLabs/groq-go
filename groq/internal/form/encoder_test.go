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
