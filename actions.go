package ds

import (
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// Option types (sealed interface pattern)
//
// The option interface is unexported with a private isOption method,
// ensuring only sseOption can satisfy it. This prevents arbitrary values
// from being misinterpreted during variadic arg partitioning in sseAction.
// ---------------------------------------------------------------------------

// option is the sealed interface for SSE action options.
type option interface {
	isOption()
}

// sseOption is a key-value pair for an SSE action options object.
type sseOption struct {
	key   string
	value string
	raw   bool // true = unquoted value, false = single-quoted value
}

func (sseOption) isOption() {}

// Opt creates an SSE action option with a single-quoted string value.
//
//	ds.Opt("requestCancellation", "disabled") // -> requestCancellation: 'disabled'
//	ds.Opt("contentType", "json")             // -> contentType: 'json'
//	ds.Opt("retry", "never")                  // -> retry: 'never'
func Opt(key, value string) option {
	return sseOption{key: key, value: value, raw: false}
}

// OptRaw creates an SSE action option with a raw (unquoted) value.
// Use for booleans, numbers, and object/regex values.
//
//	ds.OptRaw("openWhenHidden", "true")              // -> openWhenHidden: true
//	ds.OptRaw("retryMaxCount", "10")                 // -> retryMaxCount: 10
//	ds.OptRaw("filterSignals", "{include: /^foo/}")  // -> filterSignals: {include: /^foo/}
func OptRaw(key, value string) option {
	return sseOption{key: key, value: value, raw: true}
}

// ===========================================================================
// SSE Action Expression Builders
//
// Each function returns a Datastar action expression string (e.g. "@get('/url')").
// URL formatting uses fmt.Sprintf when non-option args are provided.
// Options are appended as a JS object when present.
//
// See https://data-star.dev/reference/actions
// ===========================================================================

// Get builds a @get SSE action expression.
//
//	ds.Get("/api/updates")
//	// -> "@get('/api/updates')"
//
//	ds.Get("/api/todos/%d", id)
//	// -> "@get('/api/todos/42')"
//
//	ds.Get("/api/updates", ds.Opt("requestCancellation", "disabled"))
//	// -> "@get('/api/updates',{requestCancellation: 'disabled'})"
//
//	ds.Get("/api/todos/%d", id, ds.Opt("openWhenHidden", "true"))
//	// -> "@get('/api/todos/42',{openWhenHidden: 'true'})"
func Get(urlFormat string, args ...any) string {
	return sseAction(actionGet, urlFormat, args)
}

// Post builds a @post SSE action expression.
//
//	ds.Post("/api/workcenters")
//	// -> "@post('/api/workcenters')"
func Post(urlFormat string, args ...any) string {
	return sseAction(actionPost, urlFormat, args)
}

// Put builds a @put SSE action expression.
//
//	ds.Put("/api/todos/%d", id)
//	// -> "@put('/api/todos/42')"
func Put(urlFormat string, args ...any) string {
	return sseAction(actionPut, urlFormat, args)
}

// Patch builds a @patch SSE action expression.
//
//	ds.Patch("/api/workcenters/pagesize")
//	// -> "@patch('/api/workcenters/pagesize')"
func Patch(urlFormat string, args ...any) string {
	return sseAction(actionPatch, urlFormat, args)
}

// Delete builds a @delete SSE action expression.
//
//	ds.Delete("/api/todos/%d", id)
//	// -> "@delete('/api/todos/42')"
func Delete(urlFormat string, args ...any) string {
	return sseAction(actionDelete, urlFormat, args)
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// sseAction is the shared builder for all SSE action expressions.
// It partitions args into fmt.Sprintf format args and sseOption values,
// then builds the expression string.
func sseAction(verb, urlFormat string, args []any) string {
	var fmtArgs []any
	var opts []sseOption
	for _, a := range args {
		if o, ok := a.(sseOption); ok {
			opts = append(opts, o)
		} else {
			fmtArgs = append(fmtArgs, a)
		}
	}
	url := fmt.Sprintf(urlFormat, fmtArgs...)
	if len(opts) == 0 {
		return fmt.Sprintf("@%s('%s')", verb, url)
	}
	return fmt.Sprintf("@%s('%s',%s)", verb, url, buildOpts(opts))
}

// buildOpts builds a JavaScript options object from sseOption values.
// Quoted values produce {key: 'val'}. Raw values produce {key: val}.
func buildOpts(opts []sseOption) string {
	b := sharedBuilderPool.Get().(*strings.Builder)
	defer func() {
		b.Reset()
		sharedBuilderPool.Put(b)
	}()

	b.WriteByte('{')
	for i, o := range opts {
		if i > 0 {
			b.WriteString(", ")
		}
		if o.raw {
			fmt.Fprintf(b, "%s: %s", o.key, o.value)
		} else {
			fmt.Fprintf(b, "%s: '%s'", o.key, o.value)
		}
	}
	b.WriteByte('}')
	return b.String()
}
