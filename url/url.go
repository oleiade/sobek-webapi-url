package url

import (
	"net/url"
	"strings"
)

// URL represents a WHATWG-style URL.
//
// It wraps Go's net/url.URL and provides WHATWG-compatible
// property accessors and setters.
type URL struct {
	// internal parsed URL
	inner *url.URL

	// searchParams is the attached URLSearchParams instance
	searchParams *URLSearchParams
}

// NewURL creates a new URL by parsing input relative to an optional base.
//
// If parsing fails, it returns an error that should be converted to a
// JavaScript TypeError when thrown.
func NewURL(input string, base string) (*URL, error) {
	var baseURL *url.URL
	var err error

	if base != "" {
		baseURL, err = url.Parse(base)
		if err != nil {
			return nil, NewError(TypeError, "Invalid URL")
		}
		// Validate that base is absolute
		if !baseURL.IsAbs() {
			return nil, NewError(TypeError, "Invalid URL")
		}
	}

	var parsed *url.URL
	if baseURL != nil {
		ref, err := url.Parse(input)
		if err != nil {
			return nil, NewError(TypeError, "Invalid URL")
		}
		parsed = baseURL.ResolveReference(ref)
	} else {
		parsed, err = url.Parse(input)
		if err != nil {
			return nil, NewError(TypeError, "Invalid URL")
		}
		// Without a base, the URL must be absolute
		if !parsed.IsAbs() {
			return nil, NewError(TypeError, "Invalid URL")
		}
	}

	// Validate scheme - reject empty scheme
	if parsed.Scheme == "" {
		return nil, NewError(TypeError, "Invalid URL")
	}

	u := &URL{inner: parsed}
	u.initSearchParams()

	return u, nil
}

// Parse attempts to parse input relative to base and returns the URL or nil.
// This is the implementation for the static URL.parse() method.
func Parse(input string, base string) *URL {
	u, err := NewURL(input, base)
	if err != nil {
		return nil
	}
	return u
}

// CanParse returns true if input can be parsed relative to base.
// This is the implementation for the static URL.canParse() method.
func CanParse(input string, base string) bool {
	_, err := NewURL(input, base)
	return err == nil
}

// initSearchParams initializes the searchParams field from the current query string.
func (u *URL) initSearchParams() {
	// Don't use NewURLSearchParamsFromString here because it strips leading '?'
	// but RawQuery might contain '?' as part of the actual query content
	u.searchParams = &URLSearchParams{
		entries: parseFormEncoded(u.inner.RawQuery),
		owner:   u,
	}
}

// syncFromSearchParams updates the URL's query string from the attached searchParams.
func (u *URL) syncFromSearchParams() {
	serialized := u.searchParams.String()
	u.inner.RawQuery = serialized
	// Clear ForceQuery when query becomes empty
	if serialized == "" {
		u.inner.ForceQuery = false
	}
}

// Href returns the full serialized URL.
func (u *URL) Href() string {
	return u.inner.String()
}

// SetHref replaces the entire URL by parsing the new href value.
func (u *URL) SetHref(href string) error {
	parsed, err := url.Parse(href)
	if err != nil {
		return NewError(TypeError, "Invalid URL")
	}
	if !parsed.IsAbs() {
		return NewError(TypeError, "Invalid URL")
	}
	u.inner = parsed
	// Update the existing searchParams object
	u.updateSearchParams(parsed.RawQuery)
	return nil
}

// Protocol returns the scheme followed by a colon (e.g., "https:").
func (u *URL) Protocol() string {
	return u.inner.Scheme + ":"
}

// SetProtocol sets the URL's scheme from a value like "https:" or "https".
func (u *URL) SetProtocol(protocol string) {
	// Strip trailing colon if present
	scheme := strings.TrimSuffix(protocol, ":")
	scheme = strings.ToLower(scheme)
	u.inner.Scheme = scheme
}

// Username returns the username portion of the URL.
func (u *URL) Username() string {
	if u.inner.User == nil {
		return ""
	}
	return u.inner.User.Username()
}

