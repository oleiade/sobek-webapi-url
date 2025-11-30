// Package url implements WHATWG URL API bindings for Sobek runtimes.
//
// This package provides the core URL and URLSearchParams types along with
// their Sobek bindings. The implementation follows the WHATWG URL Standard
// (https://url.spec.whatwg.org/) with the following supported features:
//
//   - URL parsing and serialization for http, https, ws, wss, ftp, and file schemes
//   - All standard URL component accessors (href, protocol, host, hostname,
//     port, pathname, search, hash, origin, username, password)
//   - URLSearchParams with full manipulation API (append, delete, get, getAll,
//     has, set, sort, forEach, entries, keys, values)
//   - Static URL.canParse() and URL.parse() methods
//   - Proper synchronization between URL.search and URL.searchParams
//   - URLSearchParams iteration via Symbol.iterator
//
// # Usage
//
// To register the URL and URLSearchParams constructors in a Sobek runtime:
//
//	rt := sobek.New()
//	if err := url.RegisterRuntime(rt); err != nil {
//	    log.Fatal(err)
//	}
//
// # Known Limitations
//
//   - Blob URLs are not supported
//   - Some edge-case Unicode/punycode behaviors may differ from browsers
//   - Origin computation for non-standard schemes returns "null"
//   - Base URL validation is more lenient than WHATWG (uses Go's net/url)
//   - Data URLs with opaque paths may not be fully supported
//   - URLSearchParams iterators are not live (don't reflect mutations during iteration)
//
// # Go API invariants
//
// The exported Go types primarily exist to back the Sobek bindings. They are
// intentionally small wrappers around Go's net/url so that they can be used in
// tests and, when needed, inside Grafana/k6 integrations. Always construct a
// URL by calling NewURL (and URLSearchParams with the constructors in this
// package); doing so guarantees that:
//
//   - URL.inner is non-nil and carries the parsed WHATWG representation.
//   - URL.searchParams is non-nil and bidirectionally synced with URL.inner.
//   - URLSearchParams.owner is either nil or points to the URL that must be
//     kept in sync.
//
// Consuming code should treat URL and URLSearchParams as opaque wrappers and
// prefer the provided accessors instead of manipulating fields directly.
package url
