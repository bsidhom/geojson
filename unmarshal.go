package geojson

import (
	"encoding/json"
	"fmt"
)

var _ json.Unmarshaler = &Point{}

func (p *Point) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	return p.unmarshalFrom(m)
}

func (p *Point) unmarshalFrom(m map[string]interface{}) error {
	typeName, err := getString("type", m)
	if err != nil {
		return nil
	}
	if typeName != "Point" {
		return fmt.Errorf("expected type \"Point\", got %q", typeName)
	}
	coordinates, err := getNumArray("coordinates", m)
	if err != nil {
		return err
	}
	length := len(coordinates)
	if length < 2 || length > 3 {
		return fmt.Errorf("Point coordinate array should have 2 or 3 elements, got %d", length)
	}

	p.X = coordinates[0]
	p.Y = coordinates[1]
	if length == 3 {
		p.Elevation = coordinates[2]
		p.HasElevation = true
	}
	p.AdditionalData = shallowCopyWithoutKeys(m, "type", "coordinates")
	return nil
}

func getString(name string, m map[string]interface{}) (string, error) {
	v, err := getField(name, m)
	if err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("field %q should have been of string, got %T", name, v)
	}
	return s, nil
}

func getNumArray(name string, m map[string]interface{}) ([]float64, error) {
	v, err := getField(name, m)
	if err != nil {
		return nil, err
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("field %q should have been []float64, got %T", name, v)
	}
	result := make([]float64, len(arr))
	for i, v := range arr {
		n, ok := v.(float64)
		if !ok {
			return nil, fmt.Errorf("expected float64 in %q[%d], got %T", name, i, v)
		}
		result[i] = n
	}
	return result, nil
}

func getField(name string, m map[string]interface{}) (interface{}, error) {
	v, ok := m[name]
	if !ok {
		return "", fmt.Errorf("missing field %q", name)
	}
	return v, nil
}

func shallowCopyWithoutKeys(m map[string]interface{}, keys ...string) map[string]interface{} {
	// NOTE: Not all keys are necessarily in map.
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		if stringInSlice(k, keys) {
			// Omit this key.
			continue
		}
		result[k] = v
	}
	// We don't want to even include an empty map if there's no additional data.
	if len(result) == 0 {
		return nil
	}
	return result
}

func stringInSlice(s string, ss []string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}
