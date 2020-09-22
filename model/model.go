package model

//Storage contains urls.
type Storage struct {
	Item []string `json:"urls"`
}

//ImageURL contains url of image.
type ImageURL struct {
	Urls struct {
		Thumbnail string `json:"full"`
	} `json:"urls"`
}
