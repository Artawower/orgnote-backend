package infrastructure

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type FileStorage struct {
	dirPath string
}

func (f *FileStorage) Upload(fileName string, file io.Reader) error {
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("file storage: upload: could not read file: %v", err)
	}

	err = ioutil.WriteFile(path.Join(f.dirPath, fileName), fileData, 0644)
	if err != nil {
		return fmt.Errorf("file storage: upload: could not write file: %v", err)
	}

	return nil
}

// Return file size in mega bytes
func (f *FileStorage) CalculateFileSize(fileName ...string) (float64, error) {
	fileSize := int64(0)

	for _, name := range fileName {
		errCh := make(chan error)
		resCh := make(chan int64)
		f.getFileSize(name, resCh, errCh)

		select {
		case err := <-errCh:
			return 0, fmt.Errorf("file storage: calculate file size: could not get file size: %v", err)
		case res := <-resCh:
			fileSize += res
		}
	}

	return float64(fileSize) / 1024 / 1024, nil
}

func (f *FileStorage) getFileSize(fileName string, resCh chan<- int64, errCh chan<- error) {
	file, err := os.Open(fileName)
	if err != nil {
		errCh <- fmt.Errorf("file storage: get file size: could not open file: %v", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		errCh <- fmt.Errorf("file storage: get file size: could not get file info: %v", err)
	}

	resCh <- fileInfo.Size()
}

func NewFileStorage(dirPath string) *FileStorage {
	return &FileStorage{
		dirPath: dirPath,
	}
}
