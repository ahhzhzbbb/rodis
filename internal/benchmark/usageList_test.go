package benchmark

import (
	"fmt"
	"rodis/internal/engine"
	"testing"
)

var args []string

func init() {
	args = make([]string, 100000)

	for i := 0; i < len(args); i++ {
		args[i] = "test" + fmt.Sprint(i)
	}
}

func BenchmarkQuickListPushBack(b *testing.B) {
	ql := engine.NewQuickList(5000, args)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ql.PushBack([]string{"test"})
	}
}

func BenchmarkZipListPushBack(b *testing.B) {
	zl := engine.NewZipList()
	for _, arg := range args {
		zl.PushBack(arg)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		zl.PushBack("test")
	}
}

func BenchmarkQuickListPushFront(b *testing.B) {
	ql := engine.NewQuickList(5000, args)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ql.PushFront([]string{"test"})
	}
}

func BenchmarkZipListPushFront(b *testing.B) {
	zl := engine.NewZipList()
	for _, arg := range args {
		zl.PushBack(arg)
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		zl.PushFront("test")
	}
}

// func BenchmarkQuickListPopFront(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		b.StopTimer()

// 		ql := engine.NewQuickList(5000, args)

// 		b.StartTimer()

// 		ql.PopFront()
// 	}
// }

// func BenchmarkZipListPopFront(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		b.StopTimer()

// 		zl := engine.NewZipList()
// 		for _, arg := range args {
// 			zl.PushBack(arg)
// 		}

// 		b.StartTimer()

// 		zl.PopFront()
// 	}
// }

func BenchmarkQuickListPopFront(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		ql := engine.NewQuickList(5000, args)

		b.StartTimer()

		for ql.Length() > 0 {
			ql.PopFront()
		}
	}
}

func BenchmarkZipListPopFront(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		zl := engine.NewZipList()
		for _, arg := range args {
			zl.PushBack(arg)
		}

		b.StartTimer()

		for zl.Length() > 0 {
			zl.PopFront()
		}
	}
}

func BenchmarkQuickListFindAndInsert(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		ql := engine.NewQuickList(5000, args)

		b.StartTimer()

		node, pos, found := ql.GetIndexOFElement("test50000")
		if !found {
			b.Fatal("element not found")
		}

		ql.Insert(node, pos, "new-value")
	}
}

func BenchmarkZipListFindAndInsert(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		zl := engine.NewZipList()
		for _, arg := range args {
			zl.PushBack(arg)
		}

		b.StartTimer()

		pos, found := zl.GetIndexOfElement("test50000")
		if !found {
			b.Fatal("element not found")
		}

		zl.Insert(pos, "new-value")
	}
}
