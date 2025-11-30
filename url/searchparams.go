package url

import (
	"sort"
	"strings"
)

// urlParam represents a single key-value pair in URLSearchParams.
type urlParam struct {
	key   string
	value string
}

// URLSearchParams represents a collection of URL query parameters.
//
// It maintains insertion order and supports the full WHATWG URLSearchParams API.
type URLSearchParams struct {
	// entries stores the parameters in insertion order
	entries []urlParam

	// owner is the URL that owns this URLSearchParams, if any.
	// When set, mutations to the params will update the owner's query string.
	owner *URL
}

// NewURLSearchParams creates an empty URLSearchParams.
func NewURLSearchParams() *URLSearchParams {
	return &URLSearchParams{
		entries: make([]urlParam, 0),
	}
}

// NewURLSearchParamsFromString parses a query string into URLSearchParams.
//
// The input may or may not have a leading "?".
func NewURLSearchParamsFromString(raw string) *URLSearchParams {
	sp := &URLSearchParams{
		entries: make([]urlParam, 0),
	}

	// Strip leading ? if present
	raw = strings.TrimPrefix(raw, "?")

	if raw == "" {
		return sp
	}

	sp.entries = parseFormEncoded(raw)
	return sp
}

// NewURLSearchParamsFromEntries creates URLSearchParams from key-value pairs.
func NewURLSearchParamsFromEntries(entries [][2]string) *URLSearchParams {
	sp := &URLSearchParams{
		entries: make([]urlParam, 0, len(entries)),
	}

	for _, entry := range entries {
		sp.entries = append(sp.entries, urlParam{
			key:   entry[0],
			value: entry[1],
		})
	}

	return sp
}

// NewURLSearchParamsFromMap creates URLSearchParams from a map.
// Note: map iteration order is not guaranteed, so the order of entries
// may vary between calls.
func NewURLSearchParamsFromMap(m map[string]string) *URLSearchParams {
	sp := &URLSearchParams{
		entries: make([]urlParam, 0, len(m)),
	}

	for k, v := range m {
		sp.entries = append(sp.entries, urlParam{
			key:   k,
			value: v,
		})
	}

	return sp
}

// Clone creates a copy of the URLSearchParams without the owner reference.
func (sp *URLSearchParams) Clone() *URLSearchParams {
	clone := &URLSearchParams{
		entries: make([]urlParam, len(sp.entries)),
	}
	copy(clone.entries, sp.entries)
	return clone
}

// syncOwner updates the owner URL's query string if one exists.
func (sp *URLSearchParams) syncOwner() {
	if sp.owner != nil {
		sp.owner.syncFromSearchParams()
	}
}

// Append adds a new key-value pair to the end of the list.
func (sp *URLSearchParams) Append(key, value string) {
	sp.entries = append(sp.entries, urlParam{key: key, value: value})
	sp.syncOwner()
}

// Delete removes all entries with the given key.
// If value is provided (non-nil), only entries with matching key AND value are removed.
func (sp *URLSearchParams) Delete(key string, value *string) {
	newEntries := make([]urlParam, 0, len(sp.entries))
	for _, entry := range sp.entries {
		if entry.key == key {
			if value == nil || entry.value == *value {
				continue // skip this entry (delete it)
			}
		}
		newEntries = append(newEntries, entry)
	}
	sp.entries = newEntries
	sp.syncOwner()
}

// Get returns the first value for the given key, or empty string if not found.
func (sp *URLSearchParams) Get(key string) (string, bool) {
	for _, entry := range sp.entries {
		if entry.key == key {
			return entry.value, true
		}
	}
	return "", false
}

// GetAll returns all values for the given key.
func (sp *URLSearchParams) GetAll(key string) []string {
	values := make([]string, 0)
	for _, entry := range sp.entries {
		if entry.key == key {
			values = append(values, entry.value)
		}
	}
	return values
}

