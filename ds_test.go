package ds_test

import (
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-scheduler/pkg/ds"
)

// ---------------------------------------------------------------------------
// Modifier helpers
// ---------------------------------------------------------------------------

func TestDuration(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want ds.Modifier
	}{
		{"zero", 0, ".0ms"},
		{"1ms", time.Millisecond, ".1ms"},
		{"500ms", 500 * time.Millisecond, ".500ms"},
		{"1s", time.Second, ".1000ms"},
		{"rounds 500us up to 1ms", 500 * time.Microsecond, ".1ms"},
		{"rounds 100us down to 0ms", 100 * time.Microsecond, ".0ms"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ds.Duration(tt.d))
		})
	}

	t.Run("panics on negative", func(t *testing.T) {
		assert.Panics(t, func() { ds.Duration(-1) })
	})
}

func TestMs(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want ds.Modifier
	}{
		{"zero", 0, ".0ms"},
		{"300", 300, ".300ms"},
		{"1000", 1000, ".1000ms"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ds.Ms(tt.n))
		})
	}

	t.Run("panics on negative", func(t *testing.T) {
		assert.Panics(t, func() { ds.Ms(-1) })
	})
}

func TestSeconds(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want ds.Modifier
	}{
		{"zero", 0, ".0s"},
		{"1", 1, ".1s"},
		{"5", 5, ".5s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ds.Seconds(tt.n))
		})
	}

	t.Run("panics on negative", func(t *testing.T) {
		assert.Panics(t, func() { ds.Seconds(-1) })
	})
}

func TestThreshold(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		want ds.Modifier
	}{
		{"0.25", 0.25, ".25"},
		{"0.50", 0.50, ".50"},
		{"0.75", 0.75, ".75"},
		{"1.0 (100%)", 1.0, ".100"},
		{"rounds 0.333 to .33", 0.333, ".33"},
		{"rounds 0.335 to .34", 0.335, ".34"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ds.Threshold(tt.val))
		})
	}

	t.Run("panics on zero", func(t *testing.T) {
		assert.Panics(t, func() { ds.Threshold(0) })
	})
	t.Run("panics on negative", func(t *testing.T) {
		assert.Panics(t, func() { ds.Threshold(-0.1) })
	})
	t.Run("panics on > 1", func(t *testing.T) {
		assert.Panics(t, func() { ds.Threshold(1.1) })
	})
}

// ---------------------------------------------------------------------------
// Merge
// ---------------------------------------------------------------------------

func TestMerge(t *testing.T) {
	t.Run("combines multiple", func(t *testing.T) {
		result := ds.Merge(
			ds.Show("$visible"),
			ds.OnClick("$open = false"),
		)
		assert.Equal(t, "$visible", result["data-show"])
		assert.Equal(t, "$open = false", result["data-on:click"])
		assert.Len(t, result, 2)
	})

	t.Run("later overrides earlier", func(t *testing.T) {
		result := ds.Merge(
			ds.Show("$first"),
			ds.Show("$second"),
		)
		assert.Equal(t, "$second", result["data-show"])
		assert.Len(t, result, 1)
	})

	t.Run("empty merge", func(t *testing.T) {
		result := ds.Merge()
		assert.Empty(t, result)
	})
}

// ---------------------------------------------------------------------------
// DOM Event Functions – data-on:{event}
// ---------------------------------------------------------------------------

func TestOnClick(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		mods    []ds.Modifier
		wantKey string
		wantVal string
	}{
		{
			name:    "simple",
			expr:    "$foo = true",
			wantKey: "data-on:click",
			wantVal: "$foo = true",
		},
		{
			name:    "with debounce",
			expr:    "$foo = true",
			mods:    []ds.Modifier{ds.ModDebounce, ds.Ms(300)},
			wantKey: "data-on:click__debounce.300ms",
			wantVal: "$foo = true",
		},
		{
			name:    "with window and throttle",
			expr:    "$x = 1",
			mods:    []ds.Modifier{ds.ModWindow, ds.ModThrottle, ds.Ms(100), ds.NoLeading},
			wantKey: "data-on:click__window__throttle.100ms.noleading",
			wantVal: "$x = 1",
		},
		{
			name:    "with prevent",
			expr:    "doStuff()",
			mods:    []ds.Modifier{ds.ModPrevent},
			wantKey: "data-on:click__prevent",
			wantVal: "doStuff()",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := ds.OnClick(tt.expr, tt.mods...)
			require.Len(t, attrs, 1)
			assert.Equal(t, tt.wantVal, attrs[tt.wantKey])
		})
	}
}

