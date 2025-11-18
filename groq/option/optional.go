package option

import (
	"encoding/json"
)

// Optional represents a value that may or may not be set.
// It is used to distinguish between a zero value and an unset value.
type Optional[T any] struct {
	Value T
	Set   bool
}

// Some creates an Optional with a value.
func Some[T any](v T) Optional[T] {
	return Optional[T]{Value: v, Set: true}
}

// None creates an unset Optional.
func None[T any]() Optional[T] {
	var zero T
	return Optional[T]{Value: zero, Set: false}
}

// Ptr returns a pointer to the Optional value.
// This is useful for structs fields that use *Optional[T] for omitempty support.
func Ptr[T any](o Optional[T]) *Optional[T] {
	return &o
}

// IsSet returns true if the value is set.
func (o Optional[T]) IsSet() bool {
	return o.Set
}

// IsZero returns true if the value is unset, allowing omitempty to work.
func (o Optional[T]) IsZero() bool {
	return !o.Set
}

// MarshalJSON implements json.Marshaler.
// If Set is false, it returns "null" (or omits if omitempty is used on the field,
// but standard library json encoder doesn't check this method for omitempty decision on structs,
// it checks the pointer or emptiness if implemented).
// Actually, if we use Optional[T] as a value receiver, it's a struct.
// `omitempty` on a struct field only works if the struct implements IsZero() (Go 1.22+) or is a pointer.
// We likely want to use pointers to Optional or rely on the user not setting it.
// But the plan says: "Optional[T] -> *T (pointer for optionality)" is what Python -> Go mapping suggested,
// BUT then "Omit Pattern" section suggested this struct.
//
// If we use this struct as a value in another struct field:
// Field Optional[int] `json:"field,omitempty"`
// Go's json encoder will NOT omit this unless we make it a pointer or it's "empty".
// But a struct with a bool is never "empty" to the encoder unless we define what empty is.
// As of Go 1.24 (which the user seems to use based on go.mod), we can't easily make `omitempty` work
// for non-pointer structs without `IsZero` (which isn't fully standard for json omitempty yet?).
// Actually, the safest way to get `omitempty` behavior is to make the field a pointer or
// make `MarshalJSON` return `null` and hope the API handles `null` as unset (which is NOT the same as Omit).
//
// The plan snippet:
// func (o Optional[T]) MarshalJSON() ([]byte, error) {
//    if !o.Set { return []byte("null"), nil } ...
// }
//
// If the API treats `null` as "unset/ignore", then returning `null` is fine.
// If the API treats `null` as "clear/delete", then we MUST omit the field entirely.
// OpenAI/Groq usually treats `null` as clear for some fields, and omit for others.
// So `omitempty` is crucial.
//
// To support `omitempty` properly with this struct, we might need to implement `IsZero() bool`?
// No, `encoding/json` doesn't use `IsZero` interface universally yet.
// The standard pattern in Go for optional fields is usually pointers `*T`.
// The plan PROPOSED this `Optional` struct.
// "Python uses Omit sentinel... Go equivalent: type Optional[T]..."
//
// If I implement this, how do I ensure it's omitted from JSON object if not set?
// If it's embedded in a struct:
// type Request struct {
//   Temp Optional[float64] `json:"temp,omitempty"`
// }
// `json.Marshal` will marshal it unless it's considered empty.
// A struct is empty if all fields are zero. `Set` bool makes it non-zero if false?
// No, `Set: false` is zero value for bool. `Value` zero value is zero.
// So `Optional[T]{}` (zero value) IS empty.
// So `omitempty` WILL work if both Value and Set are zero.
// BUT `Some(0)` -> `Set: true`. Not empty. Correct.
// `None()` -> `Set: false`. If `Value` is also zero, then it is empty.
// So `Optional[int]{Value: 0, Set: false}` is empty.
// `Optional[int]{Value: 1, Set: false}` is NOT empty.
// So as long as we ensure `None` has zero value for T, `omitempty` works!
//
// So:
// func None[T any]() Optional[T] { return Optional[T]{} } // Zero value
//
// Then MarshalJSON is only needed if we want to customize serialization of Value?
// Or if we want `Some(0)` to serialize as `0`.
// `json.Marshal` on struct calls `MarshalJSON` if present.
// If we implement `MarshalJSON`, `omitempty` checks usually happen BEFORE calling MarshalJSON?
// No, for structs, `omitempty` checks if the struct value is equal to zero value of that struct type.
// So `Optional[T]{}` is zero value.
// If we implement MarshalJSON, does it affect omitempty check? No.
//
// So `omitempty` tag works for `Optional[T]` IF `Set` is false AND `Value` is zero.
// Our `None()` should ensure `Value` is zero.
//
// What if we want to send explicit `null`?
// Then we need `Set: true` but `Value` is ... wait.
// `Optional` here is "Set vs Unset".
// If we need "Set to Null", that's different.
// `Optional[T]` implies T is the type. If T can be null (pointer), then `Some(nil)` is explicit null.
// If T is int, `Some(0)` is 0. We can't send `null` for int unless T is `*int`.
//
// So `Optional[T]` is strictly for "Is this field present in the JSON or not?".
//
// Let's follow the plan's implementation but ensure `MarshalJSON` does the right thing.
// If `Set` is true, marshal value.
// If `Set` is false, marshal "null"?
// If `omitempty` is used, `None` (zero value) will be omitted.
// If `omitempty` is NOT used, `None` will be marshaled.
// In that case, `null` is probably correct default for "not set but present".

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.Set {
		return []byte("null"), nil
	}
	return json.Marshal(o.Value)
}

// UnmarshalJSON implements json.Unmarshaler
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.Set = false
		return nil
	}
	o.Set = true
	return json.Unmarshal(data, &o.Value)
}
