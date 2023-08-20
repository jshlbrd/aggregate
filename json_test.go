package aggregate

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

// TestJSONTimeout tests that the timeout is respected by configuring
// a timeout of 1ms and sleeping for 1ms between each Add() call.
func TestJSONTimeout(t *testing.T) {
	var tests = []struct {
		data     []interface{}
		expected int
	}{
		{
			[]interface{}{
				struct {
					Foo string `json:"foo"`
					Baz int    `json:"baz"`
					Qux bool   `json:"qux"`
				}{
					Foo: "bar",
					Baz: 1,
					Qux: true,
				},
				struct {
					Foo string `json:"foo"`
					Baz int    `json:"baz"`
					Qux bool   `json:"qux"`
				}{
					Foo: "bar",
					Baz: 100,
					Qux: false,
				},
			},
			1,
		},
	}

	for _, test := range tests {
		dur, err := time.ParseDuration("1ms")
		if err != nil {
			t.Logf("error parsing duration: %v", err)
			t.Fail()
		}

		agg := JSON{}
		agg.New(100, 100, dur)

		for _, data := range test.data {
			agg.Add(data)
			time.Sleep(1 * time.Millisecond)
		}

		if agg.Count() != test.expected {
			t.Logf("expected %v, got %v", test.expected, agg.Count())
			t.Fail()
		}
	}
}
func TestJSONCount(t *testing.T) {
	var tests = []struct {
		data     []interface{}
		expected int
	}{
		{
			[]interface{}{
				struct {
					Foo string `json:"foo"`
					Baz int    `json:"baz"`
					Qux bool   `json:"qux"`
				}{
					Foo: "bar",
					Baz: 1,
					Qux: true,
				},
				struct {
					Foo string `json:"foo"`
					Baz int    `json:"baz"`
					Qux bool   `json:"qux"`
				}{
					Foo: "bar",
					Baz: 100,
					Qux: false,
				},
			},
			2,
		},
	}

	for _, test := range tests {
		dur, err := time.ParseDuration("1ms")
		if err != nil {
			t.Logf("error parsing duration: %v", err)
			t.Fail()
		}

		agg := JSON{}
		agg.New(100, 100, dur)
		for _, data := range test.data {
			agg.Add(data)
		}

		if agg.Count() != test.expected {
			t.Logf("expected %v, got %v", test.expected, agg.Count())
			t.Fail()
		}
	}
}

func TestJSONSize(t *testing.T) {
	var tests = []struct {
		data     []interface{}
		expected int
	}{
		{
			[]interface{}{
				struct {
					Foo string `json:"foo"`
					Baz int    `json:"baz"`
					Qux bool   `json:"qux"`
				}{
					Foo: "bar",
					Baz: 1,
					Qux: true,
				},
			},
			32,
		},
	}

	for _, test := range tests {
		dur, err := time.ParseDuration("1ms")
		if err != nil {
			t.Logf("error parsing duration: %v", err)
			t.Fail()
		}

		agg := JSON{}
		agg.New(100, 100, dur)
		for _, data := range test.data {
			agg.Add(data)
		}

		if agg.Size() != test.expected {
			t.Logf("expected %v, got %v", test.expected, agg.count)
			t.Fail()
		}
	}
}

func TestJSONGet(t *testing.T) {
	var tests = []struct {
		data     []interface{}
		expected [][]byte
	}{
		{
			[]interface{}{
				struct {
					Foo string `json:"foo"`
					Baz int    `json:"baz"`
					Qux bool   `json:"qux"`
				}{
					Foo: "bar",
					Baz: 1,
					Qux: true,
				},
				struct {
					Foo string `json:"foo"`
					Baz int    `json:"baz"`
					Qux bool   `json:"qux"`
				}{
					Foo: "bar",
					Baz: 100,
					Qux: false,
				},
			},
			[][]byte{
				[]byte(`{"foo":"bar","baz":1,"qux":true}`),
				[]byte(`{"foo":"bar","baz":100,"qux":false}`),
			},
		},
	}

	for _, test := range tests {
		dur, err := time.ParseDuration("1ms")
		if err != nil {
			t.Logf("error parsing duration: %v", err)
			t.Fail()
		}

		agg := JSON{}
		agg.New(100, 100, dur)
		for _, data := range test.data {
			agg.Add(data)
		}

		payload := agg.Get()
		for i, p := range payload {
			b, _ := json.Marshal(p)

			if bytes.Compare(b, test.expected[i]) != 0 {
				t.Logf("expected %v, got %v", string(test.expected[i]), string(b))
				t.Fail()
			}
		}
	}
}

func TestJSONReset(t *testing.T) {
	var tests = []struct {
		data     []interface{}
		expected int
	}{
		{
			[]interface{}{
				struct {
					Foo string `json:"foo"`
					Baz int    `json:"baz"`
					Qux bool   `json:"qux"`
				}{
					Foo: "bar",
					Baz: 1,
					Qux: true,
				},
				struct {
					Foo string `json:"foo"`
					Baz int    `json:"baz"`
					Qux bool   `json:"qux"`
				}{
					Foo: "bar",
					Baz: 100,
					Qux: false,
				},
			},
			0,
		},
	}

	for _, test := range tests {
		dur, err := time.ParseDuration("1ms")
		if err != nil {
			t.Logf("error parsing duration: %v", err)
			t.Fail()
		}

		agg := JSON{}
		agg.New(100, 100, dur)
		for _, data := range test.data {
			agg.Add(data)
		}

		agg.Reset()
		if agg.Size() != test.expected {
			t.Logf("expected %v, got %v", test.expected, agg.Size())
			t.Fail()
		}
	}
}

func benchmarkJSON(b *testing.B, data interface{}) {
	dur, _ := time.ParseDuration("1ms")

	agg := JSON{}
	agg.New(10000, 10000, dur)

	for i := 0; i < b.N; i++ {
		agg.Add(data)
	}
}

func BenchmarkJSON(b *testing.B) {
	var tests = []struct {
		name string
		data []interface{}
	}{
		{
			"json",
			[]interface{}{
				[]interface{}{
					struct {
						Foo string `json:"foo"`
						Baz string `json:"baz"`
					}{
						Foo: "bar",
						Baz: "qux",
					},
					struct {
						Foo string `json:"foo"`
						Baz string `json:"baz"`
					}{
						Foo: "bar",
						Baz: "qux",
					},
				},
			},
		},
	}

	for _, test := range tests {
		b.Run(string(test.name),
			func(b *testing.B) {
				benchmarkJSON(b, test.data)
			},
		)
	}
}