// TestAllDOMEvents verifies every typed DOM event function produces the
// correct "data-on:{event}" attribute name.
func TestAllDOMEvents(t *testing.T) {
	type eventFunc = func(string, ...ds.Modifier) templ.Attributes

	tests := []struct {
		name    string
		fn      eventFunc
		wantKey string
	}{
		// Mouse
		{"OnClick", ds.OnClick, "data-on:click"},
		{"OnDblClick", ds.OnDblClick, "data-on:dblclick"},
		{"OnMouseDown", ds.OnMouseDown, "data-on:mousedown"},
		{"OnMouseUp", ds.OnMouseUp, "data-on:mouseup"},
		{"OnMouseOver", ds.OnMouseOver, "data-on:mouseover"},
		{"OnMouseOut", ds.OnMouseOut, "data-on:mouseout"},
		{"OnMouseMove", ds.OnMouseMove, "data-on:mousemove"},
		{"OnMouseEnter", ds.OnMouseEnter, "data-on:mouseenter"},
		{"OnMouseLeave", ds.OnMouseLeave, "data-on:mouseleave"},
		{"OnContextMenu", ds.OnContextMenu, "data-on:contextmenu"},
		// Keyboard
		{"OnKeyDown", ds.OnKeyDown, "data-on:keydown"},
		{"OnKeyUp", ds.OnKeyUp, "data-on:keyup"},
		{"OnKeyPress", ds.OnKeyPress, "data-on:keypress"},
		// Focus
		{"OnFocus", ds.OnFocus, "data-on:focus"},
		{"OnBlur", ds.OnBlur, "data-on:blur"},
		{"OnFocusIn", ds.OnFocusIn, "data-on:focusin"},
		{"OnFocusOut", ds.OnFocusOut, "data-on:focusout"},
		// Form
		{"OnSubmit", ds.OnSubmit, "data-on:submit"},
		{"OnReset", ds.OnReset, "data-on:reset"},
		{"OnInput", ds.OnInput, "data-on:input"},
		{"OnChange", ds.OnChange, "data-on:change"},
		{"OnInvalid", ds.OnInvalid, "data-on:invalid"},
		{"OnSelect", ds.OnSelect, "data-on:select"},
		// Drag
		{"OnDrag", ds.OnDrag, "data-on:drag"},
		{"OnDragStart", ds.OnDragStart, "data-on:dragstart"},
		{"OnDragEnd", ds.OnDragEnd, "data-on:dragend"},
		{"OnDragOver", ds.OnDragOver, "data-on:dragover"},
		{"OnDragEnter", ds.OnDragEnter, "data-on:dragenter"},
		{"OnDragLeave", ds.OnDragLeave, "data-on:dragleave"},
		{"OnDrop", ds.OnDrop, "data-on:drop"},
		// Touch
		{"OnTouchStart", ds.OnTouchStart, "data-on:touchstart"},
		{"OnTouchEnd", ds.OnTouchEnd, "data-on:touchend"},
		{"OnTouchMove", ds.OnTouchMove, "data-on:touchmove"},
		{"OnTouchCancel", ds.OnTouchCancel, "data-on:touchcancel"},
		// Pointer
		{"OnPointerDown", ds.OnPointerDown, "data-on:pointerdown"},
		{"OnPointerUp", ds.OnPointerUp, "data-on:pointerup"},
		{"OnPointerMove", ds.OnPointerMove, "data-on:pointermove"},
		{"OnPointerOver", ds.OnPointerOver, "data-on:pointerover"},
		{"OnPointerOut", ds.OnPointerOut, "data-on:pointerout"},
		{"OnPointerEnter", ds.OnPointerEnter, "data-on:pointerenter"},
		{"OnPointerLeave", ds.OnPointerLeave, "data-on:pointerleave"},
		{"OnPointerCancel", ds.OnPointerCancel, "data-on:pointercancel"},
		{"OnGotPointerCapture", ds.OnGotPointerCapture, "data-on:gotpointercapture"},
		{"OnLostPointerCapture", ds.OnLostPointerCapture, "data-on:lostpointercapture"},
		// Scroll/Wheel
		{"OnScroll", ds.OnScroll, "data-on:scroll"},
		{"OnWheel", ds.OnWheel, "data-on:wheel"},
		// Animation/Transition
		{"OnAnimationStart", ds.OnAnimationStart, "data-on:animationstart"},
		{"OnAnimationEnd", ds.OnAnimationEnd, "data-on:animationend"},
		{"OnAnimationIteration", ds.OnAnimationIteration, "data-on:animationiteration"},
		{"OnTransitionEnd", ds.OnTransitionEnd, "data-on:transitionend"},
		// Media
		{"OnLoad", ds.OnLoad, "data-on:load"},
		{"OnError", ds.OnError, "data-on:error"},
		// Clipboard
		{"OnCopy", ds.OnCopy, "data-on:copy"},
		{"OnCut", ds.OnCut, "data-on:cut"},
		{"OnPaste", ds.OnPaste, "data-on:paste"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := tt.fn("expr")
			require.Len(t, attrs, 1)
			assert.Equal(t, "expr", attrs[tt.wantKey])
		})
	}
}

