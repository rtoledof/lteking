package uploader

import "io"

type Uploader interface {
	Upload(file io.Reader, path string) error
}
