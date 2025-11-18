package querystring

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"
)

// Values maps a string key to a list of values.
// It is typically used for query parameters and form values.
type Values map[string][]string

// ToURLValues converts Values to url.Values
func (v Values) ToURLValues() url.Values {
	uv := make(url.Values)
	for k, vals := range v {
		uv[k] = vals
	}
	return uv
}

// Stringify converts a map of query parameters to a URL-encoded string.
// It supports specific formatting for arrays (comma-separated).
// It handles:
// - string, int, float, bool, time.Time
// - slices (comma separated)
// - maps (flattened? No, usually query params are flat or simple nested. Groq seems to use flat or simple)
// The plan says: "Serializes query params via Querystring.stringify() with comma array format"
func Stringify(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}

	values := make(Values)
	if err := encode(values, v, ""); err != nil {
		return "", err
	}

	if len(values) == 0 {
		return "", nil
	}

	// Convert to comma-separated string for arrays
	var buf strings.Builder
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}

		// Join array values with commas
		valStr := strings.Join(values[k], ",")

		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(valStr))
	}

	return buf.String(), nil
}

func encode(values Values, v interface{}, prefix string) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		iter := val.MapRange()
		keys := make([]string, 0, val.Len())
		tempMap := make(map[string]reflect.Value)

		for iter.Next() {
			k := iter.Key().String()
			keys = append(keys, k)
			tempMap[k] = iter.Value()
		}
		sort.Strings(keys)

		for _, k := range keys {
			keyName := k
			if prefix != "" {
				keyName = prefix + "[" + k + "]"
			}
			if err := encode(values, tempMap[k].Interface(), keyName); err != nil {
				return err
			}
		}
	case reflect.Struct:
		if t, ok := v.(time.Time); ok {
			if prefix == "" {
				return fmt.Errorf("time at top level not supported")
			}
			values[prefix] = append(values[prefix], t.Format(time.RFC3339))
			return nil
		}
		return fmt.Errorf("struct encoding not yet implemented")
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			// For arrays, we collect them to join later
			// But wait, the recursive call structure needs care.
			// If we are at the top level, 'prefix' is empty, which is invalid for a slice.
			// Slices usually come from a key in a map.
			if prefix == "" {
				return fmt.Errorf("slice at top level not supported")
			}

			// We add to the existing key in values
			s := fmt.Sprintf("%v", val.Index(i).Interface())
			values[prefix] = append(values[prefix], s)
		}
	default:
		if prefix == "" {
			return fmt.Errorf("primitive at top level not supported")
		}
		s := fmt.Sprint(v)
		// Handle specific types if needed (e.g. bool lowercase?)
		if val.Kind() == reflect.Bool {
			s = fmt.Sprintf("%t", v) // "true" / "false"
		}
		values[prefix] = append(values[prefix], s)
	}
	return nil
}
