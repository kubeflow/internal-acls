package gcs

import (
	"github.com/pkg/errors"
	"io"
	"os"
)

type LocalFileHelper struct {}

// NewReader creates a new Reader for local file.
func (h *LocalFileHelper) NewReader(uri string) (io.Reader, error) {
	reader, err := os.Open(uri)

	if err != nil {
		return nil, errors.WithStack(errors.Wrapf(err, "Clould not read: %v", uri))
	}

	return reader, nil
}

// NewWriter creates a new Writer for the local file.
//
// TODO(jlewi): Can we add options to control filemode?
func (h *LocalFileHelper) NewWriter(uri string) (io.Writer, error) {
	_, err :=  os.Stat(uri)

	if err == nil || !os.IsNotExist(err) {
		return nil, errors.WithStack(errors.Errorf("Can't write %v; It already exists", uri ))
	}

	writer, err := os.Create(uri)

	if err != nil {
		return nil, errors.WithStack(errors.Wrapf(err, "Clould not write: %v", uri))
	}

	return writer, nil
}

// Exists checks whether the file exists.
func (h *LocalFileHelper) Exists(uri string) (bool, error) {
	_, err := os.Stat(uri)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}