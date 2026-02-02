// Package ds constants. All Datastar attribute names, DOM event names, prefixes,
// and modifier constants are defined here as the single source of truth.
//
// Unexported constants are used internally by the helper functions.
// Exported Modifier constants are part of the public API.
package ds

// ---------------------------------------------------------------------------
// Prefixes and separators (unexported)
// ---------------------------------------------------------------------------

const (
	// prefix is the common prefix for all Datastar attributes.
	prefix = "data-"
	// prefixOn is the prefix for DOM event handler attributes (colon-separated).
	prefixOn = "data-on:"
	// sepColon separates attribute names from keys in keyed attributes.
	sepColon = ":"
)

// ---------------------------------------------------------------------------
// Attribute name segments (unexported)
//
// These are the Datastar attribute identifiers without the "data-" prefix.
// Used by the internal helpers (plugin, keyed, val) to build full attribute names.
// ---------------------------------------------------------------------------

const (
	attrAttr              = "attr"
	attrBind              = "bind"
	attrClass             = "class"
	attrComputed          = "computed"
	attrEffect            = "effect"
	attrIgnore            = "ignore"
	attrIgnoreMorph       = "ignore-morph"
	attrIndicator         = "indicator"
	attrInit              = "init"
	attrJSONSignals       = "json-signals"
	attrOnIntersect       = "on-intersect"
	attrOnInterval        = "on-interval"
	attrOnSignalPatch     = "on-signal-patch"
	attrOnSignalPatchFilt = "on-signal-patch-filter"
	attrPreserveAttr      = "preserve-attr"
	attrRef               = "ref"
	attrShow              = "show"
	attrSignals           = "signals"
	attrStyle             = "style"
	attrText              = "text"
)

// ---------------------------------------------------------------------------
// SSE action verbs (unexported)
// ---------------------------------------------------------------------------

const (
	actionGet    = "get"
	actionPost   = "post"
	actionPut    = "put"
	actionPatch  = "patch"
	actionDelete = "delete"
)

// ---------------------------------------------------------------------------
// DOM event names (unexported)
//
// Standard browser events used by the On* functions. Kept as constants to
// prevent typos and centralise naming.
// ---------------------------------------------------------------------------

// Mouse events.
const (
	eventClick       = "click"
	eventDblClick    = "dblclick"
	eventMouseDown   = "mousedown"
	eventMouseUp     = "mouseup"
	eventMouseOver   = "mouseover"
	eventMouseOut    = "mouseout"
	eventMouseMove   = "mousemove"
	eventMouseEnter  = "mouseenter"
	eventMouseLeave  = "mouseleave"
	eventContextMenu = "contextmenu"
)

// Keyboard events.
const (
	eventKeyDown  = "keydown"
	eventKeyUp    = "keyup"
	eventKeyPress = "keypress"
)

// Focus events.
const (
	eventFocus    = "focus"
	eventBlur     = "blur"
	eventFocusIn  = "focusin"
	eventFocusOut = "focusout"
)

// Form events.
const (
	eventSubmit  = "submit"
	eventReset   = "reset"
	eventInput   = "input"
	eventChange  = "change"
	eventInvalid = "invalid"
	eventSelect  = "select"
)

// Drag events.
const (
	eventDrag      = "drag"
	eventDragStart = "dragstart"
	eventDragEnd   = "dragend"
	eventDragOver  = "dragover"
	eventDragEnter = "dragenter"
	eventDragLeave = "dragleave"
	eventDrop      = "drop"
)

// Touch events.
const (
	eventTouchStart  = "touchstart"
	eventTouchEnd    = "touchend"
	eventTouchMove   = "touchmove"
	eventTouchCancel = "touchcancel"
)

// Pointer events.
const (
	eventPointerDown        = "pointerdown"
	eventPointerUp          = "pointerup"
	eventPointerMove        = "pointermove"
	eventPointerOver        = "pointerover"
	eventPointerOut         = "pointerout"
	eventPointerEnter       = "pointerenter"
	eventPointerLeave       = "pointerleave"
	eventPointerCancel      = "pointercancel"
	eventGotPointerCapture  = "gotpointercapture"
	eventLostPointerCapture = "lostpointercapture"
)

// Scroll / wheel events.
const (
	eventScroll = "scroll"
	eventWheel  = "wheel"
)

// Animation / transition events.
const (
	eventAnimationStart     = "animationstart"
	eventAnimationEnd       = "animationend"
	eventAnimationIteration = "animationiteration"
	eventTransitionEnd      = "transitionend"
)

// Media events.
const (
	eventLoad  = "load"
	eventError = "error"
)

// Clipboard events.
const (
	eventCopy  = "copy"
	eventCut   = "cut"
	eventPaste = "paste"
)

// ---------------------------------------------------------------------------
// Modifier constants – double-underscore (exported)
// ---------------------------------------------------------------------------

// Double-underscore modifiers control event listener behavior and timing.
const (
	ModCapture        Modifier = "__capture"
	ModCase           Modifier = "__case"
	ModDebounce       Modifier = "__debounce"
	ModDelay          Modifier = "__delay"
	ModDuration       Modifier = "__duration"
	ModExit           Modifier = "__exit"
	ModFull           Modifier = "__full"
	ModHalf           Modifier = "__half"
	ModIfMissing      Modifier = "__ifmissing"
	ModOnce           Modifier = "__once"
	ModOutside        Modifier = "__outside"
	ModPassive        Modifier = "__passive"
	ModPrevent        Modifier = "__prevent"
	ModSelf           Modifier = "__self"
	ModStop           Modifier = "__stop"
	ModTerse          Modifier = "__terse"
	ModThreshold      Modifier = "__threshold"
	ModThrottle       Modifier = "__throttle"
	ModViewTransition Modifier = "__viewtransition"
	ModWindow         Modifier = "__window"
)

// ---------------------------------------------------------------------------
// Modifier constants – dot-tags (exported)
// ---------------------------------------------------------------------------

// Dot-tag modifiers specify case conversion and timing behavior.
const (
	Camel      Modifier = ".camel"
	Kebab      Modifier = ".kebab"
	Snake      Modifier = ".snake"
	Pascal     Modifier = ".pascal"
	Leading    Modifier = ".leading"
	NoLeading  Modifier = ".noleading"
	NoTrailing Modifier = ".notrailing"
	Trailing   Modifier = ".trailing"
)
