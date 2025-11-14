package common

import (
	"io"
	"net/http"
	"sync"
)

type ImageDataService struct {
	cache map[string][]byte
}

type ImageData struct {
	Url  string
	Data []byte
}

func NewImageDataService() *ImageDataService {
	return &ImageDataService{cache: make(map[string][]byte, 0)}
}

func (s *ImageDataService) GetImageData(urls []string) map[string][]byte {
	result := make(map[string][]byte, len(urls))

	for i := 0; i < len(urls); i++ {
		if s.cache[urls[i]] != nil && len(s.cache[urls[i]]) > 0 {
			result[urls[i]] = s.cache[urls[i]]
		}
	}

	if len(urls) == len(result) {
		return result
	}

	var waitGroup sync.WaitGroup
	channel := make(chan ImageData, len(urls)-len(result))

	for _, url := range urls {
		if url == "" || result[url] != nil {
			continue
		}

		waitGroup.Add(1)
		go getImageDataFromUrl(url, &waitGroup, channel)

	}

	waitGroup.Wait()
	close(channel)

	for res := range channel {
		if len(res.Data) == 0 {
			continue
		}

		result[res.Url] = res.Data
		s.cache[res.Url] = res.Data
	}

	return result
}

func getImageDataFromUrl(url string, wg *sync.WaitGroup, results chan<- ImageData) {
	defer wg.Done()

	imageResponse, err := http.Get(url)
	if err != nil {
		results <- ImageData{Url: url, Data: make([]byte, 0)}
	}
	defer imageResponse.Body.Close()

	data, err := io.ReadAll(imageResponse.Body)
	if err != nil {
		results <- ImageData{Url: url, Data: make([]byte, 0)}
	}

	results <- ImageData{Url: url, Data: data}
}
