package aggregate

// Strings is an intermediary structure for storing strings.
type Strings struct {
	count, maxCount int
	size, maxSize   int
	payload         []string
}

/*
New initializes a new Strings aggregate with these settings:
	maxSize:
		the maximum size of all strings stored in the aggregate; when this value is reached, no more strings can be added to the payload.
	maxCount:
		the maximum number of strings stored in the aggregate; when this value is reached, no more strings can be added to the payload.
*/
func (a *Strings) New(maxSize, maxCount int) {
	a.count, a.size = 0, 0
	a.maxCount = maxCount
	a.maxSize = maxSize

	slice := make([]string, 0, a.maxCount)
	a.payload = slice
}

// Reset resets a Strings aggregate to its initialized settings.
func (a *Strings) Reset() {
	a.count, a.size = 0, 0
	a.payload = a.payload[:0]
}

/*
Add adds a string to the aggregate payload, returning true if the add succeeded and false if the add failed.

If an add attempt fails and the payload is not empty, then the payload should be retrieved (see Get), the aggregate reset (see Reset), and the failed string should be reattempted.

If an add attempt fails and the payload is empty, then the string being added exceeds the configured limits of the aggregate and should not be reattempted.
*/
func (a *Strings) Add(data string) (bool, error) {
	newCount := a.count + 1
	if newCount > a.maxCount {
		return false, nil
	}

	newSize := a.size + len(data)
	if newSize > a.maxSize {
		return false, nil
	}

	a.payload = append(a.payload, data)
	a.size = newSize
	a.count = newCount

	return true, nil
}

// Get returns the aggregate payload.
func (a *Strings) Get() []string {
	return a.payload
}

// Count returns the number of strings in the aggregate payload.
func (a *Strings) Count() int {
	return a.count
}

// Size returns the total size of the strings in the aggregate payload.
func (a *Strings) Size() int {
	return a.size
}
