# sobek-webapi-url

A WHATWG-style `URL` and `URLSearchParams` Web API implementation for the
[Sobek](https://github.com/grafana/sobek) JavaScript runtime (used by
[k6](https://github.com/grafana/k6)).

This package lets JavaScript running inside Sobek use the standard `URL` and
`URLSearchParams` globals, and is validated against a large subset of the
upstream [Web Platform Tests](https://github.com/web-platform-tests/wpt)
for the URL standard.

## Installation

```bash
go get github.com/oleiade/sobek-webapi-url
```

## Usage

Register the URL Web API in your Sobek runtime and then use it from JavaScript:

```go
package main

import (
	"github.com/grafana/sobek"
	sobekurl "github.com/oleiade/sobek-webapi-url"
)

func main() {
	rt := sobek.New()

	// Register URL and URLSearchParams globally
	if err := sobekurl.RegisterGlobally(rt); err != nil {
		panic(err)
	}

	// Now you can use URL and URLSearchParams in JavaScript
	rt.RunString(`
        const url = new URL('https://example.com/path?foo=bar#hash');
        console.log(url.hostname);  // "example.com"
        console.log(url.searchParams.get('foo'));  // "bar"

        url.searchParams.set('baz', 'qux');
        console.log(url.href);  // "https://example.com/path?foo=bar&baz=qux#hash"
    `)
}
```

Once registered, any JavaScript code executed in the runtime can construct
and manipulate URLs using the familiar browser-style API.

## Supported Features

### URL

- **Constructor**: `new URL(url, base?)`
- **Static methods**: `URL.canParse(url, base?)`, `URL.parse(url, base?)`
- **Properties** (read/write):
  - `href`, `protocol`, `username`, `password`
  - `host`, `hostname`, `port`, `pathname`
  - `search`, `hash`
- **Properties** (read-only):
  - `origin`, `searchParams`
- **Methods**: `toString()`, `toJSON()`

### URLSearchParams

- **Constructor**: `new URLSearchParams(init?)`
  - Accepts: string, URLSearchParams, array of pairs, or object
- **Methods**:
  - `append(name, value)`
  - `delete(name, value?)`
  - `get(name)` / `getAll(name)`
  - `has(name, value?)`
  - `set(name, value)`
  - `sort()`
  - `toString()`
  - `forEach(callback, thisArg?)`
  - `entries()` / `keys()` / `values()`
- **Properties**: `size`
- **Iterable**: supports `for...of` loops

## Known Limitations

This implementation uses Go's `net/url` package under the hood, which has some differences from the WHATWG URL Standard:

1. **Base URL validation**: Go's URL parser is more lenient than WHATWG. For example, `aaa:b` is considered a valid absolute URL in Go but not in WHATWG (which requires a path separator after non-special schemes).

2. **Opaque paths**: Data URLs and other URLs with opaque paths may not be fully supported.

3. **Live iterators**: The URLSearchParams iterator does not reflect mutations made during iteration (the WHATWG spec requires live iterators).

4. **Punycode/IDNA**: International domain name handling may differ from browser implementations.

## Testing

This implementation is validated against Web Platform Tests (WPT) for the URL
and URLSearchParams APIs.

```bash
go test -v ./url/...
```

### Test Status

| Test Suite | Status |
|------------|--------|
| URLSearchParams.append | ✅ Pass |
| URLSearchParams.delete | ⚠️ Partial (data: URL edge cases) |
| URLSearchParams.get | ✅ Pass |
| URLSearchParams.getAll | ✅ Pass |
| URLSearchParams.has | ✅ Pass |
| URLSearchParams.set | ✅ Pass |
| URLSearchParams.sort | ✅ Pass |
| URLSearchParams.size | ✅ Pass |
| URLSearchParams.stringifier | ✅ Pass |
| URLSearchParams.forEach | ⚠️ Partial (live iterator edge cases) |
| URLSearchParams constructor | ⚠️ Partial (DOMException branding) |
| URL.searchParams integration | ✅ Pass |
| URL.canParse | ⚠️ Partial (WHATWG spec differences) |
| URL.parse | ⚠️ Partial (WHATWG spec differences) |
| URL.toJSON | ✅ Pass |

## WPT test files and `wptsync`

The WPT test files are vendored in the `wpt/` directory. The `wpt.json` file
specifies which tests are included, the upstream WPT commit, and any patches
applied.

The test set is managed with a small helper tool called `wptsync`, which reads
`wpt.json` and (re)populates the `wpt/` directory from the upstream WPT
repository. To change which tests are synced, edit `wpt.json` and then use
`wptsync` to refresh the local copy (see the `wptsync` documentation for
installation and usage details).

## License

This project is licensed under the same terms as the Sobek project.