// SetUsername sets the username portion of the URL.
func (u *URL) SetUsername(username string) {
	if u.inner.User == nil {
		u.inner.User = url.User(username)
	} else {
		password, hasPassword := u.inner.User.Password()
		if hasPassword {
			u.inner.User = url.UserPassword(username, password)
		} else {
			u.inner.User = url.User(username)
		}
	}
}

// Password returns the password portion of the URL.
func (u *URL) Password() string {
	if u.inner.User == nil {
		return ""
	}
	password, _ := u.inner.User.Password()
	return password
}

// SetPassword sets the password portion of the URL.
func (u *URL) SetPassword(password string) {
	username := ""
	if u.inner.User != nil {
		username = u.inner.User.Username()
	}
	u.inner.User = url.UserPassword(username, password)
}

// Host returns the host and port (if non-default) combined.
func (u *URL) Host() string {
	return u.inner.Host
}

// SetHost sets the host (and optionally port) of the URL.
func (u *URL) SetHost(host string) {
	u.inner.Host = host
}

// Hostname returns just the hostname portion (without port).
func (u *URL) Hostname() string {
	return u.inner.Hostname()
}

// SetHostname sets the hostname portion without affecting the port.
func (u *URL) SetHostname(hostname string) {
	port := u.inner.Port()
	if port != "" {
		u.inner.Host = hostname + ":" + port
	} else {
		u.inner.Host = hostname
	}
}

// Port returns the port as a string, or empty if not specified.
func (u *URL) Port() string {
	return u.inner.Port()
}

// SetPort sets the port portion of the URL.
func (u *URL) SetPort(port string) {
	hostname := u.inner.Hostname()
	if port == "" {
		u.inner.Host = hostname
	} else {
		u.inner.Host = hostname + ":" + port
	}
}

// Pathname returns the path portion of the URL.
func (u *URL) Pathname() string {
	path := u.inner.Path
	if path == "" {
		return "/"
	}
	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

// SetPathname sets the path portion of the URL.
func (u *URL) SetPathname(pathname string) {
	u.inner.Path = pathname
}

// Search returns the query string including the leading "?" if non-empty.
func (u *URL) Search() string {
	if u.inner.RawQuery == "" {
		return ""
	}
	return "?" + u.inner.RawQuery
}

// SetSearch sets the query string (with or without leading "?").
func (u *URL) SetSearch(search string) {
	// Strip leading ? if present
	search = strings.TrimPrefix(search, "?")
	u.inner.RawQuery = search
	// Clear ForceQuery when query becomes empty
	if search == "" {
		u.inner.ForceQuery = false
	}
	// Update the existing searchParams object instead of creating a new one
	u.updateSearchParams(search)
}

// updateSearchParams updates the existing searchParams with new query string.
func (u *URL) updateSearchParams(query string) {
	// Clear existing entries
	u.searchParams.entries = u.searchParams.entries[:0]
	// Parse new query and add entries
	if query != "" {
		newEntries := parseFormEncoded(query)
		u.searchParams.entries = append(u.searchParams.entries, newEntries...)
	}
}

// SearchParams returns the URLSearchParams object for this URL.
func (u *URL) SearchParams() *URLSearchParams {
	return u.searchParams
}

// Hash returns the fragment including the leading "#" if non-empty.
func (u *URL) Hash() string {
	if u.inner.Fragment == "" {
		return ""
	}
	return "#" + u.inner.Fragment
}

// SetHash sets the fragment (with or without leading "#").
func (u *URL) SetHash(hash string) {
	// Strip leading # if present
	u.inner.Fragment = strings.TrimPrefix(hash, "#")
}

// Origin returns the origin of the URL.
//
// For http, https, ws, wss schemes, this returns "scheme://host".
// For file scheme, this returns "null" per spec.
// For other schemes, this returns "null".
func (u *URL) Origin() string {
	switch u.inner.Scheme {
	case "http", "https", "ws", "wss":
		return u.inner.Scheme + "://" + u.inner.Host
	case "ftp":
		return u.inner.Scheme + "://" + u.inner.Host
	default:
		// file: and other schemes return "null"
		return "null"
	}
}

// String returns the serialized URL (same as Href).
func (u *URL) String() string {
	return u.Href()
}

// ToJSON returns the serialized URL (same as Href).
func (u *URL) ToJSON() string {
	return u.Href()
}

