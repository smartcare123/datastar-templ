// Package ds provides type-safe Datastar attribute helpers for templ templates.
//
// Every function returns templ.Attributes, so you can spread directly in templ:
//
//	<button { ds.OnClick("$open = true")... }>Open</button>
//	<div { ds.Show("$visible")... }>Content</div>
//	<input { ds.Bind("name")... } />
//
// For multiple attributes on one element, use Merge:
//
//	<div { ds.Merge(ds.Show("$open"), ds.OnClick("$open = false"))... }>
//
// See https://data-star.dev/reference/attributes
package ds

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/a-h/templ"
)

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// Modifier is a Datastar attribute modifier suffix.
// Double-underscore modifiers (e.g. __debounce) and dot-tag modifiers (e.g. .leading)
// are concatenated onto attribute names.
type Modifier string

// Filter is used by attributes that accept include/exclude regex patterns.
type Filter struct {
	Include string
	Exclude string
}

// ---------------------------------------------------------------------------
// Modifier helpers
// ---------------------------------------------------------------------------

// Duration returns a ".{N}ms" modifier tag, rounded to the nearest millisecond.
//
// Panics if the duration is negative.
//
// Example:
//
//	ds.OnClick("handler()", ds.ModDebounce, ds.Duration(300*time.Millisecond))
func Duration(d time.Duration) Modifier {
	if d < 0 {
		panic(fmt.Sprintf("ds: duration must not be negative, got %v", d))
	}
	return Modifier(fmt.Sprintf(".%dms", d.Round(time.Millisecond).Milliseconds()))
}

// Ms returns a ".{n}ms" modifier tag. Shorthand for Duration when you have a
// raw millisecond value.
//
// Panics if n is negative.
//
// Example:
//
//	ds.OnInput("search()", ds.ModDebounce, ds.Ms(300))
func Ms(n int) Modifier {
	if n < 0 {
		panic(fmt.Sprintf("ds: milliseconds must not be negative, got %d", n))
	}
	return Modifier(fmt.Sprintf(".%dms", n))
}

// Seconds returns a ".{n}s" modifier tag.
//
// Panics if n is negative.
//
// Example:
//
//	ds.OnInterval("poll()", ds.Seconds(5))
func Seconds(n int) Modifier {
	if n < 0 {
		panic(fmt.Sprintf("ds: seconds must not be negative, got %d", n))
	}
	return Modifier(fmt.Sprintf(".%ds", n))
}

// Threshold returns a visibility percentage modifier tag for the __threshold modifier.
// The value must be between 0.0 (exclusive) and 1.0 (inclusive).
//
// Panics if t is <= 0 or > 1.
//
// Example:
//
//	ds.OnIntersect("loadMore()", ds.ModThreshold, ds.Threshold(0.5))  // 50% visible
func Threshold(t float64) Modifier {
	if t <= 0 || t > 1 {
		panic(fmt.Sprintf("ds: threshold must be between 0.0 (exclusive) and 1.0 (inclusive), got %v", t))
	}
	if t == 1 {
		return Modifier(".100")
	}
	return Modifier(strings.TrimPrefix(fmt.Sprintf("%.2f", t), "0"))
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// mods concatenates modifiers into a single string.
func mods(modifiers []Modifier) string {
	if len(modifiers) == 0 {
		return ""
	}

	b := sharedBuilderPool.Get().(*strings.Builder)
	defer func() {
		b.Reset()
		sharedBuilderPool.Put(b)
	}()

	for _, m := range modifiers {
		b.WriteString(string(m))
	}
	return b.String()
}

// on builds "data-on:{event}{modifiers}".
func on(event string, modifiers []Modifier) string {
	return prefixOn + event + mods(modifiers)
}

// plugin builds "data-{name}{modifiers}" for plugin-based attributes (hyphenated, not colon).
func plugin(name string, modifiers []Modifier) string {
	return prefix + name + mods(modifiers)
}

// keyed builds "data-{name}:{key}{modifiers}".
func keyed(name, key string, modifiers []Modifier) string {
	return prefix + name + sepColon + key + mods(modifiers)
}

// val returns templ.Attributes with a single key-value pair.
func val(name, value string) templ.Attributes {
	return templ.Attributes{name: value}
}

// boolAttr returns templ.Attributes with a boolean (valueless) attribute.
func boolAttr(name string) templ.Attributes {
	return templ.Attributes{name: true}
}

// toPairs builds a JavaScript object literal from key-value string pairs.
// The valueFmt function transforms each value (e.g. identity for plain objects,
// arrow-function wrapper for computed signals).
// Panics if pairs is odd-length.
func toPairs(pairs []string, panicMsg string, valueFmt func(string) string) string {
	if len(pairs)%2 == 1 {
		panic(panicMsg)
	}
	var b strings.Builder
	b.WriteByte('{')
	for i := 0; i < len(pairs); i += 2 {
		if i > 0 {
			b.WriteString(", ")
		}
		// Quote the key to support hyphens and special characters in JavaScript object literals
		fmt.Fprintf(&b, "'%s': %s", pairs[i], valueFmt(pairs[i+1]))
	}
	b.WriteByte('}')
	return b.String()
}

// toObject builds a JavaScript object literal: "{k1: v1, k2: v2}".
func toObject(pairs []string) string {
	return toPairs(pairs, "ds: each key must have a value", func(v string) string { return v })
}

// toComputed builds a JavaScript object with arrow functions:
// "{k1: () => v1, k2: () => v2}".
func toComputed(pairs []string) string {
	return toPairs(pairs, "ds: each computed signal name must have an expression", func(v string) string {
		return "() => " + v
	})
}

// toFilter builds a filter object: "{include: /re/, exclude: /re/}".
func toFilter(f Filter) string {
	var b strings.Builder
	b.WriteByte('{')
	if f.Include != "" {
		fmt.Fprintf(&b, "include: %s", f.Include)
		if f.Exclude != "" {
			b.WriteString(", ")
		}
	}
	if f.Exclude != "" {
		fmt.Fprintf(&b, "exclude: %s", f.Exclude)
	}
	b.WriteByte('}')
	return b.String()
}

// toSignals JSON-marshals a signal map.
func toSignals(signals map[string]any) string {
	data, err := json.Marshal(signals)
	if err != nil {
		panic(fmt.Sprintf("ds: failed to marshal signals: %v", err))
	}
	return string(data)
}

// ---------------------------------------------------------------------------
// Merge
// ---------------------------------------------------------------------------

// Merge combines multiple templ.Attributes into one. Use when you need
// multiple ds attributes on a single element:
//
//	<div { ds.Merge(ds.Show("$open"), ds.OnClick("$open = false"))... }>
func Merge(attrs ...templ.Attributes) templ.Attributes {
	size := 0
	for _, a := range attrs {
		size += len(a)
	}
	m := make(templ.Attributes, size)
	for _, a := range attrs {
		for k, v := range a {
			m[k] = v
		}
	}
	return m
}