// Has returns true if a parameter with the given key exists.
// If value is provided (non-nil), returns true only if a matching key-value pair exists.
func (sp *URLSearchParams) Has(key string, value *string) bool {
	for _, entry := range sp.entries {
		if entry.key == key {
			if value == nil {
				return true
			}
			if entry.value == *value {
				return true
			}
		}
	}
	return false
}

// Set sets the value for the given key, replacing any existing values.
// If the key doesn't exist, it appends a new entry.
func (sp *URLSearchParams) Set(key, value string) {
	found := false
	newEntries := make([]urlParam, 0, len(sp.entries))

	for _, entry := range sp.entries {
		if entry.key == key {
			if !found {
				// Keep the first occurrence but update its value
				newEntries = append(newEntries, urlParam{key: key, value: value})
				found = true
			}
			// Skip subsequent occurrences (effectively deleting them)
		} else {
			newEntries = append(newEntries, entry)
		}
	}

	if !found {
		newEntries = append(newEntries, urlParam{key: key, value: value})
	}

	sp.entries = newEntries
	sp.syncOwner()
}

// Sort sorts all entries by their keys using stable sort.
// Per WHATWG URL spec, sorting is done by comparing code units (UTF-16).
func (sp *URLSearchParams) Sort() {
	sort.SliceStable(sp.entries, func(i, j int) bool {
		return compareByCodeUnits(sp.entries[i].key, sp.entries[j].key) < 0
	})
	sp.syncOwner()
}

// compareByCodeUnits compares two strings by their UTF-16 code units.
// This matches JavaScript's default string comparison behavior.
func compareByCodeUnits(a, b string) int {
	// Convert to UTF-16 code units for comparison
	aRunes := []rune(a)
	bRunes := []rune(b)

	minLen := len(aRunes)
	if len(bRunes) < minLen {
		minLen = len(bRunes)
	}

	for i := 0; i < minLen; i++ {
		// For characters in the BMP (< 0x10000), the code unit equals the code point
		// For characters >= 0x10000, we need to compare as surrogate pairs
		aUnits := runeToCodeUnits(aRunes[i])
		bUnits := runeToCodeUnits(bRunes[i])

		// Compare first code unit
		if aUnits[0] != bUnits[0] {
			if aUnits[0] < bUnits[0] {
				return -1
			}
			return 1
		}

		// If first code unit is equal and there's a second one, compare it
		if len(aUnits) > 1 && len(bUnits) > 1 {
			if aUnits[1] != bUnits[1] {
				if aUnits[1] < bUnits[1] {
					return -1
				}
				return 1
			}
		} else if len(aUnits) > 1 {
			// a has surrogate pair, b doesn't - a is "longer" at this position
			return 1
		} else if len(bUnits) > 1 {
			// b has surrogate pair, a doesn't
			return -1
		}
	}

	// All compared characters are equal, compare by length
	if len(aRunes) < len(bRunes) {
		return -1
	}
	if len(aRunes) > len(bRunes) {
		return 1
	}
	return 0
}

// runeToCodeUnits converts a rune to its UTF-16 code unit(s).
func runeToCodeUnits(r rune) []uint16 {
	if r < 0x10000 {
		return []uint16{uint16(r)}
	}
	// Surrogate pair
	r -= 0x10000
	high := uint16(0xD800 + (r >> 10))
	low := uint16(0xDC00 + (r & 0x3FF))
	return []uint16{high, low}
}

// Size returns the number of entries.
func (sp *URLSearchParams) Size() int {
	return len(sp.entries)
}

// String returns the serialized query string (without leading "?").
func (sp *URLSearchParams) String() string {
	return encodeFormEncoded(sp.entries)
}

// ForEach calls the callback function for each entry.
func (sp *URLSearchParams) ForEach(callback func(value, key string)) {
	for _, entry := range sp.entries {
		callback(entry.value, entry.key)
	}
}

// Entries returns an iterator-like slice of [key, value] pairs.
func (sp *URLSearchParams) Entries() [][2]string {
	result := make([][2]string, len(sp.entries))
	for i, entry := range sp.entries {
		result[i] = [2]string{entry.key, entry.value}
	}
	return result
}