func TestOnEvent(t *testing.T) {
	t.Run("custom event", func(t *testing.T) {
		attrs := ds.OnEvent("table-select", "$selected = evt.detail.ids")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$selected = evt.detail.ids", attrs["data-on:table-select"])
	})

	t.Run("custom event with modifiers", func(t *testing.T) {
		attrs := ds.OnEvent("my-event", "handle()", ds.ModOnce, ds.ModCapture)
		require.Len(t, attrs, 1)
		assert.Equal(t, "handle()", attrs["data-on:my-event__once__capture"])
	})
}

func TestOnEventModifiers(t *testing.T) {
	t.Run("debounce with duration and leading", func(t *testing.T) {
		attrs := ds.OnInput("doSearch()", ds.ModDebounce, ds.Duration(500*time.Millisecond), ds.Leading)
		require.Len(t, attrs, 1)
		assert.Equal(t, "doSearch()", attrs["data-on:input__debounce.500ms.leading"])
	})

	t.Run("window outside", func(t *testing.T) {
		attrs := ds.OnClick("close()", ds.ModWindow, ds.ModOutside)
		require.Len(t, attrs, 1)
		assert.Equal(t, "close()", attrs["data-on:click__window__outside"])
	})

	t.Run("once passive capture", func(t *testing.T) {
		attrs := ds.OnScroll("track()", ds.ModOnce, ds.ModPassive, ds.ModCapture)
		require.Len(t, attrs, 1)
		assert.Equal(t, "track()", attrs["data-on:scroll__once__passive__capture"])
	})

	t.Run("stop prevent", func(t *testing.T) {
		attrs := ds.OnSubmit("save()", ds.ModStop, ds.ModPrevent)
		require.Len(t, attrs, 1)
		assert.Equal(t, "save()", attrs["data-on:submit__stop__prevent"])
	})

	t.Run("view transition", func(t *testing.T) {
		attrs := ds.OnClick("navigate()", ds.ModViewTransition)
		require.Len(t, attrs, 1)
		assert.Equal(t, "navigate()", attrs["data-on:click__viewtransition"])
	})

	t.Run("case modifier", func(t *testing.T) {
		attrs := ds.OnEvent("myEvent", "handle()", ds.ModCase, ds.Kebab)
		require.Len(t, attrs, 1)
		assert.Equal(t, "handle()", attrs["data-on:myEvent__case.kebab"])
	})
}

// ---------------------------------------------------------------------------
// Datastar Attributes
// ---------------------------------------------------------------------------

