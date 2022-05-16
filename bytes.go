package aggregate

// Bytes is an intermediary structure for storing bytes.
type Bytes struct {
	count, maxCount int
	size, maxSize   int
	items           [][]byte
}

/*
New initializes a new Bytes aggregate with these settings:
	maxSize:
		the maximum size of all strings stored in the aggregate; when this value is reached, no more strings can be added to the payload.
	maxCount:
		the maximum number of strings stored in the aggregate; when this value is reached, no more strings can be added to the payload.
*/
func (a *Bytes) New(maxSize, maxCount int) {
	a.count, a.size = 0, 0
	a.maxCount = maxCount
	a.maxSize = maxSize

	slice := make([][]byte, 0, a.maxCount)
	a.items = slice
}

// Reset resets a Bytes aggregate to its initialized settings.
func (a *Bytes) Reset() {
	a.count, a.size = 0, 0
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

	a.items = append(a.items, data)
	a.size = newSize
	a.count = newCount

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
