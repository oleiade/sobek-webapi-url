// Package sobekurl registers the URL Web API with Sobek runtimes.
//
// This package provides WHATWG-style URL and URLSearchParams globals
// for use in Sobek JavaScript runtimes (as used by k6). It supports
// the common subset of the URL standard including HTTP(S), WS(S), FTP,
// and file schemes with standard parsing, serialization, and manipulation.
//
// # Usage
//
//	rt := sobek.New()
//	if err := sobekurl.RegisterGlobally(rt); err != nil {
//	    log.Fatal(err)
//	}
//	// URL and URLSearchParams are now available in the runtime
//
// After registration, JavaScript code can use the standard URL API:
//
//	const url = new URL('https://example.com/path?query=value#hash');
//	console.log(url.hostname);  // "example.com"
//	console.log(url.searchParams.get('query'));  // "value"
//
//	const params = new URLSearchParams('foo=1&bar=2');
//	params.append('baz', '3');
//	console.log(params.toString());  // "foo=1&bar=2&baz=3"
//
// # Supported Features
//
//   - URL constructor with optional base URL
//   - URL.canParse() and URL.parse() static methods
//   - All standard URL properties (href, protocol, host, hostname, port,
//     pathname, search, hash, origin, username, password, searchParams)
//   - URLSearchParams with append, delete, get, getAll, has, set, sort,
//     forEach, entries, keys, values, and size
//   - Proper bidirectional synchronization between URL.search and URL.searchParams
//
// # Known Limitations
//
//   - Blob URLs are not supported
//   - Some WHATWG edge cases may differ (uses Go's net/url internally)
//   - Data URLs with opaque paths may not be fully supported
//
// For more details, see the url subpackage documentation.
package sobekurl

import (
	"github.com/grafana/sobek"

	"github.com/oleiade/sobek-webapi-url/url"
)

// URL is a re-export of url.URL for consumers such as k6 modules.
type URL = url.URL

var (
	// ExtractURL extracts a url.URL from a Sobek Value.
	//nolint:gochecknoglobals // Re-exported for convenience
	ExtractURL = url.ExtractURL
	// ParseURLArgument parses a URL argument from a Sobek Value.
	//nolint:gochecknoglobals // Re-exported for convenience
	ParseURLArgument = url.ParseURLArgument
)

// RegisterGlobally exposes the URL and URLSearchParams constructors
// in the provided sobek runtime.
func RegisterGlobally(rt *sobek.Runtime) error {
	return url.RegisterRuntime(rt)
}
