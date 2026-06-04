package urlshortner

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

const (
	urlLen = 6
)

var (
	_randomStringGenerator = GenerateRandomString
)

type Store interface {
	GetObject(string) (string, error)
	SetObject(string, string) error
}

type URlShortner struct {
	store Store
}

func NewUrlShortner(redisStore Store) *URlShortner {
	return &URlShortner{
		store: redisStore,
	}
}

func GenerateRandomString(len int) (string, error) {

	randomBytes := make([]byte, len)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

func (u *URlShortner) ShortenUrl(urlString string) (string, error) {

	urlModified, err := _randomStringGenerator(urlLen)
	if err != nil {
		return "", err
	}
	shortenedUrl := fmt.Sprintf("https://%s.agh", urlModified)
	u.store.SetObject(urlModified, urlString)
	return shortenedUrl, nil
}

func (u *URlShortner) GetFullUrl(shortUrl string) (string, error) {
	urlObj, err := url.Parse(shortUrl)
	if err != nil {
		return "", err
	}
	host := urlObj.Host
	parsedUrl := strings.Split(host, ".")
	if len(parsedUrl) < 2 {
		return "", err
	}
	shortenedUrl := parsedUrl[0]
	if fullUrl, err := u.store.GetObject(shortenedUrl); err == nil {
		return fullUrl, nil
	}
	return "", fmt.Errorf("requested url not found ")
}
