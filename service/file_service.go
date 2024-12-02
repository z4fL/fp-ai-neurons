package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	filename := "uploaded_data-series.csv"

	var finalContent string

	if s.Repo.FileExists(filename) {
		existingContent, err := s.Repo.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		// Bandingkan isi file server dengan file yang diupload
		if strings.TrimSpace(fileContent) == strings.TrimSpace(string(existingContent)) {
			finalContent = string(existingContent)
		} else {
			// Jika berbeda, gunakan data baru dan simpan
			finalContent = fileContent
			err := s.Repo.SaveFile(filename, []byte(fileContent))
			if err != nil {
				return nil, fmt.Errorf("error saving file: %v", err)
			}
		}
	} else {
		// Jika file tidak ada, simpan file baru
		err := s.Repo.SaveFile(filename, []byte(fileContent))
		if err != nil {
			return nil, fmt.Errorf("error saving file: %v", err)
		}
		finalContent = fileContent
	}

	// Cek jika file kosong
	if strings.TrimSpace(finalContent) == "" {
		return nil, errors.New("file content is empty")
	}

	parsedData, err := s.ParseCSV(fileContent)
	if err != nil {
		return nil, fmt.Errorf("error parsing CSV: %v", err)
	}

	return parsedData, nil
}

func (s *FileService) ParseCSV(fileContent string) (map[string][]string, error) {
	reader := csv.NewReader(strings.NewReader(fileContent))

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	if len(records) <= 1 {
		return nil, fmt.Errorf("CSV does not contain data")
	}

	parsedData := make(map[string][]string)
	headers := records[0]

	for i := 1; i < len(records); i++ {
		if len(records[i]) != len(headers) {
			return nil, fmt.Errorf("invalid CSV data at line %d", i+1)
		}
		for j, header := range headers {
			parsedData[header] = append(parsedData[header], records[i][j])
		}
	}

	return parsedData, nil
}
