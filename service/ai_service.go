package service

import (
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	url := "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq"
	requestData := model.AIRequest{
		Inputs: model.Inputs{
			Table: table,
			Query: query,
		},
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-wait-for-model", "true")

	res, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New("failed to get a valid response from the AI model")
	}

	var tapasRes model.TapasResponse
	if err := json.NewDecoder(res.Body).Decode(&tapasRes); err != nil {
		return "", err
	}

	return tapasRes.Answer, nil
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	return model.ChatResponse{}, nil
}

func (s *AIService) AnalyzeFile(table map[string][]string, queries []string, token string) (string, error) {
	results := make([]string, 0, len(queries))

	for _, query := range queries {
		result, err := s.AnalyzeData(table, query, token)
		if err != nil {
			return "", err
		}
		results = append(results, result)
	}

	answer := fmt.Sprintf("From the provided data, here are the Least Electricity: %s and the Most Electricity: %s.", results[0], results[1])

	return answer, nil
}
