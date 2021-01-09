package log

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func getstr() interface{} {
	return "1111111111111111111111"
}

func getint() interface{} {
	return 1 << 16
}

func gettime() interface{} {
	return time.Now()
}

func fmtstr(b []byte, v string) []byte {
	return b
}

func fmtint(b []byte, v interface{}) []byte {
	return b
}

func fmttime(b []byte, v interface{}) []byte {
	return b
}

// func BenchmarkFieldStr(b *testing.B) {
// 	fmt := func(b []byte, v interface{}) []byte {
// 		return fmtstr(b, v.(string))
// 	}

// 	buf := make([]byte, 8)
// 	val := getstr()

// 	for i := 0; i < b.N; i++ {
// 		fmt(buf, val)
// 	}
// }

// func BenchmarkFieldInt(b *testing.B) {
// 	fmt := func(b []byte, v interface{}) []byte {
// 		return fmtint(b, v.(int))
// 	}

// 	buf := make([]byte, 8)
// 	val := getint()

// 	for i := 0; i < b.N; i++ {
// 		fmt(buf, val)
// 	}
// }

// func BenchmarkFieldTime(b *testing.B) {
// 	fmt := func(b []byte, v interface{}) []byte {
// 		return fmttime(b, v.(time.Time))
// 	}

// 	buf := make([]byte, 8)
// 	val := gettime()

// 	for i := 0; i < b.N; i++ {
// 		fmt(buf, val)
// 	}
// }

// type encoder interface {
// 	Encode(k string, v interface{})
// }

// type strencoder int

// func (strencoder) Encode(k string, v interface{}) {
// }

// func BenchmarkIntefaceCast(b *testing.B) {
// 	var f interface{} = strencoder(0)

// 	for i := 0; i < b.N; i++ {
// 		f.(encoder).Encode("key", "value")
// 	}
// }

// func BenchmarkInterfaceIndex(b *testing.B) {
// 	fns := make([]encoder, 1)
// 	fns[0] = strencoder(0)

// 	for i := 0; i < b.N; i++ {
// 		fns[0].Encode("key", "value")
// 	}
// }

// func BenchmarkFunctionIndex(b *testing.B) {
// 	fns := make([]func(string, interface{}), 1)
// 	fns[0] = func(string, interface{}) {}

// 	for i := 0; i < b.N; i++ {
// 		fns[0]("key", "value")
// 	}
// }

func BenchmarkZeroLog(b *testing.B) {
	l := zerolog.New(ioutil.Discard)
	l.Level(zerolog.InfoLevel)

	b.Run("0", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Info().Msg("info message")
		}
	})

	b.Run("1", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Info().
				Str("k1", "v1").
				Msg("info message")
		}
	})

	b.Run("2", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Info().
				Str("k1", "v1").
				Str("k2", "v2").
				Msg("info message")
		}
	})

	b.Run("3", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Info().
				Str("k", "v1").
				Str("k2", "v2").
				Str("k3", "v3").
				Msg("info message")
		}
	})

	b.Run("10", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Info().
				Str("k1", "v1").
				Str("k2", "v2").
				Str("k3", "v3").
				Str("k4", "v4").
				Str("k5", "v5").
				Str("k6", "v6").
				Str("k7", "v7").
				Str("k8", "v8").
				Str("k9", "v9").
				Str("k10", "v10").
				Msg("info message")
		}
	})

	b.Run("10x", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Info().
					Str("k1", "v1").
					Str("k2", "v2").
					Str("k3", "v3").
					Str("k4", "v4").
					Str("k5", "v5").
					Str("k6", "v6").
					Str("k7", "v7").
					Str("k8", "v8").
					Str("k9", "v9").
					Str("k10", "v10").
					Msg("info message")
			}
		})
	})
}

func BenchmarkLogFields(b *testing.B) {
	s := ioutil.Discard
	w := &memwriter{entries: make([]mementry, 1<<10)}

	l := New(Config{
		Level:  INFO,
		Stream: s,
		Writer: w,
	})

	b.Run("0", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Info("info message")
		}
	})

	b.Run("1", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Info("info message",
				Field{Key: "k1", Val: "v1"},
			)
		}
	})

	b.Run("2", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Info("info message",
				Field{Key: "k1", Val: "v1"},
				Field{Key: "k2", Val: "v2"},
			)
		}
	})

	b.Run("3", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Info("info message",
				Field{Key: "k1", Val: "v1"},
				Field{Key: "k2", Val: "v2"},
				Field{Key: "k3", Val: "v3"},
			)
		}
	})

	b.Run("10", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Info("info message",
				Field{Key: "k1", Val: "v1"},
				Field{Key: "k2", Val: "v2"},
				Field{Key: "k3", Val: "v3"},
				Field{Key: "k4", Val: "v4"},
				Field{Key: "k5", Val: "v5"},
				Field{Key: "k6", Val: "v6"},
				Field{Key: "k7", Val: "v7"},
				Field{Key: "k8", Val: "v8"},
				Field{Key: "k9", Val: "v9"},
				Field{Key: "k10", Val: "v10"},
			)
		}
	})

	b.Run("10x", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Info("info message",
					Field{Key: "k1", Val: "v1"},
					Field{Key: "k2", Val: "v2"},
					Field{Key: "k3", Val: "v3"},
					Field{Key: "k4", Val: "v4"},
					Field{Key: "k5", Val: "v5"},
					Field{Key: "k6", Val: "v6"},
					Field{Key: "k7", Val: "v7"},
					Field{Key: "k8", Val: "v8"},
					Field{Key: "k9", Val: "v9"},
					Field{Key: "k10", Val: "v10"},
				)
			}
		})
	})
}
