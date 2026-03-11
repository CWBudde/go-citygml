package gml

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cwbudde/go-citygml/types"
)

func BenchmarkParsePos_3D(b *testing.B) {
	input := "500123.456 5700234.789 42.5"

	b.ReportAllocs()

	for b.Loop() {
		_, _, err := ParsePos(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParsePosList_15Coords(b *testing.B) {
	input := "0 0 0 10 0 0 10 10 0 0 10 0 0 0 0"

	b.ReportAllocs()

	for b.Loop() {
		_, _, err := ParsePosList(input, types.Dim3D)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParsePosList_Large(b *testing.B) {
	// 100 3D points = 300 values.
	vals := make([]string, 300)
	for i := range vals {
		vals[i] = fmt.Sprintf("%d.%d", i*7, i%10)
	}

	input := strings.Join(vals, " ")

	b.ReportAllocs()

	for b.Loop() {
		_, _, err := ParsePosList(input, types.Dim3D)
		if err != nil {
			b.Fatal(err)
		}
	}
}
