package uploader

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	Directory string
}

func NewLocalStorage(directory string) *LocalStorage {
	return &LocalStorage{
		Directory: directory,
	}
}

func (s *LocalStorage) Upload(file io.Reader, path string) error {
	outputPath := filepath.Join(s.Directory, path)
	outputDir := filepath.Dir(outputPath)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, file)
	return err
}
