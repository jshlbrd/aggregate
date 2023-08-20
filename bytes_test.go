package aggregate

import (
	"bytes"
	"testing"
	"time"
)

// TestBytesTimeout tests that the timeout is respected by configuring
// a timeout of 1ms and sleeping for 1ms between each Add() call.
func TestBytesTimeout(t *testing.T) {
	var tests = []struct {
		data     [][]byte
		expected int
	}{
		{
			[][]byte{
				[]byte("foo"),
				[]byte("bar"),
				[]byte("baz"),
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

		agg := Bytes{}
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

func TestBytesCount(t *testing.T) {
	var tests = []struct {
		data     [][]byte
		expected int
	}{
		{
			[][]byte{
				[]byte("foo"),
				[]byte("bar"),
				[]byte("baz"),
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

		agg := Bytes{}
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

func TestBytesSize(t *testing.T) {
	var tests = []struct {
		data     [][]byte
		expected int
	}{
		{
			[][]byte{
				[]byte("foo"),
				[]byte("bar"),
				[]byte("baz"),
			},
			9,
		},
		{
			[][]byte{
				[]byte("foo"),
				[]byte("bar"),
				[]byte("baz"),
				[]byte("qux"),
				[]byte("quux"),
			},
			16,
		},
	}

	for _, test := range tests {
		dur, err := time.ParseDuration("1ms")
		if err != nil {
			t.Logf("error parsing duration: %v", err)
			t.Fail()
		}

		agg := Bytes{}
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

func TestBytesGet(t *testing.T) {
	var tests = []struct {
		data     [][]byte
		expected [][]byte
	}{
		{
			[][]byte{
				[]byte("foo"),
				[]byte("bar"),
				[]byte("baz"),
			},
			[][]byte{
				[]byte("foo"),
				[]byte("bar"),
				[]byte("baz"),
			},
		},
	}

	for _, test := range tests {
		dur, err := time.ParseDuration("1ms")
		if err != nil {
			t.Logf("error parsing duration: %v", err)
			t.Fail()
		}

		agg := Bytes{}
		agg.New(100, 100, dur)
		for _, data := range test.data {
			agg.Add(data)
		}

		payload := agg.Get()
		for i, p := range payload {
			if bytes.Compare(p, test.expected[i]) != 0 {
				t.Logf("expected %v, got %v", string(test.expected[i]), string(p))
				t.Fail()
			}
		}
	}
}

func TestBytesReset(t *testing.T) {
	var tests = []struct {
		data     [][]byte
		expected int
	}{
		{
			[][]byte{
				[]byte("foo"),
				[]byte("bar"),
				[]byte("baz"),
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

		agg := Bytes{}
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

func benchmarkBytes(b *testing.B, data []byte) {
	dur, _ := time.ParseDuration("1ms")

	agg := Bytes{}
	agg.New(10000, 10000, dur)

	for i := 0; i < b.N; i++ {
		agg.Add(data)
	}
}

func BenchmarkBytes(b *testing.B) {
	var tests = []struct {
		name string
		data []byte
	}{
		{
			"bytes",
			[]byte("foo"),
		},
	}

	for _, test := range tests {
		b.Run(string(test.name),
			func(b *testing.B) {
				benchmarkBytes(b, test.data)
			},
		)
	}
}