func TestShow(t *testing.T) {
	attrs := ds.Show("$visible")
	require.Len(t, attrs, 1)
	assert.Equal(t, "$visible", attrs["data-show"])
}

func TestText(t *testing.T) {
	attrs := ds.Text("$count + ' items'")
	require.Len(t, attrs, 1)
	assert.Equal(t, "$count + ' items'", attrs["data-text"])
}

func TestEffect(t *testing.T) {
	attrs := ds.Effect("$total = $price * $qty")
	require.Len(t, attrs, 1)
	assert.Equal(t, "$total = $price * $qty", attrs["data-effect"])
}

func TestInit(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.Init("$count = 1")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$count = 1", attrs["data-init"])
	})

	t.Run("with delay", func(t *testing.T) {
		attrs := ds.Init("$count = 1", ds.ModDelay, ds.Ms(500))
		require.Len(t, attrs, 1)
		assert.Equal(t, "$count = 1", attrs["data-init__delay.500ms"])
	})

	t.Run("with view transition", func(t *testing.T) {
		attrs := ds.Init("navigate()", ds.ModViewTransition)
		require.Len(t, attrs, 1)
		assert.Equal(t, "navigate()", attrs["data-init__viewtransition"])
	})
}

func TestBind(t *testing.T) {
	t.Run("keyed syntax", func(t *testing.T) {
		attrs := ds.Bind("name")
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-bind:name"])
	})

	t.Run("nested signal", func(t *testing.T) {
		attrs := ds.Bind("table.search")
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-bind:table.search"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.Bind("my-field", ds.ModCase, ds.Kebab)
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-bind:my-field__case.kebab"])
	})
}

func TestBindExpr(t *testing.T) {
	attrs := ds.BindExpr("name")
	require.Len(t, attrs, 1)
	assert.Equal(t, "name", attrs["data-bind"])
}

func TestSignals(t *testing.T) {
	t.Run("simple map", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{"foo": 1})
		require.Len(t, attrs, 1)
		assert.Equal(t, `{"foo":1}`, attrs["data-signals"])
	})

	t.Run("with ifmissing modifier", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{"foo": 1}, ds.ModIfMissing)
		require.Len(t, attrs, 1)
		assert.Equal(t, `{"foo":1}`, attrs["data-signals__ifmissing"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.Signals(map[string]any{"foo": 1}, ds.ModCase, ds.Kebab)
		require.Len(t, attrs, 1)
		assert.Equal(t, `{"foo":1}`, attrs["data-signals__case.kebab"])
	})
}

func TestSignalsJSON(t *testing.T) {
	attrs := ds.SignalsJSON(`{foo: {bar: 1}}`)
	require.Len(t, attrs, 1)
	assert.Equal(t, `{foo: {bar: 1}}`, attrs["data-signals"])
}

func TestSignalKey(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.SignalKey("foo", "1")
		require.Len(t, attrs, 1)
		assert.Equal(t, "1", attrs["data-signals:foo"])
	})

	t.Run("with ifmissing", func(t *testing.T) {
		attrs := ds.SignalKey("foo", "1", ds.ModIfMissing)
		require.Len(t, attrs, 1)
		assert.Equal(t, "1", attrs["data-signals:foo__ifmissing"])
	})
}

func TestComputed(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		attrs := ds.Computed("total", "$price * $qty")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'total': () => $price * $qty}", attrs["data-computed"])
	})

	t.Run("multiple", func(t *testing.T) {
		attrs := ds.Computed("total", "$price * $qty", "tax", "$total * 0.1")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'total': () => $price * $qty, 'tax': () => $total * 0.1}", attrs["data-computed"])
	})

	t.Run("panics on odd pairs", func(t *testing.T) {
		assert.Panics(t, func() { ds.Computed("name") }) //nolint:staticcheck // intentionally testing panic on odd args
	})
}

func TestComputedKey(t *testing.T) {
	attrs := ds.ComputedKey("total", "$price * $qty")
	require.Len(t, attrs, 1)
	assert.Equal(t, "$price * $qty", attrs["data-computed:total"])
}

