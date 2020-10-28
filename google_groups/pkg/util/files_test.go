package util

import "testing"
import  "github.com/google/go-cmp/cmp"

func TestTransformFiles(t *testing.T) {
	type testCase struct {
		InputPattern string
		OutputPattern string
		Files []string
		Expected map[string]string
	}

	cases := []testCase {
		{
			InputPattern: `gs://input/(?P<name>.*)\.pdf`,
			OutputPattern: `gs://output/{{.name}}.csv`,
			Files: []string{
			  "gs://input/file1.pdf",
			  "gs://input/file2.pdf",
			},
			Expected: map[string]string {
				"gs://input/file1.pdf": "gs://output/file1.csv",
				"gs://input/file2.pdf": "gs://output/file2.csv",
			},
		},
	}

	for _, c := range cases {
		results, err := TransformFiles(c.Files, c.InputPattern, c.OutputPattern)

		if err != nil {
			t.Errorf("TransformFiles returned error; %v", err)
			continue
		}

		if diff := cmp.Diff(c.Expected, results); diff != "" {
			t.Errorf("TransformFiles() mismatch (-want +got):\n%s", diff)
		}
	}
}
