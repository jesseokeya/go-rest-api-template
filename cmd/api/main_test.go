package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMain tests the servers main method
func TestMain(t *testing.T) {
	tests := map[string]struct {
		input  string
		output int
		err    error
	}{
		"successful conversion": {input: "1", output: 1, err: nil},
		"invalid integer":       {input: "a text", output: 0, err: &strconv.NumError{}},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		output, err := strconv.Atoi(test.input)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.output, output)
	}
}
