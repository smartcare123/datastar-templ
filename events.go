package ds

import "github.com/a-h/templ"

// ===========================================================================
// DOM Event Functions – data-on:{event}
//
// Each returns templ.Attributes{"data-on:{event}{mods}": expr}.
// For custom or uncommon events, use OnEvent.
// See https://data-star.dev/reference/attributes#data-on
// ===========================================================================

// ---------------------------------------------------------------------------
// Mouse events
// ---------------------------------------------------------------------------

// OnClick handles the "click" event.
func OnClick(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventClick, modifiers), expr)
}

// OnDblClick handles the "dblclick" event.
func OnDblClick(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventDblClick, modifiers), expr)
}

// OnMouseDown handles the "mousedown" event.
func OnMouseDown(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventMouseDown, modifiers), expr)
}

// OnMouseUp handles the "mouseup" event.
func OnMouseUp(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventMouseUp, modifiers), expr)
}

// OnMouseOver handles the "mouseover" event.
func OnMouseOver(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventMouseOver, modifiers), expr)
}

// OnMouseOut handles the "mouseout" event.
func OnMouseOut(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventMouseOut, modifiers), expr)
}

// OnMouseMove handles the "mousemove" event.
func OnMouseMove(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventMouseMove, modifiers), expr)
}

// OnMouseEnter handles the "mouseenter" event.
func OnMouseEnter(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventMouseEnter, modifiers), expr)
}

// OnMouseLeave handles the "mouseleave" event.
func OnMouseLeave(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventMouseLeave, modifiers), expr)
}

// OnContextMenu handles the "contextmenu" event.
func OnContextMenu(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventContextMenu, modifiers), expr)
}

// ---------------------------------------------------------------------------
// Keyboard events
// ---------------------------------------------------------------------------

// OnKeyDown handles the "keydown" event.
func OnKeyDown(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventKeyDown, modifiers), expr)
}

// OnKeyUp handles the "keyup" event.
func OnKeyUp(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventKeyUp, modifiers), expr)
}

// OnKeyPress handles the "keypress" event (deprecated but still used).
func OnKeyPress(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventKeyPress, modifiers), expr)
}

// ---------------------------------------------------------------------------
// Focus events
// ---------------------------------------------------------------------------

// OnFocus handles the "focus" event.
func OnFocus(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventFocus, modifiers), expr)
}

// OnBlur handles the "blur" event.
func OnBlur(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventBlur, modifiers), expr)
}

// OnFocusIn handles the "focusin" event.
func OnFocusIn(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventFocusIn, modifiers), expr)
}

// OnFocusOut handles the "focusout" event.
func OnFocusOut(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventFocusOut, modifiers), expr)
}

// ---------------------------------------------------------------------------
// Form events
// ---------------------------------------------------------------------------

// OnSubmit handles the "submit" event. Datastar automatically prevents default.
func OnSubmit(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventSubmit, modifiers), expr)
}

// OnReset handles the "reset" event.
func OnReset(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventReset, modifiers), expr)
}

// OnInput handles the "input" event.
func OnInput(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventInput, modifiers), expr)
}

// OnChange handles the "change" event.
func OnChange(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventChange, modifiers), expr)
}

// OnInvalid handles the "invalid" event.
func OnInvalid(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventInvalid, modifiers), expr)
}

// OnSelect handles the "select" event.
func OnSelect(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventSelect, modifiers), expr)
}

// ---------------------------------------------------------------------------
// Drag events
// ---------------------------------------------------------------------------

// OnDrag handles the "drag" event.
func OnDrag(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventDrag, modifiers), expr)
}

// OnDragStart handles the "dragstart" event.
func OnDragStart(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventDragStart, modifiers), expr)
}

// OnDragEnd handles the "dragend" event.
func OnDragEnd(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventDragEnd, modifiers), expr)
}

// OnDragOver handles the "dragover" event.
func OnDragOver(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventDragOver, modifiers), expr)
}

// OnDragEnter handles the "dragenter" event.
func OnDragEnter(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventDragEnter, modifiers), expr)
}

// OnDragLeave handles the "dragleave" event.
func OnDragLeave(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventDragLeave, modifiers), expr)
}

