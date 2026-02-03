# Migration Guide: V1 to V2

This guide helps you migrate from datastar-templ V1 (map-based API) to V2 (type-safe variadic API).

## Why Migrate?

V2 offers significant improvements:

- **2-3x faster performance** (~200-300ns vs ~780ns)
- **~60% less memory allocation** (~400B vs ~1200B)
- **Type safety at compile time** - catch errors before runtime
- **Eliminates odd-pair bugs** - no more "forgot a string" panics
- **Better IDE support** - autocomplete for signal types

## Breaking Changes

All pair-based functions now use typed pairs instead of variadic strings or maps.

### 1. Signals() - From map to typed helpers

**Before (V1):**
```go
ds.Signals(map[string]any{
    "count": 0,
    "message": "Hello",
    "isOpen": true,
})
```

**After (V2):**
```go
ds.Signals(
    ds.Int("count", 0),
    ds.String("message", "Hello"),
    ds.Bool("isOpen", true),
)
```

### Signal Type Helpers

| V1 Type | V2 Helper | Example |
|---------|-----------|---------|
| `int` | `ds.Int(key, value)` | `ds.Int("count", 42)` |
| `string` | `ds.String(key, value)` | `ds.String("name", "John")` |
| `bool` | `ds.Bool(key, value)` | `ds.Bool("enabled", true)` |
| `float64` | `ds.Float(key, value)` | `ds.Float("price", 19.99)` |
| `any` (complex) | `ds.JSON(key, value)` | `ds.JSON("user", userData)` |

### 2. Class() - From flat strings to typed pairs

**Before (V1):**
```go
ds.Class("hidden", "$isHidden", "font-bold", "$isBold")
```

**After (V2):**
```go
ds.Class(
    ds.C("hidden", "$isHidden"),
    ds.C("font-bold", "$isBold"),
)
```

### 3. Computed() - From flat strings to typed pairs

**Before (V1):**
```go
ds.Computed("total", "$price * $qty", "tax", "$total * 0.1")
```

**After (V2):**
```go
ds.Computed(
    ds.Comp("total", "$price * $qty"),
    ds.Comp("tax", "$total * 0.1"),
)
```

### 4. Attr() - From flat strings to typed pairs

**Before (V1):**
```go
ds.Attr("title", "$tooltip", "disabled", "$loading")
```

**After (V2):**
```go
ds.Attr(
    ds.A("title", "$tooltip"),
    ds.A("disabled", "$loading"),
)
```

### 5. Style() - From flat strings to typed pairs

**Before (V1):**
```go
ds.Style("display", "$hiding && 'none'", "color", "'red'")
```

**After (V2):**
```go
ds.Style(
    ds.S("display", "$hiding && 'none'"),
    ds.S("color", "'red'"),
)
```

## What Stays The Same

These APIs are **unchanged**:

- ✅ All event handlers: `OnClick()`, `OnInput()`, etc.
- ✅ HTTP actions: `Get()`, `Post()`, `Put()`, `Patch()`, `Delete()`
- ✅ DOM helpers: `Text()`, `Show()`, `Bind()`, `BindExpr()`
- ✅ Keyed variants: `SignalKey()`, `ClassKey()`, `ComputedKey()`, etc.
- ✅ Watchers: `OnIntersect()`, `OnInterval()`, `OnSignalPatch()`
- ✅ Utilities: `Merge()`, `Ref()`, `Indicator()`, `Init()`, `Effect()`
- ✅ Modifiers: `ModDebounce`, `ModThrottle`, `Ms()`, etc.

## Migration Checklist

### Step 1: Update Signal Declarations

Find all uses of `ds.Signals(map[string]any{...})` and convert to typed helpers.

**Search pattern:**
```regex
ds\.Signals\(map\[string\]any\{
```

**Replace with:**
```go
ds.Signals(
    // Add typed helpers here
)
```

### Step 2: Update Class Bindings

Find all uses of `ds.Class("...", "...", ...)` with flat strings.

**Search pattern:**
```regex
ds\.Class\("([^"]+)", "([^"]+)"
```

**Replace with:**
```go
ds.Class(ds.C("$1", "$2")
```

### Step 3: Update Computed Signals

Find all uses of `ds.Computed("...", "...", ...)`.

**Search pattern:**
```regex
ds\.Computed\("([^"]+)", "([^"]+)"
```

**Replace with:**
```go
ds.Computed(ds.Comp("$1", "$2")
```

