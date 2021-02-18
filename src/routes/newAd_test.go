package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"adv-backend-trainee-assignment/src/db"
	"github.com/stretchr/testify/assert"
)

func TestAPIServer_NewAd(t *testing.T) {
	tests := []struct {
		server                 APIServer
		name                   string
		body                   string
		expectedOutputCode     int
		expectedOutputEncoding string
	}{
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Correct insert",
			body:                   `{"title":"title","description":"description","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":50658783}`,
			expectedOutputCode:     http.StatusOK,
			expectedOutputEncoding: "application/json; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Too small title",
			body:                   `{"title":"","description":"desription","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":50658783}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Too big title",
			body:                   `{"title":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","description":"n","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":50658783}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Too small description",
			body:                   `{"title":"title","description":"","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":50658783}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Too big description",
			body:                   `{"title":"title","description":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":50658783}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Too small price",
			body:                   `{"title":"title","description":"description","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":0}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Too small price",
			body:                   `{"title":"title","description":"description","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":-1}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Too small price",
			body:                   `{"title":"title","description":"description","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":-9223372036854775809}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Too big price",
			body:                   `{"title":"title","description":"description","photoLinks":["https://ya.ru","http://google.com","https://example.com"],"price":9223372036854775808}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Empty photo links",
			body:                   `{"title":"title","description":"description","photoLinks":[],"price":50658783}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Too many photo links",
			body:                   `{"title":"title","description":"description","photoLinks":["https://ya.ru","http://google.com","https://example.com", "https://yandex.ru"],"price":50658783}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Bad JSON given",
			body:                   `{"title":"title","description":"description","price":50658783}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Bad JSON given",
			body:                   `{"title":"title","price":50658783}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Bad JSON given",
			body:                   `{"title":"title"}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
		{
			server:                 APIServer{db.NewMockedDBManager()},
			name:                   "Bad JSON given",
			body:                   `{}`,
			expectedOutputCode:     http.StatusBadRequest,
			expectedOutputEncoding: "text/plain; charset=utf-8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.DBManager.Close()
			request, err := http.NewRequest(http.MethodPost, "/ad", bytes.NewBuffer([]byte(tt.body)))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.server.NewAd(rr, request)
			assert.Equal(t, tt.expectedOutputCode, rr.Code, fmt.Sprintf("unexpected http code: got %v expected %v", rr.Code, tt.expectedOutputCode))
			assert.Equal(t, tt.expectedOutputEncoding, rr.Header().Get("Content-Type"), fmt.Sprintf("unexpected http content type: got %v expected %v", rr.Header().Get("Content-Type"), tt.expectedOutputEncoding))
			if tt.expectedOutputCode == http.StatusOK {
				var tmpData map[string]string
				err = json.Unmarshal(rr.Body.Bytes(), &tmpData)
				if adID, ok := tmpData["ad_id"]; ok {
					assert.Regexp(t, regexp.MustCompile(`(?i:[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12})`), adID, "returned not uuid")
				} else {
					t.Fatal("returned not valid json")
				}
			}
		})
	}
}
