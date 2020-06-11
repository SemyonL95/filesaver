package filestorage

import (
	"fmt"
	"io"
	"os"
)

type LocalStorage struct {
	StoragePath string
}

func NewLocalStorage(storagePath string) *LocalStorage {
	return &LocalStorage{
		StoragePath: storagePath,
	}
}

func (s *LocalStorage) OpenEntry(filename string) (io.ReadCloser, error) {
	filepath := fmt.Sprintf("%s/%s", s.StoragePath, filename)

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *LocalStorage) CreateEntry(filename string, file io.Reader) (io.WriteCloser, error) {
	filepath := fmt.Sprintf("%s/%s", s.StoragePath, filename)

	out, err := os.Create(filepath)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return out, nil
}
