package pocketstream

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/fspasovski/pocketstream-app/config"
	"github.com/fspasovski/pocketstream-app/model"
)

type PocketstreamService struct {
	httpClient      http.Client
	apiUrl          string
	thumbnailWidth  string
	thumbnailHeight string
}

func NewPocketstreamService(cfg *config.Config) (service *PocketstreamService) {
	return &PocketstreamService{
		apiUrl:          cfg.PocketstreamApiUrl,
		thumbnailWidth:  strconv.Itoa(int(cfg.UI.StreamsUiConfig.ThumbnailWidth)),
		thumbnailHeight: strconv.Itoa(int(cfg.UI.StreamsUiConfig.ThumbnailHeight)),
		httpClient:      http.Client{Timeout: 10 * time.Second},
	}
}

func (s *PocketstreamService) GetStreams(userLogins []string) []model.Stream {
	result := make([]model.Stream, 0)
	targetUrl, err := url.Parse(s.apiUrl + "/streams")
	if err != nil {
		log.Printf("Error occurred while parsing streams api url: %v, %v", s.apiUrl, err)
		return result
	}

	queryParams := url.Values{"user_login": userLogins}
	queryParams.Add("thumbnail_width", s.thumbnailWidth)
	queryParams.Add("thumbnail_height", s.thumbnailHeight)
	targetUrl.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", targetUrl.String(), nil)

	if err != nil {
		log.Printf("Error occurred while creating streams request for user logins: %v, %v", userLogins, err)
		return result
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("Error occurred while fetching streams for user logins: %v, %v", userLogins, err)
		return result
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading streams response body: %v", err)
		return result
	}

	var streamsResponse StreamsResponse
	err = json.Unmarshal(body, &streamsResponse)
	if err != nil {
		log.Printf("Error unmarshalling json: %v", err)
		return result
	}

	if streamsResponse.Data == nil {
		return result
	}

	for i := 0; i < len(streamsResponse.Data); i++ {
		result = append(result, model.Stream{
			Id:              streamsResponse.Data[i].Id,
			Title:           streamsResponse.Data[i].Title,
			PreviewImageURL: streamsResponse.Data[i].PreviewImageURL,
			ViewersCount:    streamsResponse.Data[i].ViewerCount,
			Broadcaster: &model.Broadcaster{
				Id:    streamsResponse.Data[i].Broadcaster.Id,
				Login: streamsResponse.Data[i].Broadcaster.Login,
			},
		})
	}
	return result
}
