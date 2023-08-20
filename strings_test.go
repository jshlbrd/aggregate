package aggregate

import (
	"testing"
	"time"
)

// TestStringsTimeout tests that the timeout is respected by configuring
// a timeout of 1ms and sleeping for 1ms between each Add() call.
func TestStringsTimeout(t *testing.T) {
	var tests = []struct {
		data     []string
		expected int
	}{
		{
			[]string{
				"foo",
				"bar",
				"baz",
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

		agg := Strings{}
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

func TestStringsCount(t *testing.T) {
	var tests = []struct {
		data     []string
		expected int
	}{
		{
			[]string{
				"foo",
				"bar",
				"baz",
			},
			3,
		},
	}

	for _, test := range tests {
		dur, err := time.ParseDuration("1ms")
		if err != nil {
			t.Logf("error parsing duration: %v", err)
			t.Fail()
		}

		agg := Strings{}
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

func TestStringsSize(t *testing.T) {
	var tests = []struct {
		data     []string
		expected int
	}{
		{
			[]string{
				"foo",
				"bar",
				"baz",
			},
			9,
		},
	}

	for _, test := range tests {
		dur, err := time.ParseDuration("1ms")
		if err != nil {
			t.Logf("error parsing duration: %v", err)
			t.Fail()
		}

		agg := Strings{}
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

func TestStringsGet(t *testing.T) {
	var tests = []struct {
		data     []string
		expected []string
	}{
		{
			[]string{
				"foo",
				"bar",
				"baz",
			},
			[]string{
				"foo",
				"bar",
				"baz",
			},
		},
	}

	for _, test := range tests {
		dur, err := time.ParseDuration("1ms")
		if err != nil {
			t.Logf("error parsing duration: %v", err)
			t.Fail()
		}

		agg := Strings{}
		agg.New(100, 100, dur)

		for _, data := range test.data {
			agg.Add(data)
		}

		payload := agg.Get()
		for i, p := range payload {
			if p != test.expected[i] {
				t.Logf("expected %v, got %v", test.expected[i], p)
				t.Fail()
			}
		}
	}
}

func TestStringsReset(t *testing.T) {
	var tests = []struct {
		data     []string
		expected int
	}{
		{
			[]string{
				"foo",
				"bar",
				"baz",
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

		agg := Strings{}
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

func benchmarkStrings(b *testing.B, data string) {
	dur, _ := time.ParseDuration("1ms")
	agg := Strings{}
	agg.New(100, 100, dur)

	for i := 0; i < b.N; i++ {
		agg.Add(data)
	}
}

func BenchmarkStrings(b *testing.B) {
	var tests = []struct {
		name string
		data string
	}{
		{
			"string",
			"foo",
		},
	}

	for _, test := range tests {
		b.Run(string(test.name),
			func(b *testing.B) {
				benchmarkStrings(b, test.data)
			},
		)
	}
}
