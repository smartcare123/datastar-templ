package ds

import (
	"encoding/json"
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
func JSON(key string, value any) Signal {
	data, err := json.Marshal(value)
	if err != nil {
		panic("ds: failed to marshal JSON signal: " + err.Error())
	}
	return Signal{key, string(data)}
}

// Pool for reusing strings.Builder instances
var signalsBuilderPool = sync.Pool{
	New: func() interface{} {
		b := new(strings.Builder)
		b.Grow(128)
		return b
	},
}

// Signals patches one or more signals using typed helpers.
//
//	{ ds.Signals(ds.Int("count", 0), ds.String("msg", "hello"))... }
//
// See https://data-star.dev/reference/attributes#data-signals
func Signals(signals ...Signal) templ.Attributes {
	b := signalsBuilderPool.Get().(*strings.Builder)
	defer func() {
		b.Reset()
		signalsBuilderPool.Put(b)
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

// ComputedPair represents a computed signal name and its expression.
type ComputedPair struct {
	name string
	expr string
}

// Comp creates a computed signal pair.
func Comp(name, expr string) ComputedPair {
	return ComputedPair{name, expr}
}

// Pool for computed signal builder
var computedBuilderPool = sync.Pool{
	New: func() interface{} {
		b := new(strings.Builder)
		b.Grow(128)
		return b
	},
}

// Computed creates read-only computed signals using typed pairs.
// Each pair is wrapped in an arrow function.
//
//	{ ds.Computed(ds.Comp("total", "$price * $qty"))... }
//
// See https://data-star.dev/reference/attributes#data-computed
func Computed(pairs ...ComputedPair) templ.Attributes {
	b := computedBuilderPool.Get().(*strings.Builder)
	defer func() {
		b.Reset()
		computedBuilderPool.Put(b)
	}()

	// Calculate exact capacity
	capacity := 2 // {}
	for i, pair := range pairs {
		if i > 0 {
			capacity += 2 // ", "
		}
		// 'name': () => expr
		capacity += 1 + len(pair.name) + 11 + len(pair.expr) // "'name': () => expr"
	}

	if b.Cap() < capacity {
		b.Grow(capacity - b.Len())
	}

	b.WriteByte('{')
	for i, pair := range pairs {
		if i > 0 {
			b.WriteByte(',')
			b.WriteByte(' ')
		}
		b.WriteByte('\'')
		b.WriteString(pair.name)
		b.WriteString("': () => ")
		b.WriteString(pair.expr)
	}
	b.WriteByte('}')

	return templ.Attributes{"data-computed": b.String()}
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

// ClassPair represents a CSS class name and its binding expression.
type ClassPair struct {
	class string
	expr  string
}

// C creates a class binding pair.
func C(class, expr string) ClassPair {
	return ClassPair{class, expr}
}

// Pool for class builder
var classBuilderPool = sync.Pool{
	New: func() interface{} {
		b := new(strings.Builder)
		b.Grow(128)
		return b
	},
}

// Class adds/removes CSS classes using typed pairs.
//
//	{ ds.Class(ds.C("hidden", "$isHidden"), ds.C("font-bold", "$isBold"))... }
//
// See https://data-star.dev/reference/attributes#data-class
func Class(pairs ...ClassPair) templ.Attributes {
	b := classBuilderPool.Get().(*strings.Builder)
	defer func() {
		b.Reset()
		classBuilderPool.Put(b)
	}()

	// Calculate exact capacity
	capacity := 2 // {}
	for i, pair := range pairs {
		if i > 0 {
			capacity += 2 // ", "
		}
		// 'class': expr
		capacity += 1 + len(pair.class) + 4 + len(pair.expr) // "'class': expr"
	}

	if b.Cap() < capacity {
		b.Grow(capacity - b.Len())
	}

	b.WriteByte('{')
	for i, pair := range pairs {
		if i > 0 {
			b.WriteByte(',')
			b.WriteByte(' ')
		}
		b.WriteByte('\'')
		b.WriteString(pair.class)
		b.WriteString("': ")
		b.WriteString(pair.expr)
	}
	b.WriteByte('}')

	return templ.Attributes{"data-class": b.String()}
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

// AttrPair represents an HTML attribute name and its binding expression.
type AttrPair struct {
	attr string
	expr string
}

// A creates an attribute binding pair.
func A(attr, expr string) AttrPair {
	return AttrPair{attr, expr}
}

// Pool for attr builder
var attrBuilderPool = sync.Pool{
	New: func() interface{} {
		b := new(strings.Builder)
		b.Grow(128)
		return b
	},
}

// Attr sets HTML attributes using typed pairs.
//
//	{ ds.Attr(ds.A("title", "$tooltip"), ds.A("disabled", "$loading"))... }
//
// See https://data-star.dev/reference/attributes#data-attr
func Attr(pairs ...AttrPair) templ.Attributes {
	b := attrBuilderPool.Get().(*strings.Builder)
	defer func() {
		b.Reset()
		attrBuilderPool.Put(b)
	}()

	// Calculate exact capacity
	capacity := 2 // {}
	for i, pair := range pairs {
		if i > 0 {
			capacity += 2 // ", "
		}
		// 'attr': expr
		capacity += 1 + len(pair.attr) + 4 + len(pair.expr) // "'attr': expr"
	}

	if b.Cap() < capacity {
		b.Grow(capacity - b.Len())
	}

	b.WriteByte('{')
	for i, pair := range pairs {
		if i > 0 {
			b.WriteByte(',')
			b.WriteByte(' ')
		}
		b.WriteByte('\'')
		b.WriteString(pair.attr)
		b.WriteString("': ")
		b.WriteString(pair.expr)
	}
	b.WriteByte('}')

	return templ.Attributes{"data-attr": b.String()}
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

// StylePair represents a CSS property and its binding expression.
type StylePair struct {
	prop string
	expr string
}

// S creates a style binding pair.
func S(prop, expr string) StylePair {
	return StylePair{prop, expr}
}

// Pool for style builder
var styleBuilderPool = sync.Pool{
	New: func() interface{} {
		b := new(strings.Builder)
		b.Grow(128)
		return b
	},
}

// Style sets inline CSS styles using typed pairs.
//
//	{ ds.Style(ds.S("display", "$hiding && 'none'"), ds.S("color", "$textColor"))... }
//
// See https://data-star.dev/reference/attributes#data-style
func Style(pairs ...StylePair) templ.Attributes {
	b := styleBuilderPool.Get().(*strings.Builder)
	defer func() {
		b.Reset()
		styleBuilderPool.Put(b)
	}()

	// Calculate exact capacity
	capacity := 2 // {}
	for i, pair := range pairs {
		if i > 0 {
			capacity += 2 // ", "
		}
		// 'prop': expr
		capacity += 1 + len(pair.prop) + 4 + len(pair.expr) // "'prop': expr"
	}

	if b.Cap() < capacity {
		b.Grow(capacity - b.Len())
	}

	b.WriteByte('{')
	for i, pair := range pairs {
		if i > 0 {
			b.WriteByte(',')
			b.WriteByte(' ')
		}
		b.WriteByte('\'')
		b.WriteString(pair.prop)
		b.WriteString("': ")
		b.WriteString(pair.expr)
	}
	b.WriteByte('}')

	return templ.Attributes{"data-style": b.String()}
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
