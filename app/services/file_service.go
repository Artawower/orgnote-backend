package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"moonbrain/app/models"
	"moonbrain/app/repositories"
	"sync"

	"github.com/rs/zerolog/log"
)

type FileStorage interface {
	Upload(fileName string, file io.Reader) error
	CalculateFileSize(fileName ...string) (float64, error)
}

type FileService struct {
	fileStorage    FileStorage
	userRepository *repositories.UserRepository
}

func NewFileService(fileStorage FileStorage, userRepository *repositories.UserRepository) *FileService {
	return &FileService{
		fileStorage:    fileStorage,
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

			file, err := fh.Open()
			defer file.Close()
			if err != nil {
				log.Err(err).Msgf("file service: upload images: could not open uploaded file: %v", fh.Filename)
			}
			err = a.fileStorage.Upload(fh.Filename, file)
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