### Step 4: Update Attr Bindings

Find all uses of `ds.Attr("...", "...", ...)`.

**Search pattern:**
```regex
ds\.Attr\("([^"]+)", "([^"]+)"
```

**Replace with:**
```go
ds.Attr(ds.A("$1", "$2")
```

### Step 5: Update Style Bindings

Find all uses of `ds.Style("...", "...", ...)`.

**Search pattern:**
```regex
ds\.Style\("([^"]+)", "([^"]+)"
```

**Replace with:**
```go
ds.Style(ds.S("$1", "$2")
```

### Step 6: Update Complex Signals

For signals containing arrays, objects, or nested data, use `ds.JSON()`:

**Before:**
```go
ds.Signals(map[string]any{
    "todos": []Todo{
        {ID: 1, Title: "Task 1"},
        {ID: 2, Title: "Task 2"},
    },
})
```

**After:**
```go
ds.Signals(
    ds.JSON("todos", []Todo{
        {ID: 1, Title: "Task 1"},
        {ID: 2, Title: "Task 2"},
    }),
)
```

### Step 7: Run Tests

```bash
go test ./...
```

The compiler will catch any remaining API mismatches.

## Example Migration

### Before (V1)
```go
templ Counter() {
    <div { ds.Signals(map[string]any{
        "count": 0,
        "step": 1,
    })... }>
        <button 
            { ds.OnClick("$count += $step")... }
            { ds.Class("active", "$count > 0")... }
        >
            Count: <span { ds.Text("$count")... }></span>
        </button>
        
        <div { ds.Computed("double", "$count * 2")... }></div>
        
        <input 
            type="number"
            { ds.Attr("value", "$step")... }
            { ds.Style("color", "$count > 10 ? 'red' : 'black'")... }
        />
    </div>
}
```

### After (V2)
```go
templ Counter() {
    <div { ds.Signals(
        ds.Int("count", 0),
        ds.Int("step", 1),
    )... }>
        <button 
            { ds.OnClick("$count += $step")... }
            { ds.Class(ds.C("active", "$count > 0"))... }
        >
            Count: <span { ds.Text("$count")... }></span>
        </button>
        
        <div { ds.Computed(ds.Comp("double", "$count * 2"))... }></div>
        
        <input 
            type="number"
            { ds.Attr(ds.A("value", "$step"))... }
            { ds.Style(ds.S("color", "$count > 10 ? 'red' : 'black'"))... }
        />
    </div>
}
```

## Common Gotchas

### 1. String Quoting

`ds.String()` automatically quotes strings for JavaScript. Don't double-quote:

**Wrong:**
```go
ds.String("name", "\"John\"")  // Produces: name: "\"John\""
```

**Right:**
```go
ds.String("name", "John")  // Produces: name: "John"
```

### 2. Empty Signals

Empty signals now use no arguments:

**Before:**
```go
ds.Signals(map[string]any{})
```

**After:**
```go
ds.Signals()  // No arguments needed
```

### 3. Multiple Pairs

Each pair must be wrapped in its helper:

**Wrong:**
```go
ds.Class(ds.C("hidden", "$isHidden", "bold", "$isBold"))  // Won't compile
```

**Right:**
```go
ds.Class(
    ds.C("hidden", "$isHidden"),
    ds.C("bold", "$isBold"),
)
```

## Benefits After Migration

### Compile-Time Safety

**Before:** Runtime panic
```go
ds.Class("hidden")  // Panics at runtime: odd number of pairs
```

**After:** Compile error
```go
ds.Class(ds.C("hidden"))  // Won't compile: missing expression
```

### Better IDE Support

V2 provides full autocomplete for type helpers:

```go
ds.Signals(
    ds.  // IDE shows: Int(), String(), Bool(), Float(), JSON()
)
```

### Performance Gains

Your application will automatically benefit from:
- ~2-3x faster signal creation
- ~60% less memory allocation
- ~50% fewer allocation operations

## Need Help?

- Review the [examples in the test files](./ds_test.go)
- Check the [API documentation](https://pkg.go.dev/github.com/Yacobolo/datastar-templ)
- Open an issue on [GitHub](https://github.com/Yacobolo/datastar-templ/issues)

## Version Support

- **V1.x**: Legacy API (map-based) - No longer maintained
- **V2.0+**: New API (type-safe variadic) - Current and recommended

We recommend migrating to V2 to benefit from performance improvements and type safety.
