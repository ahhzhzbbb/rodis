package benchmark

import (
	"fmt"
	"rodis/internal/engine"
	"testing"
)

func BenchmarkQuickList(b *testing.B) {
	args := make([]string, 10000000)
	for i := range 10000000 {
		args[i] = "test" + fmt.Sprint(i)
	}

	for b.Loop() {
		engine.NewQuickList(10000000000, args)
	}
}

func BenchmarkZipList(b *testing.B) {
	args := make([]string, 10000000)
	for i := range 10000000 {
		args[i] = "test" + fmt.Sprint(i)
	}

	for b.Loop() {
		zl := engine.NewZipList()
		for _, arg := range args {
			zl.PushBack(arg)
		}
	}
}
