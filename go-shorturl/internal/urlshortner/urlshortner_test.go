package urlshortner

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

type MockStore struct {
	StoreMap map[string]string
}

func (s *MockStore) GetObject(key string) (string, error) {
	if val := s.StoreMap[key]; val != "" {
		return val, nil
	}
	return "", fmt.Errorf("object not found in the db")
}

func (s *MockStore) SetObject(key string, val string) error {
	if keyVal := s.StoreMap[key]; keyVal != "" {
		fmt.Println("rewriting the key object")
	}
	s.StoreMap[key] = val
	return nil
}

func TestShortenUrl(t *testing.T) {
	urlShortnerSvc := NewUrlShortner(&MockStore{
		StoreMap: make(map[string]string),
	})
	originalRandomStringFunc := _randomStringGenerator
	defer func() {
		_randomStringGenerator = originalRandomStringFunc
	}()
	_randomStringGenerator = func(len int) (string, error) {
		return base64.URLEncoding.EncodeToString([]byte("mellow")), nil
	}

	input := "https://google.com"
	expectedOutput := "https://bWVsbG93.agh"
	receivedOutput, err := urlShortnerSvc.ShortenUrl(input)
	if err != nil {
		t.Fatal("failed to shorten url", err)
	}

	if strings.Compare(expectedOutput, receivedOutput) != 0 {
		t.Errorf("expected and received out put differ expected:%s recieved:%s", expectedOutput, receivedOutput)
	}
}

func TestRetrievingUrl(t *testing.T) {
	urlShortnerSvc := NewUrlShortner(&MockStore{
		StoreMap: make(map[string]string),
	})
	originalRandomStringFunc := _randomStringGenerator
	defer func() {
		_randomStringGenerator = originalRandomStringFunc
	}()
	_randomStringGenerator = func(len int) (string, error) {
		return base64.URLEncoding.EncodeToString([]byte("mellow")), nil
	}

	input := "https://google.com"
	expectedOutput := input
	shortenedUrl, err := urlShortnerSvc.ShortenUrl(input)
	if err != nil {
		t.Fatal("failed to shorten url", err)
	}
	receivedOutput, err := urlShortnerSvc.GetFullUrl(shortenedUrl)

	if err != nil {
		t.Fatal("failed to retrieve url", err)
	}

	if strings.Compare(expectedOutput, receivedOutput) != 0 {
		t.Errorf("expected and received out put differ expected:%s recieved:%s", expectedOutput, receivedOutput)
	}
}
