package ds_test

import (
	"context"
	"io"
	"strings"
	"testing"
)

// ===========================================================================
// Benchmarks
// Templates are defined in benchmark_poc.templ
// ===========================================================================

func BenchmarkPlain(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		templatePlain(42, "barr").Render(ctx, io.Discard)
	}
}

// Note: V1 (map-based) removed in V2.0 rewrite

func BenchmarkV2(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		templateV2(42, "barr").Render(ctx, io.Discard)
	}
}

// ===========================================================================
// Unit Tests - Verify Output Correctness
// ===========================================================================

func TestOutputEquivalence(t *testing.T) {
	tests := []struct {
		name      string
		foo       int
		bar       string
		wantPlain string
		wantV2    string
	}{
		{
			name:      "basic integers and strings",
			foo:       42,
			bar:       "barr",
			wantPlain: `<div data-signals="{foo: 42, bar: &#34;barr&#34;}"></div>`,
			wantV2:    `<div data-signals="{foo: 42, bar: &#34;barr&#34;}"></div>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Test Plain
			var bufPlain strings.Builder
			if err := templatePlain(tt.foo, tt.bar).Render(ctx, &bufPlain); err != nil {
				t.Fatalf("Plain render failed: %v", err)
			}
			gotPlain := bufPlain.String()

			// Test V2
			var bufV2 strings.Builder
			if err := templateV2(tt.foo, tt.bar).Render(ctx, &bufV2); err != nil {
				t.Fatalf("V2 render failed: %v", err)
			}
			gotV2 := bufV2.String()

			t.Logf("Plain: %s", gotPlain)
			t.Logf("V2:    %s", gotV2)

			// Both should produce the same output
			if gotPlain != gotV2 {
				t.Errorf("Plain and V2 output differs:\nPlain: %s\nV2:    %s", gotPlain, gotV2)
			}
		})
	}
}

// ===========================================================================
// Additional Benchmarks - Complex Scenarios
// Templates for complex scenarios are in benchmark_poc.templ
// ===========================================================================

// Note: V1Complex (map-based) removed in V2.0 rewrite

func BenchmarkV2Complex(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		templateV2Complex(42, "hello", true, 19.99).Render(ctx, io.Discard)
	}
}
