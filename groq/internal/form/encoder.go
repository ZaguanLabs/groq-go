package form

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// Encoder encodes a struct into multipart/form-data
type Encoder struct {
	w *multipart.Writer
	b *bytes.Buffer
}

// NewEncoder returns a new Encoder
func NewEncoder() *Encoder {
	b := &bytes.Buffer{}
	return &Encoder{
		w: multipart.NewWriter(b),
		b: b,
	}
}

// Encode struct to multipart
func (e *Encoder) Encode(v interface{}) (string, io.Reader, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("form encode: expected struct, got %v", val.Kind())
	}

	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		value := val.Field(i)

		tag := field.Tag.Get("json") // Reuse json tags for simplicity
		name := strings.Split(tag, ",")[0]
		if name == "-" || name == "" {
			continue
		}

		// Handle Optional
		if value.Kind() == reflect.Ptr && strings.Contains(value.Type().String(), "Optional") {
			if value.IsNil() {
				continue
			}
			// Dereference *Optional[T]
			opt := value.Elem()
			// Check Set field
			setField := opt.FieldByName("Set")
			if !setField.Bool() {
				continue
			}
			// Use Value field
			value = opt.FieldByName("Value")
		}

		if err := e.writeField(name, value.Interface()); err != nil {
			return "", nil, err
		}
	}

	if err := e.w.Close(); err != nil {
		return "", nil, err
	}

	return e.w.FormDataContentType(), e.b, nil
}

func (e *Encoder) writeField(name string, v interface{}) error {
	switch val := v.(type) {
	case *os.File:
		part, err := e.w.CreateFormFile(name, filepath.Base(val.Name()))
		if err != nil {
			return err
		}
		_, err = io.Copy(part, val)
		return err

	case io.Reader:
		// It's a file, but we need filename.
		// For now assume the struct has "File" field and handled specially or
		// assume we don't know filename.
		// API usually requires filename for file uploads.
		// If passed as io.Reader directly, we might default to "file.bin".
		part, err := e.w.CreateFormFile(name, "file.bin")
		if err != nil {
			return err
		}
		_, err = io.Copy(part, val)
		return err

	default:
		// Primitive
		return e.w.WriteField(name, fmt.Sprint(val))
	}
}
