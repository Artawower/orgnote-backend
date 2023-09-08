package services

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"moonbrain/app/models"
	"moonbrain/app/repositories"
	"path"
	"sync"

	"github.com/rs/zerolog/log"
)

type FileService struct {
	fileDir        string
	userRepository *repositories.UserRepository
}

func NewFileService(fileDir string, userRepository *repositories.UserRepository) *FileService {
	return &FileService{
		fileDir:        fileDir,
		userRepository: userRepository,
	}
}

func (a *FileService) UploadFiles(user *models.User, fileHeaders []*multipart.FileHeader) error {
	fileNames := []string{}

	wg := sync.WaitGroup{}
	for _, fh := range fileHeaders {
		fileNames = append(fileNames, fh.Filename)
		go func(fh *multipart.FileHeader) {
			wg.Add(1)
			defer wg.Done()

			err := a.UploadFile(fh)
			if err != nil {
				log.Err(err).Msg("file service: upload images: could not upload image")
				// TODO: add aggregation of errors
				return
			}
		}(fh)
		wg.Wait()
	}

	err := a.userRepository.AddFiles(user.ID.Hex(), fileNames)
	if err != nil {
		return fmt.Errorf("file service: upload images: could not add files to user: %v", err)
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
