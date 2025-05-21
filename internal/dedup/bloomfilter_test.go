package dedup

import (
	"testing"
)

func TestBloomFilter(t *testing.T) {
	// Create a bloom filter with reasonable size for testing
	bf := NewBloomFilter(1000, 5)

	// Test empty filter
	if bf.Contains("https://example.com") {
		t.Errorf("Empty bloom filter should not contain any items")
	}

	// Add a URL
	bf.Add("https://example.com")

	// Test URL presence
	if !bf.Contains("https://example.com") {
		t.Errorf("Bloom filter should contain the added URL")
	}

	// Test URL absence
	if bf.Contains("https://different-example.com") {
		t.Errorf("Bloom filter should not contain URLs that were not added")
	}

	// Test multiple additions
	urls := []string{
		"https://example.org",
		"https://test.com",
		"https://golang.org",
		"https://github.com",
	}

	for _, url := range urls {
		bf.Add(url)
	}

	for _, url := range urls {
		if !bf.Contains(url) {
			t.Errorf("Bloom filter should contain all added URLs: %s", url)
		}
	}
}

func TestOptimalParameters(t *testing.T) {
	// Test calculation of optimal size
	size := CalculateOptimalSize(10000, 0.01) // 10K items, 1% false positive rate
	if size < 95000 || size > 100000 {        // Approximate range for these parameters
		t.Errorf("Optimal size calculation is incorrect, got %d", size)
	}

	// Test calculation of optimal hash functions
	hashFuncs := CalculateOptimalHashFunctions(size, 10000)
	if hashFuncs < 6 || hashFuncs > 8 { // Approximate range for these parameters
		t.Errorf("Optimal hash function calculation is incorrect, got %d", hashFuncs)
	}
}

func BenchmarkBloomFilter(b *testing.B) {
	// Create a bloom filter
	bf := NewBloomFilter(100000, 7)

	// Pre-populate with some URLs
	for i := 0; i < 1000; i++ {
		bf.Add("https://example.com/" + string(rune(i)))
	}

	b.ResetTimer()

	// Benchmark Add operation
	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bf.Add("https://benchmark.com/" + string(rune(i%1000)))
		}
	})

	// Benchmark Contains operation
	b.Run("Contains", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bf.Contains("https://example.com/" + string(rune(i%1000)))
		}
	})
}
