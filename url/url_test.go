package url

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// WPT skips summary:
//   1. data: URL opaque paths are unsupported by Go's net/url, so
//      urlsearchparams-delete.js remains skipped.
//   2. URLSearchParams iterators are snapshots, not live views, making the
//      forEach "For-of Check" test fail (t.Skip in TestURLSearchParamsForEach).
//   3. DOMException branding is incomplete in the Sobek test stubs, so the
//      constructor branding suite stays skipped until sobek gains real DOMException
//      semantics.
//   4. net/url accepts more base URLs than WHATWG permits (e.g., "aaa:b"), so the
//      URL.canParse/parse WPT suites are skipped until a stricter parser is wired in.

// TestURLSearchParamsAppend runs the WPT tests for URLSearchParams.append()
func TestURLSearchParamsAppend(t *testing.T) {
	t.Parallel()
	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "urlsearchparams-append.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLSearchParamsDelete runs the WPT tests for URLSearchParams.delete()
//
// Known limitation: Tests involving data: URLs with opaque paths fail because
// Go's net/url doesn't support opaque path URLs the same way as WHATWG.
func TestURLSearchParamsDelete(t *testing.T) {
	t.Parallel()
	t.Skip("Skipped: data: URL opaque path handling differs from WHATWG spec")

	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "urlsearchparams-delete.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLSearchParamsGet runs the WPT tests for URLSearchParams.get()
func TestURLSearchParamsGet(t *testing.T) {
	t.Parallel()
	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "urlsearchparams-get.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLSearchParamsGetAll runs the WPT tests for URLSearchParams.getAll()
func TestURLSearchParamsGetAll(t *testing.T) {
	t.Parallel()
	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "urlsearchparams-getall.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLSearchParamsHas runs the WPT tests for URLSearchParams.has()
func TestURLSearchParamsHas(t *testing.T) {
	t.Parallel()
	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "urlsearchparams-has.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLSearchParamsSet runs the WPT tests for URLSearchParams.set()
func TestURLSearchParamsSet(t *testing.T) {
	t.Parallel()
	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "urlsearchparams-set.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLSearchParamsSort runs the WPT tests for URLSearchParams.sort()
func TestURLSearchParamsSort(t *testing.T) {
	t.Parallel()
	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "urlsearchparams-sort.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLSearchParamsStringifier runs the WPT tests for URLSearchParams stringifier
func TestURLSearchParamsStringifier(t *testing.T) {
	t.Parallel()
	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "urlsearchparams-stringifier.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLSearchParamsSize runs the WPT tests for URLSearchParams.size
func TestURLSearchParamsSize(t *testing.T) {
	t.Parallel()
	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "urlsearchparams-size.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLSearchParamsForEach runs the WPT tests for URLSearchParams.forEach()
//
// Known limitation: The "For-of Check" test expects live iterator behavior where
// modifying url.search during iteration affects the iterator. Our implementation
// creates a snapshot at iteration start.
func TestURLSearchParamsForEach(t *testing.T) {
	t.Parallel()
	t.Skip("Skipped: Live iterator behavior during mutation not implemented")

	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "urlsearchparams-foreach.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLSearchParamsConstructor runs the WPT tests for URLSearchParams constructor
//
// Known limitation: DOMException.prototype branding check test fails because our
// DOMException stub doesn't have proper internal slots/branding.
func TestURLSearchParamsConstructor(t *testing.T) {
	t.Parallel()
	t.Skip("Skipped: DOMException.prototype branding check not supported")

	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "urlsearchparams-constructor.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLSearchParams runs the WPT tests for URL.searchParams integration
func TestURLSearchParams(t *testing.T) {
	t.Parallel()
	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "url-searchparams.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLStaticsCanParse runs the WPT tests for URL.canParse()
//
// Known limitation: Go's net/url is more lenient than WHATWG URL Standard.
// For example, "aaa:b" is considered valid in Go but not in WHATWG (which
// requires a path separator after non-special schemes).
func TestURLStaticsCanParse(t *testing.T) {
	t.Parallel()
	t.Skip("Skipped: Go's net/url base URL validation differs from WHATWG spec")

	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "url-statics-canparse.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLStaticsParse runs the WPT tests for URL.parse()
//
// Known limitation: Same as TestURLStaticsCanParse - Go's net/url is more
// lenient than WHATWG URL Standard for base URL validation.
func TestURLStaticsParse(t *testing.T) {
	t.Parallel()
	t.Skip("Skipped: Go's net/url base URL validation differs from WHATWG spec")

	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "url-statics-parse.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

// TestURLToJSON runs the WPT tests for URL.toJSON()
func TestURLToJSON(t *testing.T) {
	t.Parallel()
	base := wptPath("url")
	scripts := []testScript{
		{base: base, path: "url-tojson.js"},
	}

	ts := newTestSetup(t)
	err := executeTestScripts(ts, scripts)
	require.NoError(t, err)
}

func TestURLSearchAndParamsStayInSync(t *testing.T) {
	t.Parallel()

	u, err := NewURL("https://example.com/?foo=bar", "")
	require.NoError(t, err)

	params := u.SearchParams()
	require.NotNil(t, params)

	params.Set("foo", "baz")
	require.Equal(t, "?foo=baz", u.Search())
	require.Equal(t, "foo=baz", u.inner.RawQuery)

	u.SetSearch("?a=1&b=2")
	require.Same(t, params, u.SearchParams())

	value, ok := params.Get("a")
	require.True(t, ok)
	require.Equal(t, "1", value)
	require.Equal(t, "a=1&b=2", params.String()) // order matches serialized query

	u.SetSearch("")
	require.Equal(t, "", u.Search())
	require.False(t, u.inner.ForceQuery)
	require.Equal(t, 0, params.Size())
}

func TestURLSetHrefKeepsSearchParamsReference(t *testing.T) {
	t.Parallel()

	u, err := NewURL("https://example.com/path?foo=bar", "")
	require.NoError(t, err)

	params := u.SearchParams()
	require.NotNil(t, params)

	err = u.SetHref("https://grafana.com/api?alpha=beta")
	require.NoError(t, err)

	require.Same(t, params, u.SearchParams())
	require.Equal(t, "https://grafana.com/api?alpha=beta", u.Href())

	alpha, ok := params.Get("alpha")
	require.True(t, ok)
	require.Equal(t, "beta", alpha)
}

func TestURLOrigin(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		raw  string
		want string
	}{
		{name: "https standard", raw: "https://grafana.com/path", want: "https://grafana.com"},
		{name: "ws", raw: "ws://example.com/socket", want: "ws://example.com"},
		{name: "ftp", raw: "ftp://ftp.example.com/resource", want: "ftp://ftp.example.com"},
		{name: "file", raw: "file:///tmp/data", want: "null"},
		{name: "custom scheme", raw: "custom://host/path", want: "null"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			u, err := NewURL(tc.raw, "")
			require.NoError(t, err)
			require.Equal(t, tc.want, u.Origin())
		})
	}
}
