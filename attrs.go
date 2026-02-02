package ds

import (
	"strings"

	"github.com/a-h/templ"
)

// ===========================================================================
// Datastar Attributes
// See https://data-star.dev/reference/attributes
// ===========================================================================

// ---------------------------------------------------------------------------
// data-signals
// ---------------------------------------------------------------------------

// Signals patches one or more signals using object notation.
// The map is JSON-marshaled.
//
//	{ ds.Signals(map[string]any{"foo": 1, "bar": "hello"})... }
//
// See https://data-star.dev/reference/attributes#data-signals
func Signals(signals map[string]any, modifiers ...Modifier) templ.Attributes {
	return val(plugin(attrSignals, modifiers), toSignals(signals))
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

// Computed creates read-only computed signals using object notation.
// Pairs are (name, expression) and are wrapped in arrow functions.
//
//	{ ds.Computed("total", "$price * $qty")... }
//
// See https://data-star.dev/reference/attributes#data-computed
func Computed(pairs ...string) templ.Attributes {
	return val(prefix+attrComputed, toComputed(pairs))
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

// Class adds/removes CSS classes using object notation.
// Pairs are (className, expression).
//
//	{ ds.Class("hidden", "$isHidden", "font-bold", "$isBold")... }
//
// See https://data-star.dev/reference/attributes#data-class
func Class(pairs ...string) templ.Attributes {
	return val(prefix+attrClass, toObject(pairs))
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

// Attr sets HTML attributes using object notation.
// Pairs are (attrName, expression).
//
//	{ ds.Attr("title", "$tooltip", "disabled", "$loading")... }
//
// See https://data-star.dev/reference/attributes#data-attr
func Attr(pairs ...string) templ.Attributes {
	return val(prefix+attrAttr, toObject(pairs))
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

// Style sets inline CSS styles using object notation.
// Pairs are (property, expression).
//
//	{ ds.Style("display", "$hiding && 'none'", "color", "$red ? 'red' : 'blue'")... }
//
// See https://data-star.dev/reference/attributes#data-style
func Style(pairs ...string) templ.Attributes {
	return val(prefix+attrStyle, toObject(pairs))
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
