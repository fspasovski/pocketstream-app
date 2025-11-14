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

func (s *Stream) WithImageData(previewImageBytes []byte, broadcasterImageBytes []byte) Stream {
	return Stream{
		Id:               s.Id,
		Title:            s.Title,
		ViewersCount:     s.ViewersCount,
		PreviewImageURL:  s.PreviewImageURL,
		Broadcaster:      s.Broadcaster.WithImageData(broadcasterImageBytes),
		PreviewImageData: previewImageBytes,
	}
}

func (b *Broadcaster) WithImageData(broadcasterImageBytes []byte) *Broadcaster {
	return &Broadcaster{
		Id:               b.Id,
		Login:            b.Login,
		DisplayName:      b.DisplayName,
		ProfileImageURL:  b.ProfileImageURL,
		ProfileImageData: broadcasterImageBytes,
	}
}
