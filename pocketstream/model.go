package pocketstream

type StreamsResponse struct {
	Data []StreamDto `json:"data"`
}

type StreamDto struct {
	Id              string         `json:"id"`
	Title           string         `json:"title"`
	ViewerCount     int            `json:"viewer_count"`
	PreviewImageURL string         `json:"preview_image_url"`
	Broadcaster     BroadcasterDto `json:"broadcaster"`
}

type BroadcasterDto struct {
	Id    string `json:"id"`
	Login string `json:"login"`
}
