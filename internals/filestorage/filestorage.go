package filestorage

import (
	"bufio"
	"fmt"
	"io"
	"log"
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

func (s *LocalStorage) Put(filename string, file io.Reader) error {
	filepath := fmt.Sprintf("%s/%s", s.StoragePath, filename)

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (s *LocalStorage) Get(filename string) (<-chan string, error) {
	filepath := fmt.Sprintf("%s/%s", s.StoragePath, filename)

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	outChan := make(chan string)

	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					close(outChan)
					file.Close()
					break
				}
			}

			log.Printf(line)

			outChan <- line
		}
	}()

	return outChan, nil
}
