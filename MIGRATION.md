# Migration Guide: V1 to V2

This guide helps you migrate from datastar-templ V1 (map-based API) to V2 (type-safe variadic API).

## Why Migrate?

V2 offers significant improvements:

- **2-3x faster performance** (~200-300ns vs ~780ns)
- **~60% less memory allocation** (~400B vs ~1200B)
- **Type safety at compile time** - catch errors before runtime
- **Eliminates odd-pair bugs** - no more "forgot a string" panics
- **Better IDE support** - autocomplete for signal types
- **Unified API** - Single `Pair()` helper for all attribute bindings

## Breaking Changes

All pair-based functions now use typed helpers instead of variadic strings or maps.

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
    ds.Pair("hidden", "$isHidden"),
    ds.Pair("font-bold", "$isBold"),
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
    ds.Pair("total", "$price * $qty"),
    ds.Pair("tax", "$total * 0.1"),
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
    ds.Pair("title", "$tooltip"),
    ds.Pair("disabled", "$loading"),
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
    ds.Pair("display", "$hiding && 'none'"),
    ds.Pair("color", "'red'"),
)
```

**Shorthand:** You can also use `ds.P()` instead of `ds.Pair()` for brevity:
```go
ds.Class(ds.P("btn-primary", "$isMain"))
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
    // Add typed helpers here based on value types
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
ds.Class(ds.Pair("\1", "\2"))
```

### Step 3: Update Computed Signals

Same pattern as Class:

**Search pattern:**
```regex
ds\.Computed\("([^"]+)", "([^"]+)"
```

**Replace with:**
```go
ds.Computed(ds.Pair("\1", "\2"))
```

### Step 4: Update Attr Bindings

**Search pattern:**
```regex
ds\.Attr\("([^"]+)", "([^"]+)"
```

**Replace with:**
```go
ds.Attr(ds.Pair("\1", "\2"))
```

### Step 5: Update Style Bindings

**Search pattern:**
```regex
ds\.Style\("([^"]+)", "([^"]+)"
```

**Replace with:**
```go
ds.Style(ds.Pair("\1", "\2"))
```

## Complete Example: Before & After

### Before (V1)

```go
templ TodoList(todos []Todo) {
    <div { ds.Signals(map[string]any{
        "count": len(todos),
        "loading": false,
    })... }>
        <ul { ds.Class("hidden", "$loading", "opacity-50", "$count == 0")... }>
            for _, todo := range todos {
                <li { ds.Attr("data-id", strconv.Itoa(todo.ID))... }>
                    { todo.Title }
                </li>
            }
        </ul>
    </div>
}
```

### After (V2)

```go
templ TodoList(todos []Todo) {
    <div { ds.Signals(
        ds.Int("count", len(todos)),
        ds.Bool("loading", false),
    )... }>
        <ul { ds.Class(
            ds.Pair("hidden", "$loading"),
            ds.Pair("opacity-50", "$count == 0"),
        )... }>
            for _, todo := range todos {
                <li { ds.Attr(ds.Pair("data-id", strconv.Itoa(todo.ID)))... }>
                    { todo.Title }
                </li>
            }
        </ul>
    </div>
}
```

## Why `Pair()` is Unified

V2 uses a single `Pair()` helper for all attribute bindings (Class, Computed, Attr, Style) for better API consistency:

**Benefits:**
- **Reduced cognitive load** - Learn ONE helper name instead of 4 different ones
- **Clear semantics** - `Pair` clearly communicates "key-value pairing"
- **API consistency** - No mix of long names (`Int`, `String`) and short cryptic names (`C`, `A`, `S`)
- **Better distinction** - Data transformation (`Int`, `String`, `Bool`) vs expression binding (`Pair`)

## Common Patterns

### Multiple Signals
```go
// V1
ds.Signals(map[string]any{
    "user": currentUser,
    "todos": todoList,
    "filter": "all",
})

// V2
ds.Signals(
    ds.JSON("user", currentUser),
    ds.JSON("todos", todoList),
    ds.String("filter", "all"),
)
```

### Conditional Classes
```go
// V1
ds.Class("active", "$isActive", "disabled", "$isDisabled")

// V2
ds.Class(
    ds.Pair("active", "$isActive"),
    ds.Pair("disabled", "$isDisabled"),
)
```

### Dynamic Styles
```go
// V1
ds.Style("color", "$textColor", "display", "$visible ? 'block' : 'none'")

// V2
ds.Style(
    ds.Pair("color", "$textColor"),
    ds.Pair("display", "$visible ? 'block' : 'none'"),
)
```

### Computed Values
```go
// V1
ds.Computed("total", "$price * $qty", "tax", "$total * 0.1")

// V2
ds.Computed(
    ds.Pair("total", "$price * $qty"),
    ds.Pair("tax", "$total * 0.1"),
)
```

## Testing Your Migration

After migrating, run your tests to ensure everything works:

```bash
# Build your project
go build ./...

# Run tests
go test ./...

# Run your templ generation
templ generate
```

## Performance Notes

V2 is significantly faster than V1:

| Metric | V1 | V2 | Improvement |
|--------|----|----|-------------|
| **Time** | ~780 ns/op | ~200-300 ns/op | **2-3x faster** |
| **Memory** | ~1200 B/op | ~400 B/op | **60% less** |
| **Allocations** | ~19 allocs/op | ~5-8 allocs/op | **60% fewer** |

The performance improvements come from:
- Eliminating map allocations
- Direct string building with `strings.Builder`
- Precise capacity pre-allocation
- `sync.Pool` for builder reuse
- No reflection or JSON marshaling for primitives

## Need Help?

If you encounter issues during migration:
1. Check the [README](README.md) for updated API examples
2. Review the [test files](ds_test.go) for comprehensive usage examples
3. Open an issue on GitHub with your specific use case