func TestClass(t *testing.T) {
	t.Run("single pair", func(t *testing.T) {
		attrs := ds.Class("hidden", "$isHidden")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'hidden': $isHidden}", attrs["data-class"])
	})

	t.Run("multiple pairs", func(t *testing.T) {
		attrs := ds.Class("hidden", "$isHidden", "font-bold", "$isBold")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'hidden': $isHidden, 'font-bold': $isBold}", attrs["data-class"])
	})

	t.Run("panics on odd pairs", func(t *testing.T) {
		assert.Panics(t, func() { ds.Class("hidden") }) //nolint:staticcheck // intentionally testing panic on odd args
	})
}

func TestClassKey(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.ClassKey("font-bold", "$isBold")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$isBold", attrs["data-class:font-bold"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.ClassKey("myClass", "$flag", ds.ModCase, ds.Camel)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$flag", attrs["data-class:myClass__case.camel"])
	})
}

func TestAttr(t *testing.T) {
	t.Run("single pair", func(t *testing.T) {
		attrs := ds.Attr("title", "$tooltip")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'title': $tooltip}", attrs["data-attr"])
	})

	t.Run("multiple pairs", func(t *testing.T) {
		attrs := ds.Attr("title", "$tooltip", "disabled", "$loading")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'title': $tooltip, 'disabled': $loading}", attrs["data-attr"])
	})

	t.Run("panics on odd pairs", func(t *testing.T) {
		assert.Panics(t, func() { ds.Attr("title") }) //nolint:staticcheck // intentionally testing panic on odd args
	})
}

func TestAttrKey(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		attrs := ds.AttrKey("disabled", "$loading")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$loading", attrs["data-attr:disabled"])
	})

	t.Run("title", func(t *testing.T) {
		attrs := ds.AttrKey("title", "'Theme: ' + $theme")
		require.Len(t, attrs, 1)
		assert.Equal(t, "'Theme: ' + $theme", attrs["data-attr:title"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.AttrKey("myAttr", "$val", ds.ModCase, ds.Camel)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$val", attrs["data-attr:myAttr__case.camel"])
	})
}

func TestStyle(t *testing.T) {
	t.Run("single pair", func(t *testing.T) {
		attrs := ds.Style("display", "$hiding && 'none'")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'display': $hiding && 'none'}", attrs["data-style"])
	})

	t.Run("multiple pairs", func(t *testing.T) {
		attrs := ds.Style("display", "'none'", "color", "'red'")
		require.Len(t, attrs, 1)
		assert.Equal(t, "{'display': 'none', 'color': 'red'}", attrs["data-style"])
	})

	t.Run("panics on odd pairs", func(t *testing.T) {
		assert.Panics(t, func() { ds.Style("display") }) //nolint:staticcheck // intentionally testing panic on odd args
	})
}

func TestStyleKey(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.StyleKey("background-color", "$red ? 'red' : 'blue'")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$red ? 'red' : 'blue'", attrs["data-style:background-color"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.StyleKey("myProp", "$val", ds.ModCase, ds.Camel)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$val", attrs["data-style:myProp__case.camel"])
	})
}

func TestRef(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.Ref("myEl")
		require.Len(t, attrs, 1)
		assert.Equal(t, "myEl", attrs["data-ref"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.Ref("myEl", ds.ModCase, ds.Kebab)
		require.Len(t, attrs, 1)
		assert.Equal(t, "myEl", attrs["data-ref__case.kebab"])
	})
}

func TestIndicator(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.Indicator("fetching")
		require.Len(t, attrs, 1)
		assert.Equal(t, "fetching", attrs["data-indicator"])
	})

	t.Run("with case modifier", func(t *testing.T) {
		attrs := ds.Indicator("fetching", ds.ModCase, ds.Camel)
		require.Len(t, attrs, 1)
		assert.Equal(t, "fetching", attrs["data-indicator__case.camel"])
	})
}

func TestIgnore(t *testing.T) {
	t.Run("no modifiers", func(t *testing.T) {
		attrs := ds.Ignore()
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-ignore"])
	})

	t.Run("with self", func(t *testing.T) {
		attrs := ds.Ignore(ds.ModSelf)
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-ignore__self"])
	})
}

