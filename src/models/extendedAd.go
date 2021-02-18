package models

type ExtendedAd struct {
	Title         string   `json:"title"`
	Price         int64    `json:"price"`
	MainPhotoLink string   `json:"mainPhotoLink"`
	Description   string   `json:"description,omitempty"`
	PhotoLinks    []string `json:"photoLinks,omitempty"`
}
