package ports

import "io"

type PictureStorage interface {
	Save(reader io.Reader, contentType string) (string, error)

	Retrieve(uuid string) (io.Reader, string, error)

	Delete(uuid string) error
}
