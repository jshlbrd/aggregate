package aggregate

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestJsonSize(t *testing.T) {
	var tests = []struct {
		data     interface{}
		expected int
	}{
		{
			struct {
				Foo string `json:"foo"`
				Baz int    `json:"baz"`
				Qux bool   `json:"qux"`
			}{
				Foo: "bar",
				Baz: 1,
				Qux: true,
			},
			32,
		},
	}

	for _, test := range tests {
		s, err := jsonSize(test.data)
		if err != nil {
			if err != nil {
				t.Logf("%v", err)
				t.Fail()
			}
		}

		if s != test.expected {
			t.Logf("expected %v, got %v", test.expected, s)
			t.Fail()
		}
	}
}

func TestLenSize(t *testing.T) {
	var tests = []struct {
		data     interface{}
		expected int
	}{
		{
			`foo`,
			3,
		},
		{
			[]byte(`foo`),
			3,
		},
	}

	for _, test := range tests {
		s, err := lenSize(test.data)
		if err != nil {
			if err != nil {
				t.Logf("%v", err)
				t.Fail()
			}
		}

		if s != test.expected {
			t.Logf("expected %v, got %v", test.expected, s)
			t.Fail()
		}
	}
}

func TestAdd(t *testing.T) {
	var tests = []struct {
		data     []interface{}
		maxSize  int
		expected int
	}{
		{
			// encoded JSON == 32
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
			1,
		},
	}

	for _, test := range tests {
		agg := Aggregate{}
		agg.New(test.maxSize, test.maxSize, "json")

		for _, data := range test.data {
			agg.Add(data)
		}

		if agg.count != test.expected {
			t.Logf("expected %v, got %v", test.expected, agg.count)
			t.Fail()
		}
	}
}

func TestGet(t *testing.T) {
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
		agg := Aggregate{}
		agg.New(100, 2, "json")

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

func TestCount(t *testing.T) {
	var tests = []struct {
		data     []interface{}
		expected int
	}{
		{
			[]interface{}{
				"foo",
				"bar",
			},
			2,
		},
	}

	for _, test := range tests {
		agg := Aggregate{}
		agg.New(1000, 1000, "json")

		for _, data := range test.data {
			agg.Add(data)
		}

		if agg.Count() != test.expected {
			t.Logf("expected %v, got %v", test.expected, agg.Count())
			t.Fail()
		}
	}
}

func TestSize(t *testing.T) {
	var tests = []struct {
		data     []interface{}
		expected int
	}{
		{
			// encoded JSON == 32
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
					Baz: 1,
					Qux: true,
				},
			},
			64,
		},
	}

	for _, test := range tests {
		agg := Aggregate{}
		agg.New(1000, len(test.data), "json")

		for _, data := range test.data {
			agg.Add(data)
		}

		if agg.Size() != test.expected {
			t.Logf("expected %v, got %v", test.expected, agg.Size())
			t.Fail()
		}
	}
}

func benchmark(b *testing.B, data interface{}, sizeType string) {
	agg := Aggregate{}
	agg.New(10000, 10000, sizeType)

	for i := 0; i < b.N; i++ {
		agg.Add(data)
	}
}

func BenchmarkSizeType(b *testing.B) {
	var tests = []struct {
		name string
		data []interface{}
	}{
		{
			"bytes",
			[]interface{}{
				[]byte("foo"),
				[]byte("bar"),
			},
		},
		{
			"string",
			[]interface{}{
				"foo",
				"bar",
			},
		},
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
				benchmark(b, test.data, test.name)
			},
		)
	}
}
