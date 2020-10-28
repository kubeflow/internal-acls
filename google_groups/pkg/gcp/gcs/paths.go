package gcs

import (
	"path"
)

// Dir should have the same semantics as Path.Dir except it should work with URIs
// e.g. file://some and gs:// as well as non URIs
func Dir(uri string) string {
	gcsPath, err := Parse(uri)

	if err != nil {
		return path.Dir(uri)
	}

	objDir := path.Dir(gcsPath.Path)

	p := &GcsPath{Path: objDir,
		Bucket: gcsPath.Bucket,
	}
	return p.ToURI()
}

// Base should have the same semantics as Path.Base except it should work with URIs
// e.g. file://some and gs:// as well as non URIs
func Base(uri string) string {
	gcsPath, err := Parse(uri)

	if err != nil {
		return path.Base(uri)
	}

	return path.Base(gcsPath.Path)
}