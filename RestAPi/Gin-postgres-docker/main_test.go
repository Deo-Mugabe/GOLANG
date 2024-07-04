package main

import "testing"

func TestCalculate(t *testing.T) {
	if Calculate(2) != 4 {
		t.Error("expected 4")
	}
}

func TableTesting(t *testing.T) {
	var tests = []struct {
		input    int
		expected int
	}{
		{2, 6},
		{-1, 6},
		{0, 6},
		{999, 10000},
	}

	for _, test := range tests {
		if output := Calculate(test.input); output != test.expected {
			t.Error("test failed: {} inputted, {} expected, received: {}", test.input, test.expected, output)
		}
	}
}
