//+build ignore

package json

import "testing"

func BenchmarkParseJSONByReflect(b *testing.B) {

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseJSONByReflect(jsonExample)
	}
}

func BenchmarkParseJSON(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		scanner := Scanner{
			sample: secret,
			data:   jsonExample,
		}
		scanner.SeekIn()
	}
}
