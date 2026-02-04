package ds

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/a-h/templ"
)

// ===========================================================================
// Datastar Attributes
// See https://data-star.dev/reference/attributes
// ===========================================================================

// ---------------------------------------------------------------------------
// data-signals
// ---------------------------------------------------------------------------

// Signal represents a typed key-value pair for signals.
//
// Do not construct Signal directly; use the type-safe helper functions instead:
// Int(), String(), Bool(), Float(), or JSON().
type Signal struct {
	key   string
	value string
}

// Int creates an integer signal.
func Int(key string, value int) Signal {
	return Signal{key, strconv.Itoa(value)}
}

// String creates a string signal (properly quoted for JavaScript).
func String(key string, value string) Signal {
	return Signal{key, strconv.Quote(value)}
}

// Bool creates a boolean signal.
func Bool(key string, value bool) Signal {
	return Signal{key, strconv.FormatBool(value)}
}

// Float creates a float signal.
func Float(key string, value float64) Signal {
	return Signal{key, strconv.FormatFloat(value, 'f', -1, 64)}
}

// JSON creates a signal from any value using JSON marshaling.
// Use this for complex types like arrays, objects, etc.
//
// Panics if the value cannot be marshaled to JSON. Common causes include:
//   - Circular references (e.g., a struct that references itself)
//   - Channels, functions, or other non-JSON types
//   - Values that implement json.Marshaler but return an error
//
// For production code with untrusted input, consider using JSONSafe instead.
//
// Example:
//
//	ds.Signals(ds.JSON("user", map[string]any{"name": "Alice", "age": 30}))
func JSON(key string, value any) Signal {
	data, err := json.Marshal(value)
	if err != nil {
		panic("ds: failed to marshal JSON signal: " + err.Error())
	}
	return Signal{key, string(data)}
}

// JSONSafe creates a signal from any value using JSON marshaling.
// Returns an error instead of panicking if marshaling fails.
//
// Use this for untrusted input or when you need graceful error handling.
//
// Example:
//
//	sig, err := ds.JSONSafe("user", userData)
//	if err != nil {
//	    return err // Handle gracefully
//	}
//	ds.Signals(sig)
func JSONSafe(key string, value any) (Signal, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return Signal{}, fmt.Errorf("ds: failed to marshal JSON signal: %w", err)
	}
	return Signal{key, string(data)}, nil
}

// PairItem represents a key-value expression binding for attributes.
// Used by Class, Computed, Attr, and Style functions.
//
// Do not construct PairItem directly; use Pair() or P() helper functions instead.
type PairItem struct {
	key  string
	expr string
}

// Pair creates a key-value expression binding.
// This is the recommended helper for all attribute bindings.
//
//	ds.Class(ds.Pair("hidden", "$isHidden"))
//	ds.Computed(ds.Pair("total", "$price * $qty"))
//	ds.Attr(ds.Pair("disabled", "$loading"))
//	ds.Style(ds.Pair("color", "$textColor"))
func Pair(key, expr string) PairItem {
	return PairItem{key, expr}
}

// P is a shorthand alias for Pair.
// Use this if you prefer more concise template code.
func P(key, expr string) PairItem {
	return PairItem{key, expr}
}

// sharedBuilderPool is a sync.Pool for reusing strings.Builder instances across
// all attribute-building functions. This reduces memory allocations and improves
// performance by recycling builders instead of creating new ones for each operation.
var sharedBuilderPool = sync.Pool{
	New: func() interface{} {
		b := new(strings.Builder)
		b.Grow(128) // Initial capacity for typical attribute values
		return b
	},
}

// buildPairs is a generic helper for building JavaScript object literals from key-value pairs.
// The valueFormatter function determines how each value is wrapped (e.g., "expr" vs "() => expr").
func buildPairs(pairs []PairItem, valueFormatter func(string) string) string {
	b := sharedBuilderPool.Get().(*strings.Builder)
	defer func() {
		b.Reset()
		sharedBuilderPool.Put(b)
	}()

	// Calculate exact capacity
	capacity := 2 // {}
	for i, pair := range pairs {
		if i > 0 {
			capacity += 2 // ", "
		}
		formatted := valueFormatter(pair.expr)
		capacity += 1 + len(pair.key) + 4 + len(formatted) // "'key': formatted"
	}

	if b.Cap() < capacity {
		b.Grow(capacity - b.Len())
	}

	// Build the object
	b.WriteByte('{')
	for i, pair := range pairs {
		if i > 0 {
			b.WriteByte(',')
			b.WriteByte(' ')
		}
		b.WriteByte('\'')
		b.WriteString(pair.key)
		b.WriteString("': ")
		b.WriteString(valueFormatter(pair.expr))
	}
	b.WriteByte('}')

	return b.String()
}

