package db

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"sync"
	"time"

	"adv-backend-trainee-assignment/src/models"
	"github.com/google/uuid"
)

type MockedDBManager struct {
	data map[string][]byte
	ctx  context.Context
	sync sync.Mutex
}

func NewMockedDBManager() *MockedDBManager {
	return &MockedDBManager{data: make(map[string][]byte), ctx: context.Background()}
}

func (mock *MockedDBManager) Close() error {
	mock.data = map[string][]byte{}
	return nil
}

func (mock *MockedDBManager) NewAd(adData models.CreatingAd) (string, error) {
	adID := uuid.New().String()
	marshalledPhotoLinks, err := json.Marshal(adData.PhotoLinks)
	if err != nil {
		return "", err
	} else {
		mock.sync.Lock()
		mock.data[adID], err = json.Marshal(map[string]string{
			"ad_id":       adID,
			"title":       adData.Title,
			"description": adData.Description,
			"price":       strconv.FormatInt(adData.Price, 10),
			"photo_links": string(marshalledPhotoLinks),
			"created_at":  strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		})
		mock.sync.Unlock()
		if err != nil {
			return "", err
		}
		return adID, nil
	}
}

func (mock *MockedDBManager) SelectAd(adID string) (*models.DbAd, error) {
	var rawData map[string]string
	if val, ok := mock.data[adID]; !ok {
		return nil, nil
	} else {
		err := json.Unmarshal(val, &rawData)
		if err != nil {
			return nil, err
		} else {
			var data = &models.DbAd{
				AdID:        rawData["ad_id"],
				Title:       rawData["title"],
				Description: rawData["description"],
			}

			data.Price, err = strconv.ParseInt(rawData["price"], 10, 64)
			if err != nil {
				return nil, err
			}
			rawCreatedAt, err := strconv.ParseInt(rawData["created_at"], 10, 64)
			if err != nil {
				return nil, err
			}
			data.CreatedAt = time.Unix(rawCreatedAt/1000000000, rawCreatedAt%1000000000)
			err = json.Unmarshal([]byte(rawData["photo_links"]), &data.PhotoLinks)
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	}
}

func (mock *MockedDBManager) GetAllAds(sortBy string, sortOrder string, page int, perPage int) ([]*models.DbAd, error) {
	offset := (page - 1) * perPage
	limit := len(mock.data) - offset
	if limit < 1 || offset < 0 {
		return []*models.DbAd{}, nil
	}
	raw := make([]*models.DbAd, len(mock.data))
	i := 0
	var rawData map[string]string
	for _, v := range mock.data {
		rawData = map[string]string{}
		err := json.Unmarshal(v, &rawData)
		if err != nil {
			return nil, err
		}
		var data = &models.DbAd{
			AdID:        rawData["ad_id"],
			Title:       rawData["title"],
			Description: rawData["description"],
		}
		data.Price, err = strconv.ParseInt(rawData["price"], 10, 64)
		if err != nil {
			return nil, err
		}
		rawCreatedAt, err := strconv.ParseInt(rawData["created_at"], 10, 64)
		if err != nil {
			return nil, err
		}
		data.CreatedAt = time.Unix(rawCreatedAt/1000000000, rawCreatedAt%1000000000)
		err = json.Unmarshal([]byte(rawData["photo_links"]), &data.PhotoLinks)
		if err != nil {
			return nil, err
		}
		raw[i] = data
		i++
	}
	if sortBy == "price" {
		if sortOrder == "asc" {
			sort.Slice(raw, func(i, j int) bool {
				return raw[i].Price < raw[j].Price
			})
		} else if sortOrder == "desc" {
			sort.Slice(raw, func(i, j int) bool {
				return raw[i].Price > raw[j].Price
			})
		}
	} else if sortBy == "created_at" {
		if sortOrder == "asc" {
			sort.Slice(raw, func(i, j int) bool {
				return raw[i].CreatedAt.UnixNano() < raw[j].CreatedAt.UnixNano()
			})
		} else if sortOrder == "desc" {
			sort.Slice(raw, func(i, j int) bool {
				return raw[i].CreatedAt.UnixNano() > raw[j].CreatedAt.UnixNano()
			})
		}
	}
	if limit > perPage {
		limit = perPage
	}
	return raw[offset : (page-1)*perPage+limit], nil
}
