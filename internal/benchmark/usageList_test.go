package benchmark

import (
	"fmt"
	"rodis/internal/engine"
	"runtime"
	"testing"
)

var args []string

func init() {
	args = make([]string, 1000000)
	// args = make([]string, 65000)

	for i := 0; i < len(args); i++ {
		args[i] = "test" + fmt.Sprint(i)
	}
}

func BenchmarkQuickListPushBack(b *testing.B) {
	ql := engine.NewQuickList(1024, args)

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

func BenchmarkLinkedListPushBack(b *testing.B) {
	ll := engine.NewDoubleLinkList()
	for _, arg := range args {
		ll.PushBack(engine.NewNode(arg))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ll.PushBack(engine.NewNode("test"))
	}
}

func BenchmarkQuickListPushFront(b *testing.B) {
	ql := engine.NewQuickList(1024, args)
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

func BenchmarkLinkedListPushFront(b *testing.B) {
	ll := engine.NewDoubleLinkList()
	for _, arg := range args {
		ll.PushBack(engine.NewNode(arg))
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ll.PushFront(engine.NewNode("test"))
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
	ql := engine.NewQuickList(1024, args)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if ql.Length() == 0 {
			b.StopTimer()
			ql = engine.NewQuickList(1024, args)
			b.StartTimer()
		}

		ql.PopFront()
	}
}

func BenchmarkZipListPopFront(b *testing.B) {
	zl := engine.NewZipList()
	for _, arg := range args {
		zl.PushBack(arg)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if zl.Length() == 0 {
			b.StopTimer()
			zl = engine.NewZipList()
			for _, arg := range args {
				zl.PushBack(arg)
			}
			b.StartTimer()
		}

		zl.PopFront()
	}
}

func BenchmarkLinkedListPopFront(b *testing.B) {
	ll := engine.NewDoubleLinkList()
	for _, arg := range args {
		ll.PushBack(engine.NewNode(arg))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if ll.Len() == 0 {
			b.StopTimer()
			ll = engine.NewDoubleLinkList()
			for _, arg := range args {
				ll.PushBack(engine.NewNode(arg))
			}
			b.StartTimer()
		}

		ll.PopFront()
	}
}

func BenchmarkQuickListInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		ql := engine.NewQuickList(1024, args)

		node, pos, found := ql.GetIndexOFElement("test33000")
		if !found {
			b.Fatalf("Element not found in QuickList")
		}

		b.StartTimer()

		ql.Insert(node, pos, "new-value")
	}
}

func BenchmarkZipListInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		zl := engine.NewZipList()

		for _, arg := range args {
			zl.PushBack(arg)
		}

		pos, found := zl.GetIndexOfElement("test33000")
		if !found {
			b.Fatalf("Element not found in ZipList")
		}

		b.StartTimer()

		zl.Insert(pos, "new-value")
	}
}

func BenchmarkLinkedListInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		ll := engine.NewDoubleLinkList()

		for _, arg := range args {
			ll.PushBack(engine.NewNode(arg))
		}

		node := ll.GetNodeByValue("test33000")
		if node == nil {
			b.Fatalf("Element not found in LinkedList")
		}

		b.StartTimer()

		ll.InsertAfter(node, "new-value")
	}
}

func TestQuickListMemoryUsage(t *testing.T) {
	var m1, m2 runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&m1)

	ql := engine.NewQuickList(10240, args)

	runtime.ReadMemStats(&m2)

	fmt.Printf(
		"QuickList Memory: %.2f MB\n",
		float64(m2.Alloc-m1.Alloc)/(1024*1024),
	)

	runtime.KeepAlive(ql)
}
func TestZipListMemoryUsage(t *testing.T) {
	var m1, m2 runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&m1)

	zl := engine.NewZipList()

	for _, arg := range args {
		zl.PushBack(arg)
	}

	runtime.ReadMemStats(&m2)

	fmt.Printf(
		"ZipList Memory: %.2f MB\n",
		float64(m2.Alloc-m1.Alloc)/(1024*1024),
	)

	runtime.KeepAlive(zl)
}

func TestLinkedListMemoryUsage(t *testing.T) {
	var m1, m2 runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&m1)

	ll := engine.NewDoubleLinkList()

	for _, arg := range args {
		ll.PushBack(engine.NewNode(arg))
	}

	runtime.ReadMemStats(&m2)

	fmt.Printf(
		"LinkedList Memory: %.2f MB\n",
		float64(m2.Alloc-m1.Alloc)/(1024*1024),
	)

	runtime.KeepAlive(ll)
}
