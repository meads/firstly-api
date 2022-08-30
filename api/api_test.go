package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/meads/firstly-api/db"
)

func TestImageAPI(t *testing.T) {

	tests := []struct {
		body              *bytes.Buffer
		callback          func(w *httptest.ResponseRecorder) bool
		method            string
		name              string
		responseCode      int
		route             string
		setupExpectations func(s *db.MockStore)
	}{
		{
			body: func() *bytes.Buffer {
				return bytes.NewBufferString("{\"data\":\"test\"}")
			}(),
			callback: func(w *httptest.ResponseRecorder) bool {
				return true
			},
			method:       http.MethodPost,
			name:         "should create image record when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/image/",
			setupExpectations: func(s *db.MockStore) {
				s.EXPECT().CreateImage(gomock.Any(), "test").Return(
					db.Image{
						ID:      1,
						Data:    "test",
						Created: time.Now().String(),
						Deleted: sql.NullInt32{Int32: 0, Valid: true},
					}, nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			gin.SetMode(gin.TestMode)

			ctrl := gomock.NewController(t)
			store := db.NewMockStore(ctrl)

			test.setupExpectations(store)

			NewFirstlyAPI(store, router, false)
			responseRecorder := httptest.NewRecorder()

			request := httptest.NewRequest(test.method, test.route, test.body)
			router.ServeHTTP(responseRecorder, request)

			assert.Equal(t, test.responseCode, responseRecorder.Code)

			response := db.Image{}
			if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
				t.Errorf("Error decoding response body: %v", err)
			}
		})
	}
}