// Keys returns all keys in order.
func (sp *URLSearchParams) Keys() []string {
	result := make([]string, len(sp.entries))
	for i, entry := range sp.entries {
		result[i] = entry.key
	}
	return result
}

// Values returns all values in order.
func (sp *URLSearchParams) Values() []string {
	result := make([]string, len(sp.entries))
	for i, entry := range sp.entries {
		result[i] = entry.value
	}
	return result
}

// percentDecode decodes a percent-encoded string, leaving invalid sequences as-is.
// This follows the WHATWG URL Standard's percent-decode algorithm.
func percentDecode(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) {
			// Try to decode the percent-encoded byte
			hi := unhex(s[i+1])
			lo := unhex(s[i+2])
			if hi >= 0 && lo >= 0 {
				// Valid hex digits
				result.WriteByte(byte(hi<<4 | lo))
				i += 2
				continue
			}
		}
		// Not a valid percent-encoded sequence, keep as-is
		result.WriteByte(s[i])
	}

	return result.String()
}

// unhex returns the value of a hex digit, or -1 if invalid.
func unhex(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c - 'a' + 10)
	case c >= 'A' && c <= 'F':
		return int(c - 'A' + 10)
	}
	return -1
}

// parseFormEncoded parses an application/x-www-form-urlencoded string.
func parseFormEncoded(s string) []urlParam {
	entries := make([]urlParam, 0)

	if s == "" {
		return entries
	}

	pairs := strings.Split(s, "&")
	for _, pair := range pairs {
		if pair == "" {
			continue
		}

		var key, value string
		if idx := strings.Index(pair, "="); idx >= 0 {
			key = pair[:idx]
			value = pair[idx+1:]
		} else {
			key = pair
			value = ""
		}

		// Decode + as space, then percent-decode
		key = strings.ReplaceAll(key, "+", " ")
		value = strings.ReplaceAll(value, "+", " ")

		// Use custom percent decoder that handles invalid sequences
		decodedKey := percentDecode(key)
		decodedValue := percentDecode(value)

		entries = append(entries, urlParam{
			key:   decodedKey,
			value: decodedValue,
		})
	}

	return entries
}

// encodeFormEncoded serializes entries to application/x-www-form-urlencoded format.
func encodeFormEncoded(entries []urlParam) string {
	if len(entries) == 0 {
		return ""
	}

	parts := make([]string, len(entries))
	for i, entry := range entries {
		// Use custom encoding that matches WHATWG spec
		// (encodes space as +, and uses specific character set)
		encodedKey := formEncode(entry.key)
		encodedValue := formEncode(entry.value)
		parts[i] = encodedKey + "=" + encodedValue
	}

	return strings.Join(parts, "&")
}

// formEncode encodes a string for application/x-www-form-urlencoded.
// This follows the WHATWG URL Standard encoding rules.
// The string is first converted to UTF-8 bytes, then each byte is encoded.
func formEncode(s string) string {
	var builder strings.Builder
	builder.Grow(len(s) * 3) // worst case: all characters need encoding

	// Convert to bytes (UTF-8)
	bytes := []byte(s)

	for _, c := range bytes {
		switch {
		case c == ' ':
			builder.WriteByte('+')
		case c == '*' || c == '-' || c == '.' || c == '_':
			// These characters are not encoded per WHATWG spec
			builder.WriteByte(c)
		case c >= '0' && c <= '9':
			builder.WriteByte(c)
		case c >= 'A' && c <= 'Z':
			builder.WriteByte(c)
		case c >= 'a' && c <= 'z':
			builder.WriteByte(c)
		default:
			// Percent-encode
			builder.WriteByte('%')
			builder.WriteByte(hexDigit(c >> 4))
			builder.WriteByte(hexDigit(c & 0x0F))
		}
	}

	return builder.String()
}

func hexDigit(n byte) byte {
	if n < 10 {
		return '0' + n
	}
	return 'A' + n - 10
}

