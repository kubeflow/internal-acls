package gcs

import (
	"cloud.google.com/go/storage"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/iterator"
	"regexp"
	"testing"
)

func TestParse(t *testing.T) {
	type testCase struct {
		Input string
		ExpectedErrRe string
		Expected *GcsPath
	}


	cases := []testCase{
		{
			Input: "gs://bucket/folder1/file.csv",
			ExpectedErrRe: "",
			Expected: &GcsPath{
				Bucket: "bucket",
				Path:   "folder1/file.csv",
			},
		},
		{
			Input: "gs://bucket",
			ExpectedErrRe: "",
			Expected: &GcsPath{
				Bucket: "bucket",
				Path:   "",
			},
		},
		{
			Input: "gs://bucket/",
			ExpectedErrRe: "",
			Expected: &GcsPath{
				Bucket: "bucket",
				Path:   "",
			},
		},
		{
			Input: "/some/path",
			ExpectedErrRe: ".*path.*doesn't.*match",
		},
	}


	for i, c := range cases {
		actual, err := Parse(c.Input)

		if c.ExpectedErrRe != "" {
			if err == nil {
				t.Errorf("Case %v: Expected error %v but no error returned", i, c.ExpectedErrRe)
				continue
			}

			if b, _ :=regexp.MatchString(c.ExpectedErrRe, err.Error()); !b {
				t.Errorf("Case %v; Got error %v; want error matching %v", i, err, c.ExpectedErrRe)
			}
		}

		if err != nil && c.ExpectedErrRe == "" {
			t.Errorf("Case %v: Parse gave unexpected error: %v", i, err)
			continue
		}

		if d := cmp.Diff(c.Expected, actual); d != "" {
			t.Errorf("Case %v: Parse() mismatch (-want +got):\n%s", i, d)
			continue
		}
	}
}


type FakeObjectIterator struct {
	results []string
	pos int
}

func (i *FakeObjectIterator) Next()(*storage.ObjectAttrs, error) {
	if i.pos >= len(i.results) {
		return nil, iterator.Done
	}

	p, err := Parse(i.results[i.pos])

	if err != nil {
		return nil, err
	}

	a := &storage.ObjectAttrs{
		Bucket: p.Bucket,
		Name: p.Path,
	}
	i.pos = i.pos + 1
	return a, nil
}

func TestFindMatches(t *testing.T) {
	type testCase struct {
		Input         string
		Results       []string
		Expected      []string
	}

	testCases := []testCase{
		{
			Input: "gs://mybucket/dirA/contract*.pdf",
			Results: []string {
				"gs://mybucket/dirA/contract-1.pdf",
				"gs://mybucket/dirA/contract-2.pdf",
				"gs://mybucket/dirA/contract-2.csv",
				"gs://mybucket/dirA/other-1.pdf",
			},
			Expected: []string {
				"gs://mybucket/dirA/contract-1.pdf",
				"gs://mybucket/dirA/contract-2.pdf",
			},
		},
	}

	for i, c := range testCases {
		o := &FakeObjectIterator{
			results: c.Results,
			pos:     0,
		}

		pattern, err := Parse(c.Input)

		if err != nil {
			t.Errorf("Could not parse %v; error %v", c.Input, err)
		}

		actual, err := findMatches(pattern, o)

		if err != nil && err != iterator.Done {
			t.Errorf("findMatches returned error %v", err)
			continue
		}

		if d := cmp.Diff(c.Expected, actual); d != "" {
			t.Errorf("Case %v: Parse() mismatch (-want +got):\n%s", i, d)
			continue
		}
	}
}