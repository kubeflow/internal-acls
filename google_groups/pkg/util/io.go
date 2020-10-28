package util

import (
	log "github.com/sirupsen/logrus"
	"io"
)

// MaybeClose will close the writer if its a Closer.
// Intended to be used with calls to defer.
func MaybeClose(writer io.Writer) {
	if closer, isCloser := writer.(io.Closer); isCloser {
		err := closer.Close()

		if err != nil {
			log.Errorf("Error closing writer; error %v", err)
		}
	}
}