// OnDrop handles the "drop" event.
func OnDrop(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventDrop, modifiers), expr)
}

// ---------------------------------------------------------------------------
// Touch events
// ---------------------------------------------------------------------------

// OnTouchStart handles the "touchstart" event.
func OnTouchStart(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventTouchStart, modifiers), expr)
}

// OnTouchEnd handles the "touchend" event.
func OnTouchEnd(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventTouchEnd, modifiers), expr)
}

// OnTouchMove handles the "touchmove" event.
func OnTouchMove(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventTouchMove, modifiers), expr)
}

// OnTouchCancel handles the "touchcancel" event.
func OnTouchCancel(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventTouchCancel, modifiers), expr)
}

// ---------------------------------------------------------------------------
// Pointer events
// ---------------------------------------------------------------------------

// OnPointerDown handles the "pointerdown" event.
func OnPointerDown(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventPointerDown, modifiers), expr)
}

// OnPointerUp handles the "pointerup" event.
func OnPointerUp(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventPointerUp, modifiers), expr)
}

// OnPointerMove handles the "pointermove" event.
func OnPointerMove(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventPointerMove, modifiers), expr)
}

// OnPointerOver handles the "pointerover" event.
func OnPointerOver(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventPointerOver, modifiers), expr)
}

// OnPointerOut handles the "pointerout" event.
func OnPointerOut(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventPointerOut, modifiers), expr)
}

// OnPointerEnter handles the "pointerenter" event.
func OnPointerEnter(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventPointerEnter, modifiers), expr)
}

// OnPointerLeave handles the "pointerleave" event.
func OnPointerLeave(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventPointerLeave, modifiers), expr)
}

// OnPointerCancel handles the "pointercancel" event.
func OnPointerCancel(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventPointerCancel, modifiers), expr)
}

// OnGotPointerCapture handles the "gotpointercapture" event.
func OnGotPointerCapture(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventGotPointerCapture, modifiers), expr)
}

// OnLostPointerCapture handles the "lostpointercapture" event.
func OnLostPointerCapture(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventLostPointerCapture, modifiers), expr)
}

// ---------------------------------------------------------------------------
// Scroll / Wheel events
// ---------------------------------------------------------------------------

// OnScroll handles the "scroll" event.
func OnScroll(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventScroll, modifiers), expr)
}

// OnWheel handles the "wheel" event.
func OnWheel(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventWheel, modifiers), expr)
}

// ---------------------------------------------------------------------------
// Animation / Transition events
// ---------------------------------------------------------------------------

// OnAnimationStart handles the "animationstart" event.
func OnAnimationStart(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventAnimationStart, modifiers), expr)
}

// OnAnimationEnd handles the "animationend" event.
func OnAnimationEnd(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventAnimationEnd, modifiers), expr)
}

// OnAnimationIteration handles the "animationiteration" event.
func OnAnimationIteration(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventAnimationIteration, modifiers), expr)
}

// OnTransitionEnd handles the "transitionend" event.
func OnTransitionEnd(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventTransitionEnd, modifiers), expr)
}

// ---------------------------------------------------------------------------
// Media events
// ---------------------------------------------------------------------------

// OnLoad handles the "load" event.
func OnLoad(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventLoad, modifiers), expr)
}

// OnError handles the "error" event.
func OnError(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventError, modifiers), expr)
}

// ---------------------------------------------------------------------------
// Clipboard events
// ---------------------------------------------------------------------------

// OnCopy handles the "copy" event.
func OnCopy(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventCopy, modifiers), expr)
}

// OnCut handles the "cut" event.
func OnCut(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventCut, modifiers), expr)
}

// OnPaste handles the "paste" event.
func OnPaste(expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(eventPaste, modifiers), expr)
}

// ---------------------------------------------------------------------------
// Custom event escape hatch
// ---------------------------------------------------------------------------

// OnEvent handles any event by name. Use for custom events or web component
// events not covered by the typed functions above.
//
//	{ ds.OnEvent("table-select", "$table.selected = evt.detail.ids")... }
//
// See https://data-star.dev/reference/attributes#data-on
func OnEvent(event, expr string, modifiers ...Modifier) templ.Attributes {
	return val(on(event, modifiers), expr)
}