// Signals patches one or more signals using typed helpers.
//
//	{ ds.Signals(ds.Int("count", 0), ds.String("msg", "hello"))... }
//
// See https://data-star.dev/reference/attributes#data-signals
func Signals(signals ...Signal) templ.Attributes {
	b := sharedBuilderPool.Get().(*strings.Builder)
	defer func() {
		b.Reset()
		sharedBuilderPool.Put(b)
	}()

	// Calculate exact capacity needed
	capacity := 2 // {}
	for i, sig := range signals {
		if i > 0 {
			capacity += 2 // ", "
		}
		capacity += len(sig.key) + 2 + len(sig.value) // "key: value"
	}

	// Grow if needed
	if b.Cap() < capacity {
		b.Grow(capacity - b.Len())
	}

	// Build string with byte-level optimizations
	b.WriteByte('{')
	for i, sig := range signals {
		if i > 0 {
			b.WriteByte(',')
			b.WriteByte(' ')
		}
		b.WriteString(sig.key)
		b.WriteByte(':')
		b.WriteByte(' ')
		b.WriteString(sig.value)
	}
	b.WriteByte('}')

	return templ.Attributes{"data-signals": b.String()}
}

// SignalsJSON patches signals using a pre-built JSON string value.
// Useful when you've already serialized the signals.
//
//	{ ds.SignalsJSON(myJSONString)... }
func SignalsJSON(jsonStr string, modifiers ...Modifier) templ.Attributes {
	return val(plugin(attrSignals, modifiers), jsonStr)
}

// SignalKey patches a single signal using keyed syntax: data-signals:{name}.
//
//	{ ds.SignalKey("foo", "1")... }
//
// See https://data-star.dev/reference/attributes#data-signals
func SignalKey(name, expr string, modifiers ...Modifier) templ.Attributes {
	return val(keyed(attrSignals, name, modifiers), expr)
}

// ---------------------------------------------------------------------------
// data-computed
// ---------------------------------------------------------------------------

// Computed creates read-only computed signals using typed pairs.
// Each pair is wrapped in an arrow function.
//
//	{ ds.Computed(ds.Pair("total", "$price * $qty"))... }
//
// See https://data-star.dev/reference/attributes#data-computed
func Computed(pairs ...PairItem) templ.Attributes {
	obj := buildPairs(pairs, func(expr string) string {
		return "() => " + expr
	})
	return templ.Attributes{"data-computed": obj}
}

// ComputedKey creates a single computed signal using keyed syntax.
//
//	{ ds.ComputedKey("total", "$price * $qty")... }
func ComputedKey(name, expr string, modifiers ...Modifier) templ.Attributes {
	return val(keyed(attrComputed, name, modifiers), expr)
}

// ---------------------------------------------------------------------------
// data-on-intersect / data-on-interval / data-on-signal-patch (plugin watchers)
// ---------------------------------------------------------------------------

// OnIntersect runs an expression when the element intersects with the viewport.
//
//	{ ds.OnIntersect("$visible = true", ds.ModOnce, ds.ModFull)... }
//
// See https://data-star.dev/reference/attributes#data-on-intersect
func OnIntersect(expr string, modifiers ...Modifier) templ.Attributes {
	return val(plugin(attrOnIntersect, modifiers), expr)
}

// OnInterval runs an expression at a regular interval (default: 1s).
//
//	{ ds.OnInterval("$count++", ds.ModDuration, ds.Ms(500))... }
//
// See https://data-star.dev/reference/attributes#data-on-interval
func OnInterval(expr string, modifiers ...Modifier) templ.Attributes {
	return val(plugin(attrOnInterval, modifiers), expr)
}

