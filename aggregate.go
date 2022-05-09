package aggregate

import (
	"encoding/json"
)

// InvalidSizeFunction is returned when an invalid size function is configured.
const InvalidSizeFunction = Error("InvalidSizeFunction")

// InvalidLenType is returned when an invalid interface type is processed by the len size function.
const InvalidLenType = Error("InvalidLenType")

// Aggregate is an intermediary structure for storing groups of interfaces (payload) and can be used to buffer data by count and size.
type Aggregate struct {
	count    int
	maxCount int
	size     int
	maxSize  int
	sizeType string
	payload  []interface{}
}

/*
New initializes a new Aggregate with these settings:
	maxSize:
		the maximum size of all items stored in the Aggregate; when this value is reached, no more items can be added to the payload
	maxCount:
		the maximum number of items stored in the Aggregate; when this value is reached, no more items can be added to the payload
	sizeType:
		the function used to calculate the size of items added to the Aggregate
		must be one of:
			gob (default)
			json
			unsafe
*/
func (a *Aggregate) New(maxSize, maxCount int, sizeType string) {
	a.count, a.size = 0, 0
	a.maxCount = maxCount
	a.maxSize = maxSize
	a.sizeType = sizeType

	slice := make([]interface{}, 0, a.maxCount)
	a.payload = slice
}

// Reset resets an Aggregate to its initialized settings.
func (a *Aggregate) Reset() {
	a.count, a.size = 0, 0

	slice := make([]interface{}, 0, a.maxCount)
	a.payload = slice
}

/*
Add adds data to the Aggregate payload, returning true if the add succeeded and false if the add failed.

If an add attempt fails and the payload is not empty, then the payload should be retrieved (see Get), the Aggregate reset (see Reset), and the failed item should be reattempted.

If an add attempt fails and the payload is empty, then the item being added exceeds the configured limits of the Aggregate and should not be reattempted.
*/
func (a *Aggregate) Add(data interface{}) (bool, error) {
	newCount := a.count + 1
	if newCount > a.maxCount {
		return false, nil
	}

	var size int
	switch s := a.sizeType; s {
	case "json":
		s, err := jsonSize(data)
		if err != nil {
			return false, err
		}
		size = s
	case "len":
		s, err := lenSize(data)
		if err != nil {
			return false, err
		}
		size = s
	default:
		return false, InvalidSizeFunction
	}

	newSize := a.size + size
	if newSize > a.maxSize {
		return false, nil
	}

	a.payload = append(a.payload, data)
	a.size = newSize
	a.count = newCount

	return true, nil
}

// Get returns the Aggregate payload.
func (a *Aggregate) Get() []interface{} {
	return a.payload
}

// Count returns the number of items in the Aggregate payload.
func (a *Aggregate) Count() int {
	return a.count
}

// Size returns the total size of the Aggregate payload.
func (a *Aggregate) Size() int {
	return a.size
}

// jsonSize calculates the size of an interface using the json package.
func jsonSize(v interface{}) (int, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

// lenSize calculates the size of an interface using length.
func lenSize(v interface{}) (int, error) {
	switch s := v.(type) {
	case string:
		return len(s), nil
	case []byte:
		return len(s), nil
	default:
		return 0, InvalidLenType
	}
}
