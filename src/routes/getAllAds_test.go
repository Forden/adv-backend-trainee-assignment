package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"adv-backend-trainee-assignment/src/db"
	"github.com/stretchr/testify/assert"
)

func TestAPIServer_GetAllAds(t *testing.T) {
	tests := []struct {
		server               APIServer
		name                 string
		url                  string
		population           []string
		expectedOutputCode   int
		expectedOutputLength int
		expectedOutputOrder  []int
	}{
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Get all no filters",
			url:    "/ads?page=1&perPage=10",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 3,
			expectedOutputOrder:  []int{3, 2, 1},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Get all order by createdAt",
			url:    "/ads?page=1&perPage=10&sortBy=createdAt",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 3,
			expectedOutputOrder:  []int{3, 2, 1},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Get all order by createdAt ASC",
			url:    "/ads?page=1&perPage=10&sortBy=createdAt&sortDirection=asc",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 3,
			expectedOutputOrder:  []int{1, 2, 3},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Get all order by createdAt DESC",
			url:    "/ads?page=1&perPage=10&sortBy=createdAt&sortDirection=desc",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 3,
			expectedOutputOrder:  []int{3, 2, 1},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Get all order by price ASC",
			url:    "/ads?page=1&perPage=10&sortBy=price&sortDirection=asc",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 3,
			expectedOutputOrder:  []int{3, 1, 2},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Get all order by price DESC",
			url:    "/ads?page=1&perPage=10&sortBy=price&sortDirection=desc",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 3,
			expectedOutputOrder:  []int{2, 1, 3},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Get 1 per page on 1st page",
			url:    "/ads?page=1&perPage=1",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 1,
			expectedOutputOrder:  []int{3},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Get 1 per page on 2nd page",
			url:    "/ads?page=2&perPage=1",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 1,
			expectedOutputOrder:  []int{2},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Exceed pages",
			url:    "/ads?page=3&perPage=2",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 0,
			expectedOutputOrder:  []int{},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "No per page parameter",
			url:    "/ads?page=3",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 0,
			expectedOutputOrder:  []int{},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "No page parameter",
			url:    "/ads",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 3,
			expectedOutputOrder:  []int{3, 2, 1},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Page parameter is not int",
			url:    "/ads?page=avb",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 3,
			expectedOutputOrder:  []int{3, 2, 1},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Per page parameter is not int",
			url:    "/ads?perPage=avb",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 3,
			expectedOutputOrder:  []int{3, 2, 1},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Per page parameter is too big",
			url:    "/ads?perPage=1000",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 3,
			expectedOutputOrder:  []int{3, 2, 1},
		},
		{
			server: APIServer{db.NewMockedDBManager()},
			name:   "Per page parameter is too small",
			url:    "/ads?perPage=-1000",
			population: []string{
				`{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
				`{"title":"title 2","description":"description 2","photoLinks":["http://google.com"],"price":123}`,
				`{"title":"title 3","description":"description 3","photoLinks":["http://google.com"],"price":15}`,
			},
			expectedOutputCode:   http.StatusOK,
			expectedOutputLength: 1,
			expectedOutputOrder:  []int{3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.DBManager.Close()
			population := make([]string, len(tt.population))
			i := 0
			for _, insertData := range tt.population {
				request, err := http.NewRequest(http.MethodPost, "/ad", bytes.NewBuffer([]byte(insertData)))
				if err != nil {
					t.Fatal(err)
				}
				rr := httptest.NewRecorder()
				tt.server.NewAd(rr, request)
				if rr.Code != http.StatusOK {
					t.Errorf("unexpected http code: got %v expected %v", rr.Code, tt.expectedOutputCode)
				} else {
					var tmpData map[string]string
					err = json.Unmarshal(rr.Body.Bytes(), &tmpData)
					if err != nil {
						t.Errorf("unexpected output: %v", err)
					}
					population[i] = tmpData["ad_id"]
					i++
				}
				time.Sleep(time.Nanosecond) // without sleep insertions are too fast
			}
			request, err := http.NewRequest(http.MethodGet, strings.ToLower(tt.url), nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.server.GetAllAds(rr, request)
			var tmpData []map[string]string
			err = json.Unmarshal(rr.Body.Bytes(), &tmpData)
			assert.Equal(t, tt.expectedOutputCode, rr.Code, fmt.Sprintf("unexpected http code: got %v expected %v", rr.Code, tt.expectedOutputCode))
			assert.Equal(t, tt.expectedOutputLength, len(tmpData), fmt.Sprintf("wrong output length: got %v expected %v", len(tmpData), tt.expectedOutputLength))
			for x, a := range tt.expectedOutputOrder {
				if population[a-1] != tmpData[x]["adID"] {
					t.Errorf("wrong elements order: got %v expected %v", tmpData[a-1]["adID"], population[a-1])
					break
				}
			}

		})
	}
}
