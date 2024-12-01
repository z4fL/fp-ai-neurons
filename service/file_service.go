package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"fmt"
	"strings"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	filename := "uploaded_data-series.csv"

	if s.Repo.FileExists(filename) {
		content, err := s.Repo.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		parsedData, err := s.ParseCSV(string(content))
		if err != nil {
			return nil, fmt.Errorf("error parsing CSV: %v", err)
		}

		return parsedData, nil
	} else {
		err := s.Repo.SaveFile(filename, []byte(fileContent))
		if err != nil {
			return nil, fmt.Errorf("error saving file: %v", err)
		}

		parsedData, err := s.ParseCSV(fileContent)
		if err != nil {
			return nil, fmt.Errorf("error parsing CSV: %v", err)
		}

		return parsedData, nil
	}
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
		for j, header := range headers {
			parsedData[header] = append(parsedData[header], records[i][j])
		}
	}

	return parsedData, nil
}
