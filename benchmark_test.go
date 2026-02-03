package ds_test

import (
	"testing"

	ds "github.com/Yacobolo/datastar-templ"
)

// ===========================================================================
// Benchmark Tests for V2 Performance
// ===========================================================================

func BenchmarkSignals(b *testing.B) {
	b.Run("simple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Signals(
				ds.Int("count", 42),
				ds.String("message", "hello"),
			)
		}
	})

	b.Run("complex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Signals(
				ds.Int("count", 42),
				ds.String("message", "hello"),
				ds.Bool("enabled", true),
				ds.Float("price", 19.99),
			)
		}
	})

	b.Run("with_json", func(b *testing.B) {
		data := []int{1, 2, 3, 4, 5}
		for i := 0; i < b.N; i++ {
			_ = ds.Signals(
				ds.Int("count", 42),
				ds.JSON("items", data),
			)
		}
	})
}

func BenchmarkClass(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Class(ds.C("hidden", "$isHidden"))
		}
	})

	b.Run("multiple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Class(
				ds.C("hidden", "$isHidden"),
				ds.C("font-bold", "$isBold"),
				ds.C("text-red-500", "$hasError"),
			)
		}
	})
}

func BenchmarkComputed(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Computed(ds.Comp("total", "$price * $qty"))
		}
	})

	b.Run("multiple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Computed(
				ds.Comp("total", "$price * $qty"),
				ds.Comp("tax", "$total * 0.1"),
			)
		}
	})
}

func BenchmarkAttr(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Attr(ds.A("title", "$tooltip"))
		}
	})

	b.Run("multiple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Attr(
				ds.A("title", "$tooltip"),
				ds.A("disabled", "$loading"),
			)
		}
	})
}

func BenchmarkStyle(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Style(ds.S("display", "$hiding && 'none'"))
		}
	})

	b.Run("multiple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Style(
				ds.S("display", "$hiding && 'none'"),
				ds.S("color", "$textColor"),
			)
		}
	})
}

func BenchmarkMerge(b *testing.B) {
	b.Run("two_attributes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Merge(
				ds.Signals(ds.Int("count", 0)),
				ds.OnClick("$count++"),
			)
		}
	})

	b.Run("five_attributes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Merge(
				ds.Signals(ds.Int("count", 0)),
				ds.Class(ds.C("hidden", "$isHidden")),
				ds.OnClick("toggle()"),
				ds.Text("$message"),
				ds.Show("$visible"),
			)
		}
	})
}

// Comparison benchmarks to show the improvement
func BenchmarkStringConversions(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Int("count", 42)
		}
	})

	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.String("message", "hello world")
		}
	})

	b.Run("bool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Bool("enabled", true)
		}
	})

	b.Run("float", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ds.Float("price", 19.99)
		}
	})
}
