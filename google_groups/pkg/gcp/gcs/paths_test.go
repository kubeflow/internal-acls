package gcs

import (
	"testing"
)

func Test_Dir(t *testing.T) {
	type testCase struct {
		input string
		output string
	}

	testCases := []testCase {
		{
			input: "gs://bucket/dirA/dirB",
			output: "gs://bucket/dirA",
		},
		{
			input: "/bucket/dirA/dirB",
			output: "/bucket/dirA",
		},
	}

	for _, c := range testCases {
		actual := Dir(c.input)

		if actual != c.output {
			t.Errorf("Input: %v; Got %v; Want %v", c.input, actual, c.output)
		}
	}
}

func Test_Base(t *testing.T) {
	type testCase struct {
		input string
		output string
	}

	testCases := []testCase {
		{
			input: "gs://bucket/dirA/dirB",
			output: "dirB",
		},
		{
			input: "/bucket/dirA/dirB",
			output: "dirB",
		},
	}

	for _, c := range testCases {
		actual := Base(c.input)

		if actual != c.output {
			t.Errorf("Input: %v; Got %v; Want %v", c.input, actual, c.output)
		}
	}
}