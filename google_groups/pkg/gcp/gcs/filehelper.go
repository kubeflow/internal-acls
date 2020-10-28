package gcs

import (
	"io"
)

// TODO(jlewi): We should implement a UnionFileHelper that will delegate to the GcsFileHelper or LocalFileHelper

// FileHelper is an interface intended to transparently handle working with GCS and local files.
// TODO(jlewi): Move into the util package?
type FileHelper interface {
	Exists(path string) (bool, error)
	NewReader(path string) (io.Reader, error)
	NewWriter(path string) (io.Writer, error)
}
