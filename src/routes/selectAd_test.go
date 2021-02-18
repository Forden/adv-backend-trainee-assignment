package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"adv-backend-trainee-assignment/src/db"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestAPIServer_SelectAd(t *testing.T) {
	type ExtendedAd struct {
		Title         string   `json:"title"`
		Price         int64    `json:"price"`
		MainPhotoLink string   `json:"mainPhotoLink"`
		Description   string   `json:"description,omitempty"`
		PhotoLinks    []string `json:"photoLinks,omitempty"`
	}

	tests := []struct {
		server             APIServer
		name               string
		urlFormat          string
		population         string
		expectedOutputCode int
		expectedOutput     ExtendedAd
	}{
		{
			server:             APIServer{db.NewMockedDBManager()},
			name:               "Get single ad",
			urlFormat:          "/ads/%v",
			population:         `{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
			expectedOutputCode: 200,
			expectedOutput:     ExtendedAd{Title: "title 1", Price: 100, MainPhotoLink: "https://ya.ru"},
		},
		{
			server:             APIServer{db.NewMockedDBManager()},
			name:               "Get single ad with description",
			urlFormat:          "/ads/%v?fields=description",
			population:         `{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
			expectedOutputCode: 200,
			expectedOutput:     ExtendedAd{Title: "title 1", Price: 100, MainPhotoLink: "https://ya.ru", Description: "description 1"},
		},
		{
			server:             APIServer{db.NewMockedDBManager()},
			name:               "Get single ad with description and photo links",
			urlFormat:          "/ads/%v?fields=description,photolinks",
			population:         `{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
			expectedOutputCode: 200,
			expectedOutput:     ExtendedAd{Title: "title 1", Price: 100, MainPhotoLink: "https://ya.ru", Description: "description 1", PhotoLinks: []string{"https://ya.ru", "http://google.com", "https://example.com"}},
		},
		{
			server:             APIServer{db.NewMockedDBManager()},
			name:               "Get single ad with description and photolinks with other parameters order",
			urlFormat:          "/ads/%v?fields=photolinks,description",
			population:         `{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
			expectedOutputCode: 200,
			expectedOutput:     ExtendedAd{Title: "title 1", Price: 100, MainPhotoLink: "https://ya.ru", Description: "description 1", PhotoLinks: []string{"https://ya.ru", "http://google.com", "https://example.com"}},
		},
		{
			server:             APIServer{db.NewMockedDBManager()},
			name:               "Get single ad with photo links",
			urlFormat:          "/ads/%v?fields=photolinks",
			population:         `{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
			expectedOutputCode: 200,
			expectedOutput:     ExtendedAd{Title: "title 1", Price: 100, MainPhotoLink: "https://ya.ru", PhotoLinks: []string{"https://ya.ru", "http://google.com", "https://example.com"}},
		},
		{
			server:             APIServer{db.NewMockedDBManager()},
			name:               "Get single ad with photo links with following comma",
			urlFormat:          "/ads/%v?fields=photolinks,",
			population:         `{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
			expectedOutputCode: 200,
			expectedOutput:     ExtendedAd{Title: "title 1", Price: 100, MainPhotoLink: "https://ya.ru", PhotoLinks: []string{"https://ya.ru", "http://google.com", "https://example.com"}},
		},
		{
			server:             APIServer{db.NewMockedDBManager()},
			name:               "Get single ad with no additional parameters but comma",
			urlFormat:          "/ads/%v?fields=,",
			population:         `{"title":"title 1","description":"description 1","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":100}`,
			expectedOutputCode: 200,
			expectedOutput:     ExtendedAd{Title: "title 1", Price: 100, MainPhotoLink: "https://ya.ru"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := mux.NewRouter()
			router.HandleFunc("/ads/{adID}", tt.server.SelectAd)
			request, err := http.NewRequest(http.MethodPost, "/ad", bytes.NewBuffer([]byte(tt.population)))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.server.NewAd(rr, request)
			assert.Equal(t, tt.expectedOutputCode, rr.Code, fmt.Sprintf("unexpected http code: got %v expected %v", rr.Code, tt.expectedOutputCode))
			var tmpData map[string]interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &tmpData)
			if err != nil {
				t.Errorf("unexpected output: %v", err)
			}
			request, err = http.NewRequest(http.MethodGet, fmt.Sprintf(tt.urlFormat, tmpData["ad_id"]), bytes.NewBuffer([]byte(tt.population)))
			if err != nil {
				t.Fatal(err)
			}
			rr = httptest.NewRecorder()
			router.ServeHTTP(rr, request)
			assert.Equal(t, tt.expectedOutputCode, rr.Code, fmt.Sprintf("unexpected http code: got %v expected %v", rr.Code, tt.expectedOutputCode))
			var tmp ExtendedAd
			err = json.Unmarshal(rr.Body.Bytes(), &tmp)
			if err != nil {
				t.Errorf("unexpected output: %v", err)
			}
			assert.Equal(t, tt.expectedOutput, tmp, fmt.Sprintf("unexpected output: got %v expected %v", tmpData, tt.expectedOutput))
			tt.server.DBManager.Close()
		})
	}
}
