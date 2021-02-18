package models

import "time"

type DbAd struct {
	AdID        string    `json:"ad_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	PhotoLinks  []string  `json:"photo_links"`
	CreatedAt   time.Time `json:"created_at"`
}
