package util

import (
	"bytes"
	"crypto/sha256"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"io"
	"os"
	"regexp"
	"sort"
)
import "github.com/pkg/errors"
import "text/template"


// TransformFiles matches the input pattern against the input files. For any matching file the corresponding
// output is generated according to outputPattern.
//
// inputPattern is a regex that should contain named groups using the syntax
// (?P<name>regex)
//
// For example:
// inputPattern: "gs://someBucket/(?p<name>.*)\.pdf"
//
// The outputPattern can reference the named patterns using the syntax {{.name}} e.g.
//
// outputPattern: "gs://outputBucket/{{.name}}.csv"
// regex groups to capture groups. The outputPattern is a go template that uses {{.g1}, {{.g2}}, ..., {{.gn}}
// to refer to the captured groups
func TransformFiles(files []string, inputPattern string, outputPattern string) (map[string]string , error) {
	results := map[string]string{}

	p, err := regexp.Compile(inputPattern)

	if err != nil {
		return nil, errors.WithStack(errors.Wrapf(err, "Error compiling regex: %v", inputPattern))
	}

	t, err := template.New("output").Parse(outputPattern)

	if err != nil {
		return nil, errors.WithStack(errors.Wrapf(err, "Error parsing template: %v", outputPattern))
	}

	l := &ArrayLister{
		files,
	}
	matches, err := FilterByRe(l, p)

	if err != nil {
		return results, errors.WithStack(errors.Wrapf(err, "Error occurred applying pattern %v", inputPattern))
	}

	for _, m := range matches {
		buf := new(bytes.Buffer)
		t.Execute(buf, m.Groups)
		results[m.Value] = buf.String()
	}
	return results, nil
}

// FileLister is an interface intended to transparently handle working with GCS and local files.
type FileLister interface {
	ListByRe(pattern string) ([]ReMatch, error)
}

// FilesHash generates a hash based on the contents of a list of files.
// This is intended to be used to detect when one or more files has changed.
func ContentHash(files[]string) ([]byte, error) {
	log := zapr.NewLogger(zap.L())

	// Sort the files.
	sort.Slice(files[:], func(i, j int) bool {
		return files[i] < files[j]
	})

	hash := sha256.New()

	for _, f := range files {
		input, err := os.Open(f)

		if err != nil {
			log.Error(err,"Could not read file", "file", f)
			return []byte{}, err
		}

		if _, err := io.Copy(hash, input); err != nil {
			log.Error(err, "Error reading file", "file", f)
			return []byte{}, err
		}
	}
	return hash.Sum(nil), nil
}