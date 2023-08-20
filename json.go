package aggregate

import (
	"encoding/json"
	"time"
)

// InvalidJSON is returned when an invalid JSON object is added to the aggregate.
const InvalidJSON = Error("InvalidJSON")

// JSON is an intermediary structure for storing structs that marshal to valid JSON.
type JSON struct {
	count, maxCount int
	size, maxSize   int
	maxDuration     time.Duration

	now   time.Time
	items []interface{}
}

/*
New initializes a new JSON aggregate with these settings:
	maxCount:
		the maximum number of JSON objects stored in the aggregate; when this value is reached, no more objects can be added to the payload.
	maxSize:
		the maximum size of all JSON objects stored in the aggregate; when this value is reached, no more objects can be added to the payload.
	maxDuration:
		the maximum duration that the aggregate will store JSON objects; when this duration is reached, no more objects can be added to the payload.
*/
func (a *JSON) New(maxCount, maxSize int, maxDuration time.Duration) {
	a.count, a.size = 0, 0
	a.maxCount = maxCount
	a.maxSize = maxSize
	a.maxDuration = maxDuration

	a.now = time.Now()
	a.items = make([]interface{}, 0, a.maxCount)
}

// Reset resets a JSON aggregate to its initialized settings.
func (a *JSON) Reset() {
	a.count, a.size = 0, 0

	a.now = time.Now()
	a.items = a.items[:0]
}

/*
Add adds a JSON object to the aggregate payload, returning true if the add succeeded and false if the add failed. If an invalid JSON object is added, then an error is returned.

If an add attempt fails and the payload is not empty, then the payload should be retrieved (see Get), the aggregate reset (see Reset), and the failed object should be reattempted.

If an add attempt fails and the payload is empty, then the object being added exceeds the configured limits of the aggregate and should not be reattempted.
*/
func (a *JSON) Add(data interface{}) (bool, error) {
	newCount := a.count + 1
	if newCount > a.maxCount {
		return false, nil
	}

	size, err := jsonSize(data)
	if err != nil {
		return false, err
	}

	newSize := a.size + size
	if newSize > a.maxSize {
		return false, nil
	}

	if time.Since(a.now) > a.maxDuration {
		return false, nil
	}

	a.size = newSize
	a.count = newCount

	a.now = time.Now()
	a.items = append(a.items, data)

	return true, nil
}

// Get returns the aggregate payload.
func (a *JSON) Get() []interface{} {
	return a.items
}

// Count returns the number of JSON objects in the aggregate payload.
func (a *JSON) Count() int {
	return a.count
}

// Size returns the total size of the JSON objects in the aggregate payload.
func (a *JSON) Size() int {
	return a.size
}

// size calculates the size of a JSON object. If the attempt to marshal the JSON fails or if the object is not a valid JSON object, then an error is returned.
func jsonSize(v interface{}) (int, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}

	if !json.Valid(b) {
		return 0, InvalidJSON
	}

	return len(b), nil
}
