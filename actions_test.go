package ds_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-scheduler/pkg/ds"
)

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

func TestGet(t *testing.T) {
	t.Run("simple url", func(t *testing.T) {
		assert.Equal(t, "@get('/api/updates')", ds.Get("/api/updates"))
	})

	t.Run("format args", func(t *testing.T) {
		assert.Equal(t, "@get('/api/todos/42')", ds.Get("/api/todos/%d", 42))
	})

	t.Run("multiple format args", func(t *testing.T) {
		assert.Equal(t, "@get('/api/users/5/todos/42')", ds.Get("/api/users/%d/todos/%d", 5, 42))
	})

	t.Run("string format arg", func(t *testing.T) {
		assert.Equal(t, "@get('/api/search?q=hello')", ds.Get("/api/search?q=%s", "hello"))
	})

	t.Run("single opt", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/updates',{requestCancellation: 'disabled'})",
			ds.Get("/api/updates", ds.Opt("requestCancellation", "disabled")),
		)
	})

	t.Run("multiple opts", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/updates',{requestCancellation: 'disabled', contentType: 'json'})",
			ds.Get("/api/updates",
				ds.Opt("requestCancellation", "disabled"),
				ds.Opt("contentType", "json"),
			),
		)
	})

	t.Run("format args and opts", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/todos/42',{openWhenHidden: 'true'})",
			ds.Get("/api/todos/%d", 42, ds.Opt("openWhenHidden", "true")),
		)
	})

	t.Run("raw opt", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/updates',{openWhenHidden: true})",
			ds.Get("/api/updates", ds.OptRaw("openWhenHidden", "true")),
		)
	})

	t.Run("mixed opt types", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/updates',{requestCancellation: 'disabled', openWhenHidden: true})",
			ds.Get("/api/updates",
				ds.Opt("requestCancellation", "disabled"),
				ds.OptRaw("openWhenHidden", "true"),
			),
		)
	})

	t.Run("raw opt with number", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/updates',{retryMaxCount: 10})",
			ds.Get("/api/updates", ds.OptRaw("retryMaxCount", "10")),
		)
	})

	t.Run("raw opt with object", func(t *testing.T) {
		assert.Equal(t,
			"@get('/api/updates',{filterSignals: {include: /^foo/}})",
			ds.Get("/api/updates", ds.OptRaw("filterSignals", "{include: /^foo/}")),
		)
	})

	t.Run("format args interleaved with opts", func(t *testing.T) {
		// Options can appear anywhere in the variadic args; they're partitioned by type
		assert.Equal(t,
			"@get('/api/users/5/todos/42',{requestCancellation: 'disabled'})",
			ds.Get("/api/users/%d/todos/%d", 5, ds.Opt("requestCancellation", "disabled"), 42),
		)
	})
}

// ---------------------------------------------------------------------------
// All verbs
// ---------------------------------------------------------------------------

func TestAllVerbs(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string, ...any) string
		verb string
	}{
		{"Get", ds.Get, "get"},
		{"Post", ds.Post, "post"},
		{"Put", ds.Put, "put"},
		{"Patch", ds.Patch, "patch"},
		{"Delete", ds.Delete, "delete"},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/simple", func(t *testing.T) {
			assert.Equal(t, "@"+tt.verb+"('/api/foo')", tt.fn("/api/foo"))
		})

		t.Run(tt.name+"/format_args", func(t *testing.T) {
			assert.Equal(t, "@"+tt.verb+"('/api/foo/42')", tt.fn("/api/foo/%d", 42))
		})

		t.Run(tt.name+"/with_opt", func(t *testing.T) {
			assert.Equal(t,
				"@"+tt.verb+"('/api/foo',{key: 'val'})",
				tt.fn("/api/foo", ds.Opt("key", "val")),
			)
		})
	}
}

// ---------------------------------------------------------------------------
// Post, Put, Patch, Delete specific tests
// ---------------------------------------------------------------------------

func TestPost(t *testing.T) {
	assert.Equal(t, "@post('/api/workcenters')", ds.Post("/api/workcenters"))
}

func TestPut(t *testing.T) {
	assert.Equal(t, "@put('/api/todos/42')", ds.Put("/api/todos/%d", 42))
}

func TestPatch(t *testing.T) {
	assert.Equal(t, "@patch('/api/workcenters/pagesize')", ds.Patch("/api/workcenters/pagesize"))
}

func TestDelete(t *testing.T) {
	assert.Equal(t, "@delete('/api/todos/42')", ds.Delete("/api/todos/%d", 42))
}

// ---------------------------------------------------------------------------
// Composition with attribute helpers
// ---------------------------------------------------------------------------

func TestActionComposition(t *testing.T) {
	t.Run("init with get and opts", func(t *testing.T) {
		attrs := ds.Init(ds.Get("/api/updates", ds.Opt("requestCancellation", "disabled")))
		require.Len(t, attrs, 1)
		assert.Equal(t,
			"@get('/api/updates',{requestCancellation: 'disabled'})",
			attrs["data-init"],
		)
	})

	t.Run("onclick with post", func(t *testing.T) {
		attrs := ds.OnClick(ds.Post("/api/workcenters"))
		require.Len(t, attrs, 1)
		assert.Equal(t, "@post('/api/workcenters')", attrs["data-on:click"])
	})

	t.Run("oninput with post and debounce", func(t *testing.T) {
		attrs := ds.OnInput(ds.Post("/api/search"), ds.ModDebounce, ds.Ms(300))
		require.Len(t, attrs, 1)
		assert.Equal(t, "@post('/api/search')", attrs["data-on:input__debounce.300ms"])
	})

	t.Run("onchange with patch", func(t *testing.T) {
		attrs := ds.OnChange(ds.Patch("/api/pagesize"))
		require.Len(t, attrs, 1)
		assert.Equal(t, "@patch('/api/pagesize')", attrs["data-on:change"])
	})

	t.Run("init with delay and get", func(t *testing.T) {
		attrs := ds.Init(ds.Get("/api/updates"), ds.ModDelay, ds.Ms(500))
		require.Len(t, attrs, 1)
		assert.Equal(t, "@get('/api/updates')", attrs["data-init__delay.500ms"])
	})
}

// ---------------------------------------------------------------------------
// BuildUpdatesInitURL equivalent
// ---------------------------------------------------------------------------

func TestBuildUpdatesInitURLEquivalent(t *testing.T) {
	// Reproduces the output of handlers.BuildUpdatesInitURL using ds helpers
	url := "/api/workcenters/updates?page=1&sortColumn=title&sortDir=asc"
	result := ds.Get(url, ds.Opt("requestCancellation", "disabled"))
	assert.Equal(t,
		"@get('/api/workcenters/updates?page=1&sortColumn=title&sortDir=asc',{requestCancellation: 'disabled'})",
		result,
	)
}
