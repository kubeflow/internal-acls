package util

import (
	"github.com/google/go-cmp/cmp"
	"regexp"
	"testing"
)

func Test_FilterByRe(t *testing.T) {
	type testCase struct {
		items    []string
		re       string
		expected []ReMatch
	}

	testCases := []testCase{
		{
			items: []string{
				"gs://a/b/doc-1.pdf",
				"gs://a/b/doc-2.pdf",
				"gs://a/b/skip.txt",
			},
			re: `gs://a/b/(?P<name>.*)\.pdf`,
			expected: []ReMatch{
				{
					Value: "gs://a/b/doc-1.pdf",
					Groups: map[string]string{
						"name": "doc-1",
					},
				},
				{
					Value: "gs://a/b/doc-2.pdf",
					Groups: map[string]string{
						"name": "doc-2",
					},
				},
			},
		},
	}

	for _, c := range testCases {
		p, err := regexp.Compile(c.re)

		if err != nil {
			t.Errorf("Could not compile %v; error %v", c.re, err)
			continue
		}

		l := &ArrayLister{
			c.items,
		}
		matches, err := FilterByRe(l, p)

		if err != nil {
			t.Errorf("FilterByRe failed; error %v", err)
			continue
		}

		if diff := cmp.Diff(c.expected, matches); diff != "" {
			t.Errorf("FilterByRe() mismatch (-want +got):\n%s", diff)
		}
	}
}
