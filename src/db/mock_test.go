package db

import (
	"context"
	"reflect"
	"regexp"
	"sync"
	"testing"
	"time"

	"adv-backend-trainee-assignment/src/models"
	"github.com/stretchr/testify/assert"
)

func TestMockedDBManager_GetAllAds(t *testing.T) {
	type args struct {
		sortBy    string
		sortOrder string
		page      int
		perPage   int
	}
	tests := []struct {
		name           string
		data           map[string][]byte
		args           args
		expectedOutput []*models.DbAd
		wantErr        bool
	}{
		{
			name: "Price descending all",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args: args{
				sortBy:    "price",
				sortOrder: "desc",
				page:      1,
				perPage:   10,
			},
			expectedOutput: []*models.DbAd{
				{AdID: "024e410d-f65d-470c-9920-7ddc69447ca5", Title: "title 3", Description: "description 3", Price: 120, PhotoLinks: []string{"https://example.com"}, CreatedAt: time.Unix(1257894000000000000/1000000000, 1257894000000000000%1000000000)},
				{AdID: "22e88a53-3c80-429d-9e84-99d217788098", Title: "title 1", Description: "description 1", Price: 100, PhotoLinks: []string{"https://ya.ru", "https://google.com", "https://example.com"}, CreatedAt: time.Unix(1257892000000000000/1000000000, 1257892000000000000%1000000000)},
				{AdID: "155d4a0d-52a7-42b0-a2a4-f58f4c953dc1", Title: "title 2", Description: "description 2", Price: 15, PhotoLinks: []string{"https://google.com", "https://ya.ru", "https://example.com"}, CreatedAt: time.Unix(1257893000000000000/1000000000, 1257893000000000000%1000000000)},
			},
			wantErr: false,
		},
		{
			name: "Price ascending all",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args: args{
				sortBy:    "price",
				sortOrder: "asc",
				page:      1,
				perPage:   10,
			},
			expectedOutput: []*models.DbAd{
				{AdID: "155d4a0d-52a7-42b0-a2a4-f58f4c953dc1", Title: "title 2", Description: "description 2", Price: 15, PhotoLinks: []string{"https://google.com", "https://ya.ru", "https://example.com"}, CreatedAt: time.Unix(1257893000000000000/1000000000, 1257893000000000000%1000000000)},
				{AdID: "22e88a53-3c80-429d-9e84-99d217788098", Title: "title 1", Description: "description 1", Price: 100, PhotoLinks: []string{"https://ya.ru", "https://google.com", "https://example.com"}, CreatedAt: time.Unix(1257892000000000000/1000000000, 1257892000000000000%1000000000)},
				{AdID: "024e410d-f65d-470c-9920-7ddc69447ca5", Title: "title 3", Description: "description 3", Price: 120, PhotoLinks: []string{"https://example.com"}, CreatedAt: time.Unix(1257894000000000000/1000000000, 1257894000000000000%1000000000)},
			},
			wantErr: false,
		},
		{
			name: "Price descending first",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args: args{
				sortBy:    "price",
				sortOrder: "desc",
				page:      1,
				perPage:   1,
			},
			expectedOutput: []*models.DbAd{
				{AdID: "024e410d-f65d-470c-9920-7ddc69447ca5", Title: "title 3", Description: "description 3", Price: 120, PhotoLinks: []string{"https://example.com"}, CreatedAt: time.Unix(1257894000000000000/1000000000, 1257894000000000000%1000000000)},
			},
			wantErr: false,
		},
		{
			name: "Price ascending first",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args: args{
				sortBy:    "price",
				sortOrder: "asc",
				page:      1,
				perPage:   1,
			},
			expectedOutput: []*models.DbAd{
				{AdID: "155d4a0d-52a7-42b0-a2a4-f58f4c953dc1", Title: "title 2", Description: "description 2", Price: 15, PhotoLinks: []string{"https://google.com", "https://ya.ru", "https://example.com"}, CreatedAt: time.Unix(1257893000000000000/1000000000, 1257893000000000000%1000000000)},
			},
			wantErr: false,
		},
		{
			name: "Price descending last",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args: args{
				sortBy:    "price",
				sortOrder: "desc",
				page:      3,
				perPage:   1,
			},
			expectedOutput: []*models.DbAd{
				{AdID: "155d4a0d-52a7-42b0-a2a4-f58f4c953dc1", Title: "title 2", Description: "description 2", Price: 15, PhotoLinks: []string{"https://google.com", "https://ya.ru", "https://example.com"}, CreatedAt: time.Unix(1257893000000000000/1000000000, 1257893000000000000%1000000000)},
			},
			wantErr: false,
		},
		{
			name: "Price ascending last",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args: args{
				sortBy:    "price",
				sortOrder: "asc",
				page:      3,
				perPage:   1,
			},
			expectedOutput: []*models.DbAd{
				{AdID: "024e410d-f65d-470c-9920-7ddc69447ca5", Title: "title 3", Description: "description 3", Price: 120, PhotoLinks: []string{"https://example.com"}, CreatedAt: time.Unix(1257894000000000000/1000000000, 1257894000000000000%1000000000)},
			},
			wantErr: false,
		},
		{
			name: "Creation date descending all",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args: args{
				sortBy:    "created_at",
				sortOrder: "desc",
				page:      1,
				perPage:   10,
			},
			expectedOutput: []*models.DbAd{
				{AdID: "024e410d-f65d-470c-9920-7ddc69447ca5", Title: "title 3", Description: "description 3", Price: 120, PhotoLinks: []string{"https://example.com"}, CreatedAt: time.Unix(1257894000000000000/1000000000, 1257894000000000000%1000000000)},
				{AdID: "155d4a0d-52a7-42b0-a2a4-f58f4c953dc1", Title: "title 2", Description: "description 2", Price: 15, PhotoLinks: []string{"https://google.com", "https://ya.ru", "https://example.com"}, CreatedAt: time.Unix(1257893000000000000/1000000000, 1257893000000000000%1000000000)},
				{AdID: "22e88a53-3c80-429d-9e84-99d217788098", Title: "title 1", Description: "description 1", Price: 100, PhotoLinks: []string{"https://ya.ru", "https://google.com", "https://example.com"}, CreatedAt: time.Unix(1257892000000000000/1000000000, 1257892000000000000%1000000000)},
			},
			wantErr: false,
		},
		{
			name: "Creation date ascending all",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args: args{
				sortBy:    "created_at",
				sortOrder: "asc",
				page:      1,
				perPage:   10,
			},
			expectedOutput: []*models.DbAd{
				{AdID: "22e88a53-3c80-429d-9e84-99d217788098", Title: "title 1", Description: "description 1", Price: 100, PhotoLinks: []string{"https://ya.ru", "https://google.com", "https://example.com"}, CreatedAt: time.Unix(1257892000000000000/1000000000, 1257892000000000000%1000000000)},
				{AdID: "155d4a0d-52a7-42b0-a2a4-f58f4c953dc1", Title: "title 2", Description: "description 2", Price: 15, PhotoLinks: []string{"https://google.com", "https://ya.ru", "https://example.com"}, CreatedAt: time.Unix(1257893000000000000/1000000000, 1257893000000000000%1000000000)},
				{AdID: "024e410d-f65d-470c-9920-7ddc69447ca5", Title: "title 3", Description: "description 3", Price: 120, PhotoLinks: []string{"https://example.com"}, CreatedAt: time.Unix(1257894000000000000/1000000000, 1257894000000000000%1000000000)},
			},
			wantErr: false,
		},
		{
			name: "Creation date descending first",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args: args{
				sortBy:    "created_at",
				sortOrder: "desc",
				page:      1,
				perPage:   1,
			},
			expectedOutput: []*models.DbAd{
				{AdID: "024e410d-f65d-470c-9920-7ddc69447ca5", Title: "title 3", Description: "description 3", Price: 120, PhotoLinks: []string{"https://example.com"}, CreatedAt: time.Unix(1257894000000000000/1000000000, 1257894000000000000%1000000000)},
			},
			wantErr: false,
		},
		{
			name: "Creation date ascending first",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args: args{
				sortBy:    "created_at",
				sortOrder: "asc",
				page:      1,
				perPage:   1,
			},
			expectedOutput: []*models.DbAd{
				{AdID: "22e88a53-3c80-429d-9e84-99d217788098", Title: "title 1", Description: "description 1", Price: 100, PhotoLinks: []string{"https://ya.ru", "https://google.com", "https://example.com"}, CreatedAt: time.Unix(1257892000000000000/1000000000, 1257892000000000000%1000000000)},
			},
			wantErr: false,
		},
		{
			name: "Creation date descending last",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args: args{
				sortBy:    "created_at",
				sortOrder: "desc",
				page:      3,
				perPage:   1,
			},
			expectedOutput: []*models.DbAd{
				{AdID: "22e88a53-3c80-429d-9e84-99d217788098", Title: "title 1", Description: "description 1", Price: 100, PhotoLinks: []string{"https://ya.ru", "https://google.com", "https://example.com"}, CreatedAt: time.Unix(1257892000000000000/1000000000, 1257892000000000000%1000000000)},
			},
			wantErr: false,
		},
		{
			name: "Creation date ascending last",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args: args{
				sortBy:    "created_at",
				sortOrder: "asc",
				page:      3,
				perPage:   1,
			},
			expectedOutput: []*models.DbAd{
				{AdID: "024e410d-f65d-470c-9920-7ddc69447ca5", Title: "title 3", Description: "description 3", Price: 120, PhotoLinks: []string{"https://example.com"}, CreatedAt: time.Unix(1257894000000000000/1000000000, 1257894000000000000%1000000000)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := MockedDBManager{
				data: tt.data,
				ctx:  context.Background(),
				sync: sync.Mutex{},
			}
			got, err := db.GetAllAds(tt.args.sortBy, tt.args.sortOrder, tt.args.page, tt.args.perPage)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllAds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedOutput) {
				t.Errorf("GetAllAds() got = %v, expectedOutput %v", got, tt.expectedOutput)
			}
		})
	}
}

