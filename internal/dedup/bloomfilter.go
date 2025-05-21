package dedup

import (
	"hash/fnv"
	"math"
	"sync"
)

// min returns the smaller of x or y
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// max returns the larger of x or y
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// BloomFilter is a space-efficient probabilistic data structure for membership testing
// It may produce false positives (indicating an element is present when it is not),
// but never false negatives (indicating an element is not present when it actually is)
type BloomFilter struct {
	bitset        []byte // Each byte represents 8 bits
	size          uint   // Size in bits
	mutex         sync.RWMutex
	hashFunctions []func(string) uint
}

// setBit sets the bit at position pos to 1
func (bf *BloomFilter) setBit(pos uint) {
	bytePos := pos / 8
	bitPos := pos % 8
	bf.mutex.Lock()
	defer bf.mutex.Unlock()
	bf.bitset[bytePos] |= 1 << bitPos
}

// getBit checks if the bit at position pos is set
func (bf *BloomFilter) getBit(pos uint) bool {
	bytePos := pos / 8
	bitPos := pos % 8
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()
	return (bf.bitset[bytePos] & (1 << bitPos)) != 0
}

// NewBloomFilter creates a new bloom filter with the given size and number of hash functions
// size: the size of the bit array
// numHashes: the number of hash functions to use (more hashes = lower false positive rate but slower)
func NewBloomFilter(size uint, numHashes int) *BloomFilter {
	bf := &BloomFilter{
		bitset:        make([]byte, (size+7)/8), // Round up to nearest byte
		size:          size,
		hashFunctions: make([]func(string) uint, numHashes),
	}

	// Create hash functions with different seeds
	for i := 0; i < numHashes; i++ {
		seed := i
		bf.hashFunctions[i] = func(s string) uint {
			h := fnv.New64a()
			h.Write([]byte(s))
			h.Write([]byte{byte(seed)})
			return uint(h.Sum64() % uint64(size))
		}
	}

	return bf
}

// Add adds a URL to the bloom filter
func (bf *BloomFilter) Add(url string) {
	hash := fnv.New64a()
	hash.Write([]byte(url))
	baseHash := hash.Sum64()

	for _, hashFn := range bf.hashFunctions {
		hashValue := hashFn(url)
		// Use both the base hash and the current hash to get better distribution
		combinedHash := uint64(hashValue) ^ baseHash
		index := uint(combinedHash) % bf.size
		bf.setBit(index)
	}
}

// Contains checks if a URL is likely in the bloom filter
// false means definitely not in the set, true means possibly in the set
func (bf *BloomFilter) Contains(url string) bool {
	hash := fnv.New64a()
	hash.Write([]byte(url))
	baseHash := hash.Sum64()

	for _, hashFn := range bf.hashFunctions {
		hashValue := hashFn(url)
		// Use both the base hash and the current hash to get better distribution
		combinedHash := uint64(hashValue) ^ baseHash
		index := uint(combinedHash) % bf.size
		if !bf.getBit(index) {
			return false
		}
	}
	return true
}

// CalculateOptimalSize calculates the optimal bloom filter size for the expected number of items
// and desired false positive rate
func CalculateOptimalSize(expectedItems int, falsePositiveRate float64) uint {
	// m = -n * ln(p) / (ln(2)^2)
	// where m is the number of bits, n is the number of expected elements, and p is the false positive probability
	if expectedItems <= 0 {
		expectedItems = 1000 // Default to avoid division by zero
	}
	if falsePositiveRate <= 0 || falsePositiveRate >= 1 {
		falsePositiveRate = 0.01 // Default to 1% false positive rate
	}
	m := -float64(expectedItems) * math.Log(falsePositiveRate) / (math.Ln2 * math.Ln2)
	return uint(m)
}

// CalculateOptimalHashFunctions calculates the optimal number of hash functions
// for the expected number of items and filter size
func CalculateOptimalHashFunctions(filterSize uint, expectedItems int) int {
	// k = (m/n) * ln(2)
	// where k is the number of hash functions, m is the filter size in bits, and n is the expected number of elements
	if expectedItems <= 0 || filterSize == 0 {
		return 3 // Default value
	}
	k := float64(filterSize) / float64(expectedItems) * math.Ln2
	hashCount := int(math.Ceil(k))

	// Ensure we have at least 1 hash function and not too many
	hashCount = max(1, hashCount)
	hashCount = min(hashCount, 20) // Cap at 20 hash functions

	return hashCount
}
