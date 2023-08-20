package aggregate

import "time"

// Strings is an intermediary structure for storing strings.
type Strings struct {
	count, maxCount int
	size, maxSize   int
	maxDuration     time.Duration

	now   time.Time
	items []string
}

/*
New initializes a new Strings aggregate with these settings:
	maxCount:
		the maximum number of strings stored in the aggregate; when this value is reached, no more strings can be added to the payload.
	maxSize:
		the maximum size of all strings stored in the aggregate; when this value is reached, no more strings can be added to the payload.
	maxDuration:
		the maximum duration that the aggregate will store strings; when this duration is reached, no more strings can be added to the payload.
*/
func (a *Strings) New(maxCount, maxSize int, maxDuration string) {
	a.count, a.size = 0, 0
	a.maxCount = maxCount
	a.maxSize = maxSize

	dur, err := time.ParseDuration(maxDuration)
	if err != nil {
		panic(err)
	}
	a.maxDuration = dur

	a.now = time.Now()
	a.items = make([]string, 0, a.maxCount)
}

// Reset resets a Strings aggregate to its initialized settings.
func (a *Strings) Reset() {
	a.count, a.size = 0, 0

	a.now = time.Now()
	a.items = a.items[:0]
}

/*
Add adds a string to the aggregate payload, returning true if the add succeeded and false if the add failed.

If an add attempt fails and the payload is not empty, then the payload should be retrieved (see Get), the aggregate reset (see Reset), and the failed string should be reattempted.

If an add attempt fails and the payload is empty, then the string being added exceeds the configured limits of the aggregate and should not be reattempted.
*/
func (a *Strings) Add(data string) bool {
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
func (a *Strings) Get() []string {
	return a.items
}

// Count returns the number of strings in the aggregate payload.
func (a *Strings) Count() int {
	return a.count
}

// Size returns the total size of the strings in the aggregate payload.
func (a *Strings) Size() int {
	return a.size
}
