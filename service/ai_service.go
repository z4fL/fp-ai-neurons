package service

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	return nil, nil // TODO: replace this
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	// TODO: answer here
}
