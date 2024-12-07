package localfs

import (
	"io"
	"mmddvg/chapar/pkg/errs"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalPictureStorage struct {
	BasePath string
}

func NewLocalPictureStorage(basePath string) *LocalPictureStorage {
	return &LocalPictureStorage{BasePath: basePath}
}

func (l *LocalPictureStorage) Save(reader io.Reader, contentType string) (string, error) {
	id := uuid.New().String()
	filename := filepath.Join(l.BasePath, id)

	file, err := os.Create(filename)
	if err != nil {
		return "", errs.NewUnexpected(err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", errs.NewUnexpected(err)
	}

	return id, nil
}

func (l *LocalPictureStorage) Retrieve(uuid string) (io.Reader, string, error) {
	filename := filepath.Join(l.BasePath, uuid)

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", errs.NewNotFound("file", uuid)
		}
		return nil, "", errs.NewUnexpected(err)
	}

	return file, "application/octet-stream", nil
}

func (l *LocalPictureStorage) Delete(uuid string) error {
	filename := filepath.Join(l.BasePath, uuid)
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return errs.NewNotFound("file", uuid)
		}
		return errs.NewUnexpected(err)
	}
	return nil
}