func TestIgnoreMorph(t *testing.T) {
	attrs := ds.IgnoreMorph()
	require.Len(t, attrs, 1)
	assert.Equal(t, true, attrs["data-ignore-morph"])
}

func TestJSONSignals(t *testing.T) {
	t.Run("no filter", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{})
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-json-signals"])
	})

	t.Run("with terse modifier", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{}, ds.ModTerse)
		require.Len(t, attrs, 1)
		assert.Equal(t, true, attrs["data-json-signals__terse"])
	})

	t.Run("with include filter", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Include: "/user/"})
		require.Len(t, attrs, 1)
		assert.Equal(t, "{include: /user/}", attrs["data-json-signals"])
	})

	t.Run("with exclude filter", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Exclude: "/temp$/"})
		require.Len(t, attrs, 1)
		assert.Equal(t, "{exclude: /temp$/}", attrs["data-json-signals"])
	})

	t.Run("with both filters", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Include: "/^app/", Exclude: "/password/"})
		require.Len(t, attrs, 1)
		assert.Equal(t, "{include: /^app/, exclude: /password/}", attrs["data-json-signals"])
	})

	t.Run("with filter and terse", func(t *testing.T) {
		attrs := ds.JSONSignals(ds.Filter{Include: "/counter/"}, ds.ModTerse)
		require.Len(t, attrs, 1)
		assert.Equal(t, "{include: /counter/}", attrs["data-json-signals__terse"])
	})
}

func TestPreserveAttr(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		attrs := ds.PreserveAttr("open")
		require.Len(t, attrs, 1)
		assert.Equal(t, "open", attrs["data-preserve-attr"])
	})

	t.Run("multiple", func(t *testing.T) {
		attrs := ds.PreserveAttr("open", "class")
		require.Len(t, attrs, 1)
		assert.Equal(t, "open class", attrs["data-preserve-attr"])
	})
}

// ---------------------------------------------------------------------------
// Plugin watchers
// ---------------------------------------------------------------------------

func TestOnIntersect(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.OnIntersect("$visible = true")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$visible = true", attrs["data-on-intersect"])
	})

	t.Run("once full", func(t *testing.T) {
		attrs := ds.OnIntersect("$visible = true", ds.ModOnce, ds.ModFull)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$visible = true", attrs["data-on-intersect__once__full"])
	})

	t.Run("with threshold", func(t *testing.T) {
		attrs := ds.OnIntersect("$visible = true", ds.ModThreshold, ds.Threshold(0.25))
		require.Len(t, attrs, 1)
		assert.Equal(t, "$visible = true", attrs["data-on-intersect__threshold.25"])
	})

	t.Run("with half", func(t *testing.T) {
		attrs := ds.OnIntersect("$half = true", ds.ModHalf)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$half = true", attrs["data-on-intersect__half"])
	})

	t.Run("with exit", func(t *testing.T) {
		attrs := ds.OnIntersect("$exited = true", ds.ModExit)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$exited = true", attrs["data-on-intersect__exit"])
	})
}

func TestOnInterval(t *testing.T) {
	t.Run("default interval", func(t *testing.T) {
		attrs := ds.OnInterval("$count++")
		require.Len(t, attrs, 1)
		assert.Equal(t, "$count++", attrs["data-on-interval"])
	})

	t.Run("custom duration", func(t *testing.T) {
		attrs := ds.OnInterval("$count++", ds.ModDuration, ds.Ms(500))
		require.Len(t, attrs, 1)
		assert.Equal(t, "$count++", attrs["data-on-interval__duration.500ms"])
	})

	t.Run("with leading", func(t *testing.T) {
		attrs := ds.OnInterval("$count++", ds.ModDuration, ds.Ms(500), ds.Leading)
		require.Len(t, attrs, 1)
		assert.Equal(t, "$count++", attrs["data-on-interval__duration.500ms.leading"])
	})
}

