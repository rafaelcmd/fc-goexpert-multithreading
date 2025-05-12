package service

import (
	"context"
	"fmt"
	"github.com/rafaelcmd/fc-goexpert-multithreading/internal/api"
	"time"
)

type CheckerService struct {
	client  api.ClientInterface
	apiURLs []string
}

func NewCheckerService(client api.ClientInterface, apiURLs []string) *CheckerService {
	return &CheckerService{
		client:  client,
		apiURLs: apiURLs,
	}
}

func (s *CheckerService) CheckZipCode(zipCode string) (map[string]interface{}, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	type result struct {
		data map[string]interface{}
		url  string
		err  error
	}

	resultCh := make(chan result, len(s.apiURLs))

	for _, baseURL := range s.apiURLs {
		go func(url string) {
			data, err := s.client.FetchZipCodeData(ctx, fmt.Sprintf(url, zipCode))
			if err != nil {
				resultCh <- result{err: err}
				return
			}
			resultCh <- result{data, baseURL, err}
		}(baseURL)
	}

	for i := 0; i < len(s.apiURLs); i++ {
		select {
		case <-ctx.Done():
			return map[string]interface{}{}, "", fmt.Errorf("request timed out: %w", ctx.Err())
		case res := <-resultCh:
			if res.err == nil {
				return res.data, res.url, nil
			}
		}
	}

	return map[string]interface{}{}, "", fmt.Errorf("all requests failed")
}
