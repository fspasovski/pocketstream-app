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
	ClientId            string
	GqlUrl              string
	UsherUrl            string
	StreamResolution    string
	TopStreamsLimit     int
	HttpClient          *http.Client
	SearchResultsSha256 string
}

// gqlErrorEnvelope is used to detect API-level errors in a GQL response body
// before attempting to parse the data payload.
type gqlErrorEnvelope struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// checkGQLErrors returns a non-nil error if the response body contains a
// top-level "errors" array. It does not consume the body bytes.
func checkGQLErrors(body []byte) error {
	var env gqlErrorEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil // not parseable as an error envelope; let the caller handle it
	}
	if len(env.Errors) == 0 {
		return nil
	}
	msgs := make([]string, 0, len(env.Errors))
	for _, e := range env.Errors {
		msgs = append(msgs, e.Message)
	}
	return fmt.Errorf("twitch GQL error: %s", strings.Join(msgs, "; "))
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
		fmt.Println("Error parsing gql response:", e)
		return nil, e
	}

	if err := checkGQLErrors(gqlBody); err != nil {
		fmt.Println("GQL API error:", err)
		return nil, err
	}

	if parsedResponse.Data == nil || parsedResponse.Data.Streams == nil {
		return nil, fmt.Errorf("unexpected empty response from Twitch: data.streams is nil")
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
		fmt.Println("Error parsing gql response:", e)
		return nil, e
	}

	if err := checkGQLErrors(gqlBody); err != nil {
		fmt.Println("GQL API error:", err)
		return nil, err
	}

	if parsedResponse.Data == nil || parsedResponse.Data.SearchFor == nil || parsedResponse.Data.SearchFor.Channels == nil {
		return nil, fmt.Errorf("unexpected empty response from Twitch: data.searchFor.channels is nil")
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
		return
	}
	defer profileImageResp.Body.Close()

	data, err := io.ReadAll(profileImageResp.Body)
	if err != nil {
		results <- SearchChannelsEdgeImageResultDto{Edge: edge, Err: err}
		return
	}

	previewImageResp, err := http.Get(edge.Item.Stream.PreviewImageURL)
	if err != nil {
		results <- SearchChannelsEdgeImageResultDto{Edge: edge, Err: err}
		return
	}
	defer previewImageResp.Body.Close()

	previewImageData, err := io.ReadAll(previewImageResp.Body)
	if err != nil {
		results <- SearchChannelsEdgeImageResultDto{Edge: edge, Err: err}
		return
	}

	results <- SearchChannelsEdgeImageResultDto{Edge: edge, Bytes: data, PreviewImageBytes: previewImageData}
}

func getImageDataFromUrl(edge *TopChannelsEdgeGqlResponse, wg *sync.WaitGroup, results chan<- TopChannelEdgeImageResultDto) {
	defer wg.Done()

	resp, err := http.Get(edge.Node.Broadcaster.ProfileImageURL)
	if err != nil {
		results <- TopChannelEdgeImageResultDto{Edge: edge, Err: err}
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		results <- TopChannelEdgeImageResultDto{Edge: edge, Err: err}
		return
	}

	// previewImageURL contains {width}x{height} placeholders; substitute a
	// concrete size (440×248 ≈ 16:9 at a comfortable resolution for the UI).
	previewURL := strings.NewReplacer("{width}", "440", "{height}", "248").Replace(edge.Node.PreviewImageURL)
	previewImageResp, err := http.Get(previewURL)
	if err != nil {
		results <- TopChannelEdgeImageResultDto{Edge: edge, Err: err}
		return
	}
	defer previewImageResp.Body.Close()

	previewImageData, err := io.ReadAll(previewImageResp.Body)
	if err != nil {
		results <- TopChannelEdgeImageResultDto{Edge: edge, Err: err}
		return
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
	// The PlaybackAccessToken operation is sent as a persisted request with its
	// operation name and variables; Twitch's GQL endpoint resolves it server-side.
	gqlRequest := &GqlPersistedRequest{
		OperationName: "PlaybackAccessToken",
		Variables: &GqlRequestVariables{
			IsLive:     true,
			Login:      channel,
			PlayerType: "embed",
		},
	}

	// Attach the inline query as well; the endpoint accepts both forms and some
	// Twitch GQL paths require it for the PlaybackAccessToken operation.
	type playbackRequest struct {
		OperationName string               `json:"operationName"`
		Query         string               `json:"query"`
		Variables     *GqlRequestVariables `json:"variables"`
	}
	fullRequest := &playbackRequest{
		OperationName: gqlRequest.OperationName,
		Query:         "query PlaybackAccessToken($login: String!, $isLive: Boolean!, $playerType: String!) { streamPlaybackAccessToken(channelName: $login, params: {platform: \"web\", playerBackend: \"mediaplayer\", playerType: $playerType}) @include(if: $isLive) { value signature } }",
		Variables:     gqlRequest.Variables,
	}

	gqlRequestJson, err := json.Marshal(fullRequest)
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
	// Send a full GQL query string instead of a persisted (hashed) operation.
	// Persisted queries require a SHA-256 hash that Twitch rotates when they
	// update their schema; the hash in use also must exactly match the variable
	// shape Twitch currently expects (e.g. they removed recommendations_context
	// from StreamOptions). Sending the query inline avoids both problems and is
	// fully supported by the public endpoint used by Twitch's own web client.
	//
	// The previewImageUrl template field {width}x{height} is filled with a
	// concrete size; browsers use ~852×480 for standard thumbnails.
	// previewImageURL returns a template string with {width}x{height} placeholders
	// that must be substituted client-side; the field takes no GQL arguments.
	query := fmt.Sprintf(`{
		streams(first: %d) {
			edges {
				node {
					id
					title
					viewersCount
					previewImageURL
					broadcaster {
						id
						login
						displayName
						profileImageURL(width: 50)
					}
				}
			}
		}
	}`, limit)

	gqlRequest := &GqlQueryRequest{Query: query}

	gqlRequestJson, err := json.Marshal(gqlRequest)
	if err != nil {
		fmt.Println("Error marshalling GQL request:", err)
		return nil, err
	}

	return bytes.NewBuffer(gqlRequestJson), nil
}

func (s *TwitchService) getSearchChannelsGqlRequest(searchValue *string) (*bytes.Buffer, error) {
	gqlRequest := &GqlPersistedRequest{
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
