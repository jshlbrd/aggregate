# go-aggregate
Buffer strings, bytes, and JSON objects in Go!

```go
package main

import (
	"github.com/jshlbrd/go-aggregate"
)

func main() {
	agg := aggregate.Strings{}
	agg.New(2, 10*10)

	// add items to the aggregate until it is full (!ok)
	for _, s := range []string{"foo", "bar", "baz"} {
		ok := agg.Add(s)
		if !ok {
			// retrieve the items, reset the aggregate, and re-add missed item
			_ = agg.Get()
			agg.Reset()
			agg.Add(s)
		}
	}

	// retrieve any remaining items
	if agg.Count() > 0 {
		_ = agg.Get()
	}
}
```
