package terminal

import "fmt"

// RingBuffer implements a fixed-size circular buffer.
// When the buffer is full, new writes overwrite the oldest data.
// For a buffer of size N, only the most recent N bytes are retained.
type RingBuffer struct {
	data        []byte // underlying storage
	size        int64  // maximum capacity of the buffer
	writeCursor int64  // current write position
	written     int64  // total number of bytes ever written
}

// NewRingBuffer creates a new RingBuffer with the given size.
// Size must be greater than 0, otherwise an error is returned.
func NewRingBuffer(size int64) (*RingBuffer, error) {
	if size <= 0 {
		return nil, fmt.Errorf("buffer size must be positive")
	}

	return &RingBuffer{
		size: size,
		data: make([]byte, size),
	}, nil
}

// Write appends data to the buffer.
// If the input is larger than the buffer, only the last `size` bytes are kept.
// Older data may be overwritten.
func (b *RingBuffer) Write(buf []byte) (int, error) {
	n := len(buf)
	b.written += int64(n)

	// If input exceeds capacity, only keep the last portion
	if int64(n) > b.size {
		buf = buf[int64(n)-b.size:]
	}

	// Copy data into buffer, possibly wrapping around
	remaining := b.size - b.writeCursor
	copy(b.data[b.writeCursor:], buf)
	if int64(len(buf)) > remaining {
		copy(b.data, buf[remaining:])
	}

	// Update write cursor
	b.writeCursor = (b.writeCursor + int64(len(buf))) % b.size
	return n, nil
}

// Size returns the capacity of the buffer.
func (b *RingBuffer) Size() int64 {
	return b.size
}

// TotalWritten returns the total number of bytes written
// since the buffer was created or last reset.
func (b *RingBuffer) TotalWritten() int64 {
	return b.written
}

// Bytes returns a slice containing the current buffer contents
// in the correct order. The returned slice must not be modified.
func (b *RingBuffer) Bytes() []byte {
	switch {
	// Buffer completely filled and cursor at start: return directly
	case b.written >= b.size && b.writeCursor == 0:
		return b.data

	// Buffer completely filled, cursor not at start: reconstruct
	case b.written > b.size:
		out := make([]byte, b.size)
		copy(out, b.data[b.writeCursor:])
		copy(out[b.size-b.writeCursor:], b.data[:b.writeCursor])
		return out

	// Buffer not yet full: return slice up to cursor
	default:
		return b.data[:b.writeCursor]
	}
}

// Reset clears the buffer contents and resets counters.
func (b *RingBuffer) Reset() {
	b.writeCursor = 0
	b.written = 0
}

// String returns the contents of the buffer as a string.
func (b *RingBuffer) String() string {
	return string(b.Bytes())
}
