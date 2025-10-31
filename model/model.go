package model

type Stream struct {
	Id               string
	Title            string
	ViewersCount     int
	PreviewImageURL  string
	Broadcaster      *Broadcaster
	PreviewImageData []byte
}

type Broadcaster struct {
	Id               string
	Login            string
	DisplayName      string
	ProfileImageURL  string
	ProfileImageData []byte
}
