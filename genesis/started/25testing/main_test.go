package main

import "testing"

func main() {

}

var tests = []struct {
	name     string
	divident float32
	divisor  float32
	expected float32
	isErr    bool
}{
	{"valid-data", 100.0, 10.0, 10, false},
	{"invalid-data", 100.0, 0.0, 0.0, true},
}

func TestDivision(t *testing.T) {

	for _, tt := range tests {
		got, err := divid(tt.divident, tt.divisor)

		if tt.isErr {
			if err == nil {
				t.Error("expected error but did not get one")
			}
		} else {
			if err != nil {
				t.Error("did not expect an error but got one")
			}
		}

		if got != tt.expected {
			t.Errorf("exected %f but got %f", tt.expected, got)
		}
	}
}
