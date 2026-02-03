<h1 align="center">datastar-templ</h1>

<p align="center">
  <strong>Type-safe Datastar attribute helpers for templ templates.</strong>
  <br>
  <a href="https://pkg.go.dev/github.com/yacobolo/datastar-templ">
    <img src="https://img.shields.io/badge/go-reference-007d9c?logo=go&logoColor=white&style=flat-square" alt="Go Reference">
  </a>
  <a href="https://goreportcard.com/report/github.com/yacobolo/datastar-templ">
    <img src="https://goreportcard.com/badge/github.com/yacobolo/datastar-templ?style=flat-square" alt="Go Report Card">
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License: MIT">
  </a>
  <a href="https://github.com/yacobolo/datastar-templ/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/yacobolo/datastar-templ/ci.yml?branch=main&style=flat-square&label=CI" alt="CI Status">
  </a>
</p>

<p align="center">
  <img src="assets/mascot.png" alt="datastar-templ mascot" width="600">
</p>

---

`datastar-templ` is a Go library that provides compile-time type safety for Datastar attributes in templ templates. It bridges the gap between Go's templ templating system and Datastar's hypermedia framework, enabling you to build reactive web applications with full IDE autocomplete and type checking.

## Features

- **Type-Safe**: Compile-time checks for Datastar attributes with full IDE support
- **Complete Coverage**: 60+ DOM events, HTTP actions, signals, and modifiers
- **templ Integration**: Native templ.Attributes for seamless template usage

## Installation

```bash
go get github.com/yacobolo/datastar-templ
```

Tested with **Datastar 1.0.0-RC.7**. [Get started with Datastar](https://data-star.dev/guide/getting-started).

## Usage

Import the package (commonly aliased as `ds`):

```go
import ds "github.com/yacobolo/datastar-templ"
```

### Quick Start Example

```go
templ TodoApp() {
    <div { ds.Signals(map[string]any{
        "todos": []Todo{},
        "newTodo": "",
        "filter": "",
    })... }>
        // Data binding
        <input 
            type="text"
            { ds.Bind("newTodo")... }
            placeholder="New todo"
        />
        
        // Event handlers with modifiers + SSE actions
        <button { ds.OnClick(
            ds.Post("/todos"),
            ds.ModDebounce,
            ds.Ms(300),
        )... }>
            Add Todo
        </button>
        
        // Conditional rendering + merging attributes
        <div { ds.Merge(
            ds.Show("$todos.length > 0"),
            ds.Class("active", "$filter !== ''"),
        )... }>
            <span { ds.Text("$todos.length + ' items'")... }></span>
        </div>
        
        // Event handlers
        <input 
            type="search"
            { ds.Bind("filter")... }
            { ds.OnInput(
                ds.Get("/search?q=$filter"),
                ds.ModDebounce,
                ds.Ms(300),
            )... }
        />
    </div>
}
```

## API Overview

See the [Go package documentation](https://pkg.go.dev/github.com/yacobolo/datastar-templ) for the complete API reference including:

- **60+ Event Handlers**: OnClick, OnInput, OnSubmit, OnKeyDown, etc.
- **HTTP Actions**: Get, Post, Put, Patch, Delete with options
- **Signal Management**: Signals, Computed, Bind, SignalKey
- **DOM Helpers**: Text, Show, Class, Style, Attr
- **Modifiers**: Debounce, Throttle, Once, Passive, Capture, etc.
- **Watchers**: OnIntersect, OnInterval, OnSignalPatch
- **Utilities**: Merge, Ref, Indicator, Init, Effect

## Development

Run tests:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

---

<p align="center">
  <a href="https://data-star.dev">Datastar</a> •
  <a href="https://templ.guide">templ</a> •
  <a href="https://pkg.go.dev/github.com/yacobolo/datastar-templ">API Reference</a>
</p>
