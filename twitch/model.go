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

type GqlRequest struct {
	OperationName string                `json:"operationName"`
	Query         string                `json:"query"`
	Variables     *GqlRequestVariables  `json:"variables"`
	Extensions    *GqlRequestExtensions `json:"extensions"`
}

type GqlRequestVariables struct {
	IsLive            bool               `json:"isLive"`
	Login             string             `json:"login"`
	PlayerType        string             `json:"playerType"`
	ImageWidth        int                `json:"imageWidth"`
	Limit             int                `json:"limit"`
	PlatformType      string             `json:"platformType"`
	Options           *GqlRequestOptions `json:"options"`
	SortTypeIsRecency bool               `json:"sortTypeIsRecency"`
	IncludeIsDJ       bool               `json:"includeIsDJ"`
	Query             string             `json:"query"`
}

type GqlRequestExtensions struct {
	PersistedQuery *GqlRequestPersistedQuery `json:"persistedQuery"`
}

type GqlRequestOptions struct {
	IncludeRestricted      []string                          `json:"includeRestricted"`
	Sort                   string                            `json:"sort"`
	FreeformTags           *string                           `json:"freeformTags"`
	Tags                   []string                          `json:"tags"`
	RecommendationsContext *GqlRequestRecommendationsContext `json:"recommendations_context"`
	RequestId              string                            `json:"requestID"`
	BroadcasterLanguages   []string                          `json:"broadcasterLanguages"`
}

type GqlRequestRecommendationsContext struct {
	Platform string `json:"platform"`
}

type GqlRequestPersistedQuery struct {
	Version    int
	Sha256Hash string
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
	Data *TopChannelsDataGqlResponse
}

type TopChannelsDataGqlResponse struct {
	Streams *TopChannelsStreamsGqlResponse
}

type TopChannelsStreamsGqlResponse struct {
	Edges []*TopChannelsEdgeGqlResponse
}

type TopChannelsEdgeGqlResponse struct {
	Node *TopChannelsNodeGqlResponse
}

type TopChannelsNodeGqlResponse struct {
	Id              string                  `json:"id"`
	Title           string                  `json:"title"`
	ViewersCount    int                     `json:"viewersCount"`
	PreviewImageURL string                  `json:"previewImageUrl"`
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