func TestMockedDBManager_NewAd(t *testing.T) {
	type args struct {
		adData models.CreatingAd
	}
	tests := []struct {
		name    string
		db      *MockedDBManager
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Creating new ad",
			db:      NewMockedDBManager(),
			args:    args{models.CreatingAd{Title: "title 1", Price: 100, Description: "description 1", PhotoLinks: []string{"https://ya.ru"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.NewAd(tt.args.adData)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Regexp(t, regexp.MustCompile(`(?i:[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12})`), got, "returned not uuid")
		})
	}
}

func TestMockedDBManager_SelectAd(t *testing.T) {
	type args struct {
		adID string
	}
	tests := []struct {
		name           string
		data           map[string][]byte
		args           args
		expectedOutput *models.DbAd
		wantErr        bool
	}{
		{
			name: "Select ad",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args:           args{adID: "22e88a53-3c80-429d-9e84-99d217788098"},
			expectedOutput: &models.DbAd{AdID: "22e88a53-3c80-429d-9e84-99d217788098", Title: "title 1", Description: "description 1", Price: 100, PhotoLinks: []string{"https://ya.ru", "https://google.com", "https://example.com"}, CreatedAt: time.Unix(1257892000000000000/1000000000, 1257892000000000000%1000000000)},
			wantErr:        false,
		},
		{
			name: "Select not existing ad",
			data: map[string][]byte{
				"22e88a53-3c80-429d-9e84-99d217788098": []byte(`{"ad_id":"22e88a53-3c80-429d-9e84-99d217788098","title":"title 1","description":"description 1","photo_links":"[\"https://ya.ru\",\"https://google.com\",\"https://example.com\"]","price":"100","created_at":"1257892000000000000"}`),
				"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1": []byte(`{"ad_id":"155d4a0d-52a7-42b0-a2a4-f58f4c953dc1","title":"title 2","description":"description 2","photo_links":"[\"https://google.com\",\"https://ya.ru\",\"https://example.com\"]","price":"15","created_at":"1257893000000000000"}`),
				"024e410d-f65d-470c-9920-7ddc69447ca5": []byte(`{"ad_id":"024e410d-f65d-470c-9920-7ddc69447ca5","title":"title 3","description":"description 3","photo_links":"[\"https://example.com\"]","price":"120","created_at":"1257894000000000000"}`),
			},
			args:           args{adID: "abcd"},
			expectedOutput: nil,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := MockedDBManager{
				data: tt.data,
				ctx:  context.Background(),
				sync: sync.Mutex{},
			}
			got, err := db.SelectAd(tt.args.adID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectAd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedOutput) {
				t.Errorf("SelectAd() got = %v, expectedOutput %v", got, tt.expectedOutput)
			}
		})
	}
}
