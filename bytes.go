package aggregate

import "time"

// Bytes is an intermediary structure for storing bytes.
type Bytes struct {
	count, maxCount int
	size, maxSize   int
	maxDuration     time.Duration

	now   time.Time
	items [][]byte
}

/*
New initializes a new Bytes aggregate with these settings:
	maxCount:
		the maximum number of strings stored in the aggregate; when this value is reached, no more strings can be added to the payload.
	maxSize:
		the maximum size of all strings stored in the aggregate; when this value is reached, no more strings can be added to the payload.
	maxDuration:
		the maximum duration that the aggregate will store bytes; when this duration is reached, no more bytes can be added to the payload.
*/
func (a *Bytes) New(maxCount, maxSize int, maxDuration time.Duration) {
	a.count, a.size = 0, 0
	a.maxCount = maxCount
	a.maxSize = maxSize
	a.maxDuration = maxDuration

	a.now = time.Now()
	a.items = make([][]byte, 0, a.maxCount)
}

// Reset resets a Bytes aggregate to its initialized settings.
func (a *Bytes) Reset() {
	a.count, a.size = 0, 0

	a.now = time.Now()
	a.items = a.items[:0]
}

/*
Add adds bytes to the aggregate payload, returning true if the add succeeded and false if the add failed.

If an add attempt fails and the payload is not empty, then the payload should be retrieved (see Get), the aggregate reset (see Reset), and the failed bytes should be reattempted.

If an add attempt fails and the payload is empty, then the bytes being added exceed the configured limits of the aggregate and should not be reattempted.
*/
func (a *Bytes) Add(data []byte) bool {
	newCount := a.count + 1
	if newCount > a.maxCount {
		return false
	}

	newSize := a.size + len(data)
	if newSize > a.maxSize {
		return false
	}

	if time.Since(a.now) > a.maxDuration {
		return false
	}

	a.size = newSize
	a.count = newCount

	a.now = time.Now()
	a.items = append(a.items, data)

	return true
}

// Get returns the aggregate payload.
func (a *Bytes) Get() [][]byte {
	return a.items
}

// Count returns the number of strings in the aggregate payload.
func (a *Bytes) Count() int {
	return a.count
}

// Size returns the total size of the strings in the aggregate payload.
func (a *Bytes) Size() int {
	return a.size
}
