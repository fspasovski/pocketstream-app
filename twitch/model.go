package twitch

type TopChannelEdgeImageResultDto struct {
	Edge              *TopChannelsEdgeGqlResponse
	Bytes             []byte
	PreviewImageBytes []byte
	Err               error
}

type SearchChannelsEdgeImageResultDto struct {
	Edge              *SearchStreamsEdgeGqlResponse
	Bytes             []byte
	PreviewImageBytes []byte
	Err               error
}

// GqlQueryRequest is used when sending a full GQL query string directly,
// bypassing the persisted query mechanism. This is more resilient to Twitch
// rotating their persisted query hashes or changing operation variable schemas.
type GqlQueryRequest struct {
	Query string `json:"query"`
}

// GqlPersistedRequest is used when sending a persisted (hashed) GQL operation.
type GqlPersistedRequest struct {
	OperationName string                `json:"operationName"`
	Variables     *GqlRequestVariables  `json:"variables"`
	Extensions    *GqlRequestExtensions `json:"extensions"`
}

type GqlRequestVariables struct {
	IsLive      bool   `json:"isLive"`
	Login       string `json:"login"`
	PlayerType  string `json:"playerType"`
	IncludeIsDJ bool   `json:"includeIsDJ"`
	Query       string `json:"query"`
}

type GqlRequestExtensions struct {
	PersistedQuery *GqlRequestPersistedQuery `json:"persistedQuery"`
}

type GqlRequestPersistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

type StreamingUrlGqlResponse struct {
	Data *StreamingUrlGqlResponseData `json:"data"`
}

type StreamingUrlGqlResponseData struct {
	StreamPlaybackAccessToken *StreamPlaybackAccessToken `json:"streamPlaybackAccessToken"`
}

type StreamPlaybackAccessToken struct {
	Value     string `json:"value"`
	Signature string `json:"signature"`
}

type TopChannelsGqlResponse struct {
	Data *TopChannelsDataGqlResponse `json:"data"`
}

type TopChannelsDataGqlResponse struct {
	Streams *TopChannelsStreamsGqlResponse `json:"streams"`
}

type TopChannelsStreamsGqlResponse struct {
	Edges []*TopChannelsEdgeGqlResponse `json:"edges"`
}

type TopChannelsEdgeGqlResponse struct {
	Node *TopChannelsNodeGqlResponse `json:"node"`
}

type TopChannelsNodeGqlResponse struct {
	Id              string                  `json:"id"`
	Title           string                  `json:"title"`
	ViewersCount    int                     `json:"viewersCount"`
	PreviewImageURL string                  `json:"previewImageURL"`
	Broadcaster     *BroadcasterGqlResponse `json:"broadcaster"`
}

type BroadcasterGqlResponse struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"displayName"`
	ProfileImageURL string `json:"profileImageURL"`
}

type SearchStreamsGqlResponse struct {
	Data *SearchStreamsDataGqlResponse `json:"data"`
}

type SearchStreamsDataGqlResponse struct {
	SearchFor *SearchStreamsSearchForGqlResponse `json:"searchFor"`
}

type SearchStreamsSearchForGqlResponse struct {
	Channels *SearchStreamsChannelsGqlResponse `json:"channels"`
}

type SearchStreamsChannelsGqlResponse struct {
	Edges []*SearchStreamsEdgeGqlResponse `json:"edges"`
}

type SearchStreamsEdgeGqlResponse struct {
	Item *SearchStreamsItemGqlResponse `json:"item"`
}

type SearchStreamsItemGqlResponse struct {
	Id                string                              `json:"id"`
	Login             string                              `json:"login"`
	DisplayName       string                              `json:"displayName"`
	ProfileImageURL   string                              `json:"profileImageURL"`
	Stream            *SearchStreamsStreamGqlResponse     `json:"stream"`
	BroadcastSettings *SearchStreamsItemBroadcastSettings `json:"broadcastSettings"`
}

type SearchStreamsItemBroadcastSettings struct {
	Title string
}

type SearchStreamsStreamGqlResponse struct {
	Id              string `json:"id"`
	Title           string `json:"title"`
	Type            string `json:"type"`
	ViewersCount    int    `json:"viewersCount"`
	PreviewImageURL string `json:"previewImageUrl"`
}
