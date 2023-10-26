package filecache

// NewValues creates a new custom Values map.
// The keysAndValues parameters expect a list of keys and values in format:
// key1, value1, key2, value2...
func NewValues(keysAndValues ...any) Values {
	values := make(Values)

	hasKey := false
	currentKey := ""

	for i, v := range keysAndValues {
		expectKey := i%2 == 0

		if expectKey {
			key, ok := v.(string)
			if ok {
				hasKey = true
				currentKey = key
			} else {
				hasKey = false
			}
		} else {
			if !hasKey {
				continue
			}

			values[currentKey] = v
		}
	}

	return values
}

// Values are a custom values map
type Values map[string]any
