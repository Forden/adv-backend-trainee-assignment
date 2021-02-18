package models

type CreatingAd struct {
	Title       string   `json:"title"`
	Price       int64    `json:"price"`
	Description string   `json:"description"`
	PhotoLinks  []string `json:"photoLinks"`
}

type CreatedAd struct {
	AdID string `json:"ad_id"`
}
