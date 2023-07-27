package services

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"path"
	"sync"

	"github.com/rs/zerolog/log"
)

type FileService struct {
	fileDir string
}

func NewFileService(fileDir string) *FileService {
	return &FileService{
		fileDir: fileDir,
	}
}

func (a *FileService) UploadFiles(fileHeaders []*multipart.FileHeader) error {
	wg := sync.WaitGroup{}
	for _, fh := range fileHeaders {
		go func(fh *multipart.FileHeader) {
			wg.Add(1)
			defer wg.Done()

			err := a.UploadFile(fh)
			if err != nil {
				log.Err(err).Msg("file service: upload images: could not upload image")
				// TODO: add aggregation of errors
			}
		}(fh)
		wg.Wait()
	}
	return nil
}

func (f *FileService) UploadFile(fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("file service: upload image: could not open uploaded file: %v", err)
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("file service: upload image: could not read uploaded file: %v", err)
	}
	err = ioutil.WriteFile(path.Join(f.fileDir, fileHeader.Filename), fileData, 0644)
	if err != nil {
		return fmt.Errorf("file service: upload image: could not write file: %v", err)
	}
	return nil
}
