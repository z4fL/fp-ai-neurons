package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	return nil, nil // TODO: replace this
}
