package dropwriter

import (
	"io"
	"sync"
	"time"
)

// Clock abstracts time for the bucket limiter.
type Clock func() time.Time

// NewBucket creates a new token bucket limiter with the specified capacity
// and refill rate (tokens per second), using the system real-time clock.
func NewBucket(capacity, refillRatePerSec int64) *Bucket {
	return NewBucketClock(capacity, refillRatePerSec, time.Now)
}

// NewBucketClock creates a new token bucket limiter with the specified
// capacity, refill rate (tokens per second), and a custom clock function.
func NewBucketClock(capacity, refillRatePerSec int64, clock Clock) *Bucket {
	return &Bucket{
		clock:      clock,
		capacity:   capacity,
		refillRate: refillRatePerSec,
	}
}

// Bucket implements a token bucket limiter.
type Bucket struct {
	clock           Clock
	capacity        int64
	refillRate      int64
	mu              sync.Mutex
	availableTokens int64
	lastTick        time.Time
}

// adjustTokens updates the number of available tokens in the bucket
// based on the elapsed time since the last tick.
func (b *Bucket) adjustTokens() {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.clock()
	defer func() {
		b.lastTick = now
	}()

	if b.lastTick.IsZero() {
		// First adjustment ever: fill the bucket to capacity
		b.availableTokens = b.capacity
		return
	}

	// Add tokens based on elapsed time
	b.availableTokens += int64(now.Sub(b.lastTick).Seconds() * float64(b.refillRate))
	if b.availableTokens > b.capacity {
		b.availableTokens = b.capacity
	}
}

// TakeAvailable attempts to remove `req` tokens from the bucket.
func (b *Bucket) TakeAvailable(req int64) int64 {
	b.adjustTokens()

	b.mu.Lock()
	defer b.mu.Unlock()

	grant := req
	if grant > b.availableTokens {
		grant = b.availableTokens
	}
	b.availableTokens -= grant

	return grant
}

// writer is an io.Writer wrapper that applies token bucket rate limiting
// by dropping excess bytes that exceed the available tokens.
type writer struct {
	w      io.Writer
	bucket *Bucket
}

// Write writes len(buf) bytes to the underlying writer, but only up to the
// number of tokens available in the bucket. Excess bytes beyond the limit are silently discarded.
func (w *writer) Write(buf []byte) (n int, err error) {
	grant := w.bucket.TakeAvailable(int64(len(buf)))
	n, err = w.w.Write(buf[:grant])
	if err != nil {
		return
	}

	// Pretend that the entire buffer was written successfully,
	// even though some bytes may have been dropped.
	n = len(buf)

	return
}

// Writer returns a new rate-limited, dropping writer that wraps the given dst.
func Writer(dst io.Writer, b *Bucket) io.Writer {
	return &writer{w: dst, bucket: b}
}