// OnSignalPatch runs an expression whenever signals are patched.
//
//	{ ds.OnSignalPatch("console.log(patch)")... }
//
// See https://data-star.dev/reference/attributes#data-on-signal-patch
func OnSignalPatch(expr string, modifiers ...Modifier) templ.Attributes {
	return val(plugin(attrOnSignalPatch, modifiers), expr)
}

// OnSignalPatchFilter filters which signals trigger data-on-signal-patch.
//
//	{ ds.OnSignalPatchFilter(ds.Filter{Include: "/^counter$/"})... }
//
// See https://data-star.dev/reference/attributes#data-on-signal-patch-filter
func OnSignalPatchFilter(filter Filter) templ.Attributes {
	return val(prefix+attrOnSignalPatchFilt, toFilter(filter))
}

// ---------------------------------------------------------------------------
// data-bind
// ---------------------------------------------------------------------------

// Bind creates a two-way data binding between a signal and an element's value.
//
//	<input { ds.Bind("name")... } />
//	<input { ds.Bind("table.search")... } />
//
// See https://data-star.dev/reference/attributes#data-bind
func Bind(name string, modifiers ...Modifier) templ.Attributes {
	return boolAttr(keyed(attrBind, name, modifiers))
}

// BindExpr creates a two-way binding using value syntax.
//
//	<input { ds.BindExpr("name")... } />
func BindExpr(name string) templ.Attributes {
	return val(prefix+attrBind, name)
}

// ---------------------------------------------------------------------------
// data-text
// ---------------------------------------------------------------------------

// Text binds the text content of an element to an expression.
//
//	<span { ds.Text("$count")... }></span>
//
// See https://data-star.dev/reference/attributes#data-text
func Text(expr string) templ.Attributes {
	return val(prefix+attrText, expr)
}

// ---------------------------------------------------------------------------
// data-show
// ---------------------------------------------------------------------------

// Show shows or hides an element based on a boolean expression.
//
//	<div { ds.Show("$visible")... }>Content</div>
//
// See https://data-star.dev/reference/attributes#data-show
func Show(expr string) templ.Attributes {
	return val(prefix+attrShow, expr)
}

// ---------------------------------------------------------------------------
// data-class
// ---------------------------------------------------------------------------

// Class adds/removes CSS classes using typed pairs.
//
//	{ ds.Class(ds.Pair("hidden", "$isHidden"), ds.Pair("font-bold", "$isBold"))... }
//
// See https://data-star.dev/reference/attributes#data-class
func Class(pairs ...PairItem) templ.Attributes {
	obj := buildPairs(pairs, func(expr string) string { return expr })
	return templ.Attributes{"data-class": obj}
}

// ClassKey adds/removes a single CSS class using keyed syntax.
//
//	{ ds.ClassKey("font-bold", "$isBold")... }
func ClassKey(name, expr string, modifiers ...Modifier) templ.Attributes {
	return val(keyed(attrClass, name, modifiers), expr)
}

// ---------------------------------------------------------------------------
// data-attr
// ---------------------------------------------------------------------------

// Attr sets HTML attributes using typed pairs.
//
//	{ ds.Attr(ds.Pair("title", "$tooltip"), ds.Pair("disabled", "$loading"))... }
//
// See https://data-star.dev/reference/attributes#data-attr
func Attr(pairs ...PairItem) templ.Attributes {
	obj := buildPairs(pairs, func(expr string) string { return expr })
	return templ.Attributes{"data-attr": obj}
}

// AttrKey sets a single HTML attribute using keyed syntax.
//
//	{ ds.AttrKey("disabled", "$loading")... }
//	{ ds.AttrKey("title", "'Theme: ' + $theme")... }
func AttrKey(name, expr string, modifiers ...Modifier) templ.Attributes {
	return val(keyed(attrAttr, name, modifiers), expr)
}

// ---------------------------------------------------------------------------
// data-style
// ---------------------------------------------------------------------------

// Style sets inline CSS styles using typed pairs.
//
//	{ ds.Style(ds.Pair("display", "$hiding && 'none'"), ds.Pair("color", "$textColor"))... }
//
// See https://data-star.dev/reference/attributes#data-style
func Style(pairs ...PairItem) templ.Attributes {
	obj := buildPairs(pairs, func(expr string) string { return expr })
	return templ.Attributes{"data-style": obj}
}

