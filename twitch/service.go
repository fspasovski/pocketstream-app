package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/fspasovski/pocketstream-app/model"
)

type TwitchConfig struct {
	ClientId               string
	GqlUrl                 string
	UsherUrl               string
	StreamResolution       string
	TopStreamsLimit        int
	HttpClient             *http.Client
	BrowsPagePopularSha256 string
	SearchResultsSha256    string
}

type TwitchService struct {
	Config TwitchConfig
}

func (s *TwitchService) GetTopStreams() ([]model.Stream, error) {
	gqlRequest, err := s.getTopChannelsGqlRequest(s.Config.TopStreamsLimit)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", s.Config.GqlUrl, gqlRequest)
	if err != nil {
		fmt.Println("Error executing gql request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", s.Config.ClientId)

	gqlResponse, err := s.Config.HttpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer gqlResponse.Body.Close()

	gqlBody, err := io.ReadAll(gqlResponse.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return nil, err
	}

	var parsedResponse TopChannelsGqlResponse
	e := json.Unmarshal(gqlBody, &parsedResponse)
	if e != nil {
		fmt.Println("Error parsing gql response:", err)
		return nil, e
	}

	topStreams := make([]model.Stream, 0, len(parsedResponse.Data.Streams.Edges))

	var imagesFetchWaitGroup sync.WaitGroup
	edgesWithImageData := make(chan TopChannelEdgeImageResultDto, len(parsedResponse.Data.Streams.Edges))

	for _, edge := range parsedResponse.Data.Streams.Edges {
		imagesFetchWaitGroup.Add(1)
		go getImageDataFromUrl(edge, &imagesFetchWaitGroup, edgesWithImageData)
	}

	imagesFetchWaitGroup.Wait()
	close(edgesWithImageData)

	for res := range edgesWithImageData {
		if res.Err != nil {
			fmt.Println("failed to fetch:", res.Edge.Node.Broadcaster.ProfileImageURL, "err:", res.Err)
			continue
		}

		topStreams = append(topStreams, model.Stream{
			Id:               res.Edge.Node.Id,
			Title:            res.Edge.Node.Title,
			ViewersCount:     res.Edge.Node.ViewersCount,
			PreviewImageURL:  res.Edge.Node.PreviewImageURL,
			PreviewImageData: res.PreviewImageBytes,
			Broadcaster: &model.Broadcaster{
				Id:               res.Edge.Node.Broadcaster.Id,
				Login:            res.Edge.Node.Broadcaster.Login,
				DisplayName:      res.Edge.Node.Broadcaster.DisplayName,
				ProfileImageURL:  res.Edge.Node.Broadcaster.ProfileImageURL,
				ProfileImageData: res.Bytes,
			},
		})
	}

	return topStreams, nil
}

func (s *TwitchService) SearchStreams(searchValue string) ([]model.Stream, error) {
	gqlRequest, err := s.getSearchChannelsGqlRequest(&searchValue)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", s.Config.GqlUrl, gqlRequest)
	if err != nil {
		fmt.Println("Error executing gql request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", s.Config.ClientId)

	gqlResponse, err := s.Config.HttpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer gqlResponse.Body.Close()

	gqlBody, err := io.ReadAll(gqlResponse.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return nil, err
	}

	var parsedResponse SearchStreamsGqlResponse
	e := json.Unmarshal(gqlBody, &parsedResponse)
	if e != nil {
		fmt.Println("Error parsing gql response:", err)
		return nil, e
	}

	streams := make([]model.Stream, 0)

	var imagesFetchWaitGroup sync.WaitGroup
	edgesWithImageData := make(chan SearchChannelsEdgeImageResultDto, len(parsedResponse.Data.SearchFor.Channels.Edges))

	for _, edge := range parsedResponse.Data.SearchFor.Channels.Edges {
		if edge.Item.Stream != nil && edge.Item.Stream.Type == "live" {
			imagesFetchWaitGroup.Add(1)
			go getSearchChannelsImageDataFromUrl(edge, &imagesFetchWaitGroup, edgesWithImageData)
		}
	}

	imagesFetchWaitGroup.Wait()
	close(edgesWithImageData)

	for res := range edgesWithImageData {
		if res.Err != nil {
			continue
		}

		streams = append(streams, model.Stream{
			Id:               res.Edge.Item.Stream.Id,
			Title:            res.Edge.Item.BroadcastSettings.Title,
			ViewersCount:     res.Edge.Item.Stream.ViewersCount,
			PreviewImageURL:  res.Edge.Item.Stream.PreviewImageURL,
			PreviewImageData: res.PreviewImageBytes,
			Broadcaster: &model.Broadcaster{
				Id:               res.Edge.Item.Id,
				Login:            res.Edge.Item.Login,
				DisplayName:      res.Edge.Item.DisplayName,
				ProfileImageURL:  res.Edge.Item.ProfileImageURL,
				ProfileImageData: res.Bytes,
			},
		})
	}

	return streams, nil
}

func getSearchChannelsImageDataFromUrl(edge *SearchStreamsEdgeGqlResponse, wg *sync.WaitGroup, results chan<- SearchChannelsEdgeImageResultDto) {
	defer wg.Done()

	profileImageResp, err := http.Get(edge.Item.ProfileImageURL)
	if err != nil {
		results <- SearchChannelsEdgeImageResultDto{Edge: edge, Err: err}
	}
	defer profileImageResp.Body.Close()

	data, err := io.ReadAll(profileImageResp.Body)
	if err != nil {
		results <- SearchChannelsEdgeImageResultDto{Edge: edge, Err: err}
	}

	previewImageResp, err := http.Get(edge.Item.Stream.PreviewImageURL)
	if err != nil {
		results <- SearchChannelsEdgeImageResultDto{Edge: edge, Err: err}
	}
	defer profileImageResp.Body.Close()

	previewImageData, err := io.ReadAll(previewImageResp.Body)
	if err != nil {
		results <- SearchChannelsEdgeImageResultDto{Edge: edge, Err: err}
	}

	results <- SearchChannelsEdgeImageResultDto{Edge: edge, Bytes: data, PreviewImageBytes: previewImageData}
}

func getImageDataFromUrl(edge *TopChannelsEdgeGqlResponse, wg *sync.WaitGroup, results chan<- TopChannelEdgeImageResultDto) {
	defer wg.Done()

	resp, err := http.Get(edge.Node.Broadcaster.ProfileImageURL)
	if err != nil {
		results <- TopChannelEdgeImageResultDto{Edge: edge, Err: err}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		results <- TopChannelEdgeImageResultDto{Edge: edge, Err: err}
	}

	previewImageResp, err := http.Get(edge.Node.PreviewImageURL)
	if err != nil {
		results <- TopChannelEdgeImageResultDto{Edge: edge, Err: err}
	}
	defer resp.Body.Close()

	previewImageData, err := io.ReadAll(previewImageResp.Body)
	if err != nil {
		results <- TopChannelEdgeImageResultDto{Edge: edge, Err: err}
	}

	results <- TopChannelEdgeImageResultDto{Edge: edge, Bytes: data, PreviewImageBytes: previewImageData}
}

func (s *TwitchService) GetStreamingUrl(channel string) (string, error) {
	gqlResponse, err := s.getStreamingUrlGqlResponse(channel)
	if err != nil {
		return "", err
	}
	usherResponse, nil := s.getUsherResponse(channel, gqlResponse)

	return s.parseStreamUrlFromUsherResponse(usherResponse), nil
}

func (s *TwitchService) getStreamingUrlGqlResponse(channel string) (*StreamingUrlGqlResponse, error) {
	gqlRequest, err := getStreamingLinkGqlRequest(channel)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", s.Config.GqlUrl, gqlRequest)
	if err != nil {
		fmt.Println("Error executing gql request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", s.Config.ClientId)

	gqlResponse, err := s.Config.HttpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer gqlResponse.Body.Close()

	gqlBody, err := io.ReadAll(gqlResponse.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return nil, err
	}

	var parsedResponse StreamingUrlGqlResponse
	e := json.Unmarshal(gqlBody, &parsedResponse)
	if e != nil {
		fmt.Println("Error parsing gql response:", err)
		return nil, e
	}

	return &parsedResponse, nil
}

func getStreamingLinkGqlRequest(channel string) (*bytes.Buffer, error) {
	gqlRequest := &GqlRequest{
		OperationName: "PlaybackAccessToken",
		Query:         "query PlaybackAccessToken($login: String!, $isLive: Boolean!, $playerType: String!) { streamPlaybackAccessToken(channelName: $login, params: {platform: \"web\", playerBackend: \"mediaplayer\", playerType: $playerType}) @include(if: $isLive) { value signature } }",
		Variables: &GqlRequestVariables{
			IsLive:     true,
			Login:      channel,
			PlayerType: "embed",
		},
	}

	gqlRequestJson, err := json.Marshal(gqlRequest)

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return bytes.NewBuffer(gqlRequestJson), nil
}

func (s *TwitchService) getUsherResponse(channel string, gqlResponse *StreamingUrlGqlResponse) (string, error) {
	encodedToken := url.QueryEscape(gqlResponse.Data.StreamPlaybackAccessToken.Value)
	requestUrl := fmt.Sprintf("%s/%s.m3u8?sig=%s&token=%s", s.Config.UsherUrl, strings.ToLower(channel), gqlResponse.Data.StreamPlaybackAccessToken.Signature, encodedToken)

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return "", err
	}

	resp, err := s.Config.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func (s *TwitchService) parseStreamUrlFromUsherResponse(response string) string {
	rows := strings.Split(response, "\n")
	for i, row := range rows {
		if strings.Contains(row, s.Config.StreamResolution) {
			return rows[i+1]
		}
	}

	return ""
}

func (s *TwitchService) getTopChannelsGqlRequest(limit int) (*bytes.Buffer, error) {
	gqlRequest := &GqlRequest{
		OperationName: "BrowsePage_Popular",
		Variables: &GqlRequestVariables{
			ImageWidth:   50,
			Limit:        limit,
			PlatformType: "all",
			Options: &GqlRequestOptions{
				IncludeRestricted: []string{"SUB_ONLY_LIVE"},
				Sort:              "VIEWER_COUNT",
				FreeformTags:      nil,
				Tags:              []string{},
				RecommendationsContext: &GqlRequestRecommendationsContext{
					Platform: "web",
				},
				RequestId:            "fst",
				BroadcasterLanguages: []string{},
			},
			SortTypeIsRecency: false,
			IncludeIsDJ:       true,
		},
		Extensions: &GqlRequestExtensions{
			PersistedQuery: &GqlRequestPersistedQuery{
				Version:    1,
				Sha256Hash: s.Config.BrowsPagePopularSha256,
			},
		},
	}

	gqlRequestJson, err := json.Marshal(gqlRequest)

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return bytes.NewBuffer(gqlRequestJson), nil
}

func (s *TwitchService) getSearchChannelsGqlRequest(searchValue *string) (*bytes.Buffer, error) {
	gqlRequest := &GqlRequest{
		OperationName: "SearchResultsPage_SearchResults",
		Variables: &GqlRequestVariables{
			Query:       *searchValue,
			IncludeIsDJ: true,
		},
		Extensions: &GqlRequestExtensions{
			PersistedQuery: &GqlRequestPersistedQuery{
				Version:    1,
				Sha256Hash: s.Config.SearchResultsSha256,
			},
		},
	}

	gqlRequestJson, err := json.Marshal(gqlRequest)

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return bytes.NewBuffer(gqlRequestJson), nil
}