func TestOnSignalPatch(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		attrs := ds.OnSignalPatch("console.log(patch)")
		require.Len(t, attrs, 1)
		assert.Equal(t, "console.log(patch)", attrs["data-on-signal-patch"])
	})

	t.Run("with debounce", func(t *testing.T) {
		attrs := ds.OnSignalPatch("doStuff()", ds.ModDebounce, ds.Ms(500))
		require.Len(t, attrs, 1)
		assert.Equal(t, "doStuff()", attrs["data-on-signal-patch__debounce.500ms"])
	})
}

func TestOnSignalPatchFilter(t *testing.T) {
	t.Run("include only", func(t *testing.T) {
		attrs := ds.OnSignalPatchFilter(ds.Filter{Include: "/^counter$/"})
		require.Len(t, attrs, 1)
		assert.Equal(t, "{include: /^counter$/}", attrs["data-on-signal-patch-filter"])
	})

	t.Run("both", func(t *testing.T) {
		attrs := ds.OnSignalPatchFilter(ds.Filter{Include: "/user/", Exclude: "/password/"})
		require.Len(t, attrs, 1)
		assert.Equal(t, "{include: /user/, exclude: /password/}", attrs["data-on-signal-patch-filter"])
	})
}

// ---------------------------------------------------------------------------
// Real-world usage patterns (from this project's templates)
// ---------------------------------------------------------------------------

func TestRealWorldPatterns(t *testing.T) {
	t.Run("workcenters create button", func(t *testing.T) {
		attrs := ds.OnClick("$createSlideoutOpen = true")
		assert.Equal(t, "$createSlideoutOpen = true", attrs["data-on:click"])
	})

	t.Run("slideout overlay close", func(t *testing.T) {
		attrs := ds.OnClick("$createSlideoutOpen = false")
		assert.Equal(t, "$createSlideoutOpen = false", attrs["data-on:click"])
	})

	t.Run("search input with debounce", func(t *testing.T) {
		attrs := ds.OnInput("@post('/api/workcenters/search')", ds.ModDebounce, ds.Ms(300))
		assert.Equal(t, "@post('/api/workcenters/search')", attrs["data-on:input__debounce.300ms"])
	})

	t.Run("page size change", func(t *testing.T) {
		attrs := ds.OnChange("@patch('/api/workcenters/pagesize')")
		assert.Equal(t, "@patch('/api/workcenters/pagesize')", attrs["data-on:change"])
	})

	t.Run("sidebar toggle", func(t *testing.T) {
		attrs := ds.OnClick("$sidebarCollapsed = !$sidebarCollapsed")
		assert.Equal(t, "$sidebarCollapsed = !$sidebarCollapsed", attrs["data-on:click"])
	})

	t.Run("theme toggle with attr binding", func(t *testing.T) {
		click := ds.OnClick("$theme = window.themeToggle.cycle()")
		title := ds.AttrKey("title", "'Theme: ' + $theme")
		merged := ds.Merge(click, title)
		assert.Equal(t, "$theme = window.themeToggle.cycle()", merged["data-on:click"])
		assert.Equal(t, "'Theme: ' + $theme", merged["data-attr:title"])
	})

	t.Run("custom web component event", func(t *testing.T) {
		attrs := ds.OnEvent("table-select", "$table.selected = evt.detail.ids")
		assert.Equal(t, "$table.selected = evt.detail.ids", attrs["data-on:table-select"])
	})

	t.Run("init with SSE connection", func(t *testing.T) {
		attrs := ds.Init("@get('/api/workcenters/updates',{requestCancellation: 'disabled'})")
		assert.Equal(t,
			"@get('/api/workcenters/updates',{requestCancellation: 'disabled'})",
			attrs["data-init"],
		)
	})

	t.Run("toast delay init", func(t *testing.T) {
		attrs := ds.Init("el.style.opacity = '0'; setTimeout(() => el.remove(), 300)", ds.ModDelay, ds.Ms(3000))
		assert.Equal(t,
			"el.style.opacity = '0'; setTimeout(() => el.remove(), 300)",
			attrs["data-init__delay.3000ms"],
		)
	})
}
