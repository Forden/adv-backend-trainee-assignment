package models

type BasicAd struct {
	AdID          string `json:"adID"`
	Title         string `json:"title"`
	MainPhotoLink string `json:"mainPhotoLink"`
	Price         int64  `json:"price"`
}