// StyleKey sets a single inline CSS style using keyed syntax.
//
//	{ ds.StyleKey("display", "$hiding && 'none'")... }
func StyleKey(prop, expr string, modifiers ...Modifier) templ.Attributes {
	return val(keyed(attrStyle, prop, modifiers), expr)
}

// ---------------------------------------------------------------------------
// data-ref
// ---------------------------------------------------------------------------

// Ref creates a signal that is a DOM reference to the element.
//
//	{ ds.Ref("myEl")... }
//
// See https://data-star.dev/reference/attributes#data-ref
func Ref(name string, modifiers ...Modifier) templ.Attributes {
	return val(plugin(attrRef, modifiers), name)
}

// ---------------------------------------------------------------------------
// data-indicator
// ---------------------------------------------------------------------------

// Indicator creates a boolean signal that is true while a fetch is in flight.
//
//	{ ds.Indicator("fetching")... }
//
// See https://data-star.dev/reference/attributes#data-indicator
func Indicator(name string, modifiers ...Modifier) templ.Attributes {
	return val(plugin(attrIndicator, modifiers), name)
}

// ---------------------------------------------------------------------------
// data-init
// ---------------------------------------------------------------------------

// Init runs an expression when the element is loaded into the DOM.
//
//	{ ds.Init("$count = 1")... }
//	{ ds.Init("@get('/updates')", ds.ModDelay, ds.Ms(500))... }
//
// See https://data-star.dev/reference/attributes#data-init
func Init(expr string, modifiers ...Modifier) templ.Attributes {
	return val(plugin(attrInit, modifiers), expr)
}

// ---------------------------------------------------------------------------
// data-effect
// ---------------------------------------------------------------------------

// Effect executes an expression on load and whenever dependency signals change.
//
//	{ ds.Effect("$total = $price * $qty")... }
//
// See https://data-star.dev/reference/attributes#data-effect
func Effect(expr string) templ.Attributes {
	return val(prefix+attrEffect, expr)
}

// ---------------------------------------------------------------------------
// data-ignore
// ---------------------------------------------------------------------------

// Ignore tells Datastar to skip processing an element and its descendants.
//
//	{ ds.Ignore()... }
//	{ ds.Ignore(ds.ModSelf)... }
//
// See https://data-star.dev/reference/attributes#data-ignore
func Ignore(modifiers ...Modifier) templ.Attributes {
	return boolAttr(plugin(attrIgnore, modifiers))
}

// ---------------------------------------------------------------------------
// data-ignore-morph
// ---------------------------------------------------------------------------

// IgnoreMorph tells PatchElements to skip morphing an element and its children.
//
//	{ ds.IgnoreMorph()... }
//
// See https://data-star.dev/reference/attributes#data-ignore-morph
func IgnoreMorph() templ.Attributes {
	return boolAttr(prefix + attrIgnoreMorph)
}

// ---------------------------------------------------------------------------
// data-json-signals
// ---------------------------------------------------------------------------

// JSONSignals sets text content to a JSON-stringified version of signals.
//
//	{ ds.JSONSignals(ds.Filter{})... }
//	{ ds.JSONSignals(ds.Filter{Include: "/user/"}, ds.ModTerse)... }
//
// See https://data-star.dev/reference/attributes#data-json-signals
func JSONSignals(filter Filter, modifiers ...Modifier) templ.Attributes {
	name := plugin(attrJSONSignals, modifiers)
	if filter.Include == "" && filter.Exclude == "" {
		return boolAttr(name)
	}
	return val(name, toFilter(filter))
}

// ---------------------------------------------------------------------------
// data-preserve-attr
// ---------------------------------------------------------------------------

// PreserveAttr preserves attribute values when morphing DOM elements.
//
//	{ ds.PreserveAttr("open")... }
//	{ ds.PreserveAttr("open", "class")... }
//
// See https://data-star.dev/reference/attributes#data-preserve-attr
func PreserveAttr(attrs ...string) templ.Attributes {
	return val(prefix+attrPreserveAttr, strings.Join(attrs, " "))
}